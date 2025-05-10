package models

import "github.com/google/uuid"

type ChatUser struct {
	ChatID uuid.UUID `gorm:"primaryKey;type:uuid"`
	UserID uuid.UUID `gorm:"primaryKey;type:uuid"`
	RoleID int

	Chat Chat     `gorm:"foreignKey:ChatID"`
	Role ChatRole `gorm:"foreignKey:RoleID"`
}

func (ChatUser) TableName() string {
	return "chat_service.chat_user"
}
