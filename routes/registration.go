package routes

import (
	"net/http"
	"space/services"

	"github.com/gin-gonic/gin"
)

type RegisterInput struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// RegisterHandler handles user registration
// @Tags auth
// @Description Register a new user
// @ID register
// @Accept json
// @Produce json
// @Param input body RegisterInput true "User registration input"
// @Success 201 {string} string "registered"
// @Failure 400 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /register [post]
func RegisterHandler(authService *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input services.RegisterInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		err := authService.RegisterUser(input)
		if err != nil {
			status := http.StatusInternalServerError
			if err.Error() == "user already exists" {
				status = http.StatusBadRequest
			}
			c.JSON(status, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
	}
}
