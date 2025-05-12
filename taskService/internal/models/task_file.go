package models

type TaskFile struct {
	TaskID int `gorm:"primaryKey"`
	FileID int `gorm:"primaryKey"`
}

func (TaskFile) TableName() string {
	return "task_service.task_files"
}
