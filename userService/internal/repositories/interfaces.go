package repositories

import (
	"github.com/google/uuid"
	"userService/internal/models"
)

// UserRepositoryInterface - интерфейс для UserRepository для возможности мокирования
type UserRepositoryInterface interface {
	CreateUser(user *models.User) error
	GetUserByEmail(email string) (*models.User, error)
	GetUserByUsername(username string) (*models.User, error)
	GetUserByID(id uuid.UUID) (*models.User, error)
	UpdateUser(user *models.User) error
	SearchUsers(query string, limit int) ([]*models.User, error)
}

// RoleRepositoryInterface - интерфейс для RoleRepository для возможности мокирования
type RoleRepositoryInterface interface {
	GetAllRoles() ([]models.Role, error)
	GetRoleByID(id int) (*models.Role, error)
	CreateRole(role *models.Role) error
	DeleteRole(id int) error
	UpdateRolePermissions(roleID int, permissionIDs []int) error
}

// PermissionRepositoryInterface - интерфейс для PermissionRepository для возможности мокирования
type PermissionRepositoryInterface interface {
	GetAllPermissions() ([]models.Permission, error)
	GetPermissionById(id int) (*models.Permission, error)
}
