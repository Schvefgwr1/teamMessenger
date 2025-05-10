package dto

import "github.com/google/uuid"

type CreateChatDTO struct {
	Name         string      `json:"name"`
	Description  *string     `json:"description"`
	AvatarFileID *int        `json:"avatarFileID"`
	OwnerID      uuid.UUID   `json:"ownerID"`
	UserIDs      []uuid.UUID `json:"userIDs"`
}
