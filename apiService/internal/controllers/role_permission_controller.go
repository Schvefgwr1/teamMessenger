package controllers

import (
	"apiService/internal/dto"
	"apiService/internal/http_clients"
	"apiService/internal/services"
	ac "common/contracts/api-chat"
	"context"
	"log"
	"time"
)

type ChatRolePermissionController struct {
	rolePermissionClient http_clients.ChatRolePermissionClient
	cacheService         *services.CacheService
}

func NewRolePermissionController(rolePermissionClient http_clients.ChatRolePermissionClient, cacheService *services.CacheService) *ChatRolePermissionController {
	return &ChatRolePermissionController{
		rolePermissionClient: rolePermissionClient,
		cacheService:         cacheService,
	}
}

// ==================== Roles ====================

func (ctrl *ChatRolePermissionController) GetAllRoles() ([]dto.ChatRoleResponseGateway, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Пытаемся получить из кеша
	var cachedRoles []dto.ChatRoleResponseGateway
	err := ctrl.cacheService.GetChatRolesCache(ctx, &cachedRoles)
	if err == nil {
		log.Printf("Chat roles found in cache")
		return cachedRoles, nil
	}

	// Получаем из сервиса
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

	// Сохраняем в кеш
	if err := ctrl.cacheService.SetChatRolesCache(ctx, result); err != nil {
		log.Printf("Failed to cache chat roles: %v", err)
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

	// Инвалидация кеша ролей
	ctx := context.Background()
	_ = ctrl.cacheService.DeleteChatRolesCache(ctx)

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
	err := ctrl.rolePermissionClient.DeleteRole(roleID)
	if err != nil {
		return err
	}

	// Инвалидация кеша ролей
	ctx := context.Background()
	_ = ctrl.cacheService.DeleteChatRolesCache(ctx)

	return nil
}

func (ctrl *ChatRolePermissionController) UpdateRolePermissions(roleID int, req *dto.UpdateChatRolePermissionsRequestGateway) (*dto.ChatRoleResponseGateway, error) {
	updateReq := &ac.UpdateRolePermissionsRequest{
		PermissionIDs: req.PermissionIDs,
	}

	role, err := ctrl.rolePermissionClient.UpdateRolePermissions(roleID, updateReq)
	if err != nil {
		return nil, err
	}

	// Инвалидация кеша ролей
	ctx := context.Background()
	_ = ctrl.cacheService.DeleteChatRolesCache(ctx)

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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Пытаемся получить из кеша
	var cachedPermissions []dto.ChatPermissionResponseGateway
	err := ctrl.cacheService.GetChatPermissionsCache(ctx, &cachedPermissions)
	if err == nil {
		log.Printf("Chat permissions found in cache")
		return cachedPermissions, nil
	}

	// Получаем из сервиса
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

	// Сохраняем в кеш
	if err := ctrl.cacheService.SetChatPermissionsCache(ctx, result); err != nil {
		log.Printf("Failed to cache chat permissions: %v", err)
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

	// Инвалидация кеша permissions
	ctx := context.Background()
	_ = ctrl.cacheService.DeleteChatPermissionsCache(ctx)

	return &dto.ChatPermissionResponseGateway{
		ID:   permission.ID,
		Name: permission.Name,
	}, nil
}

func (ctrl *ChatRolePermissionController) DeletePermission(permissionID int) error {
	err := ctrl.rolePermissionClient.DeletePermission(permissionID)
	if err != nil {
		return err
	}

	// Инвалидация кеша permissions и ролей (т.к. роли содержат permissions)
	ctx := context.Background()
	_ = ctrl.cacheService.DeleteChatPermissionsCache(ctx)
	_ = ctrl.cacheService.DeleteChatRolesCache(ctx)

	return nil
}
