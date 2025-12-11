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

// TestTaskController_GetTaskByID_Integration_NotFound тестирует получение несуществующей задачи
func TestTaskController_GetTaskByID_Integration_NotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange - создаем мок сервер, который возвращает 404
	taskServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/v1/tasks/") && r.Method == http.MethodGet {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer taskServer.Close()

	redisClient := setupTestRedis(t)
	cacheService := services.NewCacheService(redisClient)
	taskClient := http_clients.NewTaskClient(taskServer.URL)
	fileClient := http_clients.NewFileClient("http://localhost:8080")
	taskController := controllers.NewTaskController(taskClient, fileClient, cacheService)

	taskID := 99999

	// Act
	task, err := taskController.GetTaskByID(taskID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, task)
}

// TestTaskController_CreateTask_Integration_ServiceError тестирует создание задачи при ошибке сервиса
func TestTaskController_CreateTask_Integration_ServiceError(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange - создаем мок сервер, который возвращает ошибку
	taskServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/tasks" && r.Method == http.MethodPost {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error": "invalid task data",
			})
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer taskServer.Close()

	redisClient := setupTestRedis(t)
	cacheService := services.NewCacheService(redisClient)
	taskClient := http_clients.NewTaskClient(taskServer.URL)
	fileClient := http_clients.NewFileClient("http://localhost:8080")
	taskController := controllers.NewTaskController(taskClient, fileClient, cacheService)

	createReq := NewTestCreateTaskRequest()

	// Act
	response, err := taskController.CreateTask(createReq, uuid.New())

	// Assert
	assert.Error(t, err)
	assert.Nil(t, response)
}

// TestTaskController_GetAllStatuses_Integration_ServiceError тестирует получение статусов при ошибке сервиса
func TestTaskController_GetAllStatuses_Integration_ServiceError(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange - создаем мок сервер, который возвращает 500
	taskServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/tasks/statuses" && r.Method == http.MethodGet {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer taskServer.Close()

	redisClient := setupTestRedis(t)
	cacheService := services.NewCacheService(redisClient)
	taskClient := http_clients.NewTaskClient(taskServer.URL)
	fileClient := http_clients.NewFileClient("http://localhost:8080")
	taskController := controllers.NewTaskController(taskClient, fileClient, cacheService)

	// Act
	statuses, err := taskController.GetAllStatuses()

	// Assert
	assert.Error(t, err)
	assert.Nil(t, statuses)
}
