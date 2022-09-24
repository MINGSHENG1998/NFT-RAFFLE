package routes

import (
	"nft-raffle/controllers"

	"github.com/gin-gonic/gin"
)

var (
	sendGridController controllers.SendGridController = controllers.NewSendGridController()
)

func SendGridMailRoutes(superRoute *gin.RouterGroup) {
	sendgridMailRouter := superRoute.Group("/send-grid")

	sendgridMailRouter.GET("/verify-verification-mail", sendGridController.VerifyVerificationMail())
	sendgridMailRouter.GET("/test", sendGridController.SendGridTesting())
}
