package repositories

import (
	"chatService/internal/models"
	"gorm.io/gorm"
)

type ChatPermissionRepository interface {
	GetAllPermissions() ([]models.ChatPermission, error)
	GetPermissionByID(id int) (*models.ChatPermission, error)
	GetPermissionByName(name string) (*models.ChatPermission, error)
	CreatePermission(permission *models.ChatPermission) error
	DeletePermission(id int) error
}

type chatPermissionRepository struct {
	db *gorm.DB
}

func NewChatPermissionRepository(db *gorm.DB) ChatPermissionRepository {
	return &chatPermissionRepository{db}
}

func (r *chatPermissionRepository) GetAllPermissions() ([]models.ChatPermission, error) {
	var permissions []models.ChatPermission
	err := r.db.Find(&permissions).Error
	return permissions, err
}

func (r *chatPermissionRepository) GetPermissionByID(id int) (*models.ChatPermission, error) {
	var permission models.ChatPermission
	err := r.db.First(&permission, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

func (r *chatPermissionRepository) GetPermissionByName(name string) (*models.ChatPermission, error) {
	var permission models.ChatPermission
	err := r.db.First(&permission, "name = ?", name).Error
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

func (r *chatPermissionRepository) CreatePermission(permission *models.ChatPermission) error {
	return r.db.Create(permission).Error
}

func (r *chatPermissionRepository) DeletePermission(id int) error {
	return r.db.Delete(&models.ChatPermission{}, id).Error
}
