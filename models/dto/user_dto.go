package dto

import (
	"space/models"
	"time"
)

type UserResponseDTO struct {
	UserID    int32     `json:"user_id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
	// Excludes Email, HashPassword
}

type UserDTO struct {
	UserID    int32     `json:"user_id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}

type ModeratorsResponse struct {
	Admin      UserDTO   `json:"admin"`
	Moderators []UserDTO `json:"moderators"`
}

func ToUserDTO(user *models.User) UserDTO {
	return UserDTO{
		UserID:    user.UserID,
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
	}
}
