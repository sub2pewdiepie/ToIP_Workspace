package dto

import "time"

type UserResponseDTO struct {
	UserID    int32     `json:"user_id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
	// Excludes Email, HashPassword
}
