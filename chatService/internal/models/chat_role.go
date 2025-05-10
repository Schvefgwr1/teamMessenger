package models

type ChatRole struct {
	ID   int    `gorm:"primaryKey;autoIncrement"`
	Name string `gorm:"size:100;not null"`

	Permissions []ChatPermission `gorm:"many2many:chat_service.chat_role_permissions"`
}

func (ChatRole) TableName() string {
	return "chat_service.chat_roles"
}
