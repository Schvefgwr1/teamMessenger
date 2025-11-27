package repositories

import (
	"chatService/internal/models"
	"fmt"

	"gorm.io/gorm"
)

type ChatRoleRepository interface {
	GetRoleByID(roleID int) (*models.ChatRole, error)
	GetRoleByName(roleName string) (*models.ChatRole, error)
	GetAllRoles() ([]models.ChatRole, error)
	CreateRole(role *models.ChatRole, permissionIDs []int) error
	DeleteRole(roleID int) error
	UpdateRolePermissions(roleID int, permissionIDs []int) error
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

func (r *chatRoleRepository) GetAllRoles() ([]models.ChatRole, error) {
	var roles []models.ChatRole
	err := r.db.Preload("Permissions").Find(&roles).Error
	return roles, err
}

func (r *chatRoleRepository) CreateRole(role *models.ChatRole, permissionIDs []int) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Создаём роль
		if err := tx.Create(role).Error; err != nil {
			return err
		}

		// Привязываем permissions если есть
		if len(permissionIDs) > 0 {
			for _, permID := range permissionIDs {
				if err := tx.Exec(
					"INSERT INTO chat_service.chat_role_permissions (chat_role_id, chat_permission_id) VALUES (?, ?)",
					role.ID, permID,
				).Error; err != nil {
					return err
				}
			}
		}

		// Загружаем permissions для ответа
		return tx.Preload("Permissions").First(role, role.ID).Error
	})
}

func (r *chatRoleRepository) DeleteRole(roleID int) error {
	return r.db.Delete(&models.ChatRole{}, roleID).Error
}

func (r *chatRoleRepository) UpdateRolePermissions(roleID int, permissionIDs []int) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Проверяем существование роли
		var role models.ChatRole
		if err := tx.First(&role, roleID).Error; err != nil {
			return err
		}

		// Удаляем все текущие permissions
		if err := tx.Exec(
			"DELETE FROM chat_service.chat_role_permissions WHERE chat_role_id = ?",
			roleID,
		).Error; err != nil {
			return err
		}

		// Добавляем новые permissions
		for _, permID := range permissionIDs {
			var permission models.ChatPermission
			if err := tx.First(&permission, permID).Error; err != nil {
				return fmt.Errorf("permission with id %d not found", permID)
			}
			if err := tx.Exec(
				"INSERT INTO chat_service.chat_role_permissions (chat_role_id, chat_permission_id) VALUES (?, ?)",
				roleID, permID,
			).Error; err != nil {
				return err
			}
		}

		return nil
	})
}
