package services

import (
	"apiService/internal/services"
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Тесты для UserChatListCache методов

func TestCacheService_SetUserChatListCache_Success(t *testing.T) {
	// Arrange
	redisClient := setupTestRedis(t)
	defer redisClient.Close()

	cacheService := services.NewCacheService(redisClient)
	ctx := context.Background()

	userID := uuid.New().String()
	chats := []map[string]string{{"id": "chat1", "name": "Chat 1"}}

	// Act
	err := cacheService.SetUserChatListCache(ctx, userID, chats)

	// Assert
	require.NoError(t, err)

	// Проверяем, что ключ создан
	key := cacheService.UserChatListCacheKey(userID)
	exists, _ := redisClient.Exists(ctx, key).Result()
	assert.Equal(t, int64(1), exists)
}

func TestCacheService_GetUserChatListCache_Success(t *testing.T) {
	// Arrange
	redisClient := setupTestRedis(t)
	defer redisClient.Close()

	cacheService := services.NewCacheService(redisClient)
	ctx := context.Background()

	userID := uuid.New().String()
	chats := []map[string]string{{"id": "chat1", "name": "Chat 1"}}

	// Сохраняем данные
	cacheService.SetUserChatListCache(ctx, userID, chats)

	// Act
	var result []map[string]string
	err := cacheService.GetUserChatListCache(ctx, userID, &result)

	// Assert
	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, chats[0]["id"], result[0]["id"])
}

func TestCacheService_DeleteUserChatListCache_Success(t *testing.T) {
	// Arrange
	redisClient := setupTestRedis(t)
	defer redisClient.Close()

	cacheService := services.NewCacheService(redisClient)
	ctx := context.Background()

	userID := uuid.New().String()
	chats := []map[string]string{{"id": "chat1", "name": "Chat 1"}}

	// Сохраняем данные
	cacheService.SetUserChatListCache(ctx, userID, chats)

	// Act
	err := cacheService.DeleteUserChatListCache(ctx, userID)

	// Assert
	require.NoError(t, err)

	// Проверяем, что ключ удален
	var result []map[string]string
	err = cacheService.GetUserChatListCache(ctx, userID, &result)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "cache miss")
}

// Тесты для ChatInfoCache методов

func TestCacheService_ChatInfoCacheKey(t *testing.T) {
	// Arrange
	redisClient := setupTestRedis(t)
	defer redisClient.Close()

	cacheService := services.NewCacheService(redisClient)

	chatID := uuid.New().String()

	// Act
	key := cacheService.ChatInfoCacheKey(chatID)

	// Assert
	assert.Contains(t, key, "chat:")
	assert.Contains(t, key, chatID)
}

func TestCacheService_SetChatInfoCache_Success(t *testing.T) {
	// Arrange
	redisClient := setupTestRedis(t)
	defer redisClient.Close()

	cacheService := services.NewCacheService(redisClient)
	ctx := context.Background()

	chatID := uuid.New().String()
	chatInfo := map[string]string{"name": "Test Chat", "description": "Test Description"}

	// Act
	err := cacheService.SetChatInfoCache(ctx, chatID, chatInfo)

	// Assert
	require.NoError(t, err)

	// Проверяем, что ключ создан
	key := cacheService.ChatInfoCacheKey(chatID)
	exists, _ := redisClient.Exists(ctx, key).Result()
	assert.Equal(t, int64(1), exists)
}

func TestCacheService_GetChatInfoCache_Success(t *testing.T) {
	// Arrange
	redisClient := setupTestRedis(t)
	defer redisClient.Close()

	cacheService := services.NewCacheService(redisClient)
	ctx := context.Background()

	chatID := uuid.New().String()
	chatInfo := map[string]string{"name": "Test Chat", "description": "Test Description"}

	// Сохраняем данные
	cacheService.SetChatInfoCache(ctx, chatID, chatInfo)

	// Act
	var result map[string]string
	err := cacheService.GetChatInfoCache(ctx, chatID, &result)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, chatInfo["name"], result["name"])
	assert.Equal(t, chatInfo["description"], result["description"])
}

func TestCacheService_DeleteChatInfoCache_Success(t *testing.T) {
	// Arrange
	redisClient := setupTestRedis(t)
	defer redisClient.Close()

	cacheService := services.NewCacheService(redisClient)
	ctx := context.Background()

	chatID := uuid.New().String()
	chatInfo := map[string]string{"name": "Test Chat"}

	// Сохраняем данные
	cacheService.SetChatInfoCache(ctx, chatID, chatInfo)

	// Act
	err := cacheService.DeleteChatInfoCache(ctx, chatID)

	// Assert
	require.NoError(t, err)

	// Проверяем, что ключ удален
	var result map[string]string
	err = cacheService.GetChatInfoCache(ctx, chatID, &result)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "cache miss")
}

// Тесты для ChatMembersCache методов

func TestCacheService_ChatMembersCacheKey(t *testing.T) {
	// Arrange
	redisClient := setupTestRedis(t)
	defer redisClient.Close()

	cacheService := services.NewCacheService(redisClient)

	chatID := uuid.New().String()

	// Act
	key := cacheService.ChatMembersCacheKey(chatID)

	// Assert
	assert.Contains(t, key, "chat_members:")
	assert.Contains(t, key, chatID)
}

func TestCacheService_DeleteChatMembersCache_Success(t *testing.T) {
	// Arrange
	redisClient := setupTestRedis(t)
	defer redisClient.Close()

	cacheService := services.NewCacheService(redisClient)
	ctx := context.Background()

	chatID := uuid.New().String()
	key := cacheService.ChatMembersCacheKey(chatID)

	// Создаем ключ
	redisClient.Set(ctx, key, "test", 5*time.Minute)

	// Act
	err := cacheService.DeleteChatMembersCache(ctx, chatID)

	// Assert
	require.NoError(t, err)

	// Проверяем, что ключ удален
	exists, _ := redisClient.Exists(ctx, key).Result()
	assert.Equal(t, int64(0), exists)
}

// Тесты для ChatUserRoleCache методов

func TestCacheService_ChatUserRoleCacheKey(t *testing.T) {
	// Arrange
	redisClient := setupTestRedis(t)
	defer redisClient.Close()

	cacheService := services.NewCacheService(redisClient)

	chatID := uuid.New().String()
	userID := uuid.New().String()

	// Act
	key := cacheService.ChatUserRoleCacheKey(chatID, userID)

	// Assert
	assert.Contains(t, key, "chat_user_role:")
	assert.Contains(t, key, chatID)
	assert.Contains(t, key, userID)
}

func TestCacheService_SetChatUserRoleCache_Success(t *testing.T) {
	// Arrange
	redisClient := setupTestRedis(t)
	defer redisClient.Close()

	cacheService := services.NewCacheService(redisClient)
	ctx := context.Background()

	chatID := uuid.New().String()
	userID := uuid.New().String()
	role := map[string]string{"id": "1", "name": "admin"}

	// Act
	err := cacheService.SetChatUserRoleCache(ctx, chatID, userID, role)

	// Assert
	require.NoError(t, err)

	// Проверяем, что ключ создан
	key := cacheService.ChatUserRoleCacheKey(chatID, userID)
	exists, _ := redisClient.Exists(ctx, key).Result()
	assert.Equal(t, int64(1), exists)
}

func TestCacheService_GetChatUserRoleCache_Success(t *testing.T) {
	// Arrange
	redisClient := setupTestRedis(t)
	defer redisClient.Close()

	cacheService := services.NewCacheService(redisClient)
	ctx := context.Background()

	chatID := uuid.New().String()
	userID := uuid.New().String()
	role := map[string]string{"id": "1", "name": "admin"}

	// Сохраняем данные
	cacheService.SetChatUserRoleCache(ctx, chatID, userID, role)

	// Act
	var result map[string]string
	err := cacheService.GetChatUserRoleCache(ctx, chatID, userID, &result)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, role["id"], result["id"])
	assert.Equal(t, role["name"], result["name"])
}

func TestCacheService_DeleteChatUserRoleCache_Success(t *testing.T) {
	// Arrange
	redisClient := setupTestRedis(t)
	defer redisClient.Close()

	cacheService := services.NewCacheService(redisClient)
	ctx := context.Background()

	chatID := uuid.New().String()
	userID := uuid.New().String()
	role := map[string]string{"id": "1", "name": "admin"}

	// Сохраняем данные
	cacheService.SetChatUserRoleCache(ctx, chatID, userID, role)

	// Act
	err := cacheService.DeleteChatUserRoleCache(ctx, chatID, userID)

	// Assert
	require.NoError(t, err)

	// Проверяем, что ключ удален
	var result map[string]string
	err = cacheService.GetChatUserRoleCache(ctx, chatID, userID, &result)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "cache miss")
}

// Тесты для TaskCache методов

func TestCacheService_TaskCacheKey(t *testing.T) {
	// Arrange
	redisClient := setupTestRedis(t)
	defer redisClient.Close()

	cacheService := services.NewCacheService(redisClient)

	taskID := 123

	// Act
	key := cacheService.TaskCacheKey(taskID)

	// Assert
	assert.Contains(t, key, "task:")
	assert.Contains(t, key, "123")
}

func TestCacheService_SetTaskCache_Success(t *testing.T) {
	// Arrange
	redisClient := setupTestRedis(t)
	defer redisClient.Close()

	cacheService := services.NewCacheService(redisClient)
	ctx := context.Background()

	taskID := 123
	task := map[string]string{"id": "123", "title": "Test Task"}

	// Act
	err := cacheService.SetTaskCache(ctx, taskID, task)

	// Assert
	require.NoError(t, err)

	// Проверяем, что ключ создан
	key := cacheService.TaskCacheKey(taskID)
	exists, _ := redisClient.Exists(ctx, key).Result()
	assert.Equal(t, int64(1), exists)
}

func TestCacheService_GetTaskCache_Success(t *testing.T) {
	// Arrange
	redisClient := setupTestRedis(t)
	defer redisClient.Close()

	cacheService := services.NewCacheService(redisClient)
	ctx := context.Background()

	taskID := 123
	task := map[string]string{"id": "123", "title": "Test Task"}

	// Сохраняем данные
	cacheService.SetTaskCache(ctx, taskID, task)

	// Act
	var result map[string]string
	err := cacheService.GetTaskCache(ctx, taskID, &result)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, task["id"], result["id"])
	assert.Equal(t, task["title"], result["title"])
}

func TestCacheService_DeleteTaskCache_Success(t *testing.T) {
	// Arrange
	redisClient := setupTestRedis(t)
	defer redisClient.Close()

	cacheService := services.NewCacheService(redisClient)
	ctx := context.Background()

	taskID := 123
	task := map[string]string{"id": "123", "title": "Test Task"}

	// Сохраняем данные
	cacheService.SetTaskCache(ctx, taskID, task)

	// Act
	err := cacheService.DeleteTaskCache(ctx, taskID)

	// Assert
	require.NoError(t, err)

	// Проверяем, что ключ удален
	var result map[string]string
	err = cacheService.GetTaskCache(ctx, taskID, &result)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "cache miss")
}

// Тесты для UserTasksCache методов

func TestCacheService_UserTasksCacheKey(t *testing.T) {
	// Arrange
	redisClient := setupTestRedis(t)
	defer redisClient.Close()

	cacheService := services.NewCacheService(redisClient)

	userID := uuid.New().String()

	// Act
	key := cacheService.UserTasksCacheKey(userID)

	// Assert
	assert.Contains(t, key, "user_tasks:")
	assert.Contains(t, key, userID)
}

func TestCacheService_SetUserTasksCache_Success(t *testing.T) {
	// Arrange
	redisClient := setupTestRedis(t)
	defer redisClient.Close()

	cacheService := services.NewCacheService(redisClient)
	ctx := context.Background()

	userID := uuid.New().String()
	tasks := []map[string]string{{"id": "1", "title": "Task 1"}}

	// Act
	err := cacheService.SetUserTasksCache(ctx, userID, tasks)

	// Assert
	require.NoError(t, err)

	// Проверяем, что ключ создан
	key := cacheService.UserTasksCacheKey(userID)
	exists, _ := redisClient.Exists(ctx, key).Result()
	assert.Equal(t, int64(1), exists)
}

func TestCacheService_GetUserTasksCache_Success(t *testing.T) {
	// Arrange
	redisClient := setupTestRedis(t)
	defer redisClient.Close()

	cacheService := services.NewCacheService(redisClient)
	ctx := context.Background()

	userID := uuid.New().String()
	tasks := []map[string]string{{"id": "1", "title": "Task 1"}}

	// Сохраняем данные
	cacheService.SetUserTasksCache(ctx, userID, tasks)

	// Act
	var result []map[string]string
	err := cacheService.GetUserTasksCache(ctx, userID, &result)

	// Assert
	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, tasks[0]["id"], result[0]["id"])
}

func TestCacheService_DeleteUserTasksCache_Success(t *testing.T) {
	// Arrange
	redisClient := setupTestRedis(t)
	defer redisClient.Close()

	cacheService := services.NewCacheService(redisClient)
	ctx := context.Background()

	userID := uuid.New().String()
	tasks := []map[string]string{{"id": "1", "title": "Task 1"}}

	// Сохраняем данные
	cacheService.SetUserTasksCache(ctx, userID, tasks)

	// Act
	err := cacheService.DeleteUserTasksCache(ctx, userID)

	// Assert
	require.NoError(t, err)

	// Проверяем, что ключ удален
	var result []map[string]string
	err = cacheService.GetUserTasksCache(ctx, userID, &result)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "cache miss")
}

// Тесты для ChatRolesCache методов

func TestCacheService_SetChatRolesCache_Success(t *testing.T) {
	// Arrange
	redisClient := setupTestRedis(t)
	defer redisClient.Close()

	cacheService := services.NewCacheService(redisClient)
	ctx := context.Background()

	roles := []map[string]string{{"id": "1", "name": "admin"}}

	// Act
	err := cacheService.SetChatRolesCache(ctx, roles)

	// Assert
	require.NoError(t, err)

	// Проверяем, что ключ создан
	exists, _ := redisClient.Exists(ctx, "chat_roles:all").Result()
	assert.Equal(t, int64(1), exists)
}

func TestCacheService_GetChatRolesCache_Success(t *testing.T) {
	// Arrange
	redisClient := setupTestRedis(t)
	defer redisClient.Close()

	cacheService := services.NewCacheService(redisClient)
	ctx := context.Background()

	roles := []map[string]string{{"id": "1", "name": "admin"}}

	// Сохраняем данные
	cacheService.SetChatRolesCache(ctx, roles)

	// Act
	var result []map[string]string
	err := cacheService.GetChatRolesCache(ctx, &result)

	// Assert
	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, roles[0]["id"], result[0]["id"])
}

func TestCacheService_DeleteChatRolesCache_Success(t *testing.T) {
	// Arrange
	redisClient := setupTestRedis(t)
	defer redisClient.Close()

	cacheService := services.NewCacheService(redisClient)
	ctx := context.Background()

	roles := []map[string]string{{"id": "1", "name": "admin"}}

	// Сохраняем данные
	cacheService.SetChatRolesCache(ctx, roles)

	// Act
	err := cacheService.DeleteChatRolesCache(ctx)

	// Assert
	require.NoError(t, err)

	// Проверяем, что ключ удален
	var result []map[string]string
	err = cacheService.GetChatRolesCache(ctx, &result)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "cache miss")
}

// Тесты для ChatPermissionsCache методов

func TestCacheService_SetChatPermissionsCache_Success(t *testing.T) {
	// Arrange
	redisClient := setupTestRedis(t)
	defer redisClient.Close()

	cacheService := services.NewCacheService(redisClient)
	ctx := context.Background()

	permissions := []map[string]string{{"id": "1", "name": "read"}}

	// Act
	err := cacheService.SetChatPermissionsCache(ctx, permissions)

	// Assert
	require.NoError(t, err)

	// Проверяем, что ключ создан
	exists, _ := redisClient.Exists(ctx, "chat_permissions:all").Result()
	assert.Equal(t, int64(1), exists)
}

func TestCacheService_GetChatPermissionsCache_Success(t *testing.T) {
	// Arrange
	redisClient := setupTestRedis(t)
	defer redisClient.Close()

	cacheService := services.NewCacheService(redisClient)
	ctx := context.Background()

	permissions := []map[string]string{{"id": "1", "name": "read"}}

	// Сохраняем данные
	cacheService.SetChatPermissionsCache(ctx, permissions)

	// Act
	var result []map[string]string
	err := cacheService.GetChatPermissionsCache(ctx, &result)

	// Assert
	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, permissions[0]["id"], result[0]["id"])
}

func TestCacheService_DeleteChatPermissionsCache_Success(t *testing.T) {
	// Arrange
	redisClient := setupTestRedis(t)
	defer redisClient.Close()

	cacheService := services.NewCacheService(redisClient)
	ctx := context.Background()

	permissions := []map[string]string{{"id": "1", "name": "read"}}

	// Сохраняем данные
	cacheService.SetChatPermissionsCache(ctx, permissions)

	// Act
	err := cacheService.DeleteChatPermissionsCache(ctx)

	// Assert
	require.NoError(t, err)

	// Проверяем, что ключ удален
	var result []map[string]string
	err = cacheService.GetChatPermissionsCache(ctx, &result)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "cache miss")
}

// Тесты для SearchCache методов

func TestCacheService_SearchCacheKey(t *testing.T) {
	// Arrange
	redisClient := setupTestRedis(t)
	defer redisClient.Close()

	cacheService := services.NewCacheService(redisClient)

	chatID := uuid.New().String()
	queryHash := "abc123"

	// Act
	key := cacheService.SearchCacheKey(chatID, queryHash)

	// Assert
	assert.Contains(t, key, "search:")
	assert.Contains(t, key, chatID)
	assert.Contains(t, key, queryHash)
}

func TestCacheService_SetSearchCache_Success(t *testing.T) {
	// Arrange
	redisClient := setupTestRedis(t)
	defer redisClient.Close()

	cacheService := services.NewCacheService(redisClient)
	ctx := context.Background()

	chatID := uuid.New().String()
	queryHash := "abc123"
	result := []map[string]string{{"id": "1", "text": "found"}}

	// Act
	err := cacheService.SetSearchCache(ctx, chatID, queryHash, result)

	// Assert
	require.NoError(t, err)

	// Проверяем, что ключ создан
	key := cacheService.SearchCacheKey(chatID, queryHash)
	exists, _ := redisClient.Exists(ctx, key).Result()
	assert.Equal(t, int64(1), exists)
}

func TestCacheService_GetSearchCache_Success(t *testing.T) {
	// Arrange
	redisClient := setupTestRedis(t)
	defer redisClient.Close()

	cacheService := services.NewCacheService(redisClient)
	ctx := context.Background()

	chatID := uuid.New().String()
	queryHash := "abc123"
	result := []map[string]string{{"id": "1", "text": "found"}}

	// Сохраняем данные
	cacheService.SetSearchCache(ctx, chatID, queryHash, result)

	// Act
	var searchResult []map[string]string
	err := cacheService.GetSearchCache(ctx, chatID, queryHash, &searchResult)

	// Assert
	require.NoError(t, err)
	assert.Len(t, searchResult, 1)
	assert.Equal(t, result[0]["id"], searchResult[0]["id"])
}

func TestCacheService_DeleteSearchCacheByChat_Success(t *testing.T) {
	// Arrange
	redisClient := setupTestRedis(t)
	defer redisClient.Close()

	cacheService := services.NewCacheService(redisClient)
	ctx := context.Background()

	chatID := uuid.New().String()
	queryHash1 := "hash1"
	queryHash2 := "hash2"
	result := []map[string]string{{"id": "1", "text": "found"}}

	// Сохраняем несколько результатов поиска для одного чата
	cacheService.SetSearchCache(ctx, chatID, queryHash1, result)
	cacheService.SetSearchCache(ctx, chatID, queryHash2, result)

	// Act
	err := cacheService.DeleteSearchCacheByChat(ctx, chatID)

	// Assert
	require.NoError(t, err)

	// Проверяем, что все ключи удалены
	var searchResult []map[string]string
	err = cacheService.GetSearchCache(ctx, chatID, queryHash1, &searchResult)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "cache miss")

	err = cacheService.GetSearchCache(ctx, chatID, queryHash2, &searchResult)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "cache miss")
}
