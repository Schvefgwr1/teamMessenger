package controllers

import (
	"taskService/internal/custom_errors"
	"taskService/internal/models"
	"taskService/internal/repositories"
)

type TaskStatusController struct {
	repo repositories.TaskStatusRepository
}

func NewTaskStatusController(repo repositories.TaskStatusRepository) *TaskStatusController {
	return &TaskStatusController{repo: repo}
}

func (c *TaskStatusController) Create(name string) (*models.TaskStatus, error) {
	existing, err := c.repo.GetByName(name)
	if err == nil && existing != nil {
		return nil, custom_errors.ErrStatusAlreadyExists
	}
	return c.repo.Create(name)
}

func (c *TaskStatusController) GetByID(id int) (*models.TaskStatus, error) {
	status, err := c.repo.GetByID(id)
	if err != nil {
		return nil, custom_errors.NewTaskStatusNotFoundError(string(rune(id)))
	}
	return status, nil
}

func (c *TaskStatusController) DeleteByID(id int) error {
	return c.repo.DeleteByID(id)
}

func (c *TaskStatusController) GetAll() ([]models.TaskStatus, error) {
	return c.repo.GetAll()
}
