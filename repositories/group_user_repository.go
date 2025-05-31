package repositories

import (
	"space/models"

	"gorm.io/gorm"
)

type GroupUserRepository struct {
	db *gorm.DB
}

func NewGroupUserRepository(db *gorm.DB) *GroupUserRepository {
	return &GroupUserRepository{db}
}

func (r *GroupUserRepository) GetByID(groupID, userID int32) (*models.GroupUser, error) {
	var groupUser models.GroupUser
	if err := r.db.Preload("Group").Preload("User").First(&groupUser, "group_id = ? AND user_id = ?", groupID, userID).Error; err != nil {
		return nil, err
	}
	return &groupUser, nil
}

func (r *GroupUserRepository) Create(groupUser *models.GroupUser) error {
	return r.db.Create(groupUser).Error
}

func (r *GroupUserRepository) Update(groupUser *models.GroupUser) error {
	return r.db.Save(groupUser).Error
}

func (r *GroupUserRepository) Delete(groupID, userID int32) error {
	return r.db.Delete(&models.GroupUser{}, "group_id = ? AND user_id = ?", groupID, userID).Error
}
func (r *GroupUserRepository) IsMember(groupID, userID int32) (bool, error) {
	var count int64
	err := r.db.Model(&models.GroupUser{}).
		Where("group_id = ? AND user_id = ?", groupID, userID).
		Count(&count).Error
	return count > 0, err
}
