package main

import (
	"net/http"
	helper "nft-raffle/helpers"

	"github.com/gin-gonic/gin"
)

func main() {
	port := helper.GetEnvVariable("PORT")

	if port == "" {
		port = "8000"
	}

	router := gin.New()
	router.Use(gin.Logger())

	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"testing": "OK"})
	})

	router.Run(":" + port)
}
