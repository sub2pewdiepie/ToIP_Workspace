package repositories

import (
	"space/models"
	"space/utils"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type GroupRepository struct {
	db       *gorm.DB
	userRepo UserRepository
}

func NewGroupRepository(db *gorm.DB, userRepo UserRepository) *GroupRepository {
	return &GroupRepository{db: db, userRepo: userRepo}
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

//	func (r *GroupRepository) Create(group *models.Group) error {
//		return r.db.Create(group).Error
//	}
func (r *GroupRepository) Create(group *models.Group) error {
	if err := r.db.Create(group).Error; err != nil {
		return err
	}
	// Preload Admin and AcademicGroup after creation
	return r.db.Preload("AcademicGroup").Preload("Admin").First(group, "id = ?", group.ID).Error
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

// Where user can apply (not a member)
func (r *GroupRepository) GetAvailable(userID int32, page, pageSize int) ([]models.Group, int64, error) {
	var groups []models.Group
	var total int64

	offset := (page - 1) * pageSize

	query := r.db.Model(&models.Group{}).
		Preload("Admin").
		Preload("AcademicGroup").
		Where("id NOT IN (?)", r.db.Model(&models.Group{}).
			Select("id").
			Where("admin_id = ?", userID)).
		Where("id NOT IN (?)", r.db.Model(&models.GroupUser{}).
			Select("group_id").
			Where("user_id = ?", userID)).
		Where("id NOT IN (?)", r.db.Model(&models.GroupModer{}).
			Select("group_id").
			Where("user_id = ?", userID))

	if err := query.Count(&total).Error; err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"error":   err,
			"user_id": userID,
		}).Error("Failed to count available groups")
		return nil, 0, err
	}

	if err := query.Offset(offset).Limit(pageSize).Find(&groups).Error; err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"error":   err,
			"user_id": userID,
		}).Error("Failed to fetch available groups")
		return nil, 0, err
	}

	utils.Logger.WithFields(logrus.Fields{
		"user_id":   userID,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	}).Debug("Fetched available groups")
	return groups, total, nil
}

func (r *GroupRepository) GetByID(id int32) (*models.Group, error) {
	var group models.Group
	if err := r.db.Preload("Admin").Preload("AcademicGroup").First(&group, id).Error; err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"error": err,
			"id":    id,
		}).Error("Failed to find group")
		return nil, err
	}
	return &group, nil
}

func (r *GroupRepository) FindUserGroups(username string, page, pageSize int) ([]*models.Group, int64, error) {
	// Get user_id from username
	user, err := r.userRepo.GetByUsername(username)
	if err != nil {
		return nil, 0, err
	}
	userID := user.UserID

	var groups []*models.Group
	var total int64

	// Count total groups (member or moderator)
	query := r.db.Model(&models.Group{}).
		Joins("LEFT JOIN group_users ON groups.id = group_users.group_id").
		Joins("LEFT JOIN group_moders ON groups.id = group_moders.group_id").
		Where("group_users.user_id = ? OR group_moders.user_id = ?", userID, userID).
		Preload("Admin").Preload("AcademicGroup")

	if err := query.Count(&total).Error; err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"error":    err,
			"username": username,
		}).Error("Failed to count user's groups")
		return nil, 0, err
	}

	// Fetch paginated groups
	err = query.
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&groups).Error
	if err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"error":     err,
			"username":  username,
			"page":      page,
			"page_size": pageSize,
		}).Error("Failed to fetch user's groups")
		return nil, 0, err
	}
	return groups, total, nil
}
