package services

import (
	"errors"
	"space/models"
	"space/models/dto"
	"space/repositories"
	"space/utils"

	"github.com/sirupsen/logrus"
)

// type SubjectService struct {
// 	subjectRepo *repositories.SubjectRepository
// }

// func NewSubjectService(subjectRepo *repositories.SubjectRepository) *SubjectService {
// 	return &SubjectService{subjectRepo}
// }

type SubjectService struct {
	subjectRepo *repositories.SubjectRepository
	groupRepo   *repositories.GroupRepository
	userRepo    repositories.UserRepository
}

func NewSubjectService(subjectRepo *repositories.SubjectRepository, groupRepo *repositories.GroupRepository, userRepo repositories.UserRepository) *SubjectService {
	return &SubjectService{
		subjectRepo: subjectRepo,
		groupRepo:   groupRepo,
		userRepo:    userRepo,
	}
}
func (s *SubjectService) GetSubjectByID(id int32) (*models.Subject, error) {
	return s.subjectRepo.GetByID(id)
}

func (s *SubjectService) CreateSubject(subject *models.Subject) error {
	if subject.Name == "" || subject.AcademicGroupID == 0 {
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

func (s *SubjectService) GetSubjectsByAcademicGroupID(academicGroupID int32, page, pageSize int) ([]dto.SubjectDTO, int64, error) {
	subjects, total, err := s.subjectRepo.FindByAcademicGroupID(academicGroupID, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	subjectDTOs := make([]dto.SubjectDTO, len(subjects))
	for i, subject := range subjects {
		subjectDTOs[i] = dto.SubjectDTO{
			ID:              subject.SubjectID,
			Name:            subject.Name,
			AcademicGroupID: subject.AcademicGroupID,
		}
	}

	return subjectDTOs, total, nil
}

func (s *SubjectService) GetUserSubjects(username string, page, pageSize int) ([]dto.SubjectDetailDTO, int64, error) {
	user, err := s.userRepo.GetByUsername(username)
	if err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"error":    err,
			"username": username,
		}).Error("User not found")
		return nil, 0, errors.New("user not found")
	}

	subjects, groups, total, err := s.subjectRepo.FindByUserGroups(user.UserID, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	subjectDTOs := make([]dto.SubjectDetailDTO, len(subjects))
	for i, subject := range subjects {
		groupDTOs := make([]dto.GroupDTO, 0)
		for _, group := range groups {
			if group.AcademicGroupID == subject.AcademicGroupID {
				groupDTOs = append(groupDTOs, dto.GroupDTO{
					ID:              group.ID,
					Name:            group.Name,
					AdminUsername:   group.Admin.Username,
					AcademicGroupID: group.AcademicGroupID,
					AcademicGroup:   group.AcademicGroup.Name,
				})
			}
		}
		subjectDTOs[i] = dto.SubjectDetailDTO{
			ID:   subject.SubjectID,
			Name: subject.Name,
			AcademicGroup: dto.AcademicGroupDTO{
				ID:   subject.AcademicGroup.AcademicGroupID,
				Name: subject.AcademicGroup.Name,
			},
			Groups: groupDTOs,
		}
	}

	return subjectDTOs, total, nil
}
