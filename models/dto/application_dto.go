package dto

type CreateApplicationRequest struct {
	GroupID int32  `json:"group_id" binding:"required"`
	Message string `json:"message"`
}
