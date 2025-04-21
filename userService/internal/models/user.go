package models

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Username     string    `gorm:"unique;not null"`
	Email        string    `gorm:"unique;not null"`
	PasswordHash string    `gorm:"not null"`
	Description  *string
	Gender       *string
	Age          *int
	AvatarFileID *int      `gorm:"column:avatar_file_id"`
	RoleID       int       `gorm:"column:role_id"`
	Role         Role      `gorm:"foreignKey:RoleID;references:ID"`
	CreatedAt    time.Time `gorm:"default:now()"`
	UpdatedAt    time.Time `gorm:"default:now()"`
}

func (User) TableName() string {
	return "user_service.users"
}
