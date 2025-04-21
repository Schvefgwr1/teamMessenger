package repositories

import (
	"gorm.io/gorm"
	"userService/internal/models"
)

type RoleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) *RoleRepository {
	return &RoleRepository{db: db}
}

func (r *RoleRepository) GetAllRoles() ([]models.Role, error) {
	var roles []models.Role
	err := r.db.Preload("Permissions").Find(&roles).Error
	return roles, err
}

func (r *RoleRepository) GetRoleByID(id int) (*models.Role, error) {
	var role models.Role
	return &role, r.db.Where("id = ?", id).First(&role).Error
}

func (r *RoleRepository) CreateRole(role *models.Role) error {
	tx := r.db.Begin()

	var permIds []int
	for _, rolePerm := range role.Permissions {
		permIds = append(permIds, rolePerm.ID)
	}

	// 1. Сохраняем саму роль
	if err := tx.Omit("Permissions").Create(role).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 2. Вставляем связи напрямую через SQL
	for _, pid := range permIds {
		if err := tx.Exec(
			`INSERT INTO user_service.role_permissions (role_id, permission_id) VALUES (?, ?)`,
			role.ID, pid,
		).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

func (r *RoleRepository) DeleteRole(id int) error {
	return r.db.Delete(&models.Role{}, id).Error
}
