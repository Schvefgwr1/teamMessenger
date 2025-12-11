//go:build integration
// +build integration

package integration

import (
	"apiService/internal/controllers"
	"apiService/internal/http_clients"
	"apiService/internal/services"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// TestUserController_GetUser_Integration_NotFound тестирует получение несуществующего пользователя
func TestUserController_GetUser_Integration_NotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange - создаем мок сервер, который возвращает 404
	userServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/v1/users/") {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer userServer.Close()

	redisClient := setupTestRedis(t)
	cacheService := services.NewCacheService(redisClient)
	sessionService := services.NewSessionService(redisClient)
	userClient := http_clients.NewUserClient(userServer.URL)
	fileClient := http_clients.NewFileClient("http://localhost:8080")
	userController := controllers.NewUserController(fileClient, userClient, cacheService, sessionService)

	userID := uuid.New()

	// Act
	user, err := userController.GetUser(userID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, user)
}

// TestUserController_GetUser_Integration_ServiceError тестирует получение пользователя при ошибке сервиса
func TestUserController_GetUser_Integration_ServiceError(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange - создаем мок сервер, который возвращает 500
	userServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/v1/users/") {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer userServer.Close()

	redisClient := setupTestRedis(t)
	cacheService := services.NewCacheService(redisClient)
	sessionService := services.NewSessionService(redisClient)
	userClient := http_clients.NewUserClient(userServer.URL)
	fileClient := http_clients.NewFileClient("http://localhost:8080")
	userController := controllers.NewUserController(fileClient, userClient, cacheService, sessionService)

	userID := uuid.New()

	// Act
	user, err := userController.GetUser(userID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, user)
}

// TestUserController_UpdateUser_Integration_ServiceError тестирует обновление пользователя при ошибке сервиса
func TestUserController_UpdateUser_Integration_ServiceError(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange - создаем мок сервер, который возвращает ошибку
	userServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/v1/users/") && r.Method == http.MethodPut {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error": "validation failed",
			})
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer userServer.Close()

	redisClient := setupTestRedis(t)
	cacheService := services.NewCacheService(redisClient)
	sessionService := services.NewSessionService(redisClient)
	userClient := http_clients.NewUserClient(userServer.URL)
	fileClient := http_clients.NewFileClient("http://localhost:8080")
	userController := controllers.NewUserController(fileClient, userClient, cacheService, sessionService)

	userID := uuid.New()
	updateReq := NewTestUpdateUserRequest()

	// Act
	response, err := userController.UpdateUser(userID, updateReq, nil)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, response)
}
