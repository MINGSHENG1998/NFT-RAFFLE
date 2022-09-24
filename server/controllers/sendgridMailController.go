package controllers

import (
	"context"
	"log"
	"net/http"
	"nft-raffle/database"
	"nft-raffle/dto"
	"nft-raffle/enums"
	"nft-raffle/models"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type SendGridController interface {
	SendGridTesting() gin.HandlerFunc
	VerifyVerificationMail() gin.HandlerFunc
}

type sendGridControllerStruct struct{}

var (
	mailCollection *mongo.Collection = database.OpenCollection(database.Client, "mail")
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

		err = database.Client.UseSession(ctx, func(sessionContext mongo.SessionContext) error {
			err := sessionContext.StartTransaction()
			if err != nil {
				return err
			}

			_, err = userCollection.UpdateOne(
				ctx,
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
			_, err = mailCollection.DeleteOne(ctx, bson.D{
				{Key: "email", Value: email},
				{Key: "type", Value: enums.MailVerification.String()},
			})

			if err != nil {
				sessionContext.AbortTransaction(sessionContext)
				return err
			}

			if err = sessionContext.CommitTransaction(sessionContext); err != nil {
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

func (s *sendGridControllerStruct) SendGridTesting() gin.HandlerFunc {
	return func(c *gin.Context) {
		fromName := "yongheng98_testing_sendgrid"
		fromEmail := "yongheng98@hotmail.com"

		tos := []*mail.Email{
			mail.NewEmail("yyhyap98", "yyhyap98@gmail.com"),
		}

		dynamicTemplateData := map[string]string{}
		dynamicTemplateData["Last_Name"] = "YYHYAP98"
		dynamicTemplateData["Verify_Mail_Link"] = "https://www.google.com"

		mailReq := &dto.MailRequest{
			FromName:            fromName,
			FromEmail:           fromEmail,
			MailType:            enums.MailVerification,
			Tos:                 tos,
			DynamicTemplateData: dynamicTemplateData,
		}

		var Body = sendGridMailService.DynamicTemplate(mailReq)
		responseCh, errCh := sendGridMailService.SendMailAsync(Body)

		select {
		case err := <-errCh:
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		case response := <-responseCh:
			c.JSON(http.StatusOK, gin.H{
				"status":  response.StatusCode,
				"body":    response.Body,
				"headers": response.Headers,
			})
			return
		}
	}
}
