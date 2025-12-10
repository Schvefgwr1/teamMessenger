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

// Тесты для SessionService.CreateSession

func TestSessionService_CreateSession_Success(t *testing.T) {
	// Arrange
	redisClient := setupTestRedis(t)
	defer redisClient.Close()

	sessionService := services.NewSessionService(redisClient)
	ctx := context.Background()

	userID := uuid.New()
	token := "test-token"
	expiresAt := time.Now().Add(24 * time.Hour)

	// Act
	err := sessionService.CreateSession(ctx, userID, token, expiresAt)

	// Assert
	require.NoError(t, err)

	// Проверяем, что сессия создана
	session, err := sessionService.GetSession(ctx, userID, token)
	require.NoError(t, err)
	assert.Equal(t, userID, session.UserID)
	assert.Equal(t, services.SessionActive, session.Status)
}

func TestSessionService_CreateSession_ExpiredTTL(t *testing.T) {
	// Arrange
	redisClient := setupTestRedis(t)
	defer redisClient.Close()

	sessionService := services.NewSessionService(redisClient)
	ctx := context.Background()

	userID := uuid.New()
	token := "test-token"
	expiresAt := time.Now().Add(-1 * time.Hour) // Прошедшее время

	// Act
	err := sessionService.CreateSession(ctx, userID, token, expiresAt)

	// Assert
	require.NoError(t, err)

	// Проверяем, что сессия создана с fallback TTL
	session, err := sessionService.GetSession(ctx, userID, token)
	require.NoError(t, err)
	assert.Equal(t, userID, session.UserID)
}

// Тесты для SessionService.GetSession

func TestSessionService_GetSession_Success(t *testing.T) {
	// Arrange
	redisClient := setupTestRedis(t)
	defer redisClient.Close()

	sessionService := services.NewSessionService(redisClient)
	ctx := context.Background()

	userID := uuid.New()
	token := "test-token"
	expiresAt := time.Now().Add(24 * time.Hour)

	// Создаем сессию
	sessionService.CreateSession(ctx, userID, token, expiresAt)

	// Act
	session, err := sessionService.GetSession(ctx, userID, token)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, session)
	assert.Equal(t, userID, session.UserID)
	assert.Equal(t, services.SessionActive, session.Status)
}

func TestSessionService_GetSession_NotFound(t *testing.T) {
	// Arrange
	redisClient := setupTestRedis(t)
	defer redisClient.Close()

	sessionService := services.NewSessionService(redisClient)
	ctx := context.Background()

	userID := uuid.New()
	token := "nonexistent-token"

	// Act
	session, err := sessionService.GetSession(ctx, userID, token)

	// Assert
	require.Error(t, err)
	assert.Nil(t, session)
	assert.Contains(t, err.Error(), "session not found")
}

func TestSessionService_GetSession_InvalidData(t *testing.T) {
	// Arrange
	redisClient := setupTestRedis(t)
	defer redisClient.Close()

	sessionService := services.NewSessionService(redisClient)
	ctx := context.Background()

	userID := uuid.New()
	token := "test-token"

	// Создаем сессию с валидными данными
	expiresAt := time.Now().Add(24 * time.Hour)
	sessionService.CreateSession(ctx, userID, token, expiresAt)

	// Повреждаем данные в Redis напрямую через Scan
	pattern := "session:" + userID.String() + ":*"
	iter := redisClient.Scan(ctx, 0, pattern, 0).Iterator()
	var key string
	if iter.Next(ctx) {
		key = iter.Val()
	}
	// Проверяем ошибку итератора
	if iter.Err() != nil {
		t.Fatalf("Failed to scan: %v", iter.Err())
	}

	if key != "" {
		// Повреждаем данные
		redisClient.Set(ctx, key, "invalid json", 24*time.Hour)

		// Act
		session, err := sessionService.GetSession(ctx, userID, token)

		// Assert
		require.Error(t, err)
		assert.Nil(t, session)
		// Ошибка должна содержать "unmarshal" или "failed to unmarshal"
		assert.True(t, err != nil)
	}
}

// Тесты для SessionService.RevokeSession

func TestSessionService_RevokeSession_Success(t *testing.T) {
	// Arrange
	redisClient := setupTestRedis(t)
	defer redisClient.Close()

	sessionService := services.NewSessionService(redisClient)
	ctx := context.Background()

	userID := uuid.New()
	token := "test-token"
	expiresAt := time.Now().Add(24 * time.Hour)

	// Создаем сессию
	sessionService.CreateSession(ctx, userID, token, expiresAt)

	// Act
	err := sessionService.RevokeSession(ctx, userID, token)

	// Assert
	require.NoError(t, err)

	// Проверяем, что сессия отозвана
	session, err := sessionService.GetSession(ctx, userID, token)
	require.NoError(t, err)
	assert.Equal(t, services.SessionRevoked, session.Status)
}

func TestSessionService_RevokeSession_NotFound(t *testing.T) {
	// Arrange
	redisClient := setupTestRedis(t)
	defer redisClient.Close()

	sessionService := services.NewSessionService(redisClient)
	ctx := context.Background()

	userID := uuid.New()
	token := "nonexistent-token"

	// Act
	err := sessionService.RevokeSession(ctx, userID, token)

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "session not found")
}

// Тесты для SessionService.RevokeAllUserSessions

func TestSessionService_RevokeAllUserSessions_Success(t *testing.T) {
	// Arrange
	redisClient := setupTestRedis(t)
	defer redisClient.Close()

	sessionService := services.NewSessionService(redisClient)
	ctx := context.Background()

	userID := uuid.New()
	token1 := "token1"
	token2 := "token2"
	expiresAt := time.Now().Add(24 * time.Hour)

	// Создаем несколько сессий
	sessionService.CreateSession(ctx, userID, token1, expiresAt)
	sessionService.CreateSession(ctx, userID, token2, expiresAt)

	// Act
	err := sessionService.RevokeAllUserSessions(ctx, userID)

	// Assert
	require.NoError(t, err)

	// Проверяем, что все сессии отозваны
	session1, _ := sessionService.GetSession(ctx, userID, token1)
	session2, _ := sessionService.GetSession(ctx, userID, token2)

	assert.Equal(t, services.SessionRevoked, session1.Status)
	assert.Equal(t, services.SessionRevoked, session2.Status)
}

func TestSessionService_RevokeAllUserSessions_NoSessions(t *testing.T) {
	// Arrange
	redisClient := setupTestRedis(t)
	defer redisClient.Close()

	sessionService := services.NewSessionService(redisClient)
	ctx := context.Background()

	userID := uuid.New()

	// Act
	err := sessionService.RevokeAllUserSessions(ctx, userID)

	// Assert
	require.NoError(t, err)
}

// Тесты для SessionService.IsSessionValid

func TestSessionService_IsSessionValid_ActiveSession(t *testing.T) {
	// Arrange
	redisClient := setupTestRedis(t)
	defer redisClient.Close()

	sessionService := services.NewSessionService(redisClient)
	ctx := context.Background()

	userID := uuid.New()
	token := "test-token"
	expiresAt := time.Now().Add(24 * time.Hour)

	// Создаем активную сессию
	sessionService.CreateSession(ctx, userID, token, expiresAt)

	// Act
	isValid, err := sessionService.IsSessionValid(ctx, userID, token)

	// Assert
	require.NoError(t, err)
	assert.True(t, isValid)
}

func TestSessionService_IsSessionValid_RevokedSession(t *testing.T) {
	// Arrange
	redisClient := setupTestRedis(t)
	defer redisClient.Close()

	sessionService := services.NewSessionService(redisClient)
	ctx := context.Background()

	userID := uuid.New()
	token := "test-token"
	expiresAt := time.Now().Add(24 * time.Hour)

	// Создаем и отзываем сессию
	sessionService.CreateSession(ctx, userID, token, expiresAt)
	sessionService.RevokeSession(ctx, userID, token)

	// Act
	isValid, err := sessionService.IsSessionValid(ctx, userID, token)

	// Assert
	require.NoError(t, err)
	assert.False(t, isValid)
}

func TestSessionService_IsSessionValid_ExpiredSession(t *testing.T) {
	// Arrange
	redisClient := setupTestRedis(t)
	defer redisClient.Close()

	sessionService := services.NewSessionService(redisClient)
	ctx := context.Background()

	userID := uuid.New()
	token := "test-token"
	expiresAt := time.Now().Add(-1 * time.Hour) // Прошедшее время

	// Создаем истеченную сессию
	sessionService.CreateSession(ctx, userID, token, expiresAt)

	// Act
	isValid, err := sessionService.IsSessionValid(ctx, userID, token)

	// Assert
	require.NoError(t, err)
	assert.False(t, isValid)

	// Проверяем, что статус обновлен на expired
	session, _ := sessionService.GetSession(ctx, userID, token)
	assert.Equal(t, services.SessionExpired, session.Status)
}

func TestSessionService_IsSessionValid_NotFound(t *testing.T) {
	// Arrange
	redisClient := setupTestRedis(t)
	defer redisClient.Close()

	sessionService := services.NewSessionService(redisClient)
	ctx := context.Background()

	userID := uuid.New()
	token := "nonexistent-token"

	// Act
	isValid, err := sessionService.IsSessionValid(ctx, userID, token)

	// Assert
	require.Error(t, err)
	assert.False(t, isValid)
}
