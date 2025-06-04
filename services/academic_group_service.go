package services

import (
	"errors"
	"space/models"
	"space/models/dto"
	"space/repositories"
	"space/utils"

	"github.com/sirupsen/logrus"
)

type AcademicGroupService struct {
	academicGroupRepo *repositories.AcademicGroupRepository
}

func NewAcademicGroupService(academicGroupRepo *repositories.AcademicGroupRepository) *AcademicGroupService {
	return &AcademicGroupService{academicGroupRepo}
}

func (s *AcademicGroupService) GetAcademicGroupByID(id int32) (*models.AcademicGroup, error) {
	return s.academicGroupRepo.GetByID(id)
}

func (s *AcademicGroupService) CreateAcademicGroup(academicGroup *models.AcademicGroup) error {
	if academicGroup.Name == "" {
		return errors.New("name is required")
	}
	return s.academicGroupRepo.Create(academicGroup)
}

func (s *AcademicGroupService) UpdateAcademicGroup(academicGroup *models.AcademicGroup) error {
	if academicGroup.Name == "" {
		return errors.New("name is required")
	}
	return s.academicGroupRepo.Update(academicGroup)
}

func (s *AcademicGroupService) DeleteAcademicGroup(id int32) error {
	return s.academicGroupRepo.Delete(id)
}

func (s *AcademicGroupService) GetAllAcademicGroups() ([]dto.AcademicGroupDTO, error) {
	groups, err := s.academicGroupRepo.FindAll()
	if err != nil {
		return nil, err
	}
	dtos := dto.ToAcademicGroupDTOs(groups)
	utils.Logger.WithFields(logrus.Fields{
		"count": len(dtos),
	}).Debug("Fetched all academic groups")
	return dtos, nil
}
