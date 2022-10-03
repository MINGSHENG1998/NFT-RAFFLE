package dto

type PasswordResetRequestDto struct {
	Email           string `validate:"required"`
	Code            string `validate:"required"`
	Password        string `validate:"required"`
	ConfirmPassword string `validate:"required"`
}
