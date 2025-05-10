package dto

import (
	"github.com/google/uuid"
	"time"
)

type ChatResponse struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	IsGroup      bool      `json:"isGroup"`
	Description  *string   `json:"description"`
	AvatarFileID *int      `json:"avatarFileID"`
	CreatedAt    time.Time `json:"createdAt"`
}
