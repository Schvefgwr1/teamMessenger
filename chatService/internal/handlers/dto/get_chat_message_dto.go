package dto

import (
	fc "common/contracts/file-contracts"
	"encoding/json"
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

	// Files - реальные данные для сериализации (скрыто от Swagger)
	Files *[]*fc.File `json:"-" swaggerignore:"true"`

	// FilesSwagger - только для Swagger документации (не участвует в JSON сериализации)
	FilesSwagger *[]FileSwagger `json:"files,omitempty" swaggertype:"array,object"`
}

// MarshalJSON кастомная сериализация для правильной работы с Files
// Использует Files для реальных данных, игнорируя FilesSwagger (который нужен только для Swagger)
func (m GetChatMessage) MarshalJSON() ([]byte, error) {
	// Создаем временную структуру для сериализации, используя только Files
	aux := struct {
		ID        uuid.UUID   `json:"id"`
		ChatID    uuid.UUID   `json:"chatID"`
		SenderID  *uuid.UUID  `json:"senderID"`
		Content   string      `json:"content"`
		UpdatedAt *time.Time  `json:"updatedAt"`
		CreatedAt time.Time   `json:"createdAt"`
		Files     *[]*fc.File `json:"files,omitempty"`
	}{
		ID:        m.ID,
		ChatID:    m.ChatID,
		SenderID:  m.SenderID,
		Content:   m.Content,
		UpdatedAt: m.UpdatedAt,
		CreatedAt: m.CreatedAt,
		Files:     m.Files,
	}
	return json.Marshal(aux)
}
