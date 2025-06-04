package repositories

import (
	"space/models"
	"space/utils"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type GroupApplicationRepository struct {
	db *gorm.DB
}

func NewGroupApplicationRepository(db *gorm.DB) *GroupApplicationRepository {
	return &GroupApplicationRepository{db}
}

func (r *GroupApplicationRepository) Create(app *models.GroupApplication) error {
	return r.db.Create(app).Error
}

func (r *GroupApplicationRepository) GetPendingByGroup(groupID int32) ([]models.GroupApplication, error) {
	var apps []models.GroupApplication
	err := r.db.
		Preload("User").
		Where("group_id = ? AND status = ?", groupID, "pending").
		Find(&apps).Error
	return apps, err
}

func (r *GroupApplicationRepository) UpdateStatus(appID int32, status string) error {
	return r.db.Model(&models.GroupApplication{}).
		Where("application_id = ?", appID).
		Update("status", status).Error
}

func (r *GroupApplicationRepository) ExistsPending(groupID, userID int32) (bool, error) {
	var count int64
	err := r.db.Model(&models.GroupApplication{}).
		Where("group_id = ? AND user_id = ? AND status = ?", groupID, userID, "pending").
		Count(&count).Error
	return count > 0, err
}
func (r *GroupApplicationRepository) GetByGroupAndUser(groupID, userID int32) (*models.GroupApplication, error) {
	var app models.GroupApplication
	err := r.db.First(&app, "group_id = ? AND user_id = ?", groupID, userID).Error
	if err != nil {
		return nil, err
	}
	return &app, nil
}

func (r *GroupApplicationRepository) UpdateStatusByGroupAndUser(groupID, userID int32, status string) error {
	return r.db.Model(&models.GroupApplication{}).
		Where("group_id = ? AND user_id = ?", groupID, userID).
		Update("status", status).Error
}

func (r *AcademicGroupRepository) FindAll() ([]*models.AcademicGroup, error) {
	var groups []*models.AcademicGroup
	if err := r.db.Find(&groups).Error; err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"error": err,
		}).Error("Failed to fetch all academic groups")
		return nil, err
	}
	return groups, nil
}
