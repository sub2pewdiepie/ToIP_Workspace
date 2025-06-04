package repositories

import (
	"space/models"
	"space/utils"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type GroupModerRepository struct {
	db *gorm.DB
}

func NewGroupModerRepository(db *gorm.DB) *GroupModerRepository {
	return &GroupModerRepository{db}
}

func (r *GroupModerRepository) GetByID(groupID, userID int32) (*models.GroupModer, error) {
	var groupModer models.GroupModer
	if err := r.db.First(&groupModer, "group_id = ? AND user_id = ?", groupID, userID).Error; err != nil {
		return nil, err
	}
	return &groupModer, nil
}

func (r *GroupModerRepository) Create(groupModer *models.GroupModer) error {
	return r.db.Create(groupModer).Error
}

func (r *GroupModerRepository) Delete(groupID, userID int32) error {
	return r.db.Delete(&models.GroupModer{}, "group_id = ? AND user_id = ?", groupID, userID).Error
}

func (r *GroupModerRepository) FindByGroupID(groupID int32) ([]*models.GroupModer, error) {
	var groupModers []*models.GroupModer
	if err := r.db.Preload("User").Where("group_id = ?", groupID).Find(&groupModers).Error; err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"error":    err,
			"group_id": groupID,
		}).Error("Failed to find group moderators")
		return nil, err
	}
	return groupModers, nil
}

func (r *GroupModerRepository) IsModerator(groupID, userID int32) (bool, error) {
	var groupModer models.GroupModer
	if err := r.db.Where("group_id = ? AND user_id = ?", groupID, userID).First(&groupModer).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		utils.Logger.WithFields(logrus.Fields{
			"error":    err,
			"group_id": groupID,
			"user_id":  userID,
		}).Error("Failed to check moderator status")
		return false, err
	}
	return true, nil
}
