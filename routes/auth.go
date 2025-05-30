package routes

import (
	"net/http"
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
// routes/auth.go
func LoginHandler(authService *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input services.LoginInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		token, err := authService.LoginUser(input)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": token})
	}
}
