package dto

import (
	fc "common/contracts/file-contracts"
	"encoding/json"
	"github.com/google/uuid"
	"time"
)

type ChatResponse struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	IsGroup      bool      `json:"isGroup"`
	Description  *string   `json:"description"`
	AvatarFileID *int      `json:"avatarFileID,omitempty"`

	// AvatarFile - реальные данные (скрыто от Swagger)
	AvatarFile *fc.File `json:"-" swaggerignore:"true"`

	// AvatarFileSwagger - только для Swagger документации
	AvatarFileSwagger *FileSwagger `json:"avatarFile,omitempty" swaggertype:"object"`

	CreatedAt time.Time `json:"createdAt"`
}

// MarshalJSON кастомная сериализация для правильной работы с AvatarFile
func (c ChatResponse) MarshalJSON() ([]byte, error) {
	// Создаем временную структуру для сериализации, используя только AvatarFile
	aux := struct {
		ID           uuid.UUID `json:"id"`
		Name         string    `json:"name"`
		IsGroup      bool      `json:"isGroup"`
		Description  *string   `json:"description"`
		AvatarFileID *int      `json:"avatarFileID,omitempty"`
		AvatarFile   *fc.File  `json:"avatarFile,omitempty"`
		CreatedAt    time.Time `json:"createdAt"`
	}{
		ID:           c.ID,
		Name:         c.Name,
		IsGroup:      c.IsGroup,
		Description:  c.Description,
		AvatarFileID: c.AvatarFileID,
		AvatarFile:   c.AvatarFile,
		CreatedAt:    c.CreatedAt,
	}
	return json.Marshal(aux)
}
