package repositories

import (
	"space/models"

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
	return r.db.Create(subject).Error
}

func (r *SubjectRepository) Update(subject *models.Subject) error {
	return r.db.Save(subject).Error
}

func (r *SubjectRepository) Delete(id int32) error {
	return r.db.Delete(&models.Subject{}, "subject_id = ?", id).Error
}
