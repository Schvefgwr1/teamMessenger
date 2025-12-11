//go:build integration
// +build integration

package integration

import (
	"apiService/internal/controllers"
	"apiService/internal/dto"
	"apiService/internal/services"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestChatRolePermissionController_GetAllRoles_Integration тестирует получение всех ролей с кешированием
func TestChatRolePermissionController_GetAllRoles_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange
	redisClient := setupTestRedis(t)
	_, chatServer, _, _, _, _, _, _, rolePermissionClient := setupTestHTTPClients(t)
	defer chatServer.Close()

	cacheService := services.NewCacheService(redisClient)
	rolePermissionController := controllers.NewRolePermissionController(rolePermissionClient, cacheService)

	// Act - первый запрос
	roles1, err := rolePermissionController.GetAllRoles()

	// Assert
	require.NoError(t, err)
	assert.NotEmpty(t, roles1)

	// Act - второй запрос (должен быть из кеша)
	roles2, err := rolePermissionController.GetAllRoles()

	// Assert
	require.NoError(t, err)
	assert.Equal(t, len(roles1), len(roles2))

	// Проверяем, что данные реально в кеше
	ctx := context.Background()
	var cachedRoles []dto.ChatRoleResponseGateway
	err = cacheService.GetChatRolesCache(ctx, &cachedRoles)
	require.NoError(t, err)
}

// TestChatRolePermissionController_GetRoleByID_Integration тестирует получение роли по ID
func TestChatRolePermissionController_GetRoleByID_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange
	redisClient := setupTestRedis(t)
	_, chatServer, _, _, _, _, _, _, rolePermissionClient := setupTestHTTPClients(t)
	defer chatServer.Close()

	cacheService := services.NewCacheService(redisClient)
	rolePermissionController := controllers.NewRolePermissionController(rolePermissionClient, cacheService)

	roleID := 1

	// Act
	role, err := rolePermissionController.GetRoleByID(roleID)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, role)
	assert.Equal(t, roleID, role.ID)
}

// TestChatRolePermissionController_CreateRole_Integration тестирует создание роли с инвалидацией кеша
func TestChatRolePermissionController_CreateRole_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange
	redisClient := setupTestRedis(t)
	_, chatServer, _, _, _, _, _, _, rolePermissionClient := setupTestHTTPClients(t)
	defer chatServer.Close()

	cacheService := services.NewCacheService(redisClient)
	rolePermissionController := controllers.NewRolePermissionController(rolePermissionClient, cacheService)

	// Создаем кеш перед созданием роли
	ctx := context.Background()
	cacheService.SetChatRolesCache(ctx, []dto.ChatRoleResponseGateway{})

	createReq := &dto.CreateChatRoleRequestGateway{
		Name:          "test_role",
		PermissionIDs: []int{1, 2},
	}

	// Act
	role, err := rolePermissionController.CreateRole(createReq)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, role)

	// Проверяем, что кеш ролей инвалидирован
	exists, _ := cacheService.Exists(ctx, "chat_roles:all")
	assert.False(t, exists, "Roles cache should be invalidated after create")
}

// TestChatRolePermissionController_DeleteRole_Integration тестирует удаление роли с инвалидацией кеша
func TestChatRolePermissionController_DeleteRole_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange
	redisClient := setupTestRedis(t)
	_, chatServer, _, _, _, _, _, _, rolePermissionClient := setupTestHTTPClients(t)
	defer chatServer.Close()

	cacheService := services.NewCacheService(redisClient)
	rolePermissionController := controllers.NewRolePermissionController(rolePermissionClient, cacheService)

	roleID := 1

	// Создаем кеш перед удалением
	ctx := context.Background()
	cacheService.SetChatRolesCache(ctx, []dto.ChatRoleResponseGateway{})

	// Act
	err := rolePermissionController.DeleteRole(roleID)

	// Assert
	require.NoError(t, err)

	// Проверяем, что кеш ролей инвалидирован
	exists, _ := cacheService.Exists(ctx, "chat_roles:all")
	assert.False(t, exists, "Roles cache should be invalidated after delete")
}

// TestChatRolePermissionController_UpdateRolePermissions_Integration тестирует обновление разрешений роли
func TestChatRolePermissionController_UpdateRolePermissions_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange
	redisClient := setupTestRedis(t)
	_, chatServer, _, _, _, _, _, _, rolePermissionClient := setupTestHTTPClients(t)
	defer chatServer.Close()

	cacheService := services.NewCacheService(redisClient)
	rolePermissionController := controllers.NewRolePermissionController(rolePermissionClient, cacheService)

	roleID := 1
	updateReq := &dto.UpdateChatRolePermissionsRequestGateway{
		PermissionIDs: []int{1, 2, 3},
	}

	// Создаем кеш перед обновлением
	ctx := context.Background()
	cacheService.SetChatRolesCache(ctx, []dto.ChatRoleResponseGateway{})

	// Act
	role, err := rolePermissionController.UpdateRolePermissions(roleID, updateReq)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, role)

	// Проверяем, что кеш ролей инвалидирован
	exists, _ := cacheService.Exists(ctx, "chat_roles:all")
	assert.False(t, exists, "Roles cache should be invalidated after update")
}

// TestChatRolePermissionController_GetAllPermissions_Integration тестирует получение всех разрешений с кешированием
func TestChatRolePermissionController_GetAllPermissions_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange
	redisClient := setupTestRedis(t)
	_, chatServer, _, _, _, _, _, _, rolePermissionClient := setupTestHTTPClients(t)
	defer chatServer.Close()

	cacheService := services.NewCacheService(redisClient)
	rolePermissionController := controllers.NewRolePermissionController(rolePermissionClient, cacheService)

	// Act - первый запрос
	permissions1, err := rolePermissionController.GetAllPermissions()

	// Assert
	require.NoError(t, err)
	assert.NotEmpty(t, permissions1)

	// Act - второй запрос (должен быть из кеша)
	permissions2, err := rolePermissionController.GetAllPermissions()

	// Assert
	require.NoError(t, err)
	assert.Equal(t, len(permissions1), len(permissions2))

	// Проверяем, что данные реально в кеше
	ctx := context.Background()
	var cachedPermissions []dto.ChatPermissionResponseGateway
	err = cacheService.GetChatPermissionsCache(ctx, &cachedPermissions)
	require.NoError(t, err)
}

// TestChatRolePermissionController_CreatePermission_Integration тестирует создание разрешения с инвалидацией кеша
func TestChatRolePermissionController_CreatePermission_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange
	redisClient := setupTestRedis(t)
	_, chatServer, _, _, _, _, _, _, rolePermissionClient := setupTestHTTPClients(t)
	defer chatServer.Close()

	cacheService := services.NewCacheService(redisClient)
	rolePermissionController := controllers.NewRolePermissionController(rolePermissionClient, cacheService)

	// Создаем кеш перед созданием разрешения
	ctx := context.Background()
	cacheService.SetChatPermissionsCache(ctx, []dto.ChatPermissionResponseGateway{})

	createReq := &dto.CreateChatPermissionRequestGateway{
		Name: "test_permission",
	}

	// Act
	permission, err := rolePermissionController.CreatePermission(createReq)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, permission)

	// Проверяем, что кеш разрешений инвалидирован
	exists, _ := cacheService.Exists(ctx, "chat_permissions:all")
	assert.False(t, exists, "Permissions cache should be invalidated after create")
}

// TestChatRolePermissionController_DeletePermission_Integration тестирует удаление разрешения с инвалидацией кеша
func TestChatRolePermissionController_DeletePermission_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange
	redisClient := setupTestRedis(t)
	_, chatServer, _, _, _, _, _, _, rolePermissionClient := setupTestHTTPClients(t)
	defer chatServer.Close()

	cacheService := services.NewCacheService(redisClient)
	rolePermissionController := controllers.NewRolePermissionController(rolePermissionClient, cacheService)

	permissionID := 1

	// Создаем кеш перед удалением
	ctx := context.Background()
	cacheService.SetChatPermissionsCache(ctx, []dto.ChatPermissionResponseGateway{})
	cacheService.SetChatRolesCache(ctx, []dto.ChatRoleResponseGateway{})

	// Act
	err := rolePermissionController.DeletePermission(permissionID)

	// Assert
	require.NoError(t, err)

	// Проверяем, что кеш разрешений и ролей инвалидирован
	exists1, _ := cacheService.Exists(ctx, "chat_permissions:all")
	exists2, _ := cacheService.Exists(ctx, "chat_roles:all")
	assert.False(t, exists1, "Permissions cache should be invalidated after delete")
	assert.False(t, exists2, "Roles cache should be invalidated after delete (roles contain permissions)")
}
