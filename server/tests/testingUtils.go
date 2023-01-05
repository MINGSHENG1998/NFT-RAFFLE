package tests

import (
	"github.com/gin-gonic/gin"
)

func GetGinEngine() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	return r
}
