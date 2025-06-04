package services

import (
	"errors"
	"space/models"
	"space/models/dto"
	"space/repositories"
	"space/utils"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type GroupService struct {
	groupRepo      *repositories.GroupRepository
	userRepo       repositories.UserRepository // интерфейс!!!
	groupuserRepo  *repositories.GroupUserRepository
	groupModerRepo *repositories.GroupModerRepository
}

func NewGroupService(groupRepo *repositories.GroupRepository, userRepo repositories.UserRepository, groupuserRepo *repositories.GroupUserRepository, groupModerRepo *repositories.GroupModerRepository) *GroupService {
	return &GroupService{groupRepo, userRepo, groupuserRepo, groupModerRepo}
}

func (s *GroupService) GetGroupByID(id int32) (*models.Group, error) {
	return s.groupRepo.GetByID(id)
}
func (s *GroupService) GetAllGroups(page, pageSize int32) ([]dto.GroupDTO, int64, error) {
	groups, total, err := s.groupRepo.GetAllWithPagination(page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	var groupDTOs []dto.GroupDTO
	for _, group := range groups {
		dto := dto.GroupDTO{
			ID:              group.ID,
			Name:            group.Name,
			AdminUsername:   group.Admin.Username,
			AcademicGroupID: group.AcademicGroupID,
			AcademicGroup:   group.AcademicGroup.Name,
		}
		groupDTOs = append(groupDTOs, dto)
	}

	return groupDTOs, total, nil
}

func (s *GroupService) CreateGroup(c *gin.Context, group *models.Group) error {
	if group.Name == "" || group.AcademicGroupID == 0 {
		utils.Logger.WithFields(logrus.Fields{}).Error("name and academic_group_id are required")
		return errors.New("name and academic_group_id are required")
	}

	// Get authenticated user's username from context
	username, exists := c.Get("username")
	if !exists {
		utils.Logger.WithFields(logrus.Fields{
			"username": username,
			"exists":   exists,
		}).Error("user not authenticated")

		return errors.New("user not authenticated")
	}

	// Fetch user to get their ID
	user, err := s.userRepo.GetByUsername(username.(string))
	if err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"username": username,
			"error":    err,
		}).Error("failed to find user")
		return errors.New("failed to find authenticated user")
	}

	// Set AdminID to the authenticated user's ID
	group.AdminID = user.UserID
	if err := s.groupRepo.Create(group); err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"error":    err,
			"username": username,
			"name":     group.Name,
		}).Error("failed to create new group")
		return err
	}
	// Add admin to groupuser relationship
	groupUser := &models.GroupUser{
		GroupID: group.ID,
		UserID:  user.UserID,
		Role:    "admin",
	}
	if err := s.groupuserRepo.Create(groupUser); err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"error":      err,
			"user":       username,
			"group_name": group.Name,
			"userID":     user.UserID,
			"place":      "GroupService.CreateGroup",
		}).Error("failed to assign to admin user groupuser relationship")
	}

	return nil
}

func (s *GroupService) UpdateGroup(c *gin.Context, group *models.Group) error {
	if group.Name == "" {
		return errors.New("name is required")
	}

	// Get authenticated user
	username, exists := c.Get("username")
	if !exists {
		return errors.New("user not authenticated")
	}
	user, err := s.userRepo.GetByUsername(username.(string))
	if err != nil {
		return errors.New("failed to find authenticated user")
	}

	// Check if user is the group admin
	existingGroup, err := s.groupRepo.GetByID(group.ID)
	if err != nil {
		return errors.New("group not found")
	}
	if existingGroup.AdminID != user.UserID {
		return errors.New("only the group admin can update the group")
	}

	return s.groupRepo.Update(group)
}

func (s *GroupService) DeleteGroup(id int32) error {
	return s.groupRepo.Delete(id)
}

func (s *GroupService) GetApplicationsForManagedGroups(userID int32) ([]models.GroupApplication, error) {
	groups, err := s.groupRepo.GetGroupsManagedBy(userID)
	if err != nil {
		return nil, err
	}

	if len(groups) == 0 {
		return []models.GroupApplication{}, nil
	}

	groupIDs := make([]int32, len(groups))
	for i, g := range groups {
		groupIDs[i] = g.ID
	}

	var applications []models.GroupApplication
	err = s.groupRepo.DB().Where("group_id IN ?", groupIDs).Preload("User").Preload("Group").Find(&applications).Error
	if err != nil {
		return nil, err
	}

	return applications, nil
}

func (s *GroupService) GetAvailableGroups(c *gin.Context, username string, page, pageSize int) ([]dto.GroupDTO, int64, error) {
	user, err := s.userRepo.GetByUsername(username)
	if err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"error":    err,
			"username": username,
		}).Error("Failed to find user")
		return nil, 0, errors.New("failed to find user")
	}

	utils.Logger.WithFields(logrus.Fields{
		"username":  username,
		"user_id":   user.UserID,
		"page":      page,
		"page_size": pageSize,
	}).Debug("Querying available groups")

	groups, total, err := s.groupRepo.GetAvailable(user.UserID, page, pageSize)
	if err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"error":   err,
			"user_id": user.UserID,
		}).Error("Failed to query available groups")
		return nil, 0, err
	}

	groupDTOs := make([]dto.GroupDTO, len(groups))
	for i, group := range groups {
		groupDTOs[i] = dto.ToGroupDTO(&group)
	}

	utils.Logger.WithFields(logrus.Fields{
		"username": username,
		"user_id":  user.UserID,
		"total":    total,
	}).Debug("Available groups retrieved")
	return groupDTOs, total, nil
}

func (s *GroupService) IsAdminOrModerator(groupID int32, username string) (bool, error) {
	group, err := s.groupRepo.GetByID(groupID)
	if err != nil {
		return false, err
	}
	user, err := s.userRepo.GetByUsername(username)
	if err != nil {
		return false, err
	}
	if group.AdminID == user.UserID {
		return true, nil
	}
	isModerator, err := s.groupModerRepo.IsModerator(groupID, user.UserID)
	if err != nil {
		return false, err
	}
	return isModerator, nil
}

func (s *GroupService) GetGroupModeratorsAndAdmin(groupID int32) (dto.ModeratorsResponse, error) {
	group, err := s.groupRepo.GetByID(groupID)
	if err != nil {
		return dto.ModeratorsResponse{}, err
	}
	moders, err := s.groupModerRepo.FindByGroupID(groupID)
	if err != nil {
		return dto.ModeratorsResponse{}, err
	}
	var moderatorDTOs []dto.UserDTO
	for _, moder := range moders {
		moderatorDTOs = append(moderatorDTOs, dto.ToUserDTO(&moder.User))
	}
	response := dto.ModeratorsResponse{
		Admin:      dto.ToUserDTO(&group.Admin),
		Moderators: moderatorDTOs,
	}
	utils.Logger.WithFields(logrus.Fields{
		"group_id":        groupID,
		"admin_id":        group.AdminID,
		"moderator_count": len(moderatorDTOs),
	}).Debug("Fetched group moderators and admin")
	return response, nil
}

func (s *GroupService) GetUserByUsername(username string) (*models.User, error) {
	user, err := s.userRepo.GetByUsername(username)
	if err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"error":    err,
			"username": username,
		}).Error("Failed to fetch user by username")
		return nil, err
	}
	return user, nil
}

func (s *GroupService) GetUserGroupIDs(username string) ([]int32, error) {
	user, err := s.userRepo.GetByUsername(username)
	if err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"error":    err,
			"username": username,
		}).Error("Failed to fetch user by username")
		return nil, err
	}
	groupUsers, err := s.groupuserRepo.FindByUserID(user.UserID)
	if err != nil {
		return nil, err
	}
	var groupIDs []int32
	for _, gu := range groupUsers {
		groupIDs = append(groupIDs, gu.GroupID)
	}
	utils.Logger.WithFields(logrus.Fields{
		"username":    username,
		"group_count": len(groupIDs),
	}).Debug("Fetched user's group IDs")
	return groupIDs, nil
}

func (s *GroupService) GetUserGroups(username string, page, pageSize int) (*dto.GetGroupsResponse, error) {
	if page < 1 || pageSize < 1 || pageSize > 100 {
		return nil, errors.New("invalid page or page_size")
	}

	groups, total, err := s.groupRepo.FindUserGroups(username, page, pageSize)
	if err != nil {
		return nil, err
	}

	var groupDTOs []dto.GroupDTO
	for _, group := range groups {
		groupDTOs = append(groupDTOs, dto.ToGroupDTO(group))
	}

	totalPages := (total + int64(pageSize-1)) / int64(pageSize)
	response := &dto.GetGroupsResponse{
		Groups: groupDTOs,
		Pagination: dto.PaginationMeta{
			Page:     page,
			PageSize: pageSize,
			Total:    total,
			Pages:    totalPages,
		},
	}

	utils.Logger.WithFields(logrus.Fields{
		"username":  username,
		"page":      page,
		"page_size": pageSize,
		"count":     len(groupDTOs),
	}).Debug("Fetched user's groups")
	return response, nil
}
