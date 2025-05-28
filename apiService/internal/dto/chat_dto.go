package dto

import (
	"github.com/google/uuid"
	"mime/multipart"
)

type CreateChatRequestGateway struct {
	Name        string                `form:"name" binding:"required"`
	Description *string               `form:"description"`
	OwnerID     string                `form:"ownerID" binding:"required"`
	UserIDs     []string              `form:"userIDs" binding:"required"`
	Avatar      *multipart.FileHeader `form:"avatar"`
}

// ParseUUIDs преобразует строковые идентификаторы в UUID
func (r *CreateChatRequestGateway) ParseUUIDs() (ownerID uuid.UUID, userIDs []uuid.UUID, err error) {
	// Парсим ownerID
	ownerID, err = uuid.Parse(r.OwnerID)
	if err != nil {
		return uuid.Nil, nil, err
	}

	// Преобразуем каждый string в UUID
	userIDs = make([]uuid.UUID, len(r.UserIDs))
	for i, idStr := range r.UserIDs {
		userIDs[i], err = uuid.Parse(idStr)
		if err != nil {
			return uuid.Nil, nil, err
		}
	}

	return ownerID, userIDs, nil
}

type SendMessageRequestGateway struct {
	Content string                  `form:"content" binding:"required"`
	Files   []*multipart.FileHeader `form:"files"`
}

type CreateChatResponse struct {
	ID           uuid.UUID   `json:"id"`
	Name         string      `json:"name"`
	Description  *string     `json:"description"`
	OwnerID      uuid.UUID   `json:"ownerID"`
	UserIDs      []uuid.UUID `json:"userIDs"`
	AvatarFileID *int        `json:"avatarFileID,omitempty"`
}
