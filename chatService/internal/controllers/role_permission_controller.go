package controllers

import (
	"chatService/internal/handlers/dto"
	"chatService/internal/models"
	"chatService/internal/repositories"
)

type RolePermissionController struct {
	roleRepo       repositories.ChatRoleRepository
	permissionRepo repositories.ChatPermissionRepository
}

func NewRolePermissionController(
	roleRepo repositories.ChatRoleRepository,
	permissionRepo repositories.ChatPermissionRepository,
) *RolePermissionController {
	return &RolePermissionController{
		roleRepo:       roleRepo,
		permissionRepo: permissionRepo,
	}
}

// ==================== Roles ====================

func (c *RolePermissionController) GetAllRoles() ([]dto.RoleResponse, error) {
	roles, err := c.roleRepo.GetAllRoles()
	if err != nil {
		return nil, err
	}

	var result []dto.RoleResponse
	for _, role := range roles {
		permissions := make([]dto.PermissionResponse, 0, len(role.Permissions))
		for _, p := range role.Permissions {
			permissions = append(permissions, dto.PermissionResponse{
				ID:   p.ID,
				Name: p.Name,
			})
		}
		result = append(result, dto.RoleResponse{
			ID:          role.ID,
			Name:        role.Name,
			Permissions: permissions,
		})
	}

	return result, nil
}

func (c *RolePermissionController) GetRoleByID(roleID int) (*dto.RoleResponse, error) {
	role, err := c.roleRepo.GetRoleByID(roleID)
	if err != nil {
		return nil, err
	}

	permissions := make([]dto.PermissionResponse, 0, len(role.Permissions))
	for _, p := range role.Permissions {
		permissions = append(permissions, dto.PermissionResponse{
			ID:   p.ID,
			Name: p.Name,
		})
	}

	return &dto.RoleResponse{
		ID:          role.ID,
		Name:        role.Name,
		Permissions: permissions,
	}, nil
}

func (c *RolePermissionController) CreateRole(req *dto.CreateRoleRequest) (*dto.RoleResponse, error) {
	role := &models.ChatRole{
		Name: req.Name,
	}

	if err := c.roleRepo.CreateRole(role, req.PermissionIDs); err != nil {
		return nil, err
	}

	permissions := make([]dto.PermissionResponse, 0, len(role.Permissions))
	for _, p := range role.Permissions {
		permissions = append(permissions, dto.PermissionResponse{
			ID:   p.ID,
			Name: p.Name,
		})
	}

	return &dto.RoleResponse{
		ID:          role.ID,
		Name:        role.Name,
		Permissions: permissions,
	}, nil
}

func (c *RolePermissionController) DeleteRole(roleID int) error {
	return c.roleRepo.DeleteRole(roleID)
}

func (c *RolePermissionController) UpdateRolePermissions(roleID int, req *dto.UpdateRolePermissionsRequest) (*dto.RoleResponse, error) {
	if err := c.roleRepo.UpdateRolePermissions(roleID, req.PermissionIDs); err != nil {
		return nil, err
	}

	return c.GetRoleByID(roleID)
}

// ==================== Permissions ====================

func (c *RolePermissionController) GetAllPermissions() ([]dto.PermissionResponse, error) {
	permissions, err := c.permissionRepo.GetAllPermissions()
	if err != nil {
		return nil, err
	}

	var result []dto.PermissionResponse
	for _, p := range permissions {
		result = append(result, dto.PermissionResponse{
			ID:   p.ID,
			Name: p.Name,
		})
	}

	return result, nil
}

func (c *RolePermissionController) CreatePermission(req *dto.CreatePermissionRequest) (*dto.PermissionResponse, error) {
	permission := &models.ChatPermission{
		Name: req.Name,
	}

	if err := c.permissionRepo.CreatePermission(permission); err != nil {
		return nil, err
	}

	return &dto.PermissionResponse{
		ID:   permission.ID,
		Name: permission.Name,
	}, nil
}

func (c *RolePermissionController) DeletePermission(permissionID int) error {
	return c.permissionRepo.DeletePermission(permissionID)
}
