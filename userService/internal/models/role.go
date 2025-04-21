package models

type Role struct {
	ID          *int   `gorm:"primaryKey"`
	Name        string `gorm:"unique;not null"`
	Description string
	Permissions []Permission `gorm:"many2many:user_service.role_permissions"`
}

func (Role) TableName() string {
	return "user_service.roles"
}
