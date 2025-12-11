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

// TestChatController_GetUserChats_Integration_ServiceError тестирует получение чатов при ошибке сервиса
func TestChatController_GetUserChats_Integration_ServiceError(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange - создаем мок сервер, который возвращает 500
	chatServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/v1/chats/user/") {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer chatServer.Close()

	redisClient := setupTestRedis(t)
	cacheService := services.NewCacheService(redisClient)
	chatClient := http_clients.NewChatClient(chatServer.URL)
	fileClient := http_clients.NewFileClient("http://localhost:8080")
	chatController := controllers.NewChatController(chatClient, fileClient, cacheService)

	userID := uuid.New()

	// Act
	chats, err := chatController.GetUserChats(userID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, chats)
}

// TestChatController_CreateChat_Integration_ServiceError тестирует создание чата при ошибке сервиса
func TestChatController_CreateChat_Integration_ServiceError(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange - создаем мок сервер, который возвращает ошибку
	chatServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/chats" && r.Method == http.MethodPost {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error": "invalid chat data",
			})
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer chatServer.Close()

	redisClient := setupTestRedis(t)
	cacheService := services.NewCacheService(redisClient)
	chatClient := http_clients.NewChatClient(chatServer.URL)
	fileClient := http_clients.NewFileClient("http://localhost:8080")
	chatController := controllers.NewChatController(chatClient, fileClient, cacheService)

	createReq := NewTestCreateChatRequest()
	ownerID := uuid.New()
	userIDs := []uuid.UUID{uuid.New(), uuid.New()}

	// Act
	response, err := chatController.CreateChat(createReq, ownerID, userIDs)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, response)
}

// TestChatController_GetChatMessages_Integration_NotFound тестирует получение сообщений несуществующего чата
func TestChatController_GetChatMessages_Integration_NotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange - создаем мок сервер, который возвращает 404
	chatServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/v1/chats/messages/") {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer chatServer.Close()

	redisClient := setupTestRedis(t)
	cacheService := services.NewCacheService(redisClient)
	chatClient := http_clients.NewChatClient(chatServer.URL)
	fileClient := http_clients.NewFileClient("http://localhost:8080")
	chatController := controllers.NewChatController(chatClient, fileClient, cacheService)

	chatID := uuid.New()
	userID := uuid.New()

	// Act
	messages, err := chatController.GetChatMessages(chatID, userID, 0, 10)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, messages)
}
