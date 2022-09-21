package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"nft-raffle/database"
	"nft-raffle/dto"
	"nft-raffle/enums"
	"nft-raffle/helpers"
	"nft-raffle/models"
	"nft-raffle/services"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthController interface {
	SignUp() gin.HandlerFunc
	Login() gin.HandlerFunc
	RefreshToken() gin.HandlerFunc
}

type authControllerStruct struct{}

var (
	userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")

	tokenHelper          helpers.TokenHelper          = helpers.NewTokenHelper()
	aesEncryptionHelper  helpers.AesEncrptionHelper   = helpers.NewAesEncryptionHelper()
	randomCodeGenerator  helpers.RandomCodeGenerator  = helpers.NewRandomCodeGenerator()
	dataValidationHelper helpers.DataValidationHelper = helpers.NewDataValidationHelper()
	passwordHelper       helpers.PasswordHelper       = helpers.NewPasswordHelper()
	dotEnvHelper         helpers.DotEnvHelper         = helpers.NewDotEnvHelper()

	sendGridMailService     services.SendGridMailService     = services.NewSendGridMailService()
	verificationMailService services.VerificationMailService = services.NewVerificationMailService()

	verifcationCodeExpiration string = dotEnvHelper.GetEnvVariable("VERIFICATION_MAIL_CODE_EXPIRATION")
	fromName                  string = dotEnvHelper.GetEnvVariable("SENDGRID_FROM_NAME")
	fromEmail                 string = dotEnvHelper.GetEnvVariable("SENDGRID_FROM_EMAIL")
	verifcationMailReturnHost string = dotEnvHelper.GetEnvVariable("VERIFICATION_MAIL_RETURN_HOST")
	verifcationMailReturnPort string = dotEnvHelper.GetEnvVariable("VERIFICATION_MAIL_RETURN_PORT")

	validate = validator.New()
)

func NewAuthController() AuthController {
	return &authControllerStruct{}
}

func (a *authControllerStruct) SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User

		if err := c.BindJSON(&user); err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(user)
		if validationErr != nil {
			log.Println(validationErr)
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		if emailValidationErr := dataValidationHelper.IsEmailValid(user.Email); emailValidationErr != nil {
			log.Println(emailValidationErr)
			c.JSON(http.StatusBadRequest, gin.H{"error": emailValidationErr.Error()})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)

		emailCount, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		defer cancel()

		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if emailCount > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user email already exists"})
			return
		}

		hashedPassword, err := passwordHelper.HashPassword(user.Password)

		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		user.Password = hashedPassword

		user.Created_at, err = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while parsing created_at"})
			return
		}

		user.Updated_at, err = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while parsing updated_at"})
			return
		}

		user.Is_email_verified = false

		user.ID = primitive.NewObjectID()
		user.User_id = user.ID.Hex()

		signedToken, signedRefreshToken, err := tokenHelper.GenerateAllTokens(
			user.Email, user.First_name, user.Last_name, user.User_id, user.User_role, user.Is_email_verified)

		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		user.Access_token = signedToken
		user.Refresh_token = signedRefreshToken

		resultInsertionNumber, insertError := userCollection.InsertOne(ctx, user)
		defer cancel()

		if insertError != nil {
			log.Println(insertError)
			c.JSON(http.StatusInternalServerError, gin.H{"error": insertError})
			return
		}

		// mail verification
		randomSixDigits := randomCodeGenerator.GenerateRandomDigits(6)
		verifcationCodeExpirationInt, err := strconv.ParseInt(verifcationCodeExpiration, 10, 64)

		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		expires_at, err := time.Parse(time.RFC3339, time.Now().Local().Add(time.Hour*time.Duration(verifcationCodeExpirationInt)).Format(time.RFC3339))

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
		dynamicTemplateData["Last_Name"] = fmt.Sprintf("%s %s", user.First_name, user.Last_name)
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

		dynamicTemplateData["Verify_Mail_Link"] = fmt.Sprintf(
			"%s/%s/api/test?email=%s&code=%s",
			verifcationMailReturnHost, verifcationMailReturnPort, encryptedEmailValue, encryptedRandomSixDigits,
		)

		mailReq := &dto.MailRequest{
			FromName:            fromName,
			FromEmail:           fromEmail,
			MailType:            enums.MailVerification,
			Tos:                 tos,
			DynamicTemplateData: dynamicTemplateData,
		}

		go sendGridMailService.SendVerificationMailAsync(mailReq)

		ctx2, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		mailCount, err := mailCollection.CountDocuments(
			ctx2,
			bson.D{
				{Key: "email", Value: user.Email},
				{Key: "type", Value: enums.MailVerification.String()},
			},
		)
		defer cancel()

		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while counting mail from mail collection in db"})
			return
		}

		if mailCount > 0 {
			// update current verification mail
			// update mail in db
			mailUpdateError := verificationMailService.UpdateVerificationEmail(user.Email, randomSixDigits, expires_at)

			if mailUpdateError != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":         "error occured while updating verification email in db",
					"error_details": mailUpdateError.Error(),
				})
				return
			}
		} else {
			// create new verification mail
			// insert mail into db
			mailInsertError := verificationMailService.CreateNewVerifcationMail(user.Email, randomSixDigits, expires_at)

			if mailInsertError != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":         "error occured while inserting new verification email into db",
					"error_details": mailInsertError.Error(),
				})
				return
			}
		}

		c.JSON(http.StatusOK, resultInsertionNumber)
	}
}

func (a *authControllerStruct) Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		var foundUser models.User

		if err := c.BindJSON(&user); err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)

		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
		defer cancel()

		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "email or password is incorrect"})
			return
		}

		passwordValidationError := passwordHelper.VerifyPassword(foundUser.Password, user.Password)

		if passwordValidationError != nil {
			log.Println(passwordValidationError)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "email or password is incorrect"})
			return
		}

		if foundUser.Email == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
			return
		}

		signedToken, signedRefreshToken, err := tokenHelper.GenerateAllTokens(
			foundUser.Email, foundUser.First_name, foundUser.Last_name, foundUser.User_id, foundUser.User_role, foundUser.Is_email_verified)

		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		err = tokenHelper.UpdateAllTokens(signedToken, signedRefreshToken, foundUser.User_id)

		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		err = userCollection.FindOne(ctx, bson.M{"user_id": foundUser.User_id}).Decode(&foundUser)
		defer cancel()

		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, foundUser)
	}
}

func (a *authControllerStruct) RefreshToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		var requestBody map[string]interface{}

		jsonData, err := io.ReadAll(c.Request.Body)

		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if err := json.Unmarshal(jsonData, &requestBody); err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		signedRefreshToken, ok := requestBody["refresh_token"].(string)

		if !ok {
			log.Println("invalid refresh token")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid refresh token"})
			return
		}

		claims, err := tokenHelper.ValidateRefreshToken(signedRefreshToken)

		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		uid := claims.Uid
		var foundUser models.User

		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)

		err = userCollection.FindOne(ctx, bson.M{"user_id": uid}).Decode(&foundUser)
		defer cancel()

		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if foundUser.Email == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
			return
		}

		signedToken, signedRefreshToken, err := tokenHelper.GenerateAllTokens(
			foundUser.Email, foundUser.First_name, foundUser.Last_name, foundUser.User_id, foundUser.User_role, foundUser.Is_email_verified)

		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		err = tokenHelper.UpdateAllTokens(signedToken, signedRefreshToken, foundUser.User_id)

		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		err = userCollection.FindOne(ctx, bson.M{"user_id": foundUser.User_id}).Decode(&foundUser)
		defer cancel()

		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, foundUser)
	}
}
