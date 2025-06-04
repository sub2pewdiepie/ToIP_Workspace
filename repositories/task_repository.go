package repositories

import (
	"space/models"
	"space/utils"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type TaskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) *TaskRepository {
	return &TaskRepository{db}
}

func (r *TaskRepository) GetByID(id int32) (*models.Task, error) {
	var task models.Task
	if err := r.db.Preload("Subject").Preload("User").First(&task, "task_id = ?", id).Error; err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *TaskRepository) Create(task *models.Task) error {
	if err := r.db.Create(task).Error; err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"error": err,
			"task":  task.Title,
		}).Error("Failed to create task")
		return err
	}
	return nil
}

func (r *TaskRepository) Update(task *models.Task) error {
	return r.db.Save(task).Error
}

func (r *TaskRepository) FindByGroupID(groupID int32) ([]*models.Task, error) {
	var tasks []*models.Task
	if err := r.db.Preload("User").Where("group_id = ?", groupID).Find(&tasks).Error; err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"error":    err,
			"group_id": groupID,
		}).Error("Failed to find tasks by group ID")
		return nil, err
	}
	return tasks, nil
}

func (r *TaskRepository) UpdateVerificationStatus(taskID int32, isVerified bool) error {
	if err := r.db.Model(&models.Task{}).Where("id = ?", taskID).Update("is_verified", isVerified).Error; err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"error":       err,
			"task_id":     taskID,
			"is_verified": isVerified,
		}).Error("Failed to update task verification status")
		return err
	}
	return nil
}

func (r *TaskRepository) FindByGroupIDs(groupIDs []int32) ([]*models.Task, error) {
	var tasks []*models.Task
	if err := r.db.Preload("User").Where("group_id IN ?", groupIDs).Find(&tasks).Error; err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"error":     err,
			"group_ids": groupIDs,
		}).Error("Failed to find tasks by group IDs")
		return nil, err
	}
	return tasks, nil
}
