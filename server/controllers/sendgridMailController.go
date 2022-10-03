package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"nft-raffle/dto"
	"nft-raffle/enums"
	"nft-raffle/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SendGridController interface {
	VerifyVerificationMail() gin.HandlerFunc
	SendPasswordResetMail() gin.HandlerFunc
	VerifyPasswordResetMail() gin.HandlerFunc
}

type sendGridControllerStruct struct{}

var (
	mailCollection *mongo.Collection = nftRaffleDb.OpenCollection(nftRaffleDbClient, "mail")

	passwordResetMailCodeExpiration string = dotEnvHelper.GetEnvVariable("PASSWORD_RESET_MAIL_CODE_EXPIRATION")
	passwordResetMailReturnHost     string = dotEnvHelper.GetEnvVariable("PASSWORD_RESET_MAIL_RETURN_HOST")
	passwordResetMailReturnPort     string = dotEnvHelper.GetEnvVariable("PASSWORD_RESET_MAIL_RETURN_PORT")
)

func NewSendGridController() SendGridController {
	return &sendGridControllerStruct{}
}

func (s *sendGridControllerStruct) VerifyVerificationMail() gin.HandlerFunc {
	return func(c *gin.Context) {
		encryptedEmail := c.Query("email")

		if encryptedEmail == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "email query param is not existing in the URL"})
			return
		}

		encryptedCode := c.Query("code")

		if encryptedCode == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "code query param is not existing in the URL"})
			return
		}

		email, err := aesEncryptionHelper.AesGCMDecrypt(encryptedEmail)

		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		code, err := aesEncryptionHelper.AesGCMDecrypt(encryptedCode)

		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var verificationMail models.Mail
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)

		err = mailCollection.FindOne(ctx, bson.D{
			{Key: "email", Value: email},
			{Key: "type", Value: enums.MailVerification.String()},
		}).Decode(&verificationMail)
		defer cancel()

		if err != nil {
			log.Print(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "email is not valid"})
			return
		}

		// check for the code
		if verificationMail.Code != code {
			log.Println("verification code does not match")
			c.JSON(http.StatusBadRequest, gin.H{"error": "verification code does not match"})
			return
		}

		// check for expires_at
		if verificationMail.Expires_at.Unix() < time.Now().Local().Unix() {
			log.Println("verification mail has expired")
			c.JSON(http.StatusBadRequest, gin.H{"error": "verification mail has expired"})
			return
		}

		// update user.Is_email_verified to true
		var updateObj bson.D

		updateObj = append(updateObj, bson.E{Key: "is_email_verified", Value: true})

		Updated_at, err := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while parsing updated_at"})
			return
		}

		updateObj = append(updateObj, bson.E{Key: "updated_at", Value: Updated_at})

		upsert := true
		filter := bson.M{"email": verificationMail.Email}
		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		err = nftRaffleDbClient.UseSession(ctx, func(sessionContext mongo.SessionContext) error {
			err := sessionContext.StartTransaction()
			if err != nil {
				return err
			}

			_, err = userCollection.UpdateOne(
				sessionContext,
				filter,
				bson.D{
					{Key: "$set", Value: updateObj},
				},
				&opt,
			)

			if err != nil {
				sessionContext.AbortTransaction(sessionContext)
				return err
			}

			// now remove verifcation mail from mailCollection
			_, err = mailCollection.DeleteOne(sessionContext, bson.D{
				{Key: "email", Value: email},
				{Key: "type", Value: enums.MailVerification.String()},
			})

			if err != nil {
				sessionContext.AbortTransaction(sessionContext)
				return err
			}

			if err := sessionContext.CommitTransaction(sessionContext); err != nil {
				return err
			}

			return nil
		})
		defer cancel()

		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// generate token and return user
		var user models.User

		err = userCollection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
		defer cancel()

		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		signedToken, signedRefreshToken, err := tokenHelper.GenerateAllTokens(
			user.Email, user.First_name, user.Last_name, user.User_id, user.User_role, user.Is_email_verified)

		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		err = tokenHelper.UpdateAllTokens(signedToken, signedRefreshToken, user.User_id)

		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		err = userCollection.FindOne(ctx, bson.M{"user_id": user.User_id}).Decode(&user)
		defer cancel()

		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, user)
	}
}

func (s *sendGridControllerStruct) SendPasswordResetMail() gin.HandlerFunc {
	return func(c *gin.Context) {
		var mailDto dto.SendPasswordResetMailRequestDto
		var user models.User

		if err := c.BindJSON(&mailDto); err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if emailValidationErr := dataValidationHelper.IsEmailValid(mailDto.Email); emailValidationErr != nil {
			log.Println(emailValidationErr)
			c.JSON(http.StatusBadRequest, gin.H{"error": emailValidationErr.Error()})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)

		err := userCollection.FindOne(ctx, bson.M{"email": mailDto.Email}).Decode(&user)
		defer cancel()

		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "email not valid"})
			return
		}

		randomSixDigits := randomCodeGenerator.GenerateRandomDigits(6)
		passwordResetMailCodeExpirationInt, err := strconv.ParseInt(passwordResetMailCodeExpiration, 10, 64)

		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		expires_at, err := time.Parse(time.RFC3339, time.Now().Local().Add(time.Hour*time.Duration(passwordResetMailCodeExpirationInt)).Format(time.RFC3339))

		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while parsing mail expires_at"})
			return
		}

		// send email
		tos := []*mail.Email{
			// hardcoded for testing
			mail.NewEmail("yyhyap98", "yyhyap98@gmail.com"),
		}

		dynamicTemplateData := map[string]string{}
		dynamicTemplateData["Full_Name"] = fmt.Sprintf("%s %s", user.First_name, user.Last_name)

		encryptedEmailValue, err := aesEncryptionHelper.AesGCMEncrypt(user.Email)

		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while encrypting user email"})
			return
		}

		encryptedRandomSixDigits, err := aesEncryptionHelper.AesGCMEncrypt(randomSixDigits)

		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while encrypting random six digits"})
			return
		}

		dynamicTemplateData["Password_Reset_Mail_Link"] = fmt.Sprintf(
			"%s:%s/api/send-grid/test?email=%s&code=%s",
			passwordResetMailReturnHost, passwordResetMailReturnPort, encryptedEmailValue, encryptedRandomSixDigits,
		)

		mailReq := &dto.MailRequest{
			FromName:            fromName,
			FromEmail:           fromEmail,
			MailType:            enums.PasswordReset,
			Tos:                 tos,
			DynamicTemplateData: dynamicTemplateData,
		}

		go sendGridMailService.SendMail(mailReq)

		mailCount, err := mailCollection.CountDocuments(
			ctx,
			bson.D{
				{Key: "email", Value: user.Email},
				{Key: "type", Value: enums.PasswordReset.String()},
			},
		)
		defer cancel()

		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while counting mail from mail collection in db"})
			return
		}

		if mailCount > 0 {
			// update current password reset mail
			// update mail in db
			mailUpdateError := sendGridMailService.UpdateEmail(enums.PasswordReset, user.Email, randomSixDigits, expires_at)

			if mailUpdateError != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":         "error occured while updating verification email in db",
					"error_details": mailUpdateError.Error(),
				})
				return
			}
		} else {
			// create new password reset mail
			// insert mail into db
			mailInsertError := sendGridMailService.CreateNewMail(enums.PasswordReset, user.Email, randomSixDigits, expires_at)

			if mailInsertError != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":         "error occured while inserting new verification email into db",
					"error_details": mailInsertError.Error(),
				})
				return
			}
		}

		c.Status(http.StatusOK)
	}
}

func (s *sendGridControllerStruct) VerifyPasswordResetMail() gin.HandlerFunc {
	return func(c *gin.Context) {
		encryptedEmail := c.Query("email")

		if encryptedEmail == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "email query param is not existing in the URL"})
			return
		}

		encryptedCode := c.Query("code")

		if encryptedCode == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "code query param is not existing in the URL"})
			return
		}

		email, err := aesEncryptionHelper.AesGCMDecrypt(encryptedEmail)

		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		code, err := aesEncryptionHelper.AesGCMDecrypt(encryptedCode)

		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var passwordResetMail models.Mail
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)

		err = mailCollection.FindOne(ctx, bson.D{
			{Key: "email", Value: email},
			{Key: "type", Value: enums.PasswordReset.String()},
		}).Decode(&passwordResetMail)
		defer cancel()

		if err != nil {
			log.Print(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "email is not valid"})
			return
		}

		// check for the code
		if passwordResetMail.Code != code {
			log.Println("verification code does not match")
			c.JSON(http.StatusBadRequest, gin.H{"error": "verification code does not match"})
			return
		}

		// check for expires_at
		if passwordResetMail.Expires_at.Unix() < time.Now().Local().Unix() {
			log.Println("verification mail has expired")
			c.JSON(http.StatusBadRequest, gin.H{"error": "verification mail has expired"})
			return
		}

		// password reset mail verified
		c.JSON(http.StatusOK, gin.H{
			"email": passwordResetMail.Email,
			"code":  passwordResetMail.Code,
		})
	}
}
