package test_helpers

import (
	"nft-raffle/helpers"
	"testing"
)

var (
	dataValidationHelper helpers.DataValidationHelper = helpers.NewDataValidationHelper()
)

func TestIsEmailValid(t *testing.T) {
	email := "email@e.com"
	err := dataValidationHelper.IsEmailValid(email)

	if err != nil {
		t.Error(err.Error())
	}
}
