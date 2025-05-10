package models

type ChatPermission struct {
	ID   int    `gorm:"primaryKey;autoIncrement"`
	Name string `gorm:"size:100;not null"`
}

func (ChatPermission) TableName() string {
	return "chat_service.chat_permissions"
}
