package user_contracts

import (
	"github.com/google/uuid"
)

// Permission - право доступа
type Permission struct {
	ID          int    `json:"ID"`
	Name        string `json:"Name"`
	Description string `json:"Description"`
}

// Role - роль пользователя с правами
type Role struct {
	ID          int          `json:"ID"`
	Name        string       `json:"Name"`
	Description string       `json:"Description"`
	Permissions []Permission `json:"Permissions"`
}

type User struct {
	ID          uuid.UUID
	Username    string
	Email       string
	Description *string
	Gender      *string
	Age         *int
	Role        Role
}
