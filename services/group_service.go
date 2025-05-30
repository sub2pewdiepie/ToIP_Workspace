package services

import (
	"errors"
	"space/models"
	"space/models/dto"
	"space/repositories"

	"github.com/gin-gonic/gin"
)

type GroupService struct {
	groupRepo *repositories.GroupRepository
	userRepo  repositories.UserRepository // интерфейс!!!
}

func NewGroupService(groupRepo *repositories.GroupRepository, userRepo repositories.UserRepository) *GroupService {
	return &GroupService{groupRepo, userRepo}
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
		return errors.New("name and academic_group_id are required")
	}

	// Get authenticated user's username from context
	username, exists := c.Get("username")
	if !exists {
		return errors.New("user not authenticated")
	}

	// Fetch user to get their ID
	user, err := s.userRepo.GetByUsername(username.(string))
	if err != nil {
		return errors.New("failed to find authenticated user")
	}

	// Set AdminID to the authenticated user's ID
	group.AdminID = user.UserID

	return s.groupRepo.Create(group)
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
