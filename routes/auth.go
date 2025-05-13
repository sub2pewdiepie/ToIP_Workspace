package routes

import (
	"net/http"
	"space/auth"
	"space/services"

	"github.com/gin-gonic/gin"
)

// Credentials structure for login input
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Simple hardcoded login credentials
// FOR TESTING PURPOSES ONLY
var users = map[string]string{
	"user": "password",
}

// LoginHandler handles user login and JWT token generation
// @Tags auth
// @Description Хендлер авторизации
// @ID login
// @Accept json
// @Produce json
// @Param input body Credentials true "credentials"
// @Success 200 {string} string "token"
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /login [post]
func LoginHandler(c *gin.Context) {
	var creds Credentials
	if err := c.ShouldBindJSON(&creds); err != nil {
		// c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		services.HandleError(c, http.StatusBadRequest, "Invalid request")
		return
	}

	// Validate credentials
	if password, ok := users[creds.Username]; !ok || password != creds.Password {
		services.HandleError(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Generate JWT
	token, err := auth.GenerateJWT(creds.Username)
	if err != nil {
		services.HandleError(c, http.StatusInternalServerError, "Could not generate token")
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
