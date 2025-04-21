package models

type Permission struct {
	ID          int    `gorm:"primaryKey"`
	Name        string `gorm:"unique;not null"`
	Description string
}

func (Permission) TableName() string {
	return "user_service.permissions"
}
