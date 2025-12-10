package services

import (
	"apiService/internal/services"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Тесты для edge cases CacheService

func TestCacheService_Get_RedisError(t *testing.T) {
	// Arrange
	// Создаем Redis клиент и закрываем его для создания ошибки
	redisClient := setupTestRedis(t)
	redisClient.Close() // Закрываем для создания ошибки

	cacheService := services.NewCacheService(redisClient)
	ctx := context.Background()

	key := "test:key"

	// Act
	var result map[string]string
	err := cacheService.Get(ctx, key, &result)

	// Assert
	require.Error(t, err)
	// Должна быть ошибка Redis, не cache miss
	assert.Contains(t, err.Error(), "failed to get cache data")
}

func TestCacheService_DeleteByPattern_IteratorError(t *testing.T) {
	// Arrange
	// Создаем Redis клиент и закрываем его для создания ошибки
	redisClient := setupTestRedis(t)
	redisClient.Close() // Закрываем для создания ошибки

	cacheService := services.NewCacheService(redisClient)
	ctx := context.Background()

	pattern := "test:*"

	// Act
	err := cacheService.DeleteByPattern(ctx, pattern)

	// Assert
	require.Error(t, err)
}
