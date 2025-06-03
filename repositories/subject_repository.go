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
