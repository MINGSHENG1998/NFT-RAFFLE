package routes

import (
	"nft-raffle/controllers"
	"nft-raffle/middleware"

	"github.com/gin-gonic/gin"
)

var (
	authController controllers.AuthController = controllers.NewAuthController()
	authMiddleware middleware.AuthMiddleware  = middleware.NewAuthMiddleware()
)

func AuthRoutes(superRoute *gin.RouterGroup) {
	authRouter := superRoute.Group("/auth")

	authRouter.POST("/signup", authController.SignUp())
	authRouter.POST("/login", authController.Login())
	authRouter.POST("/refresh-token", authController.RefreshToken())
	authRouter.POST("/reset-user-password", authController.ResetUserPassword())
	authRouter.GET("/test-redis", authMiddleware.Authenticate(), authController.TestRedis())
}
