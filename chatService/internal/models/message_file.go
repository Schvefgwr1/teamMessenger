package models

import "github.com/google/uuid"

type MessageFile struct {
	MessageID uuid.UUID `gorm:"primaryKey;type:uuid"`
	FileID    int       `gorm:"primaryKey"`
}

func (MessageFile) TableName() string {
	return "chat_service.message_files"
}
