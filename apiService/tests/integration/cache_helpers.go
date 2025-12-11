//go:build integration
// +build integration

package integration

import (
	"apiService/internal/services"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// assertCacheExists проверяет, что ключ существует в кеше
func assertCacheExists(t *testing.T, cacheService *services.CacheService, key string) {
	t.Helper()
	ctx := context.Background()
	exists, err := cacheService.Exists(ctx, key)
	require.NoError(t, err)
	assert.True(t, exists, "Cache key %s should exist", key)
}

// assertCacheNotExists проверяет, что ключ не существует в кеше
func assertCacheNotExists(t *testing.T, cacheService *services.CacheService, key string) {
	t.Helper()
	ctx := context.Background()
	exists, err := cacheService.Exists(ctx, key)
	require.NoError(t, err)
	assert.False(t, exists, "Cache key %s should not exist", key)
}

// assertCacheValue проверяет значение в кеше
func assertCacheValue(t *testing.T, cacheService *services.CacheService, key string, expected interface{}) {
	t.Helper()
	ctx := context.Background()
	var actual interface{}
	err := cacheService.Get(ctx, key, &actual)
	require.NoError(t, err)
	assert.Equal(t, expected, actual, "Cache value for key %s should match", key)
}

// assertCacheTTL проверяет TTL ключа в кеше (с допуском ±1 секунда)
func assertCacheTTL(t *testing.T, cacheService *services.CacheService, key string, expectedTTL time.Duration) {
	t.Helper()
	ctx := context.Background()
	ttl, err := cacheService.GetTTL(ctx, key)
	require.NoError(t, err)

	// Проверяем, что TTL находится в пределах ожидаемого значения ±1 секунда
	minTTL := expectedTTL - time.Second
	maxTTL := expectedTTL + time.Second
	assert.True(t, ttl >= minTTL && ttl <= maxTTL,
		"Cache TTL for key %s should be between %v and %v, got %v", key, minTTL, maxTTL, ttl)
}

// getCacheTTL возвращает TTL ключа в кеше
func getCacheTTL(t *testing.T, cacheService *services.CacheService, key string) time.Duration {
	t.Helper()
	ctx := context.Background()
	ttl, err := cacheService.GetTTL(ctx, key)
	require.NoError(t, err)
	return ttl
}
