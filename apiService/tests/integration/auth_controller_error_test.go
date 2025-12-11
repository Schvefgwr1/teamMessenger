//go:build integration
// +build integration

package integration

import (
	"apiService/internal/controllers"
	"apiService/internal/http_clients"
	"apiService/internal/services"
	au "common/contracts/api-user"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestAuthController_Login_Integration_InvalidCredentials тестирует логин с неверными учетными данными
func TestAuthController_Login_Integration_InvalidCredentials(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange - создаем мок сервер, который возвращает ошибку
	userServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/auth/login" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error": "invalid credentials",
			})
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer userServer.Close()

	redisClient := setupTestRedis(t)
	sessionService := services.NewSessionService(redisClient)
	userClient := http_clients.NewUserClient(userServer.URL)
	fileClient := http_clients.NewFileClient("http://localhost:8080")
	authController := controllers.NewAuthController(fileClient, userClient, sessionService)

	loginData := &au.Login{
		Login:    "invalid_user",
		Password: "wrong_password",
	}

	// Act
	ctx := context.Background()
	token, userID, err := authController.Login(ctx, loginData)

	// Assert
	assert.Error(t, err)
	assert.Empty(t, token)
	assert.Equal(t, uuid.Nil, userID)
}

// TestAuthController_Login_Integration_ServiceUnavailable тестирует логин при недоступности сервиса
func TestAuthController_Login_Integration_ServiceUnavailable(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange - создаем мок сервер, который возвращает 500
	userServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/auth/login" {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer userServer.Close()

	redisClient := setupTestRedis(t)
	sessionService := services.NewSessionService(redisClient)
	userClient := http_clients.NewUserClient(userServer.URL)
	fileClient := http_clients.NewFileClient("http://localhost:8080")
	authController := controllers.NewAuthController(fileClient, userClient, sessionService)

	loginData := &au.Login{
		Login:    "testuser",
		Password: "password123",
	}

	// Act
	ctx := context.Background()
	token, userID, err := authController.Login(ctx, loginData)

	// Assert
	assert.Error(t, err)
	assert.Empty(t, token)
	assert.Equal(t, uuid.Nil, userID)
}

// TestAuthController_Logout_Integration_InvalidSession тестирует выход с несуществующей сессией
func TestAuthController_Logout_Integration_InvalidSession(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange
	redisClient := setupTestRedis(t)
	sessionService := services.NewSessionService(redisClient)
	_, _, _, _, userClient, _, _, fileClient, _ := setupTestHTTPClients(t)
	authController := controllers.NewAuthController(fileClient, userClient, sessionService)

	userID := uuid.New()
	invalidToken := "invalid_token_12345"

	// Act
	ctx := context.Background()
	err := authController.Logout(ctx, userID, invalidToken)

	// Assert - RevokeSession возвращает ошибку для несуществующей сессии
	// Это ожидаемое поведение, так как метод сначала получает сессию через GetSession
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "session not found", "Error should indicate session not found")
}

// TestAuthController_Register_Integration_ServiceError тестирует регистрацию при ошибке сервиса
func TestAuthController_Register_Integration_ServiceError(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange - создаем мок сервер, который возвращает ошибку
	userServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/auth/register" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error": "email already exists",
			})
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer userServer.Close()

	redisClient := setupTestRedis(t)
	sessionService := services.NewSessionService(redisClient)
	userClient := http_clients.NewUserClient(userServer.URL)
	fileClient := http_clients.NewFileClient("http://localhost:8080")
	authController := controllers.NewAuthController(fileClient, userClient, sessionService)

	registerReq := NewTestRegisterRequest()
	registerReq.Email = "existing@example.com"

	// Act
	response := authController.Register(registerReq, nil)

	// Assert
	require.NotNil(t, response)
	assert.NotNil(t, response.Error)
}
