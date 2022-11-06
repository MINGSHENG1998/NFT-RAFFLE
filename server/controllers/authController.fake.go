package controllers

import (
	"net/http"
	"nft-raffle/logger"
	"nft-raffle/models"

	"github.com/gin-gonic/gin"
)

var (
	FakeAuthController IFakeAuthController = NewFakeAuthController()
)

type IFakeAuthController interface {
	FakeLogin() gin.HandlerFunc
}

func NewFakeAuthController() IFakeAuthController {
	return &authControllerStruct{}
}

func (a *authControllerStruct) FakeLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		var foundUser models.User

		if err := c.BindJSON(&user); err != nil {
			logger.Logger.Error(err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// find user from mongodb
		foundUser = models.User{
			Email:    "testingaaa@gmail.com",
			Password: "$2a$14$X7pxIBiQtS/SFhyOHo1aIO6PFTEY5.w2xHR84e.0nOi.kqwdiTylm",
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

		// generate all tokens

		// update all tokens

		// find one user
		foundUser = models.User{
			Email:      "testingaaa@gmail.com",
			First_name: "aaa",
			Last_name:  "bbb",
		}

		// loc, _ := time.LoadLocation("Asia/Singapore")
		// logger.Logger.Debug(fmt.Sprintf("local date time %v", foundUser.Updated_at.In(loc)))
		c.JSON(http.StatusOK, foundUser)
	}
}
