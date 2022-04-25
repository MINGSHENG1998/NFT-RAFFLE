package routes

import (
	"nft-raffle/controllers"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(superRoute *gin.RouterGroup) {
	authRouter := superRoute.Group("/auth")

	authRouter.POST("/signup", controllers.SignUp())
	authRouter.POST("/login", controllers.Login())
}
