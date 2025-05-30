package services

import (
	"errors"
	"space/models"
	"space/repositories"
)

// type GroupApplicationService struct {
// 	repo          *repositories.GroupApplicationRepository
// 	groupRepo     *repositories.GroupRepository
// 	groupUserRepo *repositories.GroupUserRepository
// }

// func NewGroupApplicationService(
// 	repo *repositories.GroupApplicationRepository,
// 	groupRepo *repositories.GroupRepository,
// 	groupUserRepo *repositories.GroupUserRepository,
// ) *GroupApplicationService {
// 	return &GroupApplicationService{repo, groupRepo, groupUserRepo}
// }

type GroupApplicationService struct {
	repo           *repositories.GroupApplicationRepository
	groupRepo      *repositories.GroupRepository
	groupModerRepo *repositories.GroupModerRepository
	userRepo       repositories.UserRepository
	groupUserRepo  *repositories.GroupUserRepository
}

func NewGroupApplicationService(repo *repositories.GroupApplicationRepository,
	groupRepo *repositories.GroupRepository,
	groupModerRepo *repositories.GroupModerRepository,
	userRepo repositories.UserRepository,
	groupUserRepo *repositories.GroupUserRepository) *GroupApplicationService {

	return &GroupApplicationService{
		repo:           repo,
		groupRepo:      groupRepo,
		groupModerRepo: groupModerRepo,
		userRepo:       userRepo,
		groupUserRepo:  groupUserRepo,
	}
}

func (s *GroupApplicationService) ApplyToGroup(userID, groupID int32, message string) error {
	// Check if already a member
	isMember, err := s.groupUserRepo.IsMember(groupID, userID)
	if err != nil {
		return err
	}
	if isMember {
		return errors.New("user is already a group member")
	}

	// Check if there's already a pending application
	exists, err := s.repo.ExistsPending(groupID, userID)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("application already submitted and pending")
	}

	// Create application
	app := &models.GroupApplication{
		GroupID: groupID,
		UserID:  userID,
		Status:  "pending",
		Message: message,
	}
	return s.repo.Create(app)
}

func (s *GroupApplicationService) GetPendingApplications(userID int32) ([]models.GroupApplication, error) {
	// Fetch groups user moderates or admins
	groupIDs, err := s.groupRepo.GetGroupsManagedBy(userID)
	if err != nil {
		return nil, err
	}

	var allApps []models.GroupApplication
	for _, gid := range groupIDs {
		apps, err := s.repo.GetPendingByGroup(gid.ID)
		if err != nil {
			return nil, err
		}
		allApps = append(allApps, apps...)
	}
	return allApps, nil
}

//	func (s *GroupApplicationService) ReviewApplication(appID int32, status string) error {
//		if status != "approved" && status != "rejected" {
//			return errors.New("invalid status")
//		}
//		return s.repo.UpdateStatus(appID, status)
//	}
func (s *GroupApplicationService) ReviewApplication(groupID, userID int32, username, status string) error {
	if status != "approved" && status != "rejected" {
		return errors.New("invalid status")
	}

	// Check permissions
	isAuthorized, err := s.groupRepo.IsAdminOrModerator(groupID, userID)
	if err != nil {
		return err
	}
	if !isAuthorized {
		return errors.New("unauthorized: not an admin or moderator")
	}

	// Update status
	return s.repo.UpdateStatusByGroupAndUser(groupID, userID, status)
}

func (s *GroupApplicationService) ApproveApplication(groupID, userID int32, actingUserID int32) error {
	isAuthorized, err := s.groupRepo.IsAdminOrModerator(groupID, actingUserID)
	if err != nil || !isAuthorized {
		return errors.New("not authorized to approve applications")
	}

	// Check application exists and is pending
	app, err := s.repo.GetByGroupAndUser(groupID, userID)
	if err != nil || app.Status != "pending" {
		return errors.New("no pending application found")
	}

	// Approve application
	err = s.repo.UpdateStatus(app.ApplicationID, "approved")
	if err != nil {
		return err
	}

	// Add to GroupUser
	return s.groupUserRepo.Create(&models.GroupUser{
		GroupID: groupID,
		UserID:  userID,
	})
}

func (s *GroupApplicationService) RejectApplication(groupID, userID int32, actingUserID int32) error {
	isAuthorized, err := s.groupRepo.IsAdminOrModerator(groupID, actingUserID)
	if err != nil || !isAuthorized {
		return errors.New("not authorized to reject applications")
	}

	app, err := s.repo.GetByGroupAndUser(groupID, userID)
	if err != nil || app.Status != "pending" {
		return errors.New("no pending application found")
	}

	return s.repo.UpdateStatus(app.ApplicationID, "rejected")
}
