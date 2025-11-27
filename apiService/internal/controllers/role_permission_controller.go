package controllers

import (
	"apiService/internal/dto"
	"apiService/internal/http_clients"
	ac "common/contracts/api-chat"
)

type ChatRolePermissionController struct {
	rolePermissionClient http_clients.ChatRolePermissionClient
}

func NewRolePermissionController(rolePermissionClient http_clients.ChatRolePermissionClient) *ChatRolePermissionController {
	return &ChatRolePermissionController{
		rolePermissionClient: rolePermissionClient,
	}
}

// ==================== Roles ====================

func (ctrl *ChatRolePermissionController) GetAllRoles() ([]dto.ChatRoleResponseGateway, error) {
	roles, err := ctrl.rolePermissionClient.GetAllRoles()
	if err != nil {
		return nil, err
	}

	result := make([]dto.ChatRoleResponseGateway, 0, len(roles))
	for _, role := range roles {
		permissions := make([]dto.ChatPermissionResponseGateway, 0, len(role.Permissions))
		for _, p := range role.Permissions {
			permissions = append(permissions, dto.ChatPermissionResponseGateway{
				ID:   p.ID,
				Name: p.Name,
			})
		}
		result = append(result, dto.ChatRoleResponseGateway{
			ID:          role.ID,
			Name:        role.Name,
			Permissions: permissions,
		})
	}

	return result, nil
}

func (ctrl *ChatRolePermissionController) GetRoleByID(roleID int) (*dto.ChatRoleResponseGateway, error) {
	role, err := ctrl.rolePermissionClient.GetRoleByID(roleID)
	if err != nil {
		return nil, err
	}

	permissions := make([]dto.ChatPermissionResponseGateway, 0, len(role.Permissions))
	for _, p := range role.Permissions {
		permissions = append(permissions, dto.ChatPermissionResponseGateway{
			ID:   p.ID,
			Name: p.Name,
		})
	}

	return &dto.ChatRoleResponseGateway{
		ID:          role.ID,
		Name:        role.Name,
		Permissions: permissions,
	}, nil
}

func (ctrl *ChatRolePermissionController) CreateRole(req *dto.CreateRoleRequestGateway) (*dto.ChatRoleResponseGateway, error) {
	createReq := &ac.CreateRoleRequest{
		Name:          req.Name,
		PermissionIDs: req.PermissionIds,
	}

	role, err := ctrl.rolePermissionClient.CreateRole(createReq)
	if err != nil {
		return nil, err
	}

	permissions := make([]dto.ChatPermissionResponseGateway, 0, len(role.Permissions))
	for _, p := range role.Permissions {
		permissions = append(permissions, dto.ChatPermissionResponseGateway{
			ID:   p.ID,
			Name: p.Name,
		})
	}

	return &dto.ChatRoleResponseGateway{
		ID:          role.ID,
		Name:        role.Name,
		Permissions: permissions,
	}, nil
}

func (ctrl *ChatRolePermissionController) DeleteRole(roleID int) error {
	return ctrl.rolePermissionClient.DeleteRole(roleID)
}

func (ctrl *ChatRolePermissionController) UpdateRolePermissions(roleID int, req *dto.UpdateChatRolePermissionsRequestGateway) (*dto.ChatRoleResponseGateway, error) {
	updateReq := &ac.UpdateRolePermissionsRequest{
		PermissionIDs: req.PermissionIDs,
	}

	role, err := ctrl.rolePermissionClient.UpdateRolePermissions(roleID, updateReq)
	if err != nil {
		return nil, err
	}

	permissions := make([]dto.ChatPermissionResponseGateway, 0, len(role.Permissions))
	for _, p := range role.Permissions {
		permissions = append(permissions, dto.ChatPermissionResponseGateway{
			ID:   p.ID,
			Name: p.Name,
		})
	}

	return &dto.ChatRoleResponseGateway{
		ID:          role.ID,
		Name:        role.Name,
		Permissions: permissions,
	}, nil
}

// ==================== Permissions ====================

func (ctrl *ChatRolePermissionController) GetAllPermissions() ([]dto.ChatPermissionResponseGateway, error) {
	permissions, err := ctrl.rolePermissionClient.GetAllPermissions()
	if err != nil {
		return nil, err
	}

	result := make([]dto.ChatPermissionResponseGateway, 0, len(permissions))
	for _, p := range permissions {
		result = append(result, dto.ChatPermissionResponseGateway{
			ID:   p.ID,
			Name: p.Name,
		})
	}

	return result, nil
}

func (ctrl *ChatRolePermissionController) CreatePermission(req *dto.CreateChatPermissionRequestGateway) (*dto.ChatPermissionResponseGateway, error) {
	createReq := &ac.CreatePermissionRequest{
		Name: req.Name,
	}

	permission, err := ctrl.rolePermissionClient.CreatePermission(createReq)
	if err != nil {
		return nil, err
	}

	return &dto.ChatPermissionResponseGateway{
		ID:   permission.ID,
		Name: permission.Name,
	}, nil
}

func (ctrl *ChatRolePermissionController) DeletePermission(permissionID int) error {
	return ctrl.rolePermissionClient.DeletePermission(permissionID)
}
