package chat_contracts

import (
	"github.com/google/uuid"
	"time"
)

type Chat struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name         string    `gorm:"size:255;not null"`
	IsGroup      bool      `gorm:"default:false"`
	Description  *string
	AvatarFileID *int
	CreatedAt    time.Time
}
