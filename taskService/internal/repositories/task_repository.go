package repositories

import (
	"gorm.io/gorm"
	"taskService/internal/custom_errors"
	"taskService/internal/handlers/dto"
	"taskService/internal/models"
)

type TaskRepository interface {
	Create(task *models.Task) error
	UpdateStatus(taskID int, statusID int) error
	GetByID(taskID int) (*models.Task, error)
	GetUserTasks(userID string, limit, offset int) (*[]dto.TaskToList, error)
}

type taskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &taskRepository{db: db}
}

func (r *taskRepository) Create(task *models.Task) error {
	return r.db.Omit("Status").Create(task).Error
}

func (r *taskRepository) UpdateStatus(taskID int, statusID int) error {
	result := r.db.Model(&models.Task{}).Where("id = ?", taskID).Update("status", statusID)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return custom_errors.NewTaskNotFoundError(taskID)
	}
	return nil
}

func (r *taskRepository) GetByID(taskID int) (*models.Task, error) {
	var task models.Task
	err := r.db.Preload("Status").Preload("Files").First(&task, taskID).Error
	return &task, err
}

func (r *taskRepository) GetUserTasks(userID string, limit, offset int) (*[]dto.TaskToList, error) {
	var tasks []dto.TaskToList

	err := r.db.
		Table("task_service.tasks AS t").
		Select("t.id, t.title, s.name AS status").
		Joins("JOIN task_service.task_statuses s ON t.status = s.id").
		Where("t.executor_id = ?", userID).
		Limit(limit).
		Offset(offset).
		Scan(&tasks).Error

	return &tasks, err
}
