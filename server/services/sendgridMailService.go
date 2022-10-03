package services

import (
	"context"
	"log"
	"nft-raffle/database"
	"nft-raffle/dto"
	"nft-raffle/enums"
	"nft-raffle/helpers"
	"nft-raffle/models"
	"time"

	"github.com/sendgrid/rest"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SendGridMailService interface {
	SendMailAsync(body []byte) (chan *rest.Response, chan error)
	DynamicTemplate(mailRequest *dto.MailRequest) []byte
	SendMail(mailRequest *dto.MailRequest)
	CreateNewMail(mailType enums.MailType, email string, randomSixDigits string, expires_at time.Time) error
	UpdateEmail(mailType enums.MailType, email string, randomSixDigits string, expires_at time.Time) error
}

type sendGridMailServiceStruct struct{}

var (
	nftRaffleDb       database.NftRaffleMongoDbConnection = database.NewNftRaffleMongoDbConnection()
	nftRaffleDbClient *mongo.Client                       = nftRaffleDb.DBClient()
	mailCollection    *mongo.Collection                   = nftRaffleDb.OpenCollection(nftRaffleDbClient, "mail")

	dotEnvHelper helpers.DotEnvHelper = helpers.NewDotEnvHelper()

	sendgridMailVerificationDynamicTemplateId  string = dotEnvHelper.GetEnvVariable("SENDGRID_MAIL_VERIFICATION_DYNAMIC_TEMPLATE_ID")
	sendgridMailPasswordResetDynamicTemplateId string = dotEnvHelper.GetEnvVariable("SENDGRID_MAIL_PASSWORD_RESET_DYNAMIC_TEMPLATE_ID")
	sendgridApiKey                             string = dotEnvHelper.GetEnvVariable("SENDGRID_API_KEY")
	sendgridApiEndPoint                        string = dotEnvHelper.GetEnvVariable("SENDGRID_API_ENDPOINT")
	sendgridApiHost                            string = dotEnvHelper.GetEnvVariable("SENDGRID_API_HOST")
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
	} else if mailRequest.MailType == enums.PasswordReset {
		m.SetTemplateID(sendgridMailPasswordResetDynamicTemplateId)
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

func (s *sendGridMailServiceStruct) SendMail(mailRequest *dto.MailRequest) {
	body := s.DynamicTemplate(mailRequest)
	responseCh, errCh := s.SendMailAsync(body)

	select {
	case err := <-errCh:
		log.Println(err)
	case response := <-responseCh:
		// do nothing
		_ = response
		for _, val := range mailRequest.Tos {
			log.Printf("%s mail for %s has been sent \n", mailRequest.MailType.String(), val.Address)
		}
	}
}

func (s *sendGridMailServiceStruct) CreateNewMail(mailType enums.MailType, email string, randomSixDigits string, expires_at time.Time) error {
	var mail models.Mail
	var err error
	mail.ID = primitive.NewObjectID()
	mail.Mail_id = mail.ID.Hex()
	mail.Email = email

	mail.Code = randomSixDigits
	mail.Type = mailType.String()

	mail.Created_at, err = time.Parse(time.RFC3339, time.Now().Local().Format(time.RFC3339))

	if err != nil {
		log.Println(err)
		return err
	}

	mail.Updated_at, err = time.Parse(time.RFC3339, time.Now().Local().Format(time.RFC3339))

	if err != nil {
		log.Println(err)
		return err
	}

	mail.Expires_at = expires_at

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	_, insertError := mailCollection.InsertOne(ctx, mail)
	defer cancel()

	if insertError != nil {
		log.Println(insertError)
		return insertError
	}

	return nil
}

func (s *sendGridMailServiceStruct) UpdateEmail(mailType enums.MailType, email string, randomSixDigits string, expires_at time.Time) error {
	var updateObj bson.D

	updateObj = append(updateObj, bson.E{Key: "code", Value: randomSixDigits})
	updateObj = append(updateObj, bson.E{Key: "expires_at", Value: expires_at})

	Updated_at, err := time.Parse(time.RFC3339, time.Now().Local().Format(time.RFC3339))

	if err != nil {
		return err
	}

	updateObj = append(updateObj, bson.E{Key: "updated_at", Value: Updated_at})

	upsert := true
	filter := bson.D{
		{Key: "email", Value: email},
		{Key: "type", Value: mailType.String()},
	}
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)

	_, updateError := mailCollection.UpdateOne(
		ctx,
		filter,
		bson.D{
			{Key: "$set", Value: updateObj},
		},
		&opt,
	)
	defer cancel()

	if updateError != nil {
		return updateError
	}

	return nil
}
