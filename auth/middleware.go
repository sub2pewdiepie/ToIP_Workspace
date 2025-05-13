package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates the JWT token

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")

		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Missing or invalid Authorization header"})
			c.Abort()
			return
		}

		// Extract the token
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		_, err := ValidateJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized: " + err.Error()})
			c.Abort()
			return
		}

		c.Next()
	}
}
