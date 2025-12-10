package controllers

import (
	"apiService/internal/controllers"
	"apiService/internal/dto"
	"apiService/internal/services"
	"context"
	"errors"
	"testing"

	ac "common/contracts/api-chat"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Тесты для ChatRolePermissionController.GetAllRoles

func TestRolePermissionController_GetAllRoles_Success_FromCache(t *testing.T) {
	// Arrange
	mockRolePermissionClient := new(MockChatRolePermissionClient)
	redisClient := setupTestRedis(t)
	defer redisClient.Close()
	cacheService := services.NewCacheService(redisClient)

	controller := controllers.NewRolePermissionController(mockRolePermissionClient, cacheService)

	cachedRoles := []dto.ChatRoleResponseGateway{
		{ID: 1, Name: "admin"},
	}

	// Сохраняем роли в кеш
	cacheService.SetChatRolesCache(context.Background(), cachedRoles)

	// Act
	result, err := controller.GetAllRoles()

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 1)

	mockRolePermissionClient.AssertNotCalled(t, "GetAllRoles")
}

func TestRolePermissionController_GetAllRoles_Success_FromService(t *testing.T) {
	// Arrange
	mockRolePermissionClient := new(MockChatRolePermissionClient)
	redisClient := setupTestRedis(t)
	defer redisClient.Close()
	cacheService := services.NewCacheService(redisClient)

	controller := controllers.NewRolePermissionController(mockRolePermissionClient, cacheService)

	expectedRoles := []ac.RoleResponse{
		{
			ID:   1,
			Name: "admin",
			Permissions: []ac.PermissionResponse{
				{ID: 1, Name: "read"},
			},
		},
	}

	mockRolePermissionClient.On("GetAllRoles").Return(expectedRoles, nil)

	// Act
	result, err := controller.GetAllRoles()

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 1)
	assert.Equal(t, expectedRoles[0].Name, result[0].Name)

	mockRolePermissionClient.AssertExpectations(t)
}

func TestRolePermissionController_GetAllRoles_ServiceError(t *testing.T) {
	// Arrange
	mockRolePermissionClient := new(MockChatRolePermissionClient)
	redisClient := setupTestRedis(t)
	defer redisClient.Close()
	cacheService := services.NewCacheService(redisClient)

	controller := controllers.NewRolePermissionController(mockRolePermissionClient, cacheService)

	serviceError := errors.New("service error")

	mockRolePermissionClient.On("GetAllRoles").Return(nil, serviceError)

	// Act
	result, err := controller.GetAllRoles()

	// Assert
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, serviceError, err)

	mockRolePermissionClient.AssertExpectations(t)
}

// Тесты для ChatRolePermissionController.GetRoleByID

func TestRolePermissionController_GetRoleByID_Success(t *testing.T) {
	// Arrange
	mockRolePermissionClient := new(MockChatRolePermissionClient)
	redisClient := setupTestRedis(t)
	defer redisClient.Close()
	cacheService := services.NewCacheService(redisClient)

	controller := controllers.NewRolePermissionController(mockRolePermissionClient, cacheService)

	roleID := 1
	expectedRole := &ac.RoleResponse{
		ID:   roleID,
		Name: "admin",
		Permissions: []ac.PermissionResponse{
			{ID: 1, Name: "read"},
		},
	}

	mockRolePermissionClient.On("GetRoleByID", roleID).Return(expectedRole, nil)

	// Act
	result, err := controller.GetRoleByID(roleID)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedRole.Name, result.Name)

	mockRolePermissionClient.AssertExpectations(t)
}

func TestRolePermissionController_GetRoleByID_ServiceError(t *testing.T) {
	// Arrange
	mockRolePermissionClient := new(MockChatRolePermissionClient)
	redisClient := setupTestRedis(t)
	defer redisClient.Close()
	cacheService := services.NewCacheService(redisClient)

	controller := controllers.NewRolePermissionController(mockRolePermissionClient, cacheService)

	roleID := 1
	serviceError := errors.New("service error")

	mockRolePermissionClient.On("GetRoleByID", roleID).Return(nil, serviceError)

	// Act
	result, err := controller.GetRoleByID(roleID)

	// Assert
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, serviceError, err)

	mockRolePermissionClient.AssertExpectations(t)
}

// Тесты для ChatRolePermissionController.CreateRole

func TestRolePermissionController_CreateRole_Success(t *testing.T) {
	// Arrange
	mockRolePermissionClient := new(MockChatRolePermissionClient)
	redisClient := setupTestRedis(t)
	defer redisClient.Close()
	cacheService := services.NewCacheService(redisClient)

	controller := controllers.NewRolePermissionController(mockRolePermissionClient, cacheService)

	req := &dto.CreateChatRoleRequestGateway{
		Name:          "moderator",
		PermissionIDs: []int{1, 2},
	}

	expectedRole := &ac.RoleResponse{
		ID:   2,
		Name: "moderator",
		Permissions: []ac.PermissionResponse{
			{ID: 1, Name: "read"},
			{ID: 2, Name: "write"},
		},
	}

	mockRolePermissionClient.On("CreateRole", mock.Anything).Return(expectedRole, nil)

	// Act
	result, err := controller.CreateRole(req)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedRole.Name, result.Name)

	mockRolePermissionClient.AssertExpectations(t)
}

func TestRolePermissionController_CreateRole_ServiceError(t *testing.T) {
	// Arrange
	mockRolePermissionClient := new(MockChatRolePermissionClient)
	redisClient := setupTestRedis(t)
	defer redisClient.Close()
	cacheService := services.NewCacheService(redisClient)

	controller := controllers.NewRolePermissionController(mockRolePermissionClient, cacheService)

	req := &dto.CreateChatRoleRequestGateway{
		Name:          "moderator",
		PermissionIDs: []int{1, 2},
	}

	serviceError := errors.New("service error")

	mockRolePermissionClient.On("CreateRole", mock.Anything).Return(nil, serviceError)

	// Act
	result, err := controller.CreateRole(req)

	// Assert
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, serviceError, err)

	mockRolePermissionClient.AssertExpectations(t)
}

// Тесты для ChatRolePermissionController.DeleteRole

func TestRolePermissionController_DeleteRole_Success(t *testing.T) {
	// Arrange
	mockRolePermissionClient := new(MockChatRolePermissionClient)
	redisClient := setupTestRedis(t)
	defer redisClient.Close()
	cacheService := services.NewCacheService(redisClient)

	controller := controllers.NewRolePermissionController(mockRolePermissionClient, cacheService)

	roleID := 1

	mockRolePermissionClient.On("DeleteRole", roleID).Return(nil)

	// Act
	err := controller.DeleteRole(roleID)

	// Assert
	require.NoError(t, err)

	mockRolePermissionClient.AssertExpectations(t)
}

func TestRolePermissionController_DeleteRole_ServiceError(t *testing.T) {
	// Arrange
	mockRolePermissionClient := new(MockChatRolePermissionClient)
	redisClient := setupTestRedis(t)
	defer redisClient.Close()
	cacheService := services.NewCacheService(redisClient)

	controller := controllers.NewRolePermissionController(mockRolePermissionClient, cacheService)

	roleID := 1
	serviceError := errors.New("service error")

	mockRolePermissionClient.On("DeleteRole", roleID).Return(serviceError)

	// Act
	err := controller.DeleteRole(roleID)

	// Assert
	require.Error(t, err)
	assert.Equal(t, serviceError, err)

	mockRolePermissionClient.AssertExpectations(t)
}

// Тесты для ChatRolePermissionController.UpdateRolePermissions

func TestRolePermissionController_UpdateRolePermissions_Success(t *testing.T) {
	// Arrange
	mockRolePermissionClient := new(MockChatRolePermissionClient)
	redisClient := setupTestRedis(t)
	defer redisClient.Close()
	cacheService := services.NewCacheService(redisClient)

	controller := controllers.NewRolePermissionController(mockRolePermissionClient, cacheService)

	roleID := 1
	req := &dto.UpdateChatRolePermissionsRequestGateway{
		PermissionIDs: []int{1, 2, 3},
	}

	expectedRole := &ac.RoleResponse{
		ID:   roleID,
		Name: "admin",
		Permissions: []ac.PermissionResponse{
			{ID: 1, Name: "read"},
			{ID: 2, Name: "write"},
			{ID: 3, Name: "delete"},
		},
	}

	mockRolePermissionClient.On("UpdateRolePermissions", roleID, mock.Anything).Return(expectedRole, nil)

	// Act
	result, err := controller.UpdateRolePermissions(roleID, req)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedRole.Name, result.Name)

	mockRolePermissionClient.AssertExpectations(t)
}

// Тесты для ChatRolePermissionController.GetAllPermissions

func TestRolePermissionController_GetAllPermissions_Success_FromCache(t *testing.T) {
	// Arrange
	mockRolePermissionClient := new(MockChatRolePermissionClient)
	redisClient := setupTestRedis(t)
	defer redisClient.Close()
	cacheService := services.NewCacheService(redisClient)

	controller := controllers.NewRolePermissionController(mockRolePermissionClient, cacheService)

	cachedPermissions := []dto.ChatPermissionResponseGateway{
		{ID: 1, Name: "read"},
		{ID: 2, Name: "write"},
	}

	// Сохраняем permissions в кеш
	cacheService.SetChatPermissionsCache(context.Background(), cachedPermissions)

	// Act
	result, err := controller.GetAllPermissions()

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 2)

	mockRolePermissionClient.AssertNotCalled(t, "GetAllPermissions")
}

func TestRolePermissionController_GetAllPermissions_Success_FromService(t *testing.T) {
	// Arrange
	mockRolePermissionClient := new(MockChatRolePermissionClient)
	redisClient := setupTestRedis(t)
	defer redisClient.Close()
	cacheService := services.NewCacheService(redisClient)

	controller := controllers.NewRolePermissionController(mockRolePermissionClient, cacheService)

	expectedPermissions := []ac.PermissionResponse{
		{ID: 1, Name: "read"},
		{ID: 2, Name: "write"},
	}

	mockRolePermissionClient.On("GetAllPermissions").Return(expectedPermissions, nil)

	// Act
	result, err := controller.GetAllPermissions()

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 2)

	mockRolePermissionClient.AssertExpectations(t)
}

func TestRolePermissionController_GetAllPermissions_ServiceError(t *testing.T) {
	// Arrange
	mockRolePermissionClient := new(MockChatRolePermissionClient)
	redisClient := setupTestRedis(t)
	defer redisClient.Close()
	cacheService := services.NewCacheService(redisClient)

	controller := controllers.NewRolePermissionController(mockRolePermissionClient, cacheService)

	serviceError := errors.New("service error")

	mockRolePermissionClient.On("GetAllPermissions").Return(nil, serviceError)

	// Act
	result, err := controller.GetAllPermissions()

	// Assert
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, serviceError, err)

	mockRolePermissionClient.AssertExpectations(t)
}

// Тесты для ChatRolePermissionController.CreatePermission

func TestRolePermissionController_CreatePermission_Success(t *testing.T) {
	// Arrange
	mockRolePermissionClient := new(MockChatRolePermissionClient)
	redisClient := setupTestRedis(t)
	defer redisClient.Close()
	cacheService := services.NewCacheService(redisClient)

	controller := controllers.NewRolePermissionController(mockRolePermissionClient, cacheService)

	req := &dto.CreateChatPermissionRequestGateway{
		Name: "delete",
	}

	expectedPermission := &ac.PermissionResponse{
		ID:   3,
		Name: "delete",
	}

	mockRolePermissionClient.On("CreatePermission", mock.Anything).Return(expectedPermission, nil)

	// Act
	result, err := controller.CreatePermission(req)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedPermission.Name, result.Name)

	mockRolePermissionClient.AssertExpectations(t)
}

func TestRolePermissionController_CreatePermission_ServiceError(t *testing.T) {
	// Arrange
	mockRolePermissionClient := new(MockChatRolePermissionClient)
	redisClient := setupTestRedis(t)
	defer redisClient.Close()
	cacheService := services.NewCacheService(redisClient)

	controller := controllers.NewRolePermissionController(mockRolePermissionClient, cacheService)

	req := &dto.CreateChatPermissionRequestGateway{
		Name: "delete",
	}

	serviceError := errors.New("service error")

	mockRolePermissionClient.On("CreatePermission", mock.Anything).Return(nil, serviceError)

	// Act
	result, err := controller.CreatePermission(req)

	// Assert
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, serviceError, err)

	mockRolePermissionClient.AssertExpectations(t)
}

// Тесты для ChatRolePermissionController.DeletePermission

func TestRolePermissionController_DeletePermission_Success(t *testing.T) {
	// Arrange
	mockRolePermissionClient := new(MockChatRolePermissionClient)
	redisClient := setupTestRedis(t)
	defer redisClient.Close()
	cacheService := services.NewCacheService(redisClient)

	controller := controllers.NewRolePermissionController(mockRolePermissionClient, cacheService)

	permissionID := 1

	mockRolePermissionClient.On("DeletePermission", permissionID).Return(nil)

	// Act
	err := controller.DeletePermission(permissionID)

	// Assert
	require.NoError(t, err)

	mockRolePermissionClient.AssertExpectations(t)
}

func TestRolePermissionController_DeletePermission_ServiceError(t *testing.T) {
	// Arrange
	mockRolePermissionClient := new(MockChatRolePermissionClient)
	redisClient := setupTestRedis(t)
	defer redisClient.Close()
	cacheService := services.NewCacheService(redisClient)

	controller := controllers.NewRolePermissionController(mockRolePermissionClient, cacheService)

	permissionID := 1
	serviceError := errors.New("service error")

	mockRolePermissionClient.On("DeletePermission", permissionID).Return(serviceError)

	// Act
	err := controller.DeletePermission(permissionID)

	// Assert
	require.Error(t, err)
	assert.Equal(t, serviceError, err)

	mockRolePermissionClient.AssertExpectations(t)
}
