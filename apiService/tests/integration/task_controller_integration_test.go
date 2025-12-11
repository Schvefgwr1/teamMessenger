//go:build integration
// +build integration

package integration

import (
	"apiService/internal/controllers"
	"apiService/internal/dto"
	"apiService/internal/services"
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestTaskController_CreateTask_Integration тестирует создание задачи с инвалидацией кеша
func TestTaskController_CreateTask_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange
	redisClient := setupTestRedis(t)
	_, _, taskServer, _, _, _, taskClient, fileClient, _ := setupTestHTTPClients(t)
	defer taskServer.Close()

	cacheService := services.NewCacheService(redisClient)
	taskController := controllers.NewTaskController(taskClient, fileClient, cacheService)

	creatorID := uuid.New()
	executorID := uuid.New()

	// Создаем кеш перед созданием задачи
	ctx := context.Background()
	cacheService.SetUserTasksCache(ctx, creatorID.String(), []interface{}{})

	createReq := &dto.CreateTaskRequestGateway{
		Title:      "test_task",
		ExecutorID: executorID.String(),
	}

	// Act
	task, err := taskController.CreateTask(createReq, creatorID)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, task)

	// Проверяем, что кеш задач инвалидирован для создателя
	exists, _ := cacheService.Exists(ctx, cacheService.UserTasksCacheKey(creatorID.String()))
	assert.False(t, exists, "User tasks cache should be invalidated after create")
}

// TestTaskController_UpdateTaskStatus_Integration тестирует обновление статуса задачи с инвалидацией кеша
func TestTaskController_UpdateTaskStatus_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange
	redisClient := setupTestRedis(t)
	_, _, taskServer, _, _, _, taskClient, fileClient, _ := setupTestHTTPClients(t)
	defer taskServer.Close()

	cacheService := services.NewCacheService(redisClient)
	taskController := controllers.NewTaskController(taskClient, fileClient, cacheService)

	taskID := 1
	statusID := 2

	// Создаем кеш перед обновлением
	ctx := context.Background()
	cacheService.SetTaskCache(ctx, taskID, map[string]interface{}{"id": taskID})

	// Act
	err := taskController.UpdateTaskStatus(taskID, statusID)

	// Assert
	require.NoError(t, err)

	// Проверяем, что кеш задачи инвалидирован
	exists, _ := cacheService.Exists(ctx, cacheService.TaskCacheKey(taskID))
	assert.False(t, exists, "Task cache should be invalidated after status update")
}

// TestTaskController_GetTaskByID_Integration тестирует получение задачи по ID с кешированием
func TestTaskController_GetTaskByID_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange
	redisClient := setupTestRedis(t)
	_, _, taskServer, _, _, _, taskClient, fileClient, _ := setupTestHTTPClients(t)
	defer taskServer.Close()

	cacheService := services.NewCacheService(redisClient)
	taskController := controllers.NewTaskController(taskClient, fileClient, cacheService)

	taskID := 1

	// Act - первый запрос (должен идти в сервис)
	task1, err := taskController.GetTaskByID(taskID)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, task1)

	// Act - второй запрос (должен быть из кеша)
	task2, err := taskController.GetTaskByID(taskID)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, task2)

	// Проверяем, что данные реально в кеше
	ctx := context.Background()
	var cachedTask interface{}
	err = cacheService.GetTaskCache(ctx, taskID, &cachedTask)
	require.NoError(t, err)
}

// TestTaskController_GetUserTasks_Integration тестирует получение задач пользователя с кешированием
func TestTaskController_GetUserTasks_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange
	redisClient := setupTestRedis(t)
	_, _, taskServer, _, _, _, taskClient, fileClient, _ := setupTestHTTPClients(t)
	defer taskServer.Close()

	cacheService := services.NewCacheService(redisClient)
	taskController := controllers.NewTaskController(taskClient, fileClient, cacheService)

	userID := uuid.New().String()

	// Act - первый запрос (offset=0, limit=20, должен кешироваться)
	tasks1, err := taskController.GetUserTasks(userID, 20, 0)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, tasks1)

	// Act - второй запрос (должен быть из кеша)
	tasks2, err := taskController.GetUserTasks(userID, 20, 0)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, tasks2)

	// Проверяем, что данные реально в кеше
	ctx := context.Background()
	var cachedTasks []interface{}
	err = cacheService.GetUserTasksCache(ctx, userID, &cachedTasks)
	require.NoError(t, err)
}

// TestTaskController_GetAllStatuses_Integration тестирует получение всех статусов с кешированием
func TestTaskController_GetAllStatuses_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange
	redisClient := setupTestRedis(t)
	_, _, taskServer, _, _, _, taskClient, fileClient, _ := setupTestHTTPClients(t)
	defer taskServer.Close()

	cacheService := services.NewCacheService(redisClient)
	taskController := controllers.NewTaskController(taskClient, fileClient, cacheService)

	// Act - первый запрос
	statuses1, err := taskController.GetAllStatuses()

	// Assert
	require.NoError(t, err)
	assert.NotEmpty(t, statuses1)

	// Act - второй запрос (должен быть из кеша)
	statuses2, err := taskController.GetAllStatuses()

	// Assert
	require.NoError(t, err)
	assert.Equal(t, len(statuses1), len(statuses2))

	// Проверяем, что данные реально в кеше с улучшенными проверками
	ctx := context.Background()
	cacheKey := "task:statuses"
	assertCacheExists(t, cacheService, cacheKey)

	var cachedStatuses []interface{}
	err = cacheService.Get(ctx, cacheKey, &cachedStatuses)
	require.NoError(t, err)
	assert.Equal(t, len(statuses1), len(cachedStatuses))

	// Проверяем TTL кеша (должен быть установлен на 1 час)
	ttl := getCacheTTL(t, cacheService, cacheKey)
	assert.Greater(t, ttl, time.Duration(0), "Cache TTL should be greater than 0")
}

// TestTaskController_CreateStatus_Integration тестирует создание статуса с инвалидацией кеша
func TestTaskController_CreateStatus_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange
	redisClient := setupTestRedis(t)
	_, _, taskServer, _, _, _, taskClient, fileClient, _ := setupTestHTTPClients(t)
	defer taskServer.Close()

	cacheService := services.NewCacheService(redisClient)
	taskController := controllers.NewTaskController(taskClient, fileClient, cacheService)

	// Создаем кеш перед созданием статуса
	ctx := context.Background()
	cacheService.Set(ctx, "task:statuses", []interface{}{}, time.Hour)

	statusName := "test_status_" + uuid.New().String()[:8]

	// Act
	status, err := taskController.CreateStatus(statusName)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, status)

	// Проверяем, что кеш статусов инвалидирован
	exists, _ := cacheService.Exists(ctx, "task:statuses")
	assert.False(t, exists, "Statuses cache should be invalidated after create")
}

// TestTaskController_GetStatusByID_Integration тестирует получение статуса по ID с кешированием
func TestTaskController_GetStatusByID_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange
	redisClient := setupTestRedis(t)
	_, _, taskServer, _, _, _, taskClient, fileClient, _ := setupTestHTTPClients(t)
	defer taskServer.Close()

	cacheService := services.NewCacheService(redisClient)
	taskController := controllers.NewTaskController(taskClient, fileClient, cacheService)

	statusID := 1

	// Act - первый запрос
	status1, err := taskController.GetStatusByID(statusID)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, status1)

	// Act - второй запрос (должен быть из кеша)
	status2, err := taskController.GetStatusByID(statusID)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, status1.ID, status2.ID)
}

// TestTaskController_DeleteStatus_Integration тестирует удаление статуса с инвалидацией кеша
func TestTaskController_DeleteStatus_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange
	redisClient := setupTestRedis(t)
	_, _, taskServer, _, _, _, taskClient, fileClient, _ := setupTestHTTPClients(t)
	defer taskServer.Close()

	cacheService := services.NewCacheService(redisClient)
	taskController := controllers.NewTaskController(taskClient, fileClient, cacheService)

	statusID := 1

	// Создаем кеш перед удалением
	ctx := context.Background()
	cacheService.Set(ctx, "task:statuses", []interface{}{}, time.Hour)
	cacheService.Set(ctx, "task:status:1", map[string]interface{}{"id": 1}, time.Hour)

	// Act
	err := taskController.DeleteStatus(statusID)

	// Assert
	require.NoError(t, err)

	// Проверяем, что кеш инвалидирован
	exists1, _ := cacheService.Exists(ctx, "task:statuses")
	exists2, _ := cacheService.Exists(ctx, "task:status:1")
	assert.False(t, exists1, "Statuses cache should be invalidated after delete")
	assert.False(t, exists2, "Status cache should be invalidated after delete")
}
