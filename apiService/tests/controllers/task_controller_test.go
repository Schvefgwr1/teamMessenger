package controllers

import (
	"apiService/internal/controllers"
	"apiService/internal/dto"
	"apiService/internal/services"
	"context"
	"errors"
	"testing"
	"time"

	at "common/contracts/api-task"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Тесты для TaskController.CreateTask

func TestTaskController_CreateTask_Success(t *testing.T) {
	// Arrange
	mockTaskClient := new(MockTaskClient)
	mockFileClient := new(MockFileClient)
	redisClient := setupTestRedis(t)
	defer redisClient.Close()
	cacheService := services.NewCacheService(redisClient)

	controller := controllers.NewTaskController(mockTaskClient, mockFileClient, cacheService)

	creatorID := uuid.New()
	executorID := uuid.New()
	req := &dto.CreateTaskRequestGateway{
		Title:       "Test Task",
		Description: stringPtr("Test Description"),
		ExecutorID:  executorID.String(),
		ChatID:      nil,
		Files:       nil,
	}

	expectedTask := &at.TaskResponse{
		ID:    1,
		Title: "Test Task",
	}

	mockTaskClient.On("CreateTask", mock.Anything).Return(expectedTask, nil)

	// Act
	result, err := controller.CreateTask(req, creatorID)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedTask.Title, result.Title)

	mockTaskClient.AssertExpectations(t)
}

func TestTaskController_CreateTask_ServiceError(t *testing.T) {
	// Arrange
	mockTaskClient := new(MockTaskClient)
	mockFileClient := new(MockFileClient)
	redisClient := setupTestRedis(t)
	defer redisClient.Close()
	cacheService := services.NewCacheService(redisClient)

	controller := controllers.NewTaskController(mockTaskClient, mockFileClient, cacheService)

	creatorID := uuid.New()
	executorID := uuid.New()
	req := &dto.CreateTaskRequestGateway{
		Title:       "Test Task",
		Description: stringPtr("Test Description"),
		ExecutorID:  executorID.String(),
		ChatID:      nil,
		Files:       nil,
	}

	serviceError := errors.New("service error")

	mockTaskClient.On("CreateTask", mock.Anything).Return(nil, serviceError)

	// Act
	result, err := controller.CreateTask(req, creatorID)

	// Assert
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "error of task client")

	mockTaskClient.AssertExpectations(t)
}

func TestTaskController_CreateTask_InvalidUUID(t *testing.T) {
	// Arrange
	mockTaskClient := new(MockTaskClient)
	mockFileClient := new(MockFileClient)
	redisClient := setupTestRedis(t)
	defer redisClient.Close()
	cacheService := services.NewCacheService(redisClient)

	controller := controllers.NewTaskController(mockTaskClient, mockFileClient, cacheService)

	creatorID := uuid.New()
	req := &dto.CreateTaskRequestGateway{
		Title:       "Test Task",
		Description: stringPtr("Test Description"),
		ExecutorID:  "invalid-uuid",
		ChatID:      nil,
		Files:       nil,
	}

	// Act
	result, err := controller.CreateTask(req, creatorID)

	// Assert
	require.Error(t, err)
	assert.Nil(t, result)

	mockTaskClient.AssertNotCalled(t, "CreateTask", mock.Anything)
}

// Тесты для TaskController.GetTaskByID

func TestTaskController_GetTaskByID_Success_FromCache(t *testing.T) {
	// Arrange
	mockTaskClient := new(MockTaskClient)
	mockFileClient := new(MockFileClient)
	redisClient := setupTestRedis(t)
	defer redisClient.Close()
	cacheService := services.NewCacheService(redisClient)

	controller := controllers.NewTaskController(mockTaskClient, mockFileClient, cacheService)

	taskID := 1
	cachedTask := &at.TaskServiceResponse{
		Task: &at.TaskResponse{
			ID:    taskID,
			Title: "Test Task",
		},
	}

	// Сохраняем задачу в кеш
	cacheService.SetTaskCache(context.Background(), taskID, cachedTask)

	// Act
	result, err := controller.GetTaskByID(taskID)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, cachedTask.Task.ID, result.Task.ID)

	mockTaskClient.AssertNotCalled(t, "GetTaskByID", mock.Anything)
}

func TestTaskController_GetTaskByID_Success_FromService(t *testing.T) {
	// Arrange
	mockTaskClient := new(MockTaskClient)
	mockFileClient := new(MockFileClient)
	redisClient := setupTestRedis(t)
	defer redisClient.Close()
	cacheService := services.NewCacheService(redisClient)

	controller := controllers.NewTaskController(mockTaskClient, mockFileClient, cacheService)

	taskID := 1
	expectedTask := &at.TaskServiceResponse{
		Task: &at.TaskResponse{
			ID:    taskID,
			Title: "Test Task",
		},
	}

	mockTaskClient.On("GetTaskByID", taskID).Return(expectedTask, nil)

	// Act
	result, err := controller.GetTaskByID(taskID)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedTask.Task.ID, result.Task.ID)

	mockTaskClient.AssertExpectations(t)
}

func TestTaskController_GetTaskByID_ServiceError(t *testing.T) {
	// Arrange
	mockTaskClient := new(MockTaskClient)
	mockFileClient := new(MockFileClient)
	redisClient := setupTestRedis(t)
	defer redisClient.Close()
	cacheService := services.NewCacheService(redisClient)

	controller := controllers.NewTaskController(mockTaskClient, mockFileClient, cacheService)

	taskID := 1
	serviceError := errors.New("service error")

	mockTaskClient.On("GetTaskByID", taskID).Return(nil, serviceError)

	// Act
	result, err := controller.GetTaskByID(taskID)

	// Assert
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, serviceError, err)

	mockTaskClient.AssertExpectations(t)
}

// Тесты для TaskController.UpdateTaskStatus

func TestTaskController_UpdateTaskStatus_Success(t *testing.T) {
	// Arrange
	mockTaskClient := new(MockTaskClient)
	mockFileClient := new(MockFileClient)
	redisClient := setupTestRedis(t)
	defer redisClient.Close()
	cacheService := services.NewCacheService(redisClient)

	controller := controllers.NewTaskController(mockTaskClient, mockFileClient, cacheService)

	taskID := 1
	statusID := 2

	mockTaskClient.On("UpdateTaskStatus", taskID, statusID).Return(nil)

	// Act
	err := controller.UpdateTaskStatus(taskID, statusID)

	// Assert
	require.NoError(t, err)

	mockTaskClient.AssertExpectations(t)
}

func TestTaskController_UpdateTaskStatus_ServiceError(t *testing.T) {
	// Arrange
	mockTaskClient := new(MockTaskClient)
	mockFileClient := new(MockFileClient)
	redisClient := setupTestRedis(t)
	defer redisClient.Close()
	cacheService := services.NewCacheService(redisClient)

	controller := controllers.NewTaskController(mockTaskClient, mockFileClient, cacheService)

	taskID := 1
	statusID := 2
	serviceError := errors.New("service error")

	mockTaskClient.On("UpdateTaskStatus", taskID, statusID).Return(serviceError)

	// Act
	err := controller.UpdateTaskStatus(taskID, statusID)

	// Assert
	require.Error(t, err)
	assert.Equal(t, serviceError, err)

	mockTaskClient.AssertExpectations(t)
}

// Тесты для TaskController.GetUserTasks

func TestTaskController_GetUserTasks_Success_FromCache(t *testing.T) {
	// Arrange
	mockTaskClient := new(MockTaskClient)
	mockFileClient := new(MockFileClient)
	redisClient := setupTestRedis(t)
	defer redisClient.Close()
	cacheService := services.NewCacheService(redisClient)

	controller := controllers.NewTaskController(mockTaskClient, mockFileClient, cacheService)

	userID := uuid.New().String()
	cachedTasks := []at.TaskToList{
		{ID: 1, Title: "Task 1"},
		{ID: 2, Title: "Task 2"},
	}

	// Сохраняем задачи в кеш
	cacheService.SetUserTasksCache(context.Background(), userID, cachedTasks)

	// Act
	result, err := controller.GetUserTasks(userID, 10, 0)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, *result, 2)

	mockTaskClient.AssertNotCalled(t, "GetUserTasks", mock.Anything, mock.Anything, mock.Anything)
}

func TestTaskController_GetUserTasks_Success_FromService(t *testing.T) {
	// Arrange
	mockTaskClient := new(MockTaskClient)
	mockFileClient := new(MockFileClient)
	redisClient := setupTestRedis(t)
	defer redisClient.Close()
	cacheService := services.NewCacheService(redisClient)

	controller := controllers.NewTaskController(mockTaskClient, mockFileClient, cacheService)

	userID := uuid.New().String()
	expectedTasks := &[]at.TaskToList{
		{ID: 1, Title: "Task 1"},
		{ID: 2, Title: "Task 2"},
	}

	mockTaskClient.On("GetUserTasks", userID, 20, 0).Return(expectedTasks, nil)

	// Act
	result, err := controller.GetUserTasks(userID, 10, 0)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, *result, 2)

	mockTaskClient.AssertExpectations(t)
}

// Тесты для TaskController.GetAllStatuses

func TestTaskController_GetAllStatuses_Success_FromCache(t *testing.T) {
	// Arrange
	mockTaskClient := new(MockTaskClient)
	mockFileClient := new(MockFileClient)
	redisClient := setupTestRedis(t)
	defer redisClient.Close()
	cacheService := services.NewCacheService(redisClient)

	controller := controllers.NewTaskController(mockTaskClient, mockFileClient, cacheService)

	cachedStatuses := []at.TaskStatus{
		{ID: 1, Name: "Todo"},
		{ID: 2, Name: "In Progress"},
	}

	// Сохраняем статусы в кеш
	cacheService.Set(context.Background(), "task:statuses", cachedStatuses, time.Hour)

	// Act
	result, err := controller.GetAllStatuses()

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 2)

	mockTaskClient.AssertNotCalled(t, "GetAllStatuses")
}

func TestTaskController_GetAllStatuses_Success_FromService(t *testing.T) {
	// Arrange
	mockTaskClient := new(MockTaskClient)
	mockFileClient := new(MockFileClient)
	redisClient := setupTestRedis(t)
	defer redisClient.Close()
	cacheService := services.NewCacheService(redisClient)

	controller := controllers.NewTaskController(mockTaskClient, mockFileClient, cacheService)

	expectedStatuses := []at.TaskStatus{
		{ID: 1, Name: "Todo"},
		{ID: 2, Name: "In Progress"},
	}

	mockTaskClient.On("GetAllStatuses").Return(expectedStatuses, nil)

	// Act
	result, err := controller.GetAllStatuses()

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 2)

	mockTaskClient.AssertExpectations(t)
}

// Тесты для TaskController.CreateStatus

func TestTaskController_CreateStatus_Success(t *testing.T) {
	// Arrange
	mockTaskClient := new(MockTaskClient)
	mockFileClient := new(MockFileClient)
	redisClient := setupTestRedis(t)
	defer redisClient.Close()
	cacheService := services.NewCacheService(redisClient)

	controller := controllers.NewTaskController(mockTaskClient, mockFileClient, cacheService)

	statusName := "New Status"
	expectedStatus := &at.TaskStatus{
		ID:   3,
		Name: statusName,
	}

	mockTaskClient.On("CreateStatus", mock.Anything).Return(expectedStatus, nil)

	// Act
	result, err := controller.CreateStatus(statusName)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedStatus.Name, result.Name)

	mockTaskClient.AssertExpectations(t)
}

// Тесты для TaskController.GetStatusByID

func TestTaskController_GetStatusByID_Success_FromCache(t *testing.T) {
	// Arrange
	mockTaskClient := new(MockTaskClient)
	mockFileClient := new(MockFileClient)
	redisClient := setupTestRedis(t)
	defer redisClient.Close()
	cacheService := services.NewCacheService(redisClient)

	controller := controllers.NewTaskController(mockTaskClient, mockFileClient, cacheService)

	statusID := 1
	cachedStatus := &at.TaskStatus{
		ID:   statusID,
		Name: "Todo",
	}

	// Сохраняем статус в кеш
	cacheService.Set(context.Background(), "task:status:1", cachedStatus, time.Hour)

	// Act
	result, err := controller.GetStatusByID(statusID)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, cachedStatus.ID, result.ID)

	mockTaskClient.AssertNotCalled(t, "GetStatusByID", mock.Anything)
}

func TestTaskController_GetStatusByID_Success_FromService(t *testing.T) {
	// Arrange
	mockTaskClient := new(MockTaskClient)
	mockFileClient := new(MockFileClient)
	redisClient := setupTestRedis(t)
	defer redisClient.Close()
	cacheService := services.NewCacheService(redisClient)

	controller := controllers.NewTaskController(mockTaskClient, mockFileClient, cacheService)

	statusID := 1
	expectedStatus := &at.TaskStatus{
		ID:   statusID,
		Name: "Todo",
	}

	mockTaskClient.On("GetStatusByID", statusID).Return(expectedStatus, nil)

	// Act
	result, err := controller.GetStatusByID(statusID)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedStatus.ID, result.ID)

	mockTaskClient.AssertExpectations(t)
}

// Тесты для TaskController.DeleteStatus

func TestTaskController_DeleteStatus_Success(t *testing.T) {
	// Arrange
	mockTaskClient := new(MockTaskClient)
	mockFileClient := new(MockFileClient)
	redisClient := setupTestRedis(t)
	defer redisClient.Close()
	cacheService := services.NewCacheService(redisClient)

	controller := controllers.NewTaskController(mockTaskClient, mockFileClient, cacheService)

	statusID := 1

	mockTaskClient.On("DeleteStatus", statusID).Return(nil)

	// Act
	err := controller.DeleteStatus(statusID)

	// Assert
	require.NoError(t, err)

	mockTaskClient.AssertExpectations(t)
}
