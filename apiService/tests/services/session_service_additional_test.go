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

// Тесты для улучшения покрытия SessionService

func TestSessionService_CreateSession_NegativeTTL(t *testing.T) {
	// Arrange
	redisClient := setupTestRedis(t)
	defer redisClient.Close()

	sessionService := services.NewSessionService(redisClient)
	ctx := context.Background()

	userID := uuid.New()
	token := "test-token"
	// Устанавливаем expiresAt в прошлом, чтобы TTL был отрицательным
	expiresAt := time.Now().Add(-time.Hour)

	// Act
	err := sessionService.CreateSession(ctx, userID, token, expiresAt)

	// Assert
	require.NoError(t, err)
	// Должен использоваться fallback TTL (24 часа)
	session, err := sessionService.GetSession(ctx, userID, token)
	require.NoError(t, err)
	assert.Equal(t, userID, session.UserID)
}

func TestSessionService_RevokeAllUserSessions_WithErrors(t *testing.T) {
	// Arrange
	redisClient := setupTestRedis(t)
	defer redisClient.Close()

	sessionService := services.NewSessionService(redisClient)
	ctx := context.Background()

	userID := uuid.New()
	token1 := "token-1"
	token2 := "token-2"
	expiresAt := time.Now().Add(time.Hour)

	// Создаем две сессии
	err := sessionService.CreateSession(ctx, userID, token1, expiresAt)
	require.NoError(t, err)
	err = sessionService.CreateSession(ctx, userID, token2, expiresAt)
	require.NoError(t, err)

	// Создаем невалидный ключ сессии для этого пользователя
	key := "session:" + userID.String() + ":invalid-hash"
	redisClient.Set(ctx, key, "invalid json", time.Hour)

	// Act
	err = sessionService.RevokeAllUserSessions(ctx, userID)

	// Assert
	// Метод должен обработать ошибки и продолжить работу
	// Проверяем, что валидные сессии отозваны
	isValid1, _ := sessionService.IsSessionValid(ctx, userID, token1)
	isValid2, _ := sessionService.IsSessionValid(ctx, userID, token2)
	assert.False(t, isValid1)
	assert.False(t, isValid2)
}

func TestSessionService_RevokeAllUserSessions_MarshalError(t *testing.T) {
	// Arrange
	redisClient := setupTestRedis(t)
	defer redisClient.Close()

	sessionService := services.NewSessionService(redisClient)
	ctx := context.Background()

	userID := uuid.New()
	token := "test-token"
	expiresAt := time.Now().Add(time.Hour)

	// Создаем сессию
	err := sessionService.CreateSession(ctx, userID, token, expiresAt)
	require.NoError(t, err)

	// Создаем невалидный ключ для теста
	key := "session:" + userID.String() + ":invalid-hash"

	// Сохраняем невалидные данные, которые нельзя распарсить
	redisClient.Set(ctx, key, "invalid json", time.Hour)

	// Act
	err = sessionService.RevokeAllUserSessions(ctx, userID)

	// Assert
	// Метод должен обработать ошибку unmarshal и продолжить
	// Не должно быть ошибки, так как ошибки обрабатываются внутри цикла
	require.NoError(t, err)
}

func TestSessionService_RevokeAllUserSessions_RedisGetError(t *testing.T) {
	// Arrange
	redisClient := setupTestRedis(t)
	defer redisClient.Close()

	sessionService := services.NewSessionService(redisClient)
	ctx := context.Background()

	userID := uuid.New()
	token := "test-token"
	expiresAt := time.Now().Add(time.Hour)

	// Создаем сессию
	err := sessionService.CreateSession(ctx, userID, token, expiresAt)
	require.NoError(t, err)

	// Закрываем Redis для создания ошибки
	redisClient.Close()

	// Act
	err = sessionService.RevokeAllUserSessions(ctx, userID)

	// Assert
	// Должна быть ошибка при сканировании
	require.Error(t, err)
}

func TestSessionService_IsSessionValid_RevokedStatus(t *testing.T) {
	// Arrange
	redisClient := setupTestRedis(t)
	defer redisClient.Close()

	sessionService := services.NewSessionService(redisClient)
	ctx := context.Background()

	userID := uuid.New()
	token := "test-token"
	expiresAt := time.Now().Add(time.Hour)

	// Создаем и отзываем сессию
	err := sessionService.CreateSession(ctx, userID, token, expiresAt)
	require.NoError(t, err)

	err = sessionService.RevokeSession(ctx, userID, token)
	require.NoError(t, err)

	// Act
	isValid, err := sessionService.IsSessionValid(ctx, userID, token)

	// Assert
	require.NoError(t, err)
	assert.False(t, isValid)
}

func TestSessionService_IsSessionValid_ExpiredStatus(t *testing.T) {
	// Arrange
	redisClient := setupTestRedis(t)
	defer redisClient.Close()

	sessionService := services.NewSessionService(redisClient)
	ctx := context.Background()

	userID := uuid.New()
	token := "test-token"
	// Создаем сессию, которая уже истекла
	expiresAt := time.Now().Add(-time.Hour)

	// Создаем сессию с истекшим временем
	err := sessionService.CreateSession(ctx, userID, token, expiresAt)
	require.NoError(t, err)

	// Act - проверяем валидность истекшей сессии
	isValid, err := sessionService.IsSessionValid(ctx, userID, token)

	// Assert
	require.NoError(t, err)
	assert.False(t, isValid)

	// Проверяем, что статус обновлен на expired
	session, err := sessionService.GetSession(ctx, userID, token)
	require.NoError(t, err)
	assert.Equal(t, services.SessionExpired, session.Status)
}
