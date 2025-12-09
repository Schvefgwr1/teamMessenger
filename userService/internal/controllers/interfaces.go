package controllers

import (
	au "common/contracts/api-user"
	fc "common/contracts/file-contracts"
	"github.com/google/uuid"
	"userService/internal/handlers/dto"
	"userService/internal/models"
)

// UserControllerInterface - интерфейс для UserController для возможности мокирования
type UserControllerInterface interface {
	GetUserProfile(id uuid.UUID) (*models.User, *fc.File, error)
	UpdateUserProfile(req *au.UpdateUserRequest, userId *uuid.UUID) error
	GetUserBrief(userID uuid.UUID, chatID string, requesterID string) (*dto.UserBriefResponse, error)
	SearchUsers(query string, limit int) (*dto.UserSearchResponse, error)
	UpdateUserRole(userID uuid.UUID, roleID int) error
}

// AuthControllerInterface - интерфейс для AuthController для возможности мокирования
type AuthControllerInterface interface {
	Register(req *au.RegisterUserRequest) (*models.User, error)
	Login(req *au.Login, ipAddress, userAgent string) (string, uuid.UUID, error)
}

// PermissionControllerInterface - интерфейс для PermissionController для возможности мокирования
type PermissionControllerInterface interface {
	GetPermissions() ([]models.Permission, error)
}

// RoleControllerInterface - интерфейс для RoleController для возможности мокирования
type RoleControllerInterface interface {
	GetRoles() ([]models.Role, error)
	CreateRole(roleDTO *dto.CreateRole) error
	DeleteRole(roleID int) error
	UpdateRolePermissions(roleID int, permissionIDs []int) error
}
