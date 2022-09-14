package controllers

import (
	"log"
	"net/http"
	"nft-raffle/database"
	"nft-raffle/dto"
	"nft-raffle/enums"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type SendGridController interface {
	SendGridTesting() gin.HandlerFunc
}

type sendGridControllerStruct struct{}

var (
	mailCollection *mongo.Collection = database.OpenCollection(database.Client, "mail")
)

func NewSendGridController() SendGridController {
	return &sendGridControllerStruct{}
}

func (s *sendGridControllerStruct) SendGridTesting() gin.HandlerFunc {
	return func(c *gin.Context) {
		fromName := "yongheng98_testing_sendgrid"
		fromEmail := "yongheng98@hotmail.com"

		tos := []*mail.Email{
			mail.NewEmail("yyhyap98", "yyhyap98@gmail.com"),
		}

		dynamicTemplateData := map[string]string{}
		dynamicTemplateData["Last_Name"] = "YYHYAP98"
		dynamicTemplateData["Verify_Mail_Link"] = "https://www.google.com"

		mailReq := &dto.MailRequest{
			FromName:            fromName,
			FromEmail:           fromEmail,
			MailType:            enums.MailVerification,
			Tos:                 tos,
			DynamicTemplateData: dynamicTemplateData,
		}

		var Body = sendGridMailService.DynamicTemplate(mailReq)
		responseCh, errCh := sendGridMailService.SendMailAsync(Body)

		select {
		case err := <-errCh:
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		case response := <-responseCh:
			c.JSON(http.StatusOK, gin.H{
				"status":  response.StatusCode,
				"body":    response.Body,
				"headers": response.Headers,
			})
			return
		}
	}
}
