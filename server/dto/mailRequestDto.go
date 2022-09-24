package dto

import (
	"nft-raffle/enums"

	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type MailRequest struct {
	FromName            string
	FromEmail           string
	MailType            enums.MailType
	Tos                 []*mail.Email
	DynamicTemplateData map[string]string
}
