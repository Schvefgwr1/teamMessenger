package dto

import (
	fc "common/contracts/file-contracts"
	"time"

	"github.com/google/uuid"
)

type GetChatMessage struct {
	ID        uuid.UUID  `json:"id"`
	ChatID    uuid.UUID  `json:"chatID"`
	SenderID  *uuid.UUID `json:"senderID"`
	Content   string     `json:"content"`
	UpdatedAt *time.Time `json:"updatedAt"`
	CreatedAt time.Time  `json:"createdAt"`

	// Files Swagger override
	Files *[]*fc.File `json:"files" swaggertype:"array,object" swaggerignore:"true"`

	// Эта часть будет видна Swagger
	FilesSwagger *[]FileSwagger `json:"files"`
}
