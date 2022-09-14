package services

import (
	"nft-raffle/dto"
	"nft-raffle/enums"
	"nft-raffle/helpers"

	"github.com/sendgrid/rest"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type SendGridMailService interface {
	SendMailAsync(body []byte) (chan *rest.Response, chan error)
	DynamicTemplate(mailRequest *dto.MailRequest) []byte
	SendVerificationMailAsync(mailRequest *dto.MailRequest) (chan *rest.Response, chan error)
}

type sendGridMailServiceStruct struct{}

var (
	dotEnvHelper helpers.DotEnvHelper = helpers.NewDotEnvHelper()

	sendgridMailVerificationDynamicTemplateId string = dotEnvHelper.GetEnvVariable("SENDGRID_MAIL_VERIFICATION_DYNAMIC_TEMPLATE_ID")
	sendgridApiKey                            string = dotEnvHelper.GetEnvVariable("SENDGRID_API_KEY")
	sendgridApiEndPoint                       string = dotEnvHelper.GetEnvVariable("SENDGRID_API_ENDPOINT")
	sendgridApiHost                           string = dotEnvHelper.GetEnvVariable("SENDGRID_API_HOST")
)

func NewSendGridMailService() SendGridMailService {
	return &sendGridMailServiceStruct{}
}

func (s *sendGridMailServiceStruct) SendMailAsync(body []byte) (chan *rest.Response, chan error) {
	request := sendgrid.GetRequest(sendgridApiKey, sendgridApiEndPoint, sendgridApiHost)
	request.Method = "POST"
	request.Body = body
	responseCh, errCh := sendgrid.MakeRequestAsync(request)
	return responseCh, errCh
}

func (s *sendGridMailServiceStruct) DynamicTemplate(mailRequest *dto.MailRequest) []byte {
	m := mail.NewV3Mail()

	from := mail.NewEmail(mailRequest.FromName, mailRequest.FromEmail)
	m.SetFrom(from)

	if mailRequest.MailType == enums.MailVerification {
		m.SetTemplateID(sendgridMailVerificationDynamicTemplateId)
	}

	p := mail.NewPersonalization()

	p.AddTos(mailRequest.Tos...)

	if mailRequest.DynamicTemplateData != nil {
		for key, val := range mailRequest.DynamicTemplateData {
			p.SetDynamicTemplateData(key, val)
		}
	}

	m.AddPersonalizations(p)

	return mail.GetRequestBody(m)
}

func (s *sendGridMailServiceStruct) SendVerificationMailAsync(mailRequest *dto.MailRequest) (chan *rest.Response, chan error) {
	body := s.DynamicTemplate(mailRequest)
	responseCh, errCh := s.SendMailAsync(body)
	return responseCh, errCh
}
