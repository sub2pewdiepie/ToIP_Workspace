package routes

import (
	"net/http"
	"space/services"

	"github.com/gin-gonic/gin"
)

// LoginHandler godoc
// @Summary Authenticate user and generate JWT token
// @Description Authenticate a user with username and password, returning a JWT token for protected endpoints.
// @Tags auth
// @ID login
// @Accept json
// @Produce json
// @Param input body services.LoginInput true "User credentials"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /login [post]
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
