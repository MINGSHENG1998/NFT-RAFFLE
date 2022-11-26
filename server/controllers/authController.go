package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"nft-raffle/database"
	"nft-raffle/dto"
	"nft-raffle/enums"
	"nft-raffle/helpers"
	"nft-raffle/logger"
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
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	AuthController IAuthController = NewAuthController()

	nftRaffleDbClient *mongo.Client                        = database.NftRaffleDbClient
	nftRaffleDb       database.INftRaffleMongoDbConnection = database.NftRaffleMongoDbConnection
	userCollection    *mongo.Collection                    = nftRaffleDb.OpenCollection(nftRaffleDbClient, "user")

	tokenHelper          helpers.ITokenHelper          = helpers.TokenHelper
	aesEncryptionHelper  helpers.IAesEncrptionHelper   = helpers.AesEncryptionHelper
	randomCodeGenerator  helpers.IRandomCodeGenerator  = helpers.RandomCodeGenerator
	dataValidationHelper helpers.IDataValidationHelper = helpers.DataValidationHelper
	passwordHelper       helpers.IPasswordHelper       = helpers.PasswordHelper
	dotEnvHelper         helpers.IDotEnvHelper         = helpers.DotEnvHelper
	timeHelper           helpers.ITimeHelper           = helpers.TimeHelper

	sendGridMailService services.ISendGridMailService = services.SendGridMailService

	verifcationCodeExpiration string = dotEnvHelper.GetEnvVariable("VERIFICATION_MAIL_CODE_EXPIRATION")
	fromName                  string = dotEnvHelper.GetEnvVariable("SENDGRID_FROM_NAME")
	fromEmail                 string = dotEnvHelper.GetEnvVariable("SENDGRID_FROM_EMAIL")
	verifcationMailReturnHost string = dotEnvHelper.GetEnvVariable("VERIFICATION_MAIL_RETURN_HOST")
	verifcationMailReturnPort string = dotEnvHelper.GetEnvVariable("VERIFICATION_MAIL_RETURN_PORT")

	validate = validator.New()
)

type IAuthController interface {
	SignUp(c *gin.Context)
	Login(c *gin.Context)
	RefreshToken(c *gin.Context)
	ResetUserPassword(c *gin.Context)
	TestRedis(c *gin.Context)
}

type authControllerStruct struct{}

func NewAuthController() IAuthController {
	return &authControllerStruct{}
}

func (a *authControllerStruct) SignUp(c *gin.Context) {
	var user models.User

	if err := c.BindJSON(&user); err != nil {
		logger.Logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	validationErr := validate.Struct(user)
	if validationErr != nil {
		logger.Logger.Error(validationErr.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
		return
	}

	if emailValidationErr := dataValidationHelper.IsEmailValid(user.Email); emailValidationErr != nil {
		logger.Logger.Error(emailValidationErr.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": emailValidationErr.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	emailCount, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})

	if err != nil {
		logger.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if emailCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user email already exists"})
		return
	}

	hashedPassword, err := passwordHelper.HashPassword(user.Password)

	if err != nil {
		logger.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user.Password = hashedPassword

	user.Created_at, err = timeHelper.GetCurrentTimeSingapore()

	if err != nil {
		logger.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while parsing created_at"})
		return
	}

	user.Updated_at, err = timeHelper.GetCurrentTimeSingapore()

	if err != nil {
		logger.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while parsing updated_at"})
		return
	}

	user.Is_email_verified = false

	user.ID = primitive.NewObjectID()
	user.User_id = user.ID.Hex()

	signedToken, signedRefreshToken, err := tokenHelper.GenerateAllTokens(
		user.Email, user.First_name, user.Last_name, user.User_id, user.User_role, user.Is_email_verified)

	if err != nil {
		logger.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user.Access_token = signedToken
	user.Refresh_token = signedRefreshToken

	resultInsertionNumber, insertError := userCollection.InsertOne(ctx, user)

	if insertError != nil {
		logger.Logger.Error(insertError.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": insertError})
		return
	}

	// mail verification
	randomSixDigits := randomCodeGenerator.GenerateRandomDigits(6)
	verifcationCodeExpirationInt, err := strconv.ParseInt(verifcationCodeExpiration, 10, 64)

	if err != nil {
		logger.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	expires_at, err := timeHelper.GetCurrentTimeSingaporeWithAdditionalDuration(time.Hour * time.Duration(verifcationCodeExpirationInt))

	if err != nil {
		logger.Logger.Error(err.Error())
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
		logger.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while encrypting user email"})
		return
	}

	encryptedRandomSixDigits, err := aesEncryptionHelper.AesGCMEncrypt(randomSixDigits)

	if err != nil {
		logger.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while encrypting random six digits"})
		return
	}

	dynamicTemplateData["Verify_Mail_Link"] = fmt.Sprintf(
		"%s:%s/api/test?email=%s&code=%s",
		verifcationMailReturnHost, verifcationMailReturnPort, encryptedEmailValue, encryptedRandomSixDigits,
	)

	mailReq := &dto.MailRequest{
		FromName:            fromName,
		FromEmail:           fromEmail,
		MailType:            enums.MailVerification,
		Tos:                 tos,
		DynamicTemplateData: dynamicTemplateData,
	}

	go sendGridMailService.SendMail(mailReq)

	mailCount, err := mailCollection.CountDocuments(
		ctx,
		bson.D{
			{Key: "email", Value: user.Email},
			{Key: "type", Value: enums.MailVerification.String()},
		},
	)

	if err != nil {
		logger.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while counting mail from mail collection in db"})
		return
	}

	if mailCount > 0 {
		// update current verification mail
		// update mail in db
		mailUpdateError := sendGridMailService.UpdateEmail(enums.MailVerification, user.Email, randomSixDigits, expires_at)

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
		mailInsertError := sendGridMailService.CreateNewMail(enums.MailVerification, user.Email, randomSixDigits, expires_at)

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

func (a *authControllerStruct) Login(c *gin.Context) {
	var user models.User
	var foundUser models.User

	if err := c.BindJSON(&user); err != nil {
		logger.Logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)

	if err != nil {
		logger.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "email or password is incorrect"})
		return
	}

	passwordValidationError := passwordHelper.VerifyPassword(foundUser.Password, user.Password)

	if passwordValidationError != nil {
		logger.Logger.Error(passwordValidationError.Error())
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
		logger.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = tokenHelper.UpdateAllTokens(signedToken, signedRefreshToken, foundUser.User_id)

	if err != nil {
		logger.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = userCollection.FindOne(ctx, bson.M{"user_id": foundUser.User_id}).Decode(&foundUser)

	if err != nil {
		logger.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// loc, _ := time.LoadLocation("Asia/Singapore")
	// logger.Logger.Debug(fmt.Sprintf("local date time %v", foundUser.Updated_at.In(loc)))
	c.JSON(http.StatusOK, foundUser)
}

func (a *authControllerStruct) RefreshToken(c *gin.Context) {
	var requestBody map[string]interface{}

	jsonData, err := io.ReadAll(c.Request.Body)

	if err != nil {
		logger.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := json.Unmarshal(jsonData, &requestBody); err != nil {
		logger.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	signedRefreshToken, ok := requestBody["refresh_token"].(string)

	if !ok {
		logger.Logger.Error("invalid refresh token")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid refresh token"})
		return
	}

	claims, err := tokenHelper.ValidateRefreshToken(signedRefreshToken)

	if err != nil {
		logger.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// check JWT blacklist
	uid := claims.Uid
	var foundUser models.User

	blacklistRefreshTokenExpiration, err := tokenHelper.GetBlacklistRefreshTokenUserId(uid)

	if err != nil {
		logger.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// found in blacklist
	if blacklistRefreshTokenExpiration >= 0 {
		if claims.ExpiresAt < blacklistRefreshTokenExpiration {
			// forced logout
			logger.Logger.Error("refresh token has expired")
			c.JSON(http.StatusBadRequest, gin.H{"error": "refresh token has expired"})
			return
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	err = userCollection.FindOne(ctx, bson.M{"user_id": uid}).Decode(&foundUser)

	if err != nil {
		logger.Logger.Error(err.Error())
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
		logger.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = tokenHelper.UpdateAllTokens(signedToken, signedRefreshToken, foundUser.User_id)

	if err != nil {
		logger.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = userCollection.FindOne(ctx, bson.M{"user_id": foundUser.User_id}).Decode(&foundUser)

	if err != nil {
		logger.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, foundUser)
}

func (a *authControllerStruct) ResetUserPassword(c *gin.Context) {
	var passwordResetRequest dto.PasswordResetRequestDto

	err := c.BindJSON(&passwordResetRequest)

	if err != nil {
		logger.Logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	validationError := validate.Struct(passwordResetRequest)

	if validationError != nil {
		logger.Logger.Error(validationError.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": validationError.Error()})
		return
	}

	if passwordResetRequest.Password != passwordResetRequest.ConfirmPassword {
		logger.Logger.Error("password and confirm password not matching")
		c.JSON(http.StatusBadRequest, gin.H{"error": "password and confirm password not matching"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	userCh := make(chan models.User)
	userErrCh := make(chan error)

	mailCh := make(chan models.Mail)
	mailErrCh := make(chan error)

	go func(email string) {
		var u models.User
		userErr := userCollection.FindOne(ctx, bson.M{"email": email}).Decode(&u)
		userErrCh <- userErr
		userCh <- u
	}(passwordResetRequest.Email)

	go func(email, emailType string) {
		var m models.Mail
		mailErr := mailCollection.FindOne(ctx, bson.D{
			{Key: "email", Value: passwordResetRequest.Email},
			{Key: "type", Value: enums.PasswordReset.String()},
		}).Decode(&m)
		mailErrCh <- mailErr
		mailCh <- m
	}(passwordResetRequest.Email, enums.PasswordReset.String())

	err = <-userErrCh

	if err != nil {
		logger.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user := <-userCh

	passwordMatchingErr := passwordHelper.VerifyPassword(user.Password, passwordResetRequest.Password)

	// new password is same as old password
	if passwordMatchingErr == nil {
		logger.Logger.Error("new password is same as old password")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "new password is same as old password"})
		return
	}

	err = <-mailErrCh

	if err != nil {
		logger.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	passwordResetMail := <-mailCh

	// check code
	if passwordResetRequest.Code != passwordResetMail.Code {
		logger.Logger.Error("password reset mail code not match")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "password reset mail code not match"})
		return
	}

	var updateObj bson.D

	hashedPassword, err := passwordHelper.HashPassword(passwordResetRequest.Password)

	if err != nil {
		logger.Logger.Error(validationError.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": validationError.Error()})
		return
	}

	updateObj = append(updateObj, bson.E{Key: "password", Value: hashedPassword})

	Updated_at, err := timeHelper.GetCurrentTimeSingapore()

	if err != nil {
		logger.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while parsing updated_at"})
		return
	}

	updateObj = append(updateObj, bson.E{Key: "updated_at", Value: Updated_at})

	upsert := true
	filter := bson.M{"user_id": user.User_id}
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	err = nftRaffleDbClient.UseSession(ctx, func(sessionContext mongo.SessionContext) error {
		err := sessionContext.StartTransaction()
		if err != nil {
			return err
		}

		redisTokenErrCh := make(chan error)

		go func(userId string) {
			redisTokenErr := tokenHelper.SetBlacklistAccessAndRefreshTokenUserId(user.User_id)
			redisTokenErrCh <- redisTokenErr
		}(user.User_id)

		_, err = userCollection.UpdateOne(
			sessionContext,
			filter,
			bson.D{
				{Key: "$set", Value: updateObj},
			},
			&opt,
		)

		if err != nil {
			logger.Logger.Error(err.Error())
			sessionContext.AbortTransaction(sessionContext)
			return err
		}

		// now remove password reset mail from mailCollection
		_, err = mailCollection.DeleteOne(sessionContext, bson.D{
			{Key: "email", Value: passwordResetMail.Email},
			{Key: "type", Value: enums.PasswordReset.String()},
		})

		if err != nil {
			logger.Logger.Error(err.Error())
			sessionContext.AbortTransaction(sessionContext)
			return err
		}

		err = <-redisTokenErrCh

		if err != nil {
			sessionContext.AbortTransaction(sessionContext)
			return err
		}

		if err := sessionContext.CommitTransaction(sessionContext); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		logger.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func (a *authControllerStruct) TestRedis(c *gin.Context) {
	userId := "123124124"
	err := tokenHelper.SetBlacklistAccessAndRefreshTokenUserId(userId)

	if err != nil {
		logger.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	blacklistAccessTokenExpiration, err := tokenHelper.GetBlacklistAccessTokenUserId(userId)

	var blacklistAccessTokenExpirationStr string
	if blacklistAccessTokenExpiration > 0 {
		blacklistAccessTokenExpirationStr = strconv.Itoa(int(blacklistAccessTokenExpiration))
	} else {
		blacklistAccessTokenExpirationStr = "empty"
	}

	if err != nil {
		logger.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var blacklistRefreshTokenExpirationStr string
	blacklistRefreshTokenExpiration, err := tokenHelper.GetBlacklistRefreshTokenUserId(userId)

	if blacklistRefreshTokenExpiration > 0 {
		blacklistRefreshTokenExpirationStr = strconv.Itoa(int(blacklistRefreshTokenExpiration))
	} else {
		blacklistRefreshTokenExpirationStr = "empty"
	}

	if err != nil {
		logger.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	uid := c.GetString("uid")

	c.JSON(http.StatusOK, gin.H{
		"user_id":                            uid,
		"blacklist_access_token_expiration":  blacklistAccessTokenExpirationStr,
		"blacklist_refresh_token_expiration": blacklistRefreshTokenExpirationStr,
	})
}
