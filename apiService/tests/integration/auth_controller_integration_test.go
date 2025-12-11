//go:build integration
// +build integration

package integration

import (
	"apiService/internal/controllers"
	"apiService/internal/services"
	au "common/contracts/api-user"
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestAuthController_Register_Integration тестирует регистрацию пользователя с реальными интеграциями
func TestAuthController_Register_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange - настройка реальных зависимостей
	redisClient := setupTestRedis(t)
	_, _, _, _, userClient, _, _, fileClient, _ := setupTestHTTPClients(t)

	sessionService := services.NewSessionService(redisClient)
	authController := controllers.NewAuthController(fileClient, userClient, sessionService)

	registerReq := NewTestRegisterRequest()

	// Act - выполнение реального сценария
	response := authController.Register(registerReq, nil)

	// Assert - проверка реального результата
	require.NotNil(t, response)
	assert.Nil(t, response.Error)
	assert.NotNil(t, response.User)
}

// TestAuthController_Register_Integration_WithFile тестирует регистрацию с загрузкой аватара
func TestAuthController_Register_Integration_WithFile(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange
	redisClient := setupTestRedis(t)
	_, _, _, _, userClient, _, _, fileClient, _ := setupTestHTTPClients(t)

	sessionService := services.NewSessionService(redisClient)
	authController := controllers.NewAuthController(fileClient, userClient, sessionService)

	registerReq := NewTestRegisterRequest()

	// Создаем тестовый файл
	file := createTestFileHeader(t, "avatar.jpg", "image/jpeg", []byte("test image content"))

	// Act
	response := authController.Register(registerReq, file)

	// Assert
	require.NotNil(t, response)
	assert.Nil(t, response.Error)
	assert.NotNil(t, response.User)
}

// TestAuthController_Login_Integration тестирует вход в систему с реальными интеграциями
func TestAuthController_Login_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange - настройка реальных зависимостей
	redisClient := setupTestRedis(t)
	_, _, _, _, userClient, _, _, fileClient, _ := setupTestHTTPClients(t)

	sessionService := services.NewSessionService(redisClient)
	authController := controllers.NewAuthController(fileClient, userClient, sessionService)

	ctx := context.Background()
	loginData := &au.Login{
		Login:    "testuser",
		Password: "password123",
	}

	// Act - выполнение реального сценария
	token, userID, err := authController.Login(ctx, loginData)

	// Assert - проверка реального результата
	require.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.NotEqual(t, uuid.Nil, userID)

	// Проверяем, что сессия реально создана в Redis
	session, err := sessionService.GetSession(ctx, userID, token)
	require.NoError(t, err)
	assert.Equal(t, services.SessionActive, session.Status)
}

// TestAuthController_Logout_Integration тестирует выход из системы с реальными интеграциями
func TestAuthController_Logout_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange
	redisClient := setupTestRedis(t)
	_, _, _, _, userClient, _, _, fileClient, _ := setupTestHTTPClients(t)

	sessionService := services.NewSessionService(redisClient)
	authController := controllers.NewAuthController(fileClient, userClient, sessionService)

	ctx := context.Background()
	userID := uuid.New()
	token := "test_token_" + uuid.New().String()

	// Создаем сессию перед выходом
	expiresAt := time.Now().Add(24 * time.Hour)
	err := sessionService.CreateSession(ctx, userID, token, expiresAt)
	require.NoError(t, err)

	// Проверяем, что сессия существует
	session, err := sessionService.GetSession(ctx, userID, token)
	require.NoError(t, err)
	assert.Equal(t, services.SessionActive, session.Status)

	// Act - выполнение реального сценария
	err = authController.Logout(ctx, userID, token)

	// Assert - проверка реального результата
	require.NoError(t, err)

	// Проверяем, что сессия реально отозвана в Redis
	session, err = sessionService.GetSession(ctx, userID, token)
	require.NoError(t, err)
	assert.Equal(t, services.SessionRevoked, session.Status)
}

// TestAuthController_Login_Integration_RevokeOldSessions тестирует отзыв старых сессий при логине
func TestAuthController_Login_Integration_RevokeOldSessions(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange
	redisClient := setupTestRedis(t)
	_, _, _, _, userClient, _, _, fileClient, _ := setupTestHTTPClients(t)

	sessionService := services.NewSessionService(redisClient)
	authController := controllers.NewAuthController(fileClient, userClient, sessionService)

	ctx := context.Background()

	// Сначала логинимся, чтобы получить userID
	loginData := &au.Login{
		Login:    "testuser",
		Password: "password123",
	}
	_, firstUserID, err := authController.Login(ctx, loginData)
	require.NoError(t, err)

	// Создаем старую сессию для этого же userID
	oldToken := "old_token_" + uuid.New().String()
	expiresAt := time.Now().Add(24 * time.Hour)
	err = sessionService.CreateSession(ctx, firstUserID, oldToken, expiresAt)
	require.NoError(t, err)

	// Act - логин должен отозвать старую сессию
	token, secondUserID, err := authController.Login(ctx, loginData)

	// Assert
	require.NoError(t, err)
	assert.NotEmpty(t, token)
	// Проверяем, что userID совпадает (тестовый сервер возвращает фиксированный userID)
	assert.Equal(t, firstUserID, secondUserID)

	// Проверяем, что старая сессия отозвана
	oldSession, err := sessionService.GetSession(ctx, firstUserID, oldToken)
	require.NoError(t, err)
	assert.Equal(t, services.SessionRevoked, oldSession.Status)
}
