package services

import (
	"space/models"
	"space/repositories"
)

type TaskService struct {
	taskRepo *repositories.TaskRepository
}

func NewTaskService(taskRepo *repositories.TaskRepository) *TaskService {
	return &TaskService{taskRepo}
}

func (s *TaskService) GetTaskByID(id int32) (*models.Task, error) {
	return s.taskRepo.GetByID(id)
}

func (s *TaskService) CreateTask(task *models.Task) error {
	// You can add validation here if needed
	return s.taskRepo.Create(task)
}

func (s *TaskService) UpdateTask(task *models.Task) error {
	// Add any business logic or validation here
	return s.taskRepo.Update(task)
}
