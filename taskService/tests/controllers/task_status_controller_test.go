package controllers

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"taskService/internal/controllers"
	"taskService/internal/custom_errors"
	"taskService/internal/models"
)

// Тесты для TaskStatusController.Create

func TestTaskStatusController_Create_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockTaskStatusRepository)
	controller := controllers.NewTaskStatusController(mockRepo)

	statusName := "new_status"
	expectedStatus := createTestTaskStatusWithID(1, statusName)

	mockRepo.On("GetByName", statusName).Return(nil, gorm.ErrRecordNotFound)
	mockRepo.On("Create", statusName).Return(expectedStatus, nil)

	// Act
	result, err := controller.Create(statusName)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expectedStatus, result)
	assert.Equal(t, statusName, result.Name)

	mockRepo.AssertExpectations(t)
}

func TestTaskStatusController_Create_StatusAlreadyExists(t *testing.T) {
	// Arrange
	mockRepo := new(MockTaskStatusRepository)
	controller := controllers.NewTaskStatusController(mockRepo)

	statusName := "existing_status"
	existingStatus := createTestTaskStatusWithID(1, statusName)

	mockRepo.On("GetByName", statusName).Return(existingStatus, nil)

	// Act
	result, err := controller.Create(statusName)

	// Assert
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, custom_errors.ErrStatusAlreadyExists, err)

	mockRepo.AssertExpectations(t)
	mockRepo.AssertNotCalled(t, "Create", mock.Anything)
}

func TestTaskStatusController_Create_RepositoryError(t *testing.T) {
	// Arrange
	mockRepo := new(MockTaskStatusRepository)
	controller := controllers.NewTaskStatusController(mockRepo)

	statusName := "new_status"
	repoError := errors.New("database error")

	mockRepo.On("GetByName", statusName).Return(nil, gorm.ErrRecordNotFound)
	mockRepo.On("Create", statusName).Return(nil, repoError)

	// Act
	result, err := controller.Create(statusName)

	// Assert
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, repoError, err)

	mockRepo.AssertExpectations(t)
}

// Тесты для TaskStatusController.GetByID

func TestTaskStatusController_GetByID_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockTaskStatusRepository)
	controller := controllers.NewTaskStatusController(mockRepo)

	statusID := 1
	expectedStatus := createTestTaskStatus()

	mockRepo.On("GetByID", statusID).Return(expectedStatus, nil)

	// Act
	result, err := controller.GetByID(statusID)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expectedStatus, result)

	mockRepo.AssertExpectations(t)
}

func TestTaskStatusController_GetByID_NotFound(t *testing.T) {
	// Arrange
	mockRepo := new(MockTaskStatusRepository)
	controller := controllers.NewTaskStatusController(mockRepo)

	statusID := 999
	mockRepo.On("GetByID", statusID).Return(nil, gorm.ErrRecordNotFound)

	// Act
	result, err := controller.GetByID(statusID)

	// Assert
	require.Error(t, err)
	assert.Nil(t, result)
	var statusNotFoundErr *custom_errors.TaskStatusNotFoundError
	assert.True(t, errors.As(err, &statusNotFoundErr))

	mockRepo.AssertExpectations(t)
}

func TestTaskStatusController_GetByID_RepositoryError(t *testing.T) {
	// Arrange
	mockRepo := new(MockTaskStatusRepository)
	controller := controllers.NewTaskStatusController(mockRepo)

	statusID := 1
	repoError := errors.New("database error")
	mockRepo.On("GetByID", statusID).Return(nil, repoError)

	// Act
	result, err := controller.GetByID(statusID)

	// Assert
	require.Error(t, err)
	assert.Nil(t, result)
	var statusNotFoundErr *custom_errors.TaskStatusNotFoundError
	assert.True(t, errors.As(err, &statusNotFoundErr))

	mockRepo.AssertExpectations(t)
}

// Тесты для TaskStatusController.DeleteByID

func TestTaskStatusController_DeleteByID_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockTaskStatusRepository)
	controller := controllers.NewTaskStatusController(mockRepo)

	statusID := 1
	mockRepo.On("DeleteByID", statusID).Return(nil)

	// Act
	err := controller.DeleteByID(statusID)

	// Assert
	require.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestTaskStatusController_DeleteByID_RepositoryError(t *testing.T) {
	// Arrange
	mockRepo := new(MockTaskStatusRepository)
	controller := controllers.NewTaskStatusController(mockRepo)

	statusID := 1
	repoError := errors.New("database error")
	mockRepo.On("DeleteByID", statusID).Return(repoError)

	// Act
	err := controller.DeleteByID(statusID)

	// Assert
	require.Error(t, err)
	assert.Equal(t, repoError, err)

	mockRepo.AssertExpectations(t)
}

// Тесты для TaskStatusController.GetAll

func TestTaskStatusController_GetAll_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockTaskStatusRepository)
	controller := controllers.NewTaskStatusController(mockRepo)

	expectedStatuses := []models.TaskStatus{
		*createTestTaskStatusWithID(1, "created"),
		*createTestTaskStatusWithID(2, "in_progress"),
		*createTestTaskStatusWithID(3, "completed"),
	}

	mockRepo.On("GetAll").Return(expectedStatuses, nil)

	// Act
	result, err := controller.GetAll()

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expectedStatuses, result)
	assert.Len(t, result, 3)

	mockRepo.AssertExpectations(t)
}

func TestTaskStatusController_GetAll_RepositoryError(t *testing.T) {
	// Arrange
	mockRepo := new(MockTaskStatusRepository)
	controller := controllers.NewTaskStatusController(mockRepo)

	repoError := errors.New("database error")
	mockRepo.On("GetAll").Return(nil, repoError)

	// Act
	result, err := controller.GetAll()

	// Assert
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, repoError, err)

	mockRepo.AssertExpectations(t)
}

func TestTaskStatusController_GetAll_EmptyList(t *testing.T) {
	// Arrange
	mockRepo := new(MockTaskStatusRepository)
	controller := controllers.NewTaskStatusController(mockRepo)

	expectedStatuses := []models.TaskStatus{}
	mockRepo.On("GetAll").Return(expectedStatuses, nil)

	// Act
	result, err := controller.GetAll()

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expectedStatuses, result)
	assert.Len(t, result, 0)

	mockRepo.AssertExpectations(t)
}
