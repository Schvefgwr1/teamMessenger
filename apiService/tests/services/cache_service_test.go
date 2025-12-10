package services

import (
	"apiService/internal/services"
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupTestRedis создает тестовый Redis клиент с miniredis
func setupTestRedis(t *testing.T) *redis.Client {
	mr, err := miniredis.Run()
	require.NoError(t, err)

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	return client
}

// Тесты для CacheService.Set

func TestCacheService_Set_Success(t *testing.T) {
	// Arrange
	redisClient := setupTestRedis(t)
	defer redisClient.Close()

	cacheService := services.NewCacheService(redisClient)
	ctx := context.Background()

	key := "test:key"
	value := map[string]string{"test": "value"}
	ttl := 5 * time.Minute

	// Act
	err := cacheService.Set(ctx, key, value, ttl)

	// Assert
	require.NoError(t, err)

	// Проверяем, что значение сохранено
	data, err := redisClient.Get(ctx, key).Result()
	require.NoError(t, err)

	var result map[string]string
	err = json.Unmarshal([]byte(data), &result)
	require.NoError(t, err)
	assert.Equal(t, value["test"], result["test"])
}

func TestCacheService_Set_MarshalError(t *testing.T) {
	// Arrange
	redisClient := setupTestRedis(t)
	defer redisClient.Close()

	cacheService := services.NewCacheService(redisClient)
	ctx := context.Background()

	key := "test:key"
	// Используем канал, который нельзя сериализовать в JSON
	value := make(chan int)
	ttl := 5 * time.Minute

	// Act
	err := cacheService.Set(ctx, key, value, ttl)

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "marshal")
}

// Тесты для CacheService.Get

func TestCacheService_Get_Success(t *testing.T) {
	// Arrange
	redisClient := setupTestRedis(t)
	defer redisClient.Close()

	cacheService := services.NewCacheService(redisClient)
	ctx := context.Background()

	key := "test:key"
	value := map[string]string{"test": "value"}
	ttl := 5 * time.Minute

	// Сохраняем значение
	data, _ := json.Marshal(value)
	redisClient.Set(ctx, key, data, ttl)

	// Act
	var result map[string]string
	err := cacheService.Get(ctx, key, &result)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, value["test"], result["test"])
}

func TestCacheService_Get_CacheMiss(t *testing.T) {
	// Arrange
	redisClient := setupTestRedis(t)
	defer redisClient.Close()

	cacheService := services.NewCacheService(redisClient)
	ctx := context.Background()

	key := "test:nonexistent"

	// Act
	var result map[string]string
	err := cacheService.Get(ctx, key, &result)

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "cache miss")
}

func TestCacheService_Get_UnmarshalError(t *testing.T) {
	// Arrange
	redisClient := setupTestRedis(t)
	defer redisClient.Close()

	cacheService := services.NewCacheService(redisClient)
	ctx := context.Background()

	key := "test:key"
	// Сохраняем невалидный JSON
	redisClient.Set(ctx, key, "invalid json", 5*time.Minute)

	// Act
	var result map[string]string
	err := cacheService.Get(ctx, key, &result)

	// Assert
	require.Error(t, err)
	// Ошибка может быть "invalid character" или "unmarshal"
	assert.True(t, err != nil)
}

// Тесты для CacheService.Delete

func TestCacheService_Delete_Success(t *testing.T) {
	// Arrange
	redisClient := setupTestRedis(t)
	defer redisClient.Close()

	cacheService := services.NewCacheService(redisClient)
	ctx := context.Background()

	key := "test:key"
	value := "test value"
	redisClient.Set(ctx, key, value, 5*time.Minute)

	// Act
	err := cacheService.Delete(ctx, key)

	// Assert
	require.NoError(t, err)

	// Проверяем, что ключ удален
	exists, _ := redisClient.Exists(ctx, key).Result()
	assert.Equal(t, int64(0), exists)
}

func TestCacheService_Delete_NonexistentKey(t *testing.T) {
	// Arrange
	redisClient := setupTestRedis(t)
	defer redisClient.Close()

	cacheService := services.NewCacheService(redisClient)
	ctx := context.Background()

	key := "test:nonexistent"

	// Act
	err := cacheService.Delete(ctx, key)

	// Assert
	// Удаление несуществующего ключа не должно возвращать ошибку
	require.NoError(t, err)
}

// Тесты для CacheService.DeleteByPattern

func TestCacheService_DeleteByPattern_Success(t *testing.T) {
	// Arrange
	redisClient := setupTestRedis(t)
	defer redisClient.Close()

	cacheService := services.NewCacheService(redisClient)
	ctx := context.Background()

	// Создаем несколько ключей с паттерном
	redisClient.Set(ctx, "test:key1", "value1", 5*time.Minute)
	redisClient.Set(ctx, "test:key2", "value2", 5*time.Minute)
	redisClient.Set(ctx, "other:key1", "value3", 5*time.Minute)

	pattern := "test:*"

	// Act
	err := cacheService.DeleteByPattern(ctx, pattern)

	// Assert
	require.NoError(t, err)

	// Проверяем, что ключи с паттерном удалены
	exists1, _ := redisClient.Exists(ctx, "test:key1").Result()
	exists2, _ := redisClient.Exists(ctx, "test:key2").Result()
	exists3, _ := redisClient.Exists(ctx, "other:key1").Result()

	assert.Equal(t, int64(0), exists1)
	assert.Equal(t, int64(0), exists2)
	assert.Equal(t, int64(1), exists3) // Этот ключ не должен быть удален
}

func TestCacheService_DeleteByPattern_NoMatches(t *testing.T) {
	// Arrange
	redisClient := setupTestRedis(t)
	defer redisClient.Close()

	cacheService := services.NewCacheService(redisClient)
	ctx := context.Background()

	pattern := "nonexistent:*"

	// Act
	err := cacheService.DeleteByPattern(ctx, pattern)

	// Assert
	require.NoError(t, err)
}

// Тесты для CacheService.Exists

func TestCacheService_Exists_KeyExists(t *testing.T) {
	// Arrange
	redisClient := setupTestRedis(t)
	defer redisClient.Close()

	cacheService := services.NewCacheService(redisClient)
	ctx := context.Background()

	key := "test:key"
	redisClient.Set(ctx, key, "value", 5*time.Minute)

	// Act
	exists, err := cacheService.Exists(ctx, key)

	// Assert
	require.NoError(t, err)
	assert.True(t, exists)
}

func TestCacheService_Exists_KeyNotExists(t *testing.T) {
	// Arrange
	redisClient := setupTestRedis(t)
	defer redisClient.Close()

	cacheService := services.NewCacheService(redisClient)
	ctx := context.Background()

	key := "test:nonexistent"

	// Act
	exists, err := cacheService.Exists(ctx, key)

	// Assert
	require.NoError(t, err)
	assert.False(t, exists)
}

// Тесты для CacheService.SetTTL

func TestCacheService_SetTTL_Success(t *testing.T) {
	// Arrange
	redisClient := setupTestRedis(t)
	defer redisClient.Close()

	cacheService := services.NewCacheService(redisClient)
	ctx := context.Background()

	key := "test:key"
	redisClient.Set(ctx, key, "value", 5*time.Minute)
	newTTL := 10 * time.Minute

	// Act
	err := cacheService.SetTTL(ctx, key, newTTL)

	// Assert
	require.NoError(t, err)

	// Проверяем TTL
	ttl, _ := redisClient.TTL(ctx, key).Result()
	assert.InDelta(t, newTTL.Seconds(), ttl.Seconds(), 1.0) // Допускаем погрешность в 1 секунду
}

// Тесты для CacheService.GetTTL

func TestCacheService_GetTTL_Success(t *testing.T) {
	// Arrange
	redisClient := setupTestRedis(t)
	defer redisClient.Close()

	cacheService := services.NewCacheService(redisClient)
	ctx := context.Background()

	key := "test:key"
	expectedTTL := 5 * time.Minute
	redisClient.Set(ctx, key, "value", expectedTTL)

	// Act
	ttl, err := cacheService.GetTTL(ctx, key)

	// Assert
	require.NoError(t, err)
	assert.InDelta(t, expectedTTL.Seconds(), ttl.Seconds(), 1.0)
}

// Тесты для специализированных методов пользователей

func TestCacheService_SetUserCache_Success(t *testing.T) {
	// Arrange
	redisClient := setupTestRedis(t)
	defer redisClient.Close()

	cacheService := services.NewCacheService(redisClient)
	ctx := context.Background()

	userID := uuid.New().String()
	userData := map[string]string{"username": "testuser"}

	// Act
	err := cacheService.SetUserCache(ctx, userID, userData)

	// Assert
	require.NoError(t, err)

	// Проверяем, что ключ создан с правильным префиксом
	key := cacheService.UserCacheKey(userID)
	exists, _ := redisClient.Exists(ctx, key).Result()
	assert.Equal(t, int64(1), exists)
}

func TestCacheService_GetUserCache_Success(t *testing.T) {
	// Arrange
	redisClient := setupTestRedis(t)
	defer redisClient.Close()

	cacheService := services.NewCacheService(redisClient)
	ctx := context.Background()

	userID := uuid.New().String()
	userData := map[string]string{"username": "testuser"}

	// Сохраняем данные
	cacheService.SetUserCache(ctx, userID, userData)

	// Act
	var result map[string]string
	err := cacheService.GetUserCache(ctx, userID, &result)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, userData["username"], result["username"])
}

func TestCacheService_DeleteUserCache_Success(t *testing.T) {
	// Arrange
	redisClient := setupTestRedis(t)
	defer redisClient.Close()

	cacheService := services.NewCacheService(redisClient)
	ctx := context.Background()

	userID := uuid.New().String()
	userData := map[string]string{"username": "testuser"}

	// Сохраняем данные
	cacheService.SetUserCache(ctx, userID, userData)

	// Act
	err := cacheService.DeleteUserCache(ctx, userID)

	// Assert
	require.NoError(t, err)

	// Проверяем, что ключ удален
	var result map[string]string
	err = cacheService.GetUserCache(ctx, userID, &result)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "cache miss")
}

// Тесты для специализированных методов чатов

func TestCacheService_SetChatMessagesCache_Success(t *testing.T) {
	// Arrange
	redisClient := setupTestRedis(t)
	defer redisClient.Close()

	cacheService := services.NewCacheService(redisClient)
	ctx := context.Background()

	chatID := uuid.New().String()
	messages := []map[string]string{{"text": "hello"}}

	// Act
	err := cacheService.SetChatMessagesCache(ctx, chatID, messages)

	// Assert
	require.NoError(t, err)

	// Проверяем, что ключ создан
	key := cacheService.ChatMessagesCacheKey(chatID)
	exists, _ := redisClient.Exists(ctx, key).Result()
	assert.Equal(t, int64(1), exists)
}

func TestCacheService_GetChatMessagesCache_Success(t *testing.T) {
	// Arrange
	redisClient := setupTestRedis(t)
	defer redisClient.Close()

	cacheService := services.NewCacheService(redisClient)
	ctx := context.Background()

	chatID := uuid.New().String()
	messages := []map[string]string{{"text": "hello"}}

	// Сохраняем данные
	cacheService.SetChatMessagesCache(ctx, chatID, messages)

	// Act
	var result []map[string]string
	err := cacheService.GetChatMessagesCache(ctx, chatID, &result)

	// Assert
	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, messages[0]["text"], result[0]["text"])
}

func TestCacheService_DeleteChatMessagesCache_Success(t *testing.T) {
	// Arrange
	redisClient := setupTestRedis(t)
	defer redisClient.Close()

	cacheService := services.NewCacheService(redisClient)
	ctx := context.Background()

	chatID := uuid.New().String()
	messages := []map[string]string{{"text": "hello"}}

	// Сохраняем данные
	cacheService.SetChatMessagesCache(ctx, chatID, messages)

	// Act
	err := cacheService.DeleteChatMessagesCache(ctx, chatID)

	// Assert
	require.NoError(t, err)

	// Проверяем, что ключ удален
	var result []map[string]string
	err = cacheService.GetChatMessagesCache(ctx, chatID, &result)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "cache miss")
}

// Тесты для генераторов ключей

func TestCacheService_UserCacheKey(t *testing.T) {
	// Arrange
	redisClient := setupTestRedis(t)
	defer redisClient.Close()

	cacheService := services.NewCacheService(redisClient)

	userID := uuid.New().String()

	// Act
	key := cacheService.UserCacheKey(userID)

	// Assert
	assert.Contains(t, key, "user:")
	assert.Contains(t, key, userID)
}

func TestCacheService_ChatMessagesCacheKey(t *testing.T) {
	// Arrange
	redisClient := setupTestRedis(t)
	defer redisClient.Close()

	cacheService := services.NewCacheService(redisClient)

	chatID := uuid.New().String()

	// Act
	key := cacheService.ChatMessagesCacheKey(chatID)

	// Assert
	assert.Contains(t, key, "messages:")
	assert.Contains(t, key, chatID)
}

func TestCacheService_UserChatListCacheKey(t *testing.T) {
	// Arrange
	redisClient := setupTestRedis(t)
	defer redisClient.Close()

	cacheService := services.NewCacheService(redisClient)

	userID := uuid.New().String()

	// Act
	key := cacheService.UserChatListCacheKey(userID)

	// Assert
	assert.Contains(t, key, "chat_list:")
	assert.Contains(t, key, userID)
}
