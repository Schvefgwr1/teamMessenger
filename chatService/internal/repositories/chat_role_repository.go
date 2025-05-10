package repositories

import (
	"chatService/internal/models"
	"gorm.io/gorm"
)

type ChatRoleRepository interface {
	GetRoleByID(roleID int) (*models.ChatRole, error)
	GetRoleByName(roleName string) (*models.ChatRole, error)
}

type chatRoleRepository struct {
	db *gorm.DB
}

func NewChatRoleRepository(db *gorm.DB) ChatRoleRepository {
	return &chatRoleRepository{db}
}

func (r *chatRoleRepository) GetRoleByID(roleID int) (*models.ChatRole, error) {
	var chatRole models.ChatRole
	err := r.db.Preload("Permissions").First(&chatRole, "id = ?", roleID).Error
	if err != nil {
		return nil, err
	}
	return &chatRole, nil
}

func (r *chatRoleRepository) GetRoleByName(roleName string) (*models.ChatRole, error) {
	var chatRole models.ChatRole
	err := r.db.Preload("Permissions").First(&chatRole, "name = ?", roleName).Error
	if err != nil {
		return nil, err
	}
	return &chatRole, nil
}
