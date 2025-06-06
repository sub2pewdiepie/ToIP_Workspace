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

func (r *TaskRepository) GetByID(taskID int32) (*models.Task, error) {
	var task models.Task
	if err := r.db.Preload("User").Preload("Subject").First(&task, taskID).Error; err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"error":   err,
			"task_id": taskID,
		}).Error("Failed to find task by ID")
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

func (r *TaskRepository) OldFindByGroupID(groupID int32) ([]*models.Task, error) {
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

func (r *TaskRepository) OldNoPagFindByGroupID(groupID int32) ([]*models.Task, error) {
	var tasks []*models.Task
	if err := r.db.Preload("User").Preload("Subject").Where("group_id = ?", groupID).Find(&tasks).Error; err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"error":    err,
			"group_id": groupID,
		}).Error("Failed to find tasks by group ID")
		return nil, err
	}
	return tasks, nil
}

func (r *TaskRepository) FindByGroupID(groupID int32, page, pageSize int) ([]*models.Task, int64, error) {
	var tasks []*models.Task
	var total int64

	query := r.db.Model(&models.Task{}).Where("group_id = ?", groupID).
		Preload("User").Preload("Subject").Preload("Group").Preload("Group.AcademicGroup")

	if err := query.Count(&total).Error; err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"error":    err,
			"group_id": groupID,
		}).Error("Failed to count tasks")
		return nil, 0, err
	}

	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&tasks).Error; err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"error":    err,
			"group_id": groupID,
		}).Error("Failed to find tasks by group ID")
		return nil, 0, err
	}

	return tasks, total, nil
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

func (r *TaskRepository) OldFindByGroupIDs(groupIDs []int32) ([]*models.Task, error) {
	var tasks []*models.Task
	if err := r.db.Preload("User").Preload("Subject").Where("group_id IN ?", groupIDs).Find(&tasks).Error; err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"error":     err,
			"group_ids": groupIDs,
		}).Error("Failed to find tasks by group IDs")
		return nil, err
	}
	return tasks, nil
}

func (r *TaskRepository) FindByGroupIDs(groupIDs []int32, page, pageSize int) ([]*models.Task, int64, error) {
	var tasks []*models.Task
	var total int64

	query := r.db.Model(&models.Task{}).Where("group_id IN ?", groupIDs).
		Preload("User").Preload("Subject").Preload("Group").Preload("Group.AcademicGroup")

	if err := query.Count(&total).Error; err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"error":     err,
			"group_ids": groupIDs,
		}).Error("Failed to count tasks")
		return nil, 0, err
	}

	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&tasks).Error; err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"error":     err,
			"group_ids": groupIDs,
		}).Error("Failed to find tasks by group IDs")
		return nil, 0, err
	}

	return tasks, total, nil
}

func (r *TaskRepository) Delete(taskID int32) error {
	if err := r.db.Delete(&models.Task{}, taskID).Error; err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"error":   err,
			"task_id": taskID,
		}).Error("Failed to delete task")
		return err
	}
	return nil
}

func (r *TaskRepository) FindBySubjectID(subjectID, groupID int32, page, pageSize int) ([]*models.Task, int64, error) {
	var tasks []*models.Task
	var total int64

	query := r.db.Model(&models.Task{}).Where("subject_id = ? AND group_id = ?", subjectID, groupID).
		Preload("User").Preload("Subject").Preload("Group").Preload("Group.AcademicGroup")

	if err := query.Count(&total).Error; err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"error":      err,
			"subject_id": subjectID,
			"group_id":   groupID,
		}).Error("Failed to count tasks")
		return nil, 0, err
	}

	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&tasks).Error; err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"error":      err,
			"subject_id": subjectID,
			"group_id":   groupID,
		}).Error("Failed to find tasks")
		return nil, 0, err
	}

	return tasks, total, nil
}
