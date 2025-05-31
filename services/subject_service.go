package services

import (
	"errors"
	"space/models"
	"space/repositories"
)

type SubjectService struct {
	subjectRepo *repositories.SubjectRepository
}

func NewSubjectService(subjectRepo *repositories.SubjectRepository) *SubjectService {
	return &SubjectService{subjectRepo}
}

func (s *SubjectService) GetSubjectByID(id int32) (*models.Subject, error) {
	return s.subjectRepo.GetByID(id)
}

func (s *SubjectService) CreateSubject(subject *models.Subject) error {
	if subject.Name == "" || subject.GroupID == 0 {
		return errors.New("name and group_id are required")
	}
	return s.subjectRepo.Create(subject)
}

func (s *SubjectService) UpdateSubject(subject *models.Subject) error {
	if subject.Name == "" {
		return errors.New("name is required")
	}
	return s.subjectRepo.Update(subject)
}

func (s *SubjectService) DeleteSubject(id int32) error {
	return s.subjectRepo.Delete(id)
}
