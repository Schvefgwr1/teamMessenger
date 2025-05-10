package models

import (
	"github.com/google/uuid"
	"time"
)

type Message struct {
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey"`
	ChatID    uuid.UUID  `gorm:"type:uuid;not null"`
	SenderID  *uuid.UUID `gorm:"type:uuid"`
	Content   string     `gorm:"type:text;not null"`
	UpdatedAt *time.Time
	CreatedAt time.Time

	Files []MessageFile `gorm:"foreignKey:MessageID"`
}

func (Message) TableName() string {
	return "chat_service.messages"
}
