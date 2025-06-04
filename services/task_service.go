package services

import (
	"space/models"
	"space/models/dto"
	"space/repositories"
	"space/utils"
	"time"

	"github.com/sirupsen/logrus"
)

type TaskService struct {
	taskRepo *repositories.TaskRepository
}

func NewTaskService(taskRepo *repositories.TaskRepository) *TaskService {
	return &TaskService{taskRepo}
}

func (s *TaskService) UpdateTask(task *models.Task) error {
	// Add any business logic or validation here
	return s.taskRepo.Update(task)
}
func (s *TaskService) CreateTask(groupID, userID int32, title, description string, deadline *time.Time, subjectID *int32) error {
	task := &models.Task{
		GroupID:     groupID,
		UserID:      userID,
		SubjectID:   subjectID,
		Title:       title,
		Description: description,
		IsVerified:  false,
		Deadline:    deadline,
	}
	if err := s.taskRepo.Create(task); err != nil {
		return err
	}
	utils.Logger.WithFields(logrus.Fields{
		"task_id":    task.ID,
		"group_id":   groupID,
		"user_id":    userID,
		"subject_id": subjectID,
	}).Info("Task created successfully")
	return nil
}
func (s *TaskService) GetGroupTasks(groupID int32) ([]dto.TaskDTO, error) {
	tasks, err := s.taskRepo.FindByGroupID(groupID)
	if err != nil {
		return nil, err
	}
	var taskDTOs []dto.TaskDTO
	for _, task := range tasks {
		taskDTOs = append(taskDTOs, dto.ToTaskDTO(task))
	}
	utils.Logger.WithFields(logrus.Fields{
		"group_id": groupID,
		"count":    len(taskDTOs),
	}).Debug("Fetched group tasks")
	return taskDTOs, nil
}

func (s *TaskService) VerifyTask(taskID int32, isVerified bool) error {
	if err := s.taskRepo.UpdateVerificationStatus(taskID, isVerified); err != nil {
		return err
	}
	utils.Logger.WithFields(logrus.Fields{
		"task_id":     taskID,
		"is_verified": isVerified,
	}).Info("Task verification status updated")
	return nil
}
func (s *TaskService) GetTaskByID(taskID int32) (*models.Task, error) {
	task, err := s.taskRepo.GetByID(taskID)
	if err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"error":   err,
			"task_id": taskID,
		}).Error("Failed to fetch task by ID")
		return nil, err
	}
	return task, nil
}
func (s *TaskService) GetTasksByGroupIDs(groupIDs []int32) ([]dto.TaskDTO, error) {
	tasks, err := s.taskRepo.FindByGroupIDs(groupIDs)
	if err != nil {
		return nil, err
	}
	var taskDTOs []dto.TaskDTO
	for _, task := range tasks {
		taskDTOs = append(taskDTOs, dto.ToTaskDTO(task))
	}
	utils.Logger.WithFields(logrus.Fields{
		"group_ids": groupIDs,
		"count":     len(taskDTOs),
	}).Debug("Fetched tasks for multiple groups")
	return taskDTOs, nil
}

func (s *TaskService) DeleteTask(taskID int32) error {
	_, err := s.taskRepo.GetByID(taskID)
	if err != nil {
		utils.Logger.WithFields(logrus.Fields{
			"error":   err,
			"task_id": taskID,
		}).Error("Task not found")
		return err
	}
	if err := s.taskRepo.Delete(taskID); err != nil {
		return err
	}
	utils.Logger.WithFields(logrus.Fields{
		"task_id": taskID,
	}).Info("Task deleted successfully")
	return nil
}
