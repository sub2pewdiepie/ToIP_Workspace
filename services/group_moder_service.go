package services

import (
	"errors"
	"space/models"
	"space/models/dto"
	"space/repositories"
	"space/utils"

	"github.com/sirupsen/logrus"
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

func (s *GroupModerService) GetModeratorsByGroupID(groupID int32) ([]dto.UserDTO, error) {
	groupModers, err := s.groupModerRepo.FindByGroupID(groupID)
	if err != nil {
		return nil, err
	}
	var moderators []dto.UserDTO
	for _, gm := range groupModers {
		moderators = append(moderators, dto.ToUserDTO(&gm.User))
	}
	utils.Logger.WithFields(logrus.Fields{
		"group_id": groupID,
		"count":    len(moderators),
	}).Debug("Fetched group moderators")
	return moderators, nil
}
