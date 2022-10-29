package helpers

import (
	"net/mail"
)

var DataValidationHelper IDataValidationHelper = NewDataValidationHelper()

type IDataValidationHelper interface {
	IsEmailValid(email string) error
}

type dataValidationHelperStruct struct{}

func NewDataValidationHelper() IDataValidationHelper {
	return &dataValidationHelperStruct{}
}

func (d *dataValidationHelperStruct) IsEmailValid(email string) error {
	_, err := mail.ParseAddress(email)
	return err
}
