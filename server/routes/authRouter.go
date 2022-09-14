package routes

import (
	"nft-raffle/controllers"

	"github.com/gin-gonic/gin"
)

var (
	authController controllers.AuthController = controllers.NewAuthController()
)

func AuthRoutes(superRoute *gin.RouterGroup) {
	authRouter := superRoute.Group("/auth")

	authRouter.POST("/signup", authController.SignUp())
	authRouter.POST("/login", authController.Login())
	authRouter.POST("/refreshToken", authController.RefreshToken())
}
