package services

import "github.com/gin-gonic/gin"

func HandleError(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, gin.H{
		"message": message,
		"error":   statusCode,
	})
}
