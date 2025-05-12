package dto

import (
	fc "common/contracts/file-contracts"
	"github.com/google/uuid"
	"time"
)

type GetChatMessage struct {
	ID        uuid.UUID   `json:"id"`
	ChatID    uuid.UUID   `json:"chatID"`
	SenderID  *uuid.UUID  `json:"senderID"`
	Content   string      `json:"content"`
	UpdatedAt *time.Time  `json:"updatedAt"`
	CreatedAt time.Time   `json:"createdAt"`
	Files     *[]*fc.File `json:"files"`
}
