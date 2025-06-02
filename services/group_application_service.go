package services

import (
	"errors"
	"space/models"
	"space/repositories"
	"space/utils"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
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

func (s *GroupApplicationService) ApplyToGroup(c *gin.Context, groupID int32, message string) error {
	username, exists := c.Get("username")
	if !exists {
		utils.Logger.Error("User not authenticated")
		return errors.New("user not authenticated")
	}

	user, err := s.userRepo.GetByUsername(username.(string))
	if err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"error":    err,
			"username": username,
		}).Error("Failed to find authenticated user")
		return errors.New("failed to find authenticated user")
	}
	utils.Logger.WithFields(logrus.Fields{
		"username": username,
		"user_id":  user.UserID,
		"group_id": groupID,
	}).Debug("Checking group membership")

	isMember, err := s.groupUserRepo.IsMember(groupID, user.UserID)
	if err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"error":    err,
			"user_id":  user.UserID,
			"group_id": groupID,
		}).Error("Failed to check group membership")
		return err
	}
	if isMember {
		utils.Logger.WithFields(logrus.Fields{
			"error":    err,
			"user_id":  user.UserID,
			"group_id": groupID,
		}).Error("Failed to check group membership")
		return errors.New("user is already a group member")
	}

	exists, err = s.repo.ExistsPending(groupID, user.UserID)
	if err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"error":    err,
			"user_id":  user.UserID,
			"group_id": groupID,
		}).Error("Failed to check pending application")
		return err
	}
	if exists {
		utils.Logger.WithFields(logrus.Fields{
			"user_id":  user.UserID,
			"group_id": groupID,
		}).Warn("Application already submitted and pending")
		return errors.New("application already submitted and pending")
	}

	app := &models.GroupApplication{
		GroupID: groupID,
		UserID:  user.UserID,
		Status:  "pending",
		Message: message,
	}
	utils.Logger.WithFields(logrus.Fields{
		"user_id":  user.UserID,
		"group_id": groupID,
	}).Debug("Creating group application")

	if err := s.repo.Create(app); err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"error":    err,
			"user_id":  user.UserID,
			"group_id": groupID,
		}).Error("Failed to create group application")
		return err
	}

	return nil

	// return s.repo.Create(app)
}

// стоит и дальше добавить логгирование
func (s *GroupApplicationService) GetPendingApplications(c *gin.Context) ([]models.GroupApplication, error) {
	// Fetch groups user moderates or admins
	username, exists := c.Get("username")
	if !exists {
		return nil, errors.New("user not authenticated")
	}

	user, err := s.userRepo.GetByUsername(username.(string))
	if err != nil {
		return nil, errors.New("failed to find authenticated user")
	}

	groupIDs, err := s.groupRepo.GetGroupsManagedBy(user.UserID)
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
func (s *GroupApplicationService) ReviewApplication(groupID int32, targetUsername, reviewerUsername, status string) error {
	utils.Logger.WithFields(logrus.Fields{
		"reviewer": reviewerUsername,
		"username": targetUsername,
		"group_id": groupID,
		"status":   status,
	}).Debug("Processing application review")

	if status != "approved" && status != "rejected" {
		utils.Logger.WithField("status", status).Error("Invalid status")
		return errors.New("invalid status")
	}

	reviewer, err := s.userRepo.GetByUsername(reviewerUsername)
	if err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"error":    err,
			"reviewer": reviewerUsername,
		}).Error("Failed to find reviewer")
		return errors.New("failed to find reviewer")
	}

	targetUser, err := s.userRepo.GetByUsername(targetUsername)
	if err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"error":    err,
			"username": targetUsername,
		}).Error("Failed to find target user")
		return errors.New("failed to find target user")
	}

	isAuthorized, err := s.groupRepo.IsAdminOrModerator(groupID, reviewer.UserID)
	if err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"error":       err,
			"reviewer_id": reviewer.UserID,
			"group_id":    groupID,
		}).Error("Failed to check authorization")
		return err
	}
	if !isAuthorized {
		utils.Logger.WithFields(logrus.Fields{
			"reviewer_id": reviewer.UserID,
			"group_id":    groupID,
		}).Error("Unauthorized: not an admin or moderator")
		return errors.New("unauthorized: not an admin or moderator")
	}

	utils.Logger.WithFields(logrus.Fields{
		"group_id": groupID,
		"user_id":  targetUser.UserID,
	}).Debug("Updating application status")

	if err := s.repo.UpdateStatusByGroupAndUser(groupID, targetUser.UserID, status); err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"error":    err,
			"group_id": groupID,
			"user_id":  targetUser.UserID,
		}).Error("Failed to update application status")
		return err
	}

	if status == "approved" {
		utils.Logger.WithFields(logrus.Fields{
			"group_id": groupID,
			"user_id":  targetUser.UserID,
		}).Debug("Creating group user for approved application")
		if err := s.groupUserRepo.Create(&models.GroupUser{
			GroupID: groupID,
			UserID:  targetUser.UserID,
		}); err != nil {
			utils.Logger.WithFields(logrus.Fields{
				"error":    err,
				"group_id": groupID,
				"user_id":  targetUser.UserID,
			}).Error("Failed to create group user")
			return err
		}
	}

	return nil
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
