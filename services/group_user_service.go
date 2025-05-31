package services

import (
	"errors"
	"space/models"
	"space/repositories"
)

type GroupUserService struct {
	groupUserRepo *repositories.GroupUserRepository
}

func NewGroupUserService(groupUserRepo *repositories.GroupUserRepository) *GroupUserService {
	return &GroupUserService{groupUserRepo}
}

func (s *GroupUserService) GetGroupUser(groupID, userID int32) (*models.GroupUser, error) {
	return s.groupUserRepo.GetByID(groupID, userID)
}

func (s *GroupUserService) CreateGroupUser(groupUser *models.GroupUser) error {
	if groupUser.GroupID == 0 || groupUser.UserID == 0 {
		return errors.New("group_id and user_id are required")
	}
	if groupUser.Role == "" {
		groupUser.Role = "member" // Default role
	}
	return s.groupUserRepo.Create(groupUser)
}

func (s *GroupUserService) UpdateGroupUser(groupUser *models.GroupUser) error {
	if groupUser.Role == "" {
		return errors.New("role is required")
	}
	return s.groupUserRepo.Update(groupUser)
}

func (s *GroupUserService) DeleteGroupUser(groupID, userID int32) error {
	return s.groupUserRepo.Delete(groupID, userID)
}
