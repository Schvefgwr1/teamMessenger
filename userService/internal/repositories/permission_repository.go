package repositories

import (
	"gorm.io/gorm"
	"userService/internal/models"
)

type PermissionRepository struct {
	db *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) *PermissionRepository {
	return &PermissionRepository{db: db}
}

func (r *PermissionRepository) GetAllPermissions() ([]models.Permission, error) {
	var permissions []models.Permission
	err := r.db.Find(&permissions).Error
	return permissions, err
}

func (r *PermissionRepository) GetPermissionById(id int) (*models.Permission, error) {
	var permission models.Permission
	return &permission, r.db.Where("id = ?", id).First(&permission).Error
}
