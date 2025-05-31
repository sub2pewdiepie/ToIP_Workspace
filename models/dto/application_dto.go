package dto

import (
	"space/models"
	"time"
)

type CreateApplicationRequest struct {
	GroupID int32  `json:"group_id" binding:"required"`
	Message string `json:"message"`
}
type GroupApplicationDTO struct {
	ApplicationID int32     `json:"application_id"`
	GroupID       int32     `json:"group_id"`
	UserID        int32     `json:"user_id"`
	Username      string    `json:"username"`
	Message       string    `json:"message"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
}

func ToGroupApplicationDTO(app *models.GroupApplication) GroupApplicationDTO {
	return GroupApplicationDTO{
		ApplicationID: app.ApplicationID,
		GroupID:       app.GroupID,
		UserID:        app.UserID,
		Username:      app.User.Username,
		Message:       app.Message,
		Status:        app.Status,
		CreatedAt:     app.CreatedAt,
	}
}
