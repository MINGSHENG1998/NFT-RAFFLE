package main

import (
	"log"
	"math/rand"
	"nft-raffle/helpers"
	"nft-raffle/routes"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	dotEnvHelper helpers.IDotEnvHelper = helpers.DotEnvHelper
)

func init() {
	log.SetPrefix("[LOG] ")
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Llongfile)
	rand.Seed(time.Now().UnixNano())
}

func main() {
	port := dotEnvHelper.GetEnvVariable("PORT")

	if port == "" {
		port = "8000"
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	routerGroup := router.Group("/api")
	routes.AddRoutes(routerGroup)

	router.Run(":" + port)
}
