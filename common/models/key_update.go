package models

import (
	"github.com/google/uuid"
	"time"
)

// PublicKeyUpdate сообщение об обновлении публичного ключа
type PublicKeyUpdate struct {
	ID           uuid.UUID `json:"id"`
	PublicKeyPEM string    `json:"public_key_pem"`
	UpdatedAt    time.Time `json:"updated_at"`
	ServiceName  string    `json:"service_name"`
	KeyVersion   int       `json:"key_version"`
}
