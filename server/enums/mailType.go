package enums

type MailType string

const (
	MailVerification MailType = "MailVerification"
)

func (m MailType) String() string {
	switch m {
	case MailVerification:
		return "MailVerification"
	}
	return "unknown"
}
