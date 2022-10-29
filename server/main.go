package main

import (
	"log"
	"nft-raffle/helpers"
	"nft-raffle/routes"

	"github.com/gin-gonic/gin"
)

var (
	dotEnvHelper helpers.IDotEnvHelper = helpers.DotEnvHelper
)

func init() {
	log.SetPrefix("[LOG] ")
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Llongfile)
}

func main() {
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
