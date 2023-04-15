package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/Hypersus/simplebank/token"
	"github.com/gin-gonic/gin"
)

const (
	authKey        = "Authorization"
	authType       = "Bearer"
	authPayloadKey = "payload"
)

func authMiddleware(token token.TokenMaker) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the token from the header
		authHeader := c.GetHeader(authKey)
		// If the token is empty, return an error
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, errMessage(errors.New("empty auth token")))
			return
		}
		// Split the token into two parts namely authentication type and the token
		fields := strings.Fields(authHeader)
		if len(fields) != 2 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, errMessage(errors.New("invalid auth token")))
			return
		}
		// Check if the authentication type is Bearer
		if fields[0] != authType {
			c.AbortWithStatusJSON(http.StatusUnauthorized, errMessage(errors.New("invalid auth type")))
			return
		}
		// Validate the token
		payload, err := token.ValidateToken(fields[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, errMessage(err))
			return
		}
		// Set the claims in the context
		c.Set(authPayloadKey, payload)
		c.Next()
	}
}
