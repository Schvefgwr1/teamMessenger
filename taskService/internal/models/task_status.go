package models

type TaskStatus struct {
	ID   int    `gorm:"primaryKey;autoIncrement"`
	Name string `gorm:"size:50;not null;unique"`
}

func (TaskStatus) TableName() string {
	return "task_service.task_statuses"
}
