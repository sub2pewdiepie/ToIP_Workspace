package repositories

import (
	"space/models"

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
