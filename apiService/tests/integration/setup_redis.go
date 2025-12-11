//go:build integration
// +build integration

package integration

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

// setupTestRedis создает подключение к тестовому Redis
func setupTestRedis(t *testing.T) *redis.Client {
	redisHost := getEnvOrDefault("REDIS_HOST", "localhost")
	redisPort := getEnvOrDefault("REDIS_PORT", "6379")
	redisAddr := fmt.Sprintf("%s:%s", redisHost, redisPort)

	logStep(t, "Подключение к Redis addr=%s", redisAddr)

	client := redis.NewClient(&redis.Options{
		Addr: redisAddr,
		DB:   15, // Используем отдельную БД для тестов
	})

	// Проверяем подключение
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		t.Fatalf("Failed to connect to Redis: %v", err)
	}

	// Очищаем тестовую БД перед тестом
	logStep(t, "Очистка тестовых данных в Redis")
	client.FlushDB(ctx)

	t.Cleanup(func() {
		logStep(t, "Очистка после теста и закрытие соединения Redis")
		client.FlushDB(context.Background())
		client.Close()
	})

	return client
}
