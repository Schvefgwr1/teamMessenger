package api_user

import (
	"github.com/google/uuid"
	"time"
)

type RegisterUserResponse struct {
	ID           *uuid.UUID `json:"ID"`
	Username     *string    `json:"Username"`
	Email        *string    `json:"Email"`
	PasswordHash *string    `json:"PasswordHash"`
	Description  *string    `json:"Description"`
	Gender       *string    `json:"Gender"`
	Age          *int       `json:"Age"`
	AvatarFileID *int       `json:"AvatarFileID"`
	CreatedAt    *time.Time `json:"CreatedAt"`
	UpdatedAt    *time.Time `json:"UpdatedAt"`
	Error        *string    `json:"error"`
}
