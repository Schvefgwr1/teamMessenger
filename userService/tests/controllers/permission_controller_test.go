package controllers

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"userService/internal/controllers"
	"userService/internal/models"
)

// MockPermissionRepository вынесен в mocks.go

// Тесты для PermissionController.GetPermissions

func TestPermissionController_GetPermissions_Success(t *testing.T) {
	// Arrange
	mockPermRepo := new(MockPermissionRepository)

	expectedPermissions := []models.Permission{
		{ID: 1, Name: "read", Description: "Read permission"},
		{ID: 2, Name: "write", Description: "Write permission"},
	}

	mockPermRepo.On("GetAllPermissions").Return(expectedPermissions, nil)

	controller := controllers.NewPermissionController(mockPermRepo)

	// Act
	permissions, err := controller.GetPermissions()

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, permissions)
	assert.Len(t, permissions, 2)
	assert.Equal(t, expectedPermissions[0].Name, permissions[0].Name)

	mockPermRepo.AssertExpectations(t)
}

func TestPermissionController_GetPermissions_EmptyList(t *testing.T) {
	// Arrange
	mockPermRepo := new(MockPermissionRepository)

	expectedPermissions := []models.Permission{}

	mockPermRepo.On("GetAllPermissions").Return(expectedPermissions, nil)

	controller := controllers.NewPermissionController(mockPermRepo)

	// Act
	permissions, err := controller.GetPermissions()

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, permissions)
	assert.Len(t, permissions, 0)

	mockPermRepo.AssertExpectations(t)
}

func TestPermissionController_GetPermissions_RepositoryError(t *testing.T) {
	// Arrange
	mockPermRepo := new(MockPermissionRepository)

	repoError := errors.New("database error")
	mockPermRepo.On("GetAllPermissions").Return(nil, repoError)

	controller := controllers.NewPermissionController(mockPermRepo)

	// Act
	permissions, err := controller.GetPermissions()

	// Assert
	require.Error(t, err)
	assert.Nil(t, permissions)
	assert.Equal(t, repoError, err)

	mockPermRepo.AssertExpectations(t)
}
