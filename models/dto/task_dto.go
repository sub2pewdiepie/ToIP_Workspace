package dto

import (
	"space/models"
	"time"
)

type TaskDTO struct {
	ID          int32     `json:"id"`
	GroupID     int32     `json:"group_id"`
	UserID      int32     `json:"user_id"`
	Username    string    `json:"username"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	IsVerified  bool      `json:"is_verified"`
	CreatedAt   time.Time `json:"created_at"`
}

func ToTaskDTO(task *models.Task) TaskDTO {
	return TaskDTO{
		ID:          task.ID,
		GroupID:     task.GroupID,
		UserID:      task.UserID,
		Username:    task.User.Username,
		Title:       task.Title,
		Description: task.Description,
		IsVerified:  task.IsVerified,
		CreatedAt:   task.CreatedAt,
	}
}

type CreateTaskRequest struct {
	GroupID     int32  `json:"group_id" binding:"required"`
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
}

type VerificationRequest struct {
	VerificationStatus bool `json:"is_verified" binding:"required"`
}
