//go:build integration
// +build integration

package integration

import (
	"apiService/internal/http_clients"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/google/uuid"
)

var (
	// Фиксированный userID для тестов логина (для проверки отзыва сессий)
	testLoginUserID     uuid.UUID
	testLoginUserIDOnce sync.Once
)

// setupUserServiceMock создает тестовый HTTP сервер для User Service
func setupUserServiceMock(t *testing.T) *httptest.Server {
	t.Helper()

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Обработка различных эндпоинтов User Service
		if r.URL.Path == "/api/v1/auth/register" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"id":    uuid.New().String(),
				"email": "test@example.com",
			})
			return
		}
		if r.URL.Path == "/api/v1/auth/login" {
			// Логин - используем фиксированный userID для тестов отзыва сессий
			testLoginUserIDOnce.Do(func() {
				testLoginUserID = uuid.New()
			})
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"token":  "test_token_" + uuid.New().String(),
				"userID": testLoginUserID.String(),
			})
			return
		}
		// Обработка /api/v1/users/search (GET) - проверяем раньше, чтобы не конфликтовать с /api/v1/users/
		if r.URL.Path == "/api/v1/users/search" && r.Method == http.MethodGet {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"users": []map[string]interface{}{
					{"id": uuid.New().String(), "username": "test_user"},
				},
			})
			return
		}
		if strings.HasPrefix(r.URL.Path, "/api/v1/users/") {
			path := r.URL.Path[len("/api/v1/users/"):]

			// Обработка /api/v1/users/{id}/role (PATCH)
			if strings.HasSuffix(path, "/role") && r.Method == http.MethodPatch {
				w.WriteHeader(http.StatusOK)
				return
			}

			// Обработка /api/v1/users/{id}/brief (GET)
			if strings.Contains(path, "/brief") && r.Method == http.MethodGet {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(map[string]interface{}{
					"userID":   uuid.New().String(),
					"username": "test_user",
					"role":     map[string]interface{}{"name": "main"},
				})
				return
			}

			// Обработка /api/v1/users/{id} (GET, PUT)
			userID := path
			if idx := strings.Index(userID, "?"); idx != -1 {
				userID = userID[:idx]
			}
			if idx := strings.Index(userID, "/"); idx != -1 {
				userID = userID[:idx]
			}

			parsedUUID, err := uuid.Parse(userID)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			if r.Method == http.MethodPut {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(map[string]interface{}{
					"user": map[string]interface{}{
						"id":       parsedUUID.String(),
						"username": "updated_user",
						"email":    "test@example.com",
					},
					"error": nil,
				})
				return
			}

			// Возвращаем тестового пользователя в формате GetUserResponse
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"file": nil,
				"user": map[string]interface{}{
					"id":       parsedUUID.String(),
					"username": "test_user_" + userID[:8],
					"email":    "test_" + userID[:8] + "@example.com",
					"role": map[string]interface{}{
						"id":   1,
						"name": "user",
					},
				},
			})
			return
		}
		if r.URL.Path == "/api/v1/permissions/" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode([]map[string]interface{}{
				{"id": 1, "name": "read"},
				{"id": 2, "name": "write"},
			})
			return
		}
		if strings.HasPrefix(r.URL.Path, "/api/v1/roles/") {
			path := r.URL.Path[len("/api/v1/roles/"):]

			if strings.Contains(path, "/permissions") && r.Method == http.MethodPatch {
				w.WriteHeader(http.StatusOK)
				return
			}

			if r.Method == http.MethodDelete {
				w.WriteHeader(http.StatusOK)
				return
			}
		}
		if r.URL.Path == "/api/v1/roles/" {
			if r.Method == http.MethodPost {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)
				json.NewEncoder(w).Encode(map[string]interface{}{
					"id":   1,
					"name": "test_role",
				})
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode([]map[string]interface{}{
				{"id": 1, "name": "user"},
				{"id": 2, "name": "admin"},
			})
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
}

// setupChatServiceMock создает тестовый HTTP сервер для Chat Service
func setupChatServiceMock(t *testing.T) *httptest.Server {
	t.Helper()

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/v1/chats/user/") {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode([]map[string]interface{}{
				{"id": uuid.New().String(), "name": "test_chat"},
			})
			return
		}
		if r.URL.Path == "/api/v1/chats" && r.Method == http.MethodPost {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"chatID": uuid.New().String(),
			})
			return
		}
		if strings.HasPrefix(r.URL.Path, "/api/v1/chats/messages/") {
			if r.Method == http.MethodPost {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)
				json.NewEncoder(w).Encode(map[string]interface{}{
					"id":      uuid.New().String(),
					"content": "test message",
				})
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode([]map[string]interface{}{
				{"id": uuid.New().String(), "content": "test message"},
			})
			return
		}
		if strings.HasPrefix(r.URL.Path, "/api/v1/chats/search/") {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"messages": []map[string]interface{}{
					{"id": uuid.New().String(), "content": "test"},
				},
				"total": 1,
			})
			return
		}
		if strings.HasPrefix(r.URL.Path, "/api/v1/chats/") && r.Method == http.MethodPut {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"chat": map[string]interface{}{
					"id":   uuid.New().String(),
					"name": "updated_chat",
				},
			})
			return
		}
		if strings.HasPrefix(r.URL.Path, "/api/v1/chats/") && r.Method == http.MethodDelete {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		if strings.Contains(r.URL.Path, "/ban/") {
			w.WriteHeader(http.StatusOK)
			return
		}
		if strings.Contains(r.URL.Path, "/roles/change") {
			w.WriteHeader(http.StatusOK)
			return
		}
		if strings.Contains(r.URL.Path, "/me/role") {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"roleId":   1,
				"roleName": "main",
				"permissions": []map[string]interface{}{
					{"id": 1, "name": "read"},
				},
			})
			return
		}
		if strings.Contains(r.URL.Path, "/members") {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode([]map[string]interface{}{
				{"userID": uuid.New().String(), "role": map[string]interface{}{"name": "main"}},
			})
			return
		}
		if strings.HasPrefix(r.URL.Path, "/api/v1/chat-roles") {
			path := r.URL.Path[len("/api/v1/chat-roles"):]

			if path != "" && path != "/" && r.Method == http.MethodGet {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(map[string]interface{}{
					"id":          1,
					"name":        "main",
					"permissions": []map[string]interface{}{},
				})
				return
			}

			if r.Method == http.MethodGet {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode([]map[string]interface{}{
					{"id": 1, "name": "main", "permissions": []map[string]interface{}{}},
				})
				return
			}
			if r.Method == http.MethodPost {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)
				json.NewEncoder(w).Encode(map[string]interface{}{
					"id":   1,
					"name": "test_role",
				})
				return
			}
			if r.Method == http.MethodDelete {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			if r.Method == http.MethodPatch {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(map[string]interface{}{
					"id":   1,
					"name": "updated_role",
				})
				return
			}
		}
		if strings.HasPrefix(r.URL.Path, "/api/v1/chat-permissions") {
			if r.Method == http.MethodGet {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode([]map[string]interface{}{
					{"id": 1, "name": "read"},
				})
				return
			}
			if r.Method == http.MethodPost {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)
				json.NewEncoder(w).Encode(map[string]interface{}{
					"id":   1,
					"name": "test_permission",
				})
				return
			}
			if r.Method == http.MethodDelete {
				w.WriteHeader(http.StatusNoContent)
				return
			}
		}
		w.WriteHeader(http.StatusNotFound)
	}))
}

// setupTaskServiceMock создает тестовый HTTP сервер для Task Service
func setupTaskServiceMock(t *testing.T) *httptest.Server {
	t.Helper()

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/tasks" && r.Method == http.MethodPost {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"id":    1,
				"title": "test_task",
			})
			return
		}
		if strings.Contains(r.URL.Path, "/status/") && r.Method == http.MethodPatch {
			w.WriteHeader(http.StatusOK)
			return
		}
		if strings.HasPrefix(r.URL.Path, "/api/v1/tasks/statuses/") {
			path := r.URL.Path[len("/api/v1/tasks/statuses/"):]
			if path != "" && r.Method == http.MethodGet {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(map[string]interface{}{
					"id":   1,
					"name": "created",
				})
				return
			}
			if path != "" && r.Method == http.MethodDelete {
				w.WriteHeader(http.StatusNoContent)
				return
			}
		}
		if r.URL.Path == "/api/v1/tasks/statuses" {
			if r.Method == http.MethodGet {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode([]map[string]interface{}{
					{"id": 1, "name": "created"},
				})
				return
			}
			if r.Method == http.MethodPost {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)
				json.NewEncoder(w).Encode(map[string]interface{}{
					"id":   1,
					"name": "test_status",
				})
				return
			}
		}
		if strings.HasPrefix(r.URL.Path, "/api/v1/tasks/") && r.Method == http.MethodGet {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"task": map[string]interface{}{
					"id":    1,
					"title": "test_task",
				},
			})
			return
		}
		if strings.Contains(r.URL.Path, "/tasks") && r.Method == http.MethodGet {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode([]map[string]interface{}{
				{"id": 1, "title": "test_task"},
			})
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
}

// setupFileServiceMock создает тестовый HTTP сервер для File Service
func setupFileServiceMock(t *testing.T) *httptest.Server {
	t.Helper()

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/files/upload" && r.Method == http.MethodPost {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"id":   1,
				"name": "test_file.txt",
				"url":  "http://example.com/files/1",
			})
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
}

// setupTestHTTPClients создает тестовые HTTP серверы для внешних сервисов
func setupTestHTTPClients(t *testing.T) (
	userServer *httptest.Server,
	chatServer *httptest.Server,
	taskServer *httptest.Server,
	fileServer *httptest.Server,
	userClient http_clients.UserClient,
	chatClient http_clients.ChatClient,
	taskClient http_clients.TaskClient,
	fileClient http_clients.FileClient,
	rolePermissionClient http_clients.ChatRolePermissionClient,
) {
	t.Helper()

	userServer = setupUserServiceMock(t)
	chatServer = setupChatServiceMock(t)
	taskServer = setupTaskServiceMock(t)
	fileServer = setupFileServiceMock(t)

	// Устанавливаем переменные окружения для HTTP клиентов
	os.Setenv("USER_SERVICE_URL", userServer.URL)
	os.Setenv("CHAT_SERVICE_URL", chatServer.URL)
	os.Setenv("TASK_SERVICE_URL", taskServer.URL)
	os.Setenv("FILE_SERVICE_URL", fileServer.URL)

	// Создаем клиенты
	userClient = http_clients.NewUserClient(userServer.URL)
	chatClient = http_clients.NewChatClient(chatServer.URL)
	taskClient = http_clients.NewTaskClient(taskServer.URL)
	fileClient = http_clients.NewFileClient(fileServer.URL)
	rolePermissionClient = http_clients.NewRolePermissionClient(chatServer.URL)

	t.Cleanup(func() {
		userServer.Close()
		chatServer.Close()
		taskServer.Close()
		fileServer.Close()
		os.Unsetenv("USER_SERVICE_URL")
		os.Unsetenv("CHAT_SERVICE_URL")
		os.Unsetenv("TASK_SERVICE_URL")
		os.Unsetenv("FILE_SERVICE_URL")
	})

	return userServer, chatServer, taskServer, fileServer, userClient, chatClient, taskClient, fileClient, rolePermissionClient
}
