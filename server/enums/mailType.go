package enums

type MailType string

const (
	MailVerification MailType = "MailVerification"
	PasswordReset    MailType = "PasswordReset"
)

func (m MailType) String() string {
	switch m {
	case MailVerification:
		return "MailVerification"
	case PasswordReset:
		return "PasswordReset"
	}
	return "unknown"
}
