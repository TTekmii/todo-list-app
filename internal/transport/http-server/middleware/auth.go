package middleware

import (
	"net/http"
	"strings"

	"github.com/TTekmii/todo-list-app/internal/app/auth"
	"github.com/gin-gonic/gin"
)

const (
	authorizationHeader = "Authorization"
	userContext         = "userID"
)

func AuthMiddleware(authService *auth.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader(authorizationHeader)
		if header == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "empty auth header"})
			return
		}

		headerParts := strings.Split(header, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid auth header format"})
			return
		}

		tokenString := headerParts[1]

		userID, err := authService.ParseToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		c.Set(userContext, userID)

		c.Next()
	}
}
