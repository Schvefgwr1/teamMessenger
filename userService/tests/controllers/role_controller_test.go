package controllers

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"userService/internal/controllers"
	"userService/internal/handlers/dto"
	"userService/internal/models"
)

// MockPermissionRepository для role_controller_test
// Используем тот же мок, что и в permission_controller_test

// Тесты для RoleController.GetRoles

func TestRoleController_GetRoles_Success(t *testing.T) {
	// Arrange
	mockRoleRepo := new(MockRoleRepository)
	mockPermRepo := new(MockPermissionRepository)

	expectedRoles := []models.Role{
		{ID: intPtr(1), Name: "admin", Description: "Admin role"},
		{ID: intPtr(2), Name: "user", Description: "User role"},
	}

	mockRoleRepo.On("GetAllRoles").Return(expectedRoles, nil)

	controller := controllers.NewRoleController(mockRoleRepo, mockPermRepo)

	// Act
	roles, err := controller.GetRoles()

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, roles)
	assert.Len(t, roles, 2)
	assert.Equal(t, expectedRoles[0].Name, roles[0].Name)

	mockRoleRepo.AssertExpectations(t)
}

func TestRoleController_GetRoles_RepositoryError(t *testing.T) {
	// Arrange
	mockRoleRepo := new(MockRoleRepository)
	mockPermRepo := new(MockPermissionRepository)

	repoError := errors.New("database error")
	mockRoleRepo.On("GetAllRoles").Return(nil, repoError)

	controller := controllers.NewRoleController(mockRoleRepo, mockPermRepo)

	// Act
	roles, err := controller.GetRoles()

	// Assert
	require.Error(t, err)
	assert.Nil(t, roles)
	assert.Equal(t, repoError, err)

	mockRoleRepo.AssertExpectations(t)
}

// Тесты для RoleController.CreateRole

func TestRoleController_CreateRole_Success(t *testing.T) {
	// Arrange
	mockRoleRepo := new(MockRoleRepository)
	mockPermRepo := new(MockPermissionRepository)

	roleDTO := &dto.CreateRole{
		Name:          "moderator",
		Description:   "Moderator role",
		PermissionIds: []int{1, 2},
	}

	perm1 := &models.Permission{ID: 1, Name: "read"}
	perm2 := &models.Permission{ID: 2, Name: "write"}

	mockPermRepo.On("GetPermissionById", 1).Return(perm1, nil)
	mockPermRepo.On("GetPermissionById", 2).Return(perm2, nil)
	mockRoleRepo.On("CreateRole", mock.MatchedBy(func(r *models.Role) bool {
		return r.Name == roleDTO.Name && len(r.Permissions) == 2
	})).Return(nil)

	controller := controllers.NewRoleController(mockRoleRepo, mockPermRepo)

	// Act
	err := controller.CreateRole(roleDTO)

	// Assert
	require.NoError(t, err)
	mockRoleRepo.AssertExpectations(t)
	mockPermRepo.AssertExpectations(t)
}

func TestRoleController_CreateRole_InvalidPermission(t *testing.T) {
	// Arrange
	mockRoleRepo := new(MockRoleRepository)
	mockPermRepo := new(MockPermissionRepository)

	roleDTO := &dto.CreateRole{
		Name:          "moderator",
		Description:   "Moderator role",
		PermissionIds: []int{1, 999}, // 999 не существует
	}

	perm1 := &models.Permission{ID: 1, Name: "read"}

	mockPermRepo.On("GetPermissionById", 1).Return(perm1, nil)
	mockPermRepo.On("GetPermissionById", 999).Return(nil, gorm.ErrRecordNotFound)
	// Контроллер пропускает несуществующие permissions и создает роль только с существующими
	mockRoleRepo.On("CreateRole", mock.MatchedBy(func(r *models.Role) bool {
		return r.Name == roleDTO.Name && len(r.Permissions) == 1 // только один permission добавлен
	})).Return(nil)

	controller := controllers.NewRoleController(mockRoleRepo, mockPermRepo)

	// Act
	err := controller.CreateRole(roleDTO)

	// Assert
	require.NoError(t, err) // Контроллер пропускает несуществующие permissions с логом
	mockPermRepo.AssertExpectations(t)
	mockRoleRepo.AssertExpectations(t)
}

func TestRoleController_CreateRole_NoPermissions(t *testing.T) {
	// Arrange
	mockRoleRepo := new(MockRoleRepository)
	mockPermRepo := new(MockPermissionRepository)

	roleDTO := &dto.CreateRole{
		Name:          "guest",
		Description:   "Guest role",
		PermissionIds: []int{},
	}

	mockRoleRepo.On("CreateRole", mock.MatchedBy(func(r *models.Role) bool {
		return r.Name == roleDTO.Name && len(r.Permissions) == 0
	})).Return(nil)

	controller := controllers.NewRoleController(mockRoleRepo, mockPermRepo)

	// Act
	err := controller.CreateRole(roleDTO)

	// Assert
	require.NoError(t, err)
	mockRoleRepo.AssertExpectations(t)
}

// Тесты для RoleController.DeleteRole

func TestRoleController_DeleteRole_Success(t *testing.T) {
	// Arrange
	mockRoleRepo := new(MockRoleRepository)
	mockPermRepo := new(MockPermissionRepository)

	roleID := 1
	role := createTestRole()

	mockRoleRepo.On("GetRoleByID", roleID).Return(role, nil)
	mockRoleRepo.On("DeleteRole", roleID).Return(nil)

	controller := controllers.NewRoleController(mockRoleRepo, mockPermRepo)

	// Act
	err := controller.DeleteRole(roleID)

	// Assert
	require.NoError(t, err)
	mockRoleRepo.AssertExpectations(t)
}

func TestRoleController_DeleteRole_NotFound(t *testing.T) {
	// Arrange
	mockRoleRepo := new(MockRoleRepository)
	mockPermRepo := new(MockPermissionRepository)

	roleID := 999

	mockRoleRepo.On("GetRoleByID", roleID).Return(nil, gorm.ErrRecordNotFound)

	controller := controllers.NewRoleController(mockRoleRepo, mockPermRepo)

	// Act
	err := controller.DeleteRole(roleID)

	// Assert
	require.Error(t, err)
	assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))

	mockRoleRepo.AssertExpectations(t)
	mockRoleRepo.AssertNotCalled(t, "DeleteRole", mock.Anything)
}

// Тесты для RoleController.UpdateRolePermissions

func TestRoleController_UpdateRolePermissions_Success(t *testing.T) {
	// Arrange
	mockRoleRepo := new(MockRoleRepository)
	mockPermRepo := new(MockPermissionRepository)

	roleID := 1
	permissionIDs := []int{1, 2, 3}
	role := createTestRole()

	perm1 := &models.Permission{ID: 1, Name: "read"}
	perm2 := &models.Permission{ID: 2, Name: "write"}
	perm3 := &models.Permission{ID: 3, Name: "delete"}

	mockRoleRepo.On("GetRoleByID", roleID).Return(role, nil)
	mockPermRepo.On("GetPermissionById", 1).Return(perm1, nil)
	mockPermRepo.On("GetPermissionById", 2).Return(perm2, nil)
	mockPermRepo.On("GetPermissionById", 3).Return(perm3, nil)
	mockRoleRepo.On("UpdateRolePermissions", roleID, permissionIDs).Return(nil)

	controller := controllers.NewRoleController(mockRoleRepo, mockPermRepo)

	// Act
	err := controller.UpdateRolePermissions(roleID, permissionIDs)

	// Assert
	require.NoError(t, err)
	mockRoleRepo.AssertExpectations(t)
	mockPermRepo.AssertExpectations(t)
}

func TestRoleController_UpdateRolePermissions_RoleNotFound(t *testing.T) {
	// Arrange
	mockRoleRepo := new(MockRoleRepository)
	mockPermRepo := new(MockPermissionRepository)

	roleID := 999
	permissionIDs := []int{1, 2}

	mockRoleRepo.On("GetRoleByID", roleID).Return(nil, gorm.ErrRecordNotFound)

	controller := controllers.NewRoleController(mockRoleRepo, mockPermRepo)

	// Act
	err := controller.UpdateRolePermissions(roleID, permissionIDs)

	// Assert
	require.Error(t, err)
	assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))

	mockRoleRepo.AssertExpectations(t)
	mockPermRepo.AssertNotCalled(t, "GetPermissionById", mock.Anything)
	mockRoleRepo.AssertNotCalled(t, "UpdateRolePermissions", mock.Anything, mock.Anything)
}

func TestRoleController_UpdateRolePermissions_PermissionNotFound(t *testing.T) {
	// Arrange
	mockRoleRepo := new(MockRoleRepository)
	mockPermRepo := new(MockPermissionRepository)

	roleID := 1
	permissionIDs := []int{1, 999} // 999 не существует
	role := createTestRole()

	perm1 := &models.Permission{ID: 1, Name: "read"}

	mockRoleRepo.On("GetRoleByID", roleID).Return(role, nil)
	mockPermRepo.On("GetPermissionById", 1).Return(perm1, nil)
	mockPermRepo.On("GetPermissionById", 999).Return(nil, gorm.ErrRecordNotFound)

	controller := controllers.NewRoleController(mockRoleRepo, mockPermRepo)

	// Act
	err := controller.UpdateRolePermissions(roleID, permissionIDs)

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "permission with id 999 not found")

	mockRoleRepo.AssertExpectations(t)
	mockPermRepo.AssertExpectations(t)
	mockRoleRepo.AssertNotCalled(t, "UpdateRolePermissions", mock.Anything, mock.Anything)
}

// Вспомогательные функции уже определены в user_controller_test.go
