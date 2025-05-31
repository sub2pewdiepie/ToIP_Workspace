package repositories

import (
	"space/models"

	"gorm.io/gorm"
)

type AcademicGroupRepository struct {
	db *gorm.DB
}

func NewAcademicGroupRepository(db *gorm.DB) *AcademicGroupRepository {
	return &AcademicGroupRepository{db}
}

func (r *AcademicGroupRepository) GetByID(id int32) (*models.AcademicGroup, error) {
	var academicGroup models.AcademicGroup
	if err := r.db.First(&academicGroup, "academic_group_id = ?", id).Error; err != nil {
		return nil, err
	}
	return &academicGroup, nil
}

func (r *AcademicGroupRepository) Create(academicGroup *models.AcademicGroup) error {
	return r.db.Create(academicGroup).Error
}

func (r *AcademicGroupRepository) Update(academicGroup *models.AcademicGroup) error {
	return r.db.Save(academicGroup).Error
}

func (r *AcademicGroupRepository) Delete(id int32) error {
	return r.db.Delete(&models.AcademicGroup{}, "academic_group_id = ?", id).Error
}
func (r *AcademicGroupRepository) GetByName(name string, academicGroup *models.AcademicGroup) error {
	return r.db.Where("name = ?", name).First(academicGroup).Error
}
