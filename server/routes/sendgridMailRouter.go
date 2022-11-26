package routes

import (
	"nft-raffle/controllers"

	"github.com/gin-gonic/gin"
)

var (
	sendGridController controllers.ISendGridController = controllers.SendGridController
)

func SendGridMailRoutes(superRoute *gin.RouterGroup) {
	sendgridMailRouter := superRoute.Group("/send-grid")

	sendgridMailRouter.POST("/send-password-reset-mail", sendGridController.SendPasswordResetMail)
	sendgridMailRouter.GET("/verify-verification-mail", sendGridController.VerifyVerificationMail)
	sendgridMailRouter.GET("/verify-password-reset-mail", sendGridController.VerifyPasswordResetMail)
}
