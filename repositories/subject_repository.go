package repositories

import (
	"space/models"
	"space/utils"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type SubjectRepository struct {
	db *gorm.DB
}

func NewSubjectRepository(db *gorm.DB) *SubjectRepository {
	return &SubjectRepository{db}
}

func (r *SubjectRepository) GetByID(id int32) (*models.Subject, error) {
	var subject models.Subject
	if err := r.db.Preload("Group").First(&subject, "subject_id = ?", id).Error; err != nil {
		return nil, err
	}
	return &subject, nil
}

func (r *SubjectRepository) Create(subject *models.Subject) error {
	if err := r.db.Create(subject).Error; err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"error":             err,
			"name":              subject.Name,
			"academic_group_id": subject.AcademicGroupID,
		}).Error("Failed to create subject")
		return err
	}
	utils.Logger.WithFields(logrus.Fields{
		"subject_id":        subject.SubjectID,
		"name":              subject.Name,
		"academic_group_id": subject.AcademicGroupID,
	}).Debug("Subject created")
	return nil
}
func (r *SubjectRepository) Update(subject *models.Subject) error {
	return r.db.Save(subject).Error
}

func (r *SubjectRepository) Delete(id int32) error {
	return r.db.Delete(&models.Subject{}, "subject_id = ?", id).Error
}

func (r *SubjectRepository) FindByName(name string) (*models.Subject, error) {
	var subject models.Subject
	if err := r.db.Where("name = ?", name).First(&subject).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, err
		}
		utils.Logger.WithFields(logrus.Fields{
			"error": err,
			"name":  name,
		}).Error("Failed to find subject by name")
		return nil, err
	}
	return &subject, nil
}

func (r *SubjectRepository) FindByAcademicGroupID(academicGroupID int32, page, pageSize int) ([]models.Subject, int64, error) {
	var subjects []models.Subject
	var total int64

	query := r.db.Model(&models.Subject{}).Where("academic_group_id = ?", academicGroupID).
		Preload("AcademicGroup")

	if err := query.Count(&total).Error; err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"error":             err,
			"academic_group_id": academicGroupID,
		}).Error("Failed to count subjects")
		return nil, 0, err
	}

	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&subjects).Error; err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"error":             err,
			"academic_group_id": academicGroupID,
		}).Error("Failed to find subjects")
		return nil, 0, err
	}

	return subjects, total, nil
}

func (r *SubjectRepository) FindByUserGroups(userID int32, page, pageSize int) ([]models.Subject, []models.Group, int64, error) {
	var subjects []models.Subject
	var groups []models.Group
	var total int64

	// Get groups the user is a member of
	if err := r.db.Model(&models.Group{}).
		Joins("JOIN group_users ON groups.id = group_users.group_id").
		Where("group_users.user_id = ?", userID).
		Find(&groups).Error; err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"error":   err,
			"user_id": userID,
		}).Error("Failed to find user groups")
		return nil, nil, 0, err
	}

	if len(groups) == 0 {
		return nil, nil, 0, nil
	}

	// Get unique academic_group_ids
	academicGroupIDs := make([]int32, 0, len(groups))
	for _, group := range groups {
		academicGroupIDs = append(academicGroupIDs, group.AcademicGroupID)
	}

	query := r.db.Model(&models.Subject{}).
		Where("academic_group_id IN ?", academicGroupIDs).
		Preload("AcademicGroup")

	if err := query.Count(&total).Error; err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"error":   err,
			"user_id": userID,
		}).Error("Failed to count user subjects")
		return nil, nil, 0, err
	}

	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&subjects).Error; err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"error":   err,
			"user_id": userID,
		}).Error("Failed to find user subjects")
		return nil, nil, 0, err
	}

	return subjects, groups, total, nil
}
