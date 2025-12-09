package controllers

import (
	"taskService/internal/handlers/dto"
	"taskService/internal/models"
)

// TaskControllerInterface - интерфейс для TaskController для возможности мокирования
type TaskControllerInterface interface {
	Create(taskDTO *dto.CreateTaskDTO) (*models.Task, error)
	UpdateStatus(taskID, statusID int) error
	GetByID(taskID int) (*dto.TaskResponse, error)
	GetUserTasks(userID string, limit, offset int) (*[]dto.TaskToList, error)
}

// TaskStatusControllerInterface - интерфейс для TaskStatusController для возможности мокирования
type TaskStatusControllerInterface interface {
	Create(name string) (*models.TaskStatus, error)
	GetByID(id int) (*models.TaskStatus, error)
	DeleteByID(id int) error
	GetAll() ([]models.TaskStatus, error)
}
