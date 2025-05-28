package api_chat

import (
	fc "common/contracts/file-contracts"
	"github.com/google/uuid"
	"time"
)

type ChatResponse struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Description  *string   `json:"description"`
	AvatarFileID *int      `json:"avatarFileID"`
	IsGroup      bool      `json:"isGroup"`
	CreatedAt    time.Time `json:"createdAt"`
}

type CreateChatRequest struct {
	Name         string      `json:"name"`
	Description  *string     `json:"description"`
	AvatarFileID *int        `json:"avatarFileID"`
	OwnerID      uuid.UUID   `json:"ownerID"`
	UserIDs      []uuid.UUID `json:"userIDs"`
}

type MessageResponse struct {
	ID        uuid.UUID      `json:"id"`
	ChatID    uuid.UUID      `json:"chatID"`
	SenderID  uuid.UUID      `json:"senderID"`
	Content   string         `json:"content"`
	Files     []*MessageFile `json:"files,omitempty"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt *time.Time     `json:"updatedAt,omitempty"`
}

type GetChatMessage struct {
	ID        uuid.UUID   `json:"id"`
	ChatID    uuid.UUID   `json:"chatID"`
	SenderID  *uuid.UUID  `json:"senderID"`
	Content   string      `json:"content"`
	UpdatedAt *time.Time  `json:"updatedAt"`
	CreatedAt time.Time   `json:"createdAt"`
	Files     *[]*fc.File `json:"files,omitempty"`
}

type GetSearchResponse struct {
	Messages *[]GetChatMessage `json:"messages"`
	Total    *int64            `json:"total"`
}

type MessageFile struct {
	MessageID uuid.UUID `json:"messageId"`
	FileID    int       `json:"fileId"`
}

type CreateMessageRequest struct {
	Content string `json:"content"`
	FileIDs []int  `json:"fileIDs,omitempty"`
}

type CreateChatServiceResponse struct {
	ChatID uuid.UUID `json:"chat_id"`
}
