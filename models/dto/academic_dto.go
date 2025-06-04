package dto

import (
	"space/models"
	"time"
)

type AcademicGroupDTO struct {
	ID        int32     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

func ToAcademicGroupDTO(group *models.AcademicGroup) AcademicGroupDTO {
	return AcademicGroupDTO{
		ID:        group.AcademicGroupID,
		Name:      group.Name,
		CreatedAt: group.CreatedAt,
	}
}

func ToAcademicGroupDTOs(groups []*models.AcademicGroup) []AcademicGroupDTO {
	var dtos []AcademicGroupDTO
	for _, group := range groups {
		dtos = append(dtos, ToAcademicGroupDTO(group))
	}
	return dtos
}
