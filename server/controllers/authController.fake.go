package controllers

import (
	"fmt"
	"net/http"
	"nft-raffle/logger"
	"nft-raffle/models"
	"sync"

	"github.com/gin-gonic/gin"
)

var (
	fakeAuthController     *fakeAuthControllerStruct
	fakeAuthControllerOnce sync.Once
)

type IFakeAuthController interface {
	FakeLogin(c *gin.Context)
}

type fakeAuthControllerStruct struct{}

func GetFakeAuthController() *fakeAuthControllerStruct {
	if fakeAuthController == nil {
		fakeAuthControllerOnce.Do(func() {
			fakeAuthController = &fakeAuthControllerStruct{}
		})
	}
	return fakeAuthController
}

func (a *fakeAuthControllerStruct) FakeLogin(c *gin.Context) {
	var user models.User
	var foundUser models.User

	if err := c.BindJSON(&user); err != nil {
		logger.Logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// find user from mongodb
	foundUser, err := fakeControllerRepository.MockFindUserByEmailSuccess(user.Email)

	if err != nil {
		logger.Logger.Error(fmt.Errorf("error while finding user: %v", err).Error())
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

	// generate all tokens

	// update all tokens

	// find one user
	foundUser, err = fakeControllerRepository.MockFindUserByIdSuccess(foundUser.User_id)

	if err != nil {
		logger.Logger.Error(fmt.Errorf("error while finding user: %v", err).Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "email or password is incorrect"})
		return
	}

	// loc, _ := time.LoadLocation("Asia/Singapore")
	// logger.Logger.Debug(fmt.Sprintf("local date time %v", foundUser.Updated_at.In(loc)))
	c.JSON(http.StatusOK, foundUser)
}
