package repositories

import (
	"gorm.io/gorm"
	"taskService/internal/models"
)

type TaskStatusRepository interface {
	GetByID(id int) (*models.TaskStatus, error)
	GetByName(name string) (*models.TaskStatus, error)
	Create(name string) (*models.TaskStatus, error)
	DeleteByID(id int) error
	GetAll() ([]models.TaskStatus, error)
}

type taskStatusRepository struct {
	db *gorm.DB
}

func NewTaskStatusRepository(db *gorm.DB) TaskStatusRepository {
	return &taskStatusRepository{db: db}
}

func (r *taskStatusRepository) GetByID(id int) (*models.TaskStatus, error) {
	var status models.TaskStatus
	if err := r.db.First(&status, id).Error; err != nil {
		return nil, err
	}
	return &status, nil
}

func (r *taskStatusRepository) GetByName(name string) (*models.TaskStatus, error) {
	var status models.TaskStatus
	err := r.db.Where("name = ?", name).First(&status).Error
	if err != nil {
		return nil, err
	}
	return &status, nil
}

func (r *taskStatusRepository) Create(name string) (*models.TaskStatus, error) {
	status := &models.TaskStatus{Name: name}
	if err := r.db.Create(status).Error; err != nil {
		return nil, err
	}
	return status, nil
}

func (r *taskStatusRepository) DeleteByID(id int) error {
	return r.db.Delete(&models.TaskStatus{}, id).Error
}

func (r *taskStatusRepository) GetAll() ([]models.TaskStatus, error) {
	var statuses []models.TaskStatus
	if err := r.db.Find(&statuses).Error; err != nil {
		return nil, err
	}
	return statuses, nil
}
