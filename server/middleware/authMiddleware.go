package middleware

import (
	"log"
	"net/http"
	"nft-raffle/helpers"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware interface {
	Authenticate() gin.HandlerFunc
}

type authMiddlewareStruct struct{}

var (
	tokenHelper helpers.TokenHelper = helpers.NewTokenHelper()
)

func NewAuthMiddleware() AuthMiddleware {
	return &authMiddlewareStruct{}
}

func (a *authMiddlewareStruct) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorizationHeader := strings.Split(c.Request.Header.Get("Authorization"), " ")

		if len(authorizationHeader) < 2 {
			log.Println("no authorization header provided")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No authorization header provided"})
			c.Abort()
			return
		}

		clientAccessToken := authorizationHeader[1]

		if clientAccessToken == "" {
			log.Println("no authorization header provided")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No authorization header provided"})
			c.Abort()
			return
		}

		claims, err := tokenHelper.ValidateAccessToken(clientAccessToken)

		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		// check Access Token blacklist
		blacklistAccessTokenExpiration, err := tokenHelper.GetBlacklistAccessTokenUserId(claims.Uid)

		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		// found in JWT blacklist
		if blacklistAccessTokenExpiration >= 0 {
			if claims.ExpiresAt < blacklistAccessTokenExpiration {
				// forced logout
				log.Println("access token has expired")
				c.JSON(http.StatusBadRequest, gin.H{"error": "access token has expired"})
				c.Abort()
				return
			}
		}

		c.Set("uid", claims.Uid)
		c.Set("email", claims.Email)
		c.Set("first_name", claims.First_name)
		c.Set("last_name", claims.Last_name)
		c.Set("user_role", claims.User_role)
		c.Set("issued_at", claims.IssuedAt)
		c.Set("subject", claims.Subject)
		c.Next()
	}
}
