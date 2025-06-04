package services

import (
	"errors"
	"space/models"
	"space/models/dto"
	"space/repositories"
	"space/utils"

	"github.com/sirupsen/logrus"
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

func (s *GroupUserService) GetUsersByGroupID(groupID int32) ([]dto.UserDTO, error) {
	groupUsers, err := s.groupUserRepo.FindByGroupID(groupID)
	if err != nil {
		return nil, err
	}
	var users []dto.UserDTO
	for _, gu := range groupUsers {
		users = append(users, dto.ToUserDTO(&gu.User))
	}
	utils.Logger.WithFields(logrus.Fields{
		"group_id": groupID,
		"count":    len(users),
	}).Debug("Fetched group users")
	return users, nil
}

func (s *GroupService) GetGroupUsers(groupID int32) ([]dto.UserDTO, error) {
	users, err := s.groupuserRepo.FindByGroupID(groupID)
	if err != nil {
		return nil, err
	}
	var userDTOs []dto.UserDTO
	for _, user := range users {
		userDTOs = append(userDTOs, dto.ToUserDTO(&user.User))
	}
	utils.Logger.WithFields(logrus.Fields{
		"group_id": groupID,
		"count":    len(userDTOs),
	}).Debug("Fetched group users")
	return userDTOs, nil
}

func (s *GroupService) IsGroupMember(groupID int32, username string) (bool, error) {
	user, err := s.userRepo.GetByUsername(username)
	if err != nil {
		return false, err
	}
	isMember, err := s.groupuserRepo.IsMember(groupID, user.UserID)
	if err != nil {
		return false, err
	}
	return isMember, nil
}

// func (s *GroupService) GetGroupModeratorsAndAdmin(groupID int32) (dto.ModeratorsResponse, error) {
// 	group, err := s.groupRepo.GetByID(groupID)
// 	if err != nil {
// 		return dto.ModeratorsResponse{}, err
// 	}
// 	moders, err := s.groupModerRepo.FindByGroupID(groupID)
// 	if err != nil {
// 		return dto.ModeratorsResponse{}, err
// 	}
// 	var moderatorDTOs []dto.UserDTO
// 	for _, moder := range moders {
// 		moderatorDTOs = append(moderatorDTOs, dto.ToUserDTO(&moder.User))
// 	}
// 	response := dto.ModeratorsResponse{
// 		Admin:      dto.ToUserDTO(&group.Admin),
// 		Moderators: moderatorDTOs,
// 	}
// 	utils.Logger.WithFields(logrus.Fields{
// 		"group_id":        groupID,
// 		"admin_id":        group.AdminID,
// 		"moderator_count": len(moderatorDTOs),
// 	}).Debug("Fetched group moderators and admin")
// 	return response, nil
// }
