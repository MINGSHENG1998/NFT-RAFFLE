package helpers

import (
	"net/mail"
)

type DataValidationHelper interface {
	IsEmailValid(email string) error
}

type dataValidationHelperStruct struct{}

func NewDataValidationHelper() DataValidationHelper {
	return &dataValidationHelperStruct{}
}

func (d *dataValidationHelperStruct) IsEmailValid(email string) error {
	_, err := mail.ParseAddress(email)
	return err
}
