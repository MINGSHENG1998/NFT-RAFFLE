package middleware

import (
	"net/http"
	"nft-raffle/helpers"
	"nft-raffle/logger"
	"strings"

	"github.com/gin-gonic/gin"
)

var (
	AuthMiddleware IAuthMiddleware = NewAuthMiddleware()

	tokenHelper helpers.ITokenHelper = helpers.TokenHelper
)

type IAuthMiddleware interface {
	Authenticate() gin.HandlerFunc
}

type authMiddlewareStruct struct{}

func NewAuthMiddleware() IAuthMiddleware {
	return &authMiddlewareStruct{}
}

func (a *authMiddlewareStruct) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorizationHeader := strings.Split(c.Request.Header.Get("Authorization"), " ")

		if len(authorizationHeader) < 2 {
			logger.Logger.Error("no authorization header provided")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No authorization header provided"})
			c.Abort()
			return
		}

		clientAccessToken := authorizationHeader[1]

		if clientAccessToken == "" {
			logger.Logger.Error("no authorization header provided")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No authorization header provided"})
			c.Abort()
			return
		}

		claims, err := tokenHelper.ValidateAccessToken(clientAccessToken)

		if err != nil {
			logger.Logger.Error(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		// check Access Token blacklist
		blacklistAccessTokenExpiration, err := tokenHelper.GetBlacklistAccessTokenUserId(claims.Uid)

		if err != nil {
			logger.Logger.Error(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		// found in JWT blacklist
		if blacklistAccessTokenExpiration >= 0 {
			if claims.ExpiresAt < blacklistAccessTokenExpiration {
				// forced logout
				logger.Logger.Error("access token has expired")
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
