package dto

import (
	"space/models"
	"time"
)

type TaskDTO struct {
	ID          int32      `json:"id"`
	GroupID     int32      `json:"group_id"`
	UserID      int32      `json:"user_id"`
	Username    string     `json:"username"`
	SubjectID   *int32     `json:"subject_id,omitempty"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	IsVerified  bool       `json:"is_verified"`
	Deadline    *time.Time `json:"deadline,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
}

func ToTaskDTO(task *models.Task) TaskDTO {
	return TaskDTO{
		ID:          task.ID,
		GroupID:     task.GroupID,
		UserID:      task.UserID,
		Username:    task.User.Username,
		SubjectID:   task.SubjectID,
		Title:       task.Title,
		Description: task.Description,
		IsVerified:  task.IsVerified,
		Deadline:    task.Deadline,
		CreatedAt:   task.CreatedAt,
	}
}

type CreateTaskRequest struct {
	GroupID     int32      `json:"group_id" binding:"required"`
	Title       string     `json:"title" binding:"required"`
	Description string     `json:"description"`
	Deadline    *time.Time `json:"deadline"`
	SubjectID   *int32     `json:"subject_id"`
}

type VerificationRequest struct {
	VerificationStatus bool `json:"is_verified" binding:"required"`
}
