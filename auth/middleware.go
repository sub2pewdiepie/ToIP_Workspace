package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates the JWT token and sets the username in the context
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

		claims, err := ValidateJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized: " + err.Error()})
			c.Abort()
			return
		}

		// Set username in context
		if claims.Username == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token: username not found"})
			c.Abort()
			return
		}
		c.Set("username", claims.Username)

		c.Next()
	}
}
