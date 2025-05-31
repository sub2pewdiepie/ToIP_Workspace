package auth

import (
	"net/http"
	"space/utils"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// AuthMiddleware validates the JWT token and sets the username in the context
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")

		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Missing or invalid Authorization header"})
			utils.Logger.WithFields(logrus.Fields{
				"method": c.Request.Method,
				"path":   c.Request.URL.Path,
			}).Error("Missing or invalid Authorization header")
			c.Abort()
			return
		}

		// Extract the token
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		claims, err := ValidateJWT(tokenString)
		if err != nil {
			utils.Logger.WithFields(logrus.Fields{
				"method": c.Request.Method,
				"path":   c.Request.URL.Path,
				"error":  err,
			}).Error("Invalid token")
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized: " + err.Error()})
			c.Abort()
			return
		}

		// Set username in context
		if claims.Username == "" {
			utils.Logger.WithFields(logrus.Fields{
				"method": c.Request.Method,
				"path":   c.Request.URL.Path,
			}).Error("Invalid token: username not found")
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token: username not found"})
			c.Abort()
			return
		}

		utils.Logger.WithFields(logrus.Fields{
			"method":   c.Request.Method,
			"path":     c.Request.URL.Path,
			"username": claims.Username,
		}).Debug("Authenticated user")
		c.Set("username", claims.Username)

		c.Next()
	}
}

// LoggingMiddleware logs HTTP requests
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		latency := time.Since(start)

		utils.Logger.WithFields(logrus.Fields{
			"method":    c.Request.Method,
			"path":      c.Request.URL.Path,
			"status":    c.Writer.Status(),
			"latency":   latency.Milliseconds(),
			"client_ip": c.ClientIP(),
		}).Info("Handled HTTP request")
	}
}
