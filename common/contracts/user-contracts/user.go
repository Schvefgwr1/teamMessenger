package user_contracts

import (
	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID
	Username     string
	Email        string
	Description  *string
	Gender       *string
	Age          *int
	AvatarFileID *int `json:"avatarFileID"`
}
