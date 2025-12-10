package controllers

import (
	"errors"
	"testing"

	"chatService/internal/controllers"
	"chatService/internal/handlers/dto"
	"chatService/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// ==================== Тесты для Roles ====================

func TestRolePermissionController_GetAllRoles_Success(t *testing.T) {
	// Arrange
	mockRoleRepo := new(MockChatRoleRepository)
	mockPermissionRepo := new(MockChatPermissionRepository)

	roles := []models.ChatRole{
		*createTestChatRoleWithPermissions(1, "owner", []models.ChatPermission{
			*createTestChatPermissionWithID(1, "send_message"),
			*createTestChatPermissionWithID(2, "delete_message"),
		}),
		*createTestChatRoleWithPermissions(2, "main", []models.ChatPermission{
			*createTestChatPermissionWithID(1, "send_message"),
		}),
	}

	mockRoleRepo.On("GetAllRoles").Return(roles, nil)

	controller := controllers.NewRolePermissionController(mockRoleRepo, mockPermissionRepo)

	// Act
	result, err := controller.GetAllRoles()

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	assert.Equal(t, "owner", result[0].Name)
	assert.Len(t, result[0].Permissions, 2)
	assert.Equal(t, "main", result[1].Name)
	assert.Len(t, result[1].Permissions, 1)

	mockRoleRepo.AssertExpectations(t)
}

func TestRolePermissionController_GetAllRoles_RepositoryError(t *testing.T) {
	// Arrange
	mockRoleRepo := new(MockChatRoleRepository)
	mockPermissionRepo := new(MockChatPermissionRepository)

	repoError := errors.New("database error")
	mockRoleRepo.On("GetAllRoles").Return(nil, repoError)

	controller := controllers.NewRolePermissionController(mockRoleRepo, mockPermissionRepo)

	// Act
	result, err := controller.GetAllRoles()

	// Assert
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, repoError, err)

	mockRoleRepo.AssertExpectations(t)
}

func TestRolePermissionController_GetRoleByID_Success(t *testing.T) {
	// Arrange
	mockRoleRepo := new(MockChatRoleRepository)
	mockPermissionRepo := new(MockChatPermissionRepository)

	roleID := 1
	role := createTestChatRoleWithPermissions(roleID, "owner", []models.ChatPermission{
		*createTestChatPermissionWithID(1, "send_message"),
		*createTestChatPermissionWithID(2, "delete_message"),
	})

	mockRoleRepo.On("GetRoleByID", roleID).Return(role, nil)

	controller := controllers.NewRolePermissionController(mockRoleRepo, mockPermissionRepo)

	// Act
	result, err := controller.GetRoleByID(roleID)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, roleID, result.ID)
	assert.Equal(t, "owner", result.Name)
	assert.Len(t, result.Permissions, 2)
	assert.Equal(t, 1, result.Permissions[0].ID)
	assert.Equal(t, "send_message", result.Permissions[0].Name)

	mockRoleRepo.AssertExpectations(t)
}

func TestRolePermissionController_GetRoleByID_NotFound(t *testing.T) {
	// Arrange
	mockRoleRepo := new(MockChatRoleRepository)
	mockPermissionRepo := new(MockChatPermissionRepository)

	roleID := 999
	mockRoleRepo.On("GetRoleByID", roleID).Return(nil, errors.New("role not found"))

	controller := controllers.NewRolePermissionController(mockRoleRepo, mockPermissionRepo)

	// Act
	result, err := controller.GetRoleByID(roleID)

	// Assert
	require.Error(t, err)
	assert.Nil(t, result)

	mockRoleRepo.AssertExpectations(t)
}

func TestRolePermissionController_CreateRole_Success(t *testing.T) {
	// Arrange
	mockRoleRepo := new(MockChatRoleRepository)
	mockPermissionRepo := new(MockChatPermissionRepository)

	createReq := &dto.CreateRoleRequest{
		Name:          "admin",
		PermissionIDs: []int{1, 2, 3},
	}

	createdRole := createTestChatRoleWithPermissions(1, "admin", []models.ChatPermission{
		*createTestChatPermissionWithID(1, "send_message"),
		*createTestChatPermissionWithID(2, "delete_message"),
		*createTestChatPermissionWithID(3, "manage_users"),
	})

	mockRoleRepo.On("CreateRole", mock.MatchedBy(func(role *models.ChatRole) bool {
		return role.Name == "admin"
	}), createReq.PermissionIDs).Return(nil).Run(func(args mock.Arguments) {
		role := args.Get(0).(*models.ChatRole)
		role.ID = 1
		role.Permissions = createdRole.Permissions
	})

	controller := controllers.NewRolePermissionController(mockRoleRepo, mockPermissionRepo)

	// Act
	result, err := controller.CreateRole(createReq)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "admin", result.Name)
	assert.Len(t, result.Permissions, 3)

	mockRoleRepo.AssertExpectations(t)
}

func TestRolePermissionController_CreateRole_RepositoryError(t *testing.T) {
	// Arrange
	mockRoleRepo := new(MockChatRoleRepository)
	mockPermissionRepo := new(MockChatPermissionRepository)

	createReq := &dto.CreateRoleRequest{
		Name:          "admin",
		PermissionIDs: []int{1, 2},
	}

	repoError := errors.New("database error")
	mockRoleRepo.On("CreateRole", mock.Anything, createReq.PermissionIDs).Return(repoError)

	controller := controllers.NewRolePermissionController(mockRoleRepo, mockPermissionRepo)

	// Act
	result, err := controller.CreateRole(createReq)

	// Assert
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, repoError, err)

	mockRoleRepo.AssertExpectations(t)
}

func TestRolePermissionController_DeleteRole_Success(t *testing.T) {
	// Arrange
	mockRoleRepo := new(MockChatRoleRepository)
	mockPermissionRepo := new(MockChatPermissionRepository)

	roleID := 1
	mockRoleRepo.On("DeleteRole", roleID).Return(nil)

	controller := controllers.NewRolePermissionController(mockRoleRepo, mockPermissionRepo)

	// Act
	err := controller.DeleteRole(roleID)

	// Assert
	require.NoError(t, err)

	mockRoleRepo.AssertExpectations(t)
}

func TestRolePermissionController_DeleteRole_RepositoryError(t *testing.T) {
	// Arrange
	mockRoleRepo := new(MockChatRoleRepository)
	mockPermissionRepo := new(MockChatPermissionRepository)

	roleID := 999
	repoError := errors.New("role not found")
	mockRoleRepo.On("DeleteRole", roleID).Return(repoError)

	controller := controllers.NewRolePermissionController(mockRoleRepo, mockPermissionRepo)

	// Act
	err := controller.DeleteRole(roleID)

	// Assert
	require.Error(t, err)
	assert.Equal(t, repoError, err)

	mockRoleRepo.AssertExpectations(t)
}

func TestRolePermissionController_UpdateRolePermissions_Success(t *testing.T) {
	// Arrange
	mockRoleRepo := new(MockChatRoleRepository)
	mockPermissionRepo := new(MockChatPermissionRepository)

	roleID := 1
	updateReq := &dto.UpdateRolePermissionsRequest{
		PermissionIDs: []int{1, 3},
	}

	updatedRole := createTestChatRoleWithPermissions(roleID, "admin", []models.ChatPermission{
		*createTestChatPermissionWithID(1, "send_message"),
		*createTestChatPermissionWithID(3, "manage_users"),
	})

	mockRoleRepo.On("UpdateRolePermissions", roleID, updateReq.PermissionIDs).Return(nil)
	mockRoleRepo.On("GetRoleByID", roleID).Return(updatedRole, nil)

	controller := controllers.NewRolePermissionController(mockRoleRepo, mockPermissionRepo)

	// Act
	result, err := controller.UpdateRolePermissions(roleID, updateReq)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, roleID, result.ID)
	assert.Len(t, result.Permissions, 2)

	mockRoleRepo.AssertExpectations(t)
}

func TestRolePermissionController_UpdateRolePermissions_UpdateError(t *testing.T) {
	// Arrange
	mockRoleRepo := new(MockChatRoleRepository)
	mockPermissionRepo := new(MockChatPermissionRepository)

	roleID := 999
	updateReq := &dto.UpdateRolePermissionsRequest{
		PermissionIDs: []int{1, 2},
	}

	repoError := errors.New("role not found")
	mockRoleRepo.On("UpdateRolePermissions", roleID, updateReq.PermissionIDs).Return(repoError)

	controller := controllers.NewRolePermissionController(mockRoleRepo, mockPermissionRepo)

	// Act
	result, err := controller.UpdateRolePermissions(roleID, updateReq)

	// Assert
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, repoError, err)

	mockRoleRepo.AssertExpectations(t)
}

func TestRolePermissionController_UpdateRolePermissions_GetRoleError(t *testing.T) {
	// Arrange
	mockRoleRepo := new(MockChatRoleRepository)
	mockPermissionRepo := new(MockChatPermissionRepository)

	roleID := 1
	updateReq := &dto.UpdateRolePermissionsRequest{
		PermissionIDs: []int{1, 2},
	}

	mockRoleRepo.On("UpdateRolePermissions", roleID, updateReq.PermissionIDs).Return(nil)
	mockRoleRepo.On("GetRoleByID", roleID).Return(nil, errors.New("role not found"))

	controller := controllers.NewRolePermissionController(mockRoleRepo, mockPermissionRepo)

	// Act
	result, err := controller.UpdateRolePermissions(roleID, updateReq)

	// Assert
	require.Error(t, err)
	assert.Nil(t, result)

	mockRoleRepo.AssertExpectations(t)
}

// ==================== Тесты для Permissions ====================

func TestRolePermissionController_GetAllPermissions_Success(t *testing.T) {
	// Arrange
	mockRoleRepo := new(MockChatRoleRepository)
	mockPermissionRepo := new(MockChatPermissionRepository)

	permissions := []models.ChatPermission{
		*createTestChatPermissionWithID(1, "send_message"),
		*createTestChatPermissionWithID(2, "delete_message"),
		*createTestChatPermissionWithID(3, "manage_users"),
	}

	mockPermissionRepo.On("GetAllPermissions").Return(permissions, nil)

	controller := controllers.NewRolePermissionController(mockRoleRepo, mockPermissionRepo)

	// Act
	result, err := controller.GetAllPermissions()

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 3)
	assert.Equal(t, 1, result[0].ID)
	assert.Equal(t, "send_message", result[0].Name)

	mockPermissionRepo.AssertExpectations(t)
}

func TestRolePermissionController_GetAllPermissions_RepositoryError(t *testing.T) {
	// Arrange
	mockRoleRepo := new(MockChatRoleRepository)
	mockPermissionRepo := new(MockChatPermissionRepository)

	repoError := errors.New("database error")
	mockPermissionRepo.On("GetAllPermissions").Return(nil, repoError)

	controller := controllers.NewRolePermissionController(mockRoleRepo, mockPermissionRepo)

	// Act
	result, err := controller.GetAllPermissions()

	// Assert
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, repoError, err)

	mockPermissionRepo.AssertExpectations(t)
}

func TestRolePermissionController_CreatePermission_Success(t *testing.T) {
	// Arrange
	mockRoleRepo := new(MockChatRoleRepository)
	mockPermissionRepo := new(MockChatPermissionRepository)

	createReq := &dto.CreatePermissionRequest{
		Name: "edit_message",
	}

	mockPermissionRepo.On("CreatePermission", mock.MatchedBy(func(p *models.ChatPermission) bool {
		return p.Name == "edit_message"
	})).Return(nil).Run(func(args mock.Arguments) {
		permission := args.Get(0).(*models.ChatPermission)
		permission.ID = 4
	})

	controller := controllers.NewRolePermissionController(mockRoleRepo, mockPermissionRepo)

	// Act
	result, err := controller.CreatePermission(createReq)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "edit_message", result.Name)
	assert.Equal(t, 4, result.ID)

	mockPermissionRepo.AssertExpectations(t)
}

func TestRolePermissionController_CreatePermission_RepositoryError(t *testing.T) {
	// Arrange
	mockRoleRepo := new(MockChatRoleRepository)
	mockPermissionRepo := new(MockChatPermissionRepository)

	createReq := &dto.CreatePermissionRequest{
		Name: "edit_message",
	}

	repoError := errors.New("permission already exists")
	mockPermissionRepo.On("CreatePermission", mock.Anything).Return(repoError)

	controller := controllers.NewRolePermissionController(mockRoleRepo, mockPermissionRepo)

	// Act
	result, err := controller.CreatePermission(createReq)

	// Assert
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, repoError, err)

	mockPermissionRepo.AssertExpectations(t)
}

func TestRolePermissionController_DeletePermission_Success(t *testing.T) {
	// Arrange
	mockRoleRepo := new(MockChatRoleRepository)
	mockPermissionRepo := new(MockChatPermissionRepository)

	permissionID := 1
	mockPermissionRepo.On("DeletePermission", permissionID).Return(nil)

	controller := controllers.NewRolePermissionController(mockRoleRepo, mockPermissionRepo)

	// Act
	err := controller.DeletePermission(permissionID)

	// Assert
	require.NoError(t, err)

	mockPermissionRepo.AssertExpectations(t)
}

func TestRolePermissionController_DeletePermission_RepositoryError(t *testing.T) {
	// Arrange
	mockRoleRepo := new(MockChatRoleRepository)
	mockPermissionRepo := new(MockChatPermissionRepository)

	permissionID := 999
	repoError := errors.New("permission not found")
	mockPermissionRepo.On("DeletePermission", permissionID).Return(repoError)

	controller := controllers.NewRolePermissionController(mockRoleRepo, mockPermissionRepo)

	// Act
	err := controller.DeletePermission(permissionID)

	// Assert
	require.Error(t, err)
	assert.Equal(t, repoError, err)

	mockPermissionRepo.AssertExpectations(t)
}
