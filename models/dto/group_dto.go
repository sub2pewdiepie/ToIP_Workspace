package dto

import "space/models"

type GroupDTO struct {
	ID              int32  `json:"id"`
	Name            string `json:"name"`
	AdminUsername   string `json:"admin_username"`
	AcademicGroupID int32  `json:"academic_group_id"`
	AcademicGroup   string `json:"academic_group_name"`
}

// GetGroupsResponse represents the paginated group list
// swagger:model
type GetGroupsResponse struct {
	Groups     []GroupDTO     `json:"groups"`
	Pagination PaginationMeta `json:"pagination"`
}

type PaginationMeta struct {
	Page     int   `json:"page" example:"1"`
	PageSize int   `json:"page_size" example:"10"`
	Total    int64 `json:"total" example:"25"`
	Pages    int64 `json:"pages" example:"3"`
}

func ToGroupDTO(group *models.Group) GroupDTO {
	return GroupDTO{
		ID:              group.ID,
		Name:            group.Name,
		AdminUsername:   group.Admin.Username,
		AcademicGroupID: group.AcademicGroupID,
		AcademicGroup:   group.AcademicGroup.Name,
	}
}

type CreateGroupRequest struct {
	Name            string `json:"name" binding:"required"`
	AcademicGroupID int32  `json:"academic_group_id" binding:"required"`
}

type UpdateGroupRequest struct {
	Name string `json:"name,omitempty"`
	// AcademicGroupID int32  `json:"academic_group_id,omitempty"`
}
