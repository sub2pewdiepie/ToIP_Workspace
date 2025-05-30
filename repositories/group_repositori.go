package repositories

import (
	"space/models"

	"gorm.io/gorm"
)

type GroupRepository struct {
	db *gorm.DB
}

func NewGroupRepository(db *gorm.DB) *GroupRepository {
	return &GroupRepository{db}
}

func (r *GroupRepository) GetByID(id int32) (*models.Group, error) {
	var group models.Group
	if err := r.db.Preload("AcademicGroup").Preload("Admin").First(&group, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &group, nil
}
func (r *GroupRepository) DB() *gorm.DB {
	return r.db
}

func (r *GroupRepository) GetAllWithPagination(page, pageSize int32) ([]models.Group, int64, error) {
	var groups []models.Group
	var total int64

	offset := (page - 1) * pageSize

	err := r.db.Preload("Admin", func(db *gorm.DB) *gorm.DB {
		return db.Select("user_id", "username")
	}).Preload("AcademicGroup").
		Limit(int(pageSize)).
		Offset(int(offset)).
		Find(&groups).
		Offset(-1).
		Limit(-1).
		Count(&total).
		Error

	if err != nil {
		return nil, 0, err
	}
	return groups, total, nil
}

func (r *GroupRepository) Create(group *models.Group) error {
	return r.db.Create(group).Error
}

func (r *GroupRepository) Update(group *models.Group) error {
	return r.db.Save(group).Error
}

func (r *GroupRepository) Delete(id int32) error {
	return r.db.Delete(&models.Group{}, "id = ?", id).Error
}

// depricated, for future reference
func (r *GroupRepository) OldGetGroupsManagedBy(userID int32) ([]models.Group, error) {
	var groups []models.Group

	err := r.db.
		Preload("Admin").
		Preload("AcademicGroup").
		Where("admin_id = ?", userID).
		Or("id IN (?)", r.db.Model(&models.GroupModer{}).Select("group_id").Where("user_id = ?", userID)).
		Find(&groups).Error

	if err != nil {
		return nil, err
	}

	return groups, nil
}

func (r *GroupRepository) GetGroupsManagedBy(userID int32) ([]models.Group, error) {
	var groups []models.Group

	err := r.db.
		Joins("LEFT JOIN group_moders gm ON groups.id = gm.group_id").
		Where("groups.admin_id = ? OR gm.user_id = ?", userID, userID).
		Preload("Admin").
		Preload("AcademicGroup").
		Find(&groups).Error

	return groups, err
}
func (r *GroupRepository) IsAdminOrModerator(groupID, userID int32) (bool, error) {
	var count int64

	// Check if user is group admin
	err := r.db.Model(&models.Group{}).
		Where("id = ? AND admin_id = ?", groupID, userID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	if count > 0 {
		return true, nil
	}

	// Check if user is a moderator
	err = r.db.Model(&models.GroupModer{}).
		Where("group_id = ? AND user_id = ?", groupID, userID).
		Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
