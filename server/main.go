package main

import (
	"log"
	"nft-raffle/helpers"
	"nft-raffle/routes"

	"github.com/gin-gonic/gin"
)

var (
	dotEnvHelper helpers.DotEnvHelper = helpers.NewDotEnvHelper()
)

func main() {
	log.SetPrefix("[LOG] ")
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Llongfile)
	port := dotEnvHelper.GetEnvVariable("PORT")

	if port == "" {
		port = "8000"
	}

	router := gin.New()
	router.Use(gin.Logger())

	routerGroup := router.Group("/api")
	routes.AddRoutes(routerGroup)

	router.Run(":" + port)
}
