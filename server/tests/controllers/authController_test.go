package tests_controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"nft-raffle/controllers"
	"nft-raffle/models"
	"nft-raffle/tests"
	"testing"

	"github.com/go-playground/assert/v2"
)

var (
	fakeAuthController controllers.IFakeAuthController = controllers.FakeAuthController
)

func TestLoginSuccess(t *testing.T) {
	r := tests.GetGinEngine()
	r.POST("/api/login", fakeAuthController.FakeLogin)

	loginUser := &models.User{
		Email:    "testingaaa@gmail.com",
		Password: "11111111",
	}

	jsonValue, _ := json.Marshal(loginUser)
	req, _ := http.NewRequest("POST", "/api/login", bytes.NewBuffer(jsonValue))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestLoginInvalidPassword(t *testing.T) {
	r := tests.GetGinEngine()
	r.POST("/api/login", fakeAuthController.FakeLogin)

	loginUser := &models.User{
		Email:    "testingaaa@gmail.com",
		Password: "11111112",
	}

	jsonValue, _ := json.Marshal(loginUser)
	req, _ := http.NewRequest("POST", "/api/login", bytes.NewBuffer(jsonValue))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
