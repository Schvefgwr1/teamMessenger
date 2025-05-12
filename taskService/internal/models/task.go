package models

import (
	"github.com/google/uuid"
	"time"
)

type Task struct {
	ID          int       `gorm:"primaryKey;autoIncrement"`
	Title       string    `gorm:"size:255;not null"`
	Description string    `gorm:"type:text"`
	StatusID    int       `gorm:"column:status"`
	CreatorID   uuid.UUID `gorm:"type:uuid"`
	ExecutorID  uuid.UUID `gorm:"type:uuid"`
	ChatID      uuid.UUID `gorm:"type:uuid"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`

	Status *TaskStatus `gorm:"foreignKey:StatusID"`
	Files  []TaskFile  `gorm:"foreignKey:TaskID"`
}

func (Task) TableName() string {
	return "task_service.tasks"
}
