package repositories

import (
	"gorm.io/gorm"
	"taskService/internal/models"
)

type TaskFileRepository interface {
	BulkCreate(taskFiles []models.TaskFile) error
}

type taskFileRepository struct {
	db *gorm.DB
}

func NewTaskFileRepository(db *gorm.DB) TaskFileRepository {
	return &taskFileRepository{db: db}
}

func (r *taskFileRepository) BulkCreate(taskFiles []models.TaskFile) error {
	return r.db.Create(&taskFiles).Error
}
