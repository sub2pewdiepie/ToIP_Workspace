package services

import (
	"errors"
	"space/models"
	"space/repositories"
)

type GroupModerService struct {
	groupModerRepo *repositories.GroupModerRepository
}

func NewGroupModerService(groupModerRepo *repositories.GroupModerRepository) *GroupModerService {
	return &GroupModerService{groupModerRepo}
}

func (s *GroupModerService) GetGroupModer(groupID, userID int32) (*models.GroupModer, error) {
	return s.groupModerRepo.GetByID(groupID, userID)
}

func (s *GroupModerService) CreateGroupModer(groupModer *models.GroupModer) error {
	if groupModer.GroupID == 0 || groupModer.UserID == 0 {
		return errors.New("group_id and user_id are required")
	}
	// Check if moderator already exists
	_, err := s.groupModerRepo.GetByID(groupModer.GroupID, groupModer.UserID)
	if err == nil {
		return errors.New("moderator already exists")
	}
	return s.groupModerRepo.Create(groupModer)
}

func (s *GroupModerService) DeleteGroupModer(groupID, userID int32) error {
	return s.groupModerRepo.Delete(groupID, userID)
}
