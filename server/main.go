package main

import (
	helper "nft-raffle/helpers"
	"nft-raffle/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	port := helper.GetEnvVariable("PORT")

	if port == "" {
		port = "8000"
	}

	router := gin.New()
	router.Use(gin.Logger())

	routerGroup := router.Group("/api")
	routes.AddRoutes(routerGroup)

	router.Run(":" + port)
}
