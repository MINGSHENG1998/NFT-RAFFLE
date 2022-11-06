package tests_helpers

import (
	"nft-raffle/helpers"
	"testing"
)

var (
	dataValidationHelper helpers.IDataValidationHelper = helpers.DataValidationHelper
)

func TestIsEmailValid(t *testing.T) {
	email := "email@e"
	err := dataValidationHelper.IsEmailValid(email)

	if err != nil {
		t.Error(err.Error())
	}
}
