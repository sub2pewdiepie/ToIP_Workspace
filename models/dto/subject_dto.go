package dto

import (
	"time"
)

type SubjectDTO struct {
	ID              int32  `json:"id"`
	Name            string `json:"name"`
	AcademicGroupID int32  `json:"academic_group_id"`
}

type SubjectDetailDTO struct {
	ID            int32            `json:"id"`
	Name          string           `json:"name"`
	AcademicGroup AcademicGroupDTO `json:"academic_group"`
	Groups        []GroupDTO       `json:"groups"` // Groups in this academic group
}

type TaskDetailDTO struct {
	ID            int32            `json:"id"`
	Title         string           `json:"title"`
	Description   string           `json:"description"`
	IsVerified    bool             `json:"is_verified"`
	Deadline      *time.Time       `json:"deadline,omitempty"`
	CreatedAt     time.Time        `json:"created_at"`
	UpdatedAt     time.Time        `json:"updated_at"`
	User          UserDTO          `json:"user"`
	Group         GroupDTO         `json:"group"`
	Subject       *SubjectDTO      `json:"subject,omitempty"`
	AcademicGroup AcademicGroupDTO `json:"academic_group"`
}

type SubjectsResponse struct {
	Subjects   []SubjectDTO   `json:"subjects"`
	Pagination PaginationMeta `json:"pagination"`
}

type SubjectsDetailResponse struct {
	Subjects   []SubjectDetailDTO `json:"subjects"`
	Pagination PaginationMeta     `json:"pagination"`
}

type TasksResponse struct {
	Tasks      []TaskDTO      `json:"tasks"`
	Pagination PaginationMeta `json:"pagination"`
}

type TasksDetailResponse struct {
	Tasks      []TaskDetailDTO `json:"tasks"`
	Pagination PaginationMeta  `json:"pagination"`
}
