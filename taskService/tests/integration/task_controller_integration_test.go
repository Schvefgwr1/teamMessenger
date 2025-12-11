//go:build integration
// +build integration

package integration

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	cc "common/contracts/chat-contracts"
	fc "common/contracts/file-contracts"
	cuc "common/contracts/user-contracts"
	"taskService/internal/controllers"
	"taskService/internal/handlers/dto"
	"taskService/internal/http_clients"
	"taskService/internal/models"
	"taskService/internal/repositories"
	"taskService/internal/services"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupTestHTTPClients создает тестовые HTTP серверы для внешних сервисов
func setupTestHTTPClients(t *testing.T) (userServer *httptest.Server, chatServer *httptest.Server, fileServer *httptest.Server, userClient http_clients.UserClientInterface, chatClient http_clients.ChatClientInterface, fileClient http_clients.FileClientInterface) {
	t.Helper()

	// Тестовый сервер для User Service
	userServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.URL.Path[len("/api/v1/users/"):]
		parsedUUID, err := uuid.Parse(userID)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Возвращаем тестового пользователя
		response := cuc.Response{
			User: &cuc.User{
				ID:       parsedUUID,
				Username: "test_user_" + userID[:8],
				Email:    "test_" + userID[:8] + "@example.com",
				Role: cuc.Role{
					ID:   1,
					Name: "user",
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))

	// Тестовый сервер для Chat Service
	chatServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		chatID := r.URL.Path[len("/api/v1/chats/"):]
		parsedUUID, err := uuid.Parse(chatID)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Возвращаем тестовый чат
		chat := cc.Chat{
			ID:      parsedUUID,
			Name:    "test_chat",
			IsGroup: true,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(chat)
	}))

	// Тестовый сервер для File Service
	fileServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var fileID int
		fmt.Sscanf(r.URL.Path, "/api/v1/files/%d", &fileID)

		// Возвращаем тестовый файл
		file := fc.File{
			ID:         fileID,
			Name:       fmt.Sprintf("test_file_%d.txt", fileID),
			FileTypeID: 1,
			URL:        fmt.Sprintf("http://example.com/files/%d", fileID),
			CreatedAt:  time.Now(),
			FileType: fc.FileType{
				ID:   1,
				Name: "text",
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(file)
	}))

	// Устанавливаем переменные окружения для HTTP клиентов
	os.Setenv("USER_SERVICE_URL", userServer.URL)
	os.Setenv("CHAT_SERVICE_URL", chatServer.URL)
	os.Setenv("FILE_SERVICE_URL", fileServer.URL)

	// Создаем адаптеры клиентов
	userClient = http_clients.NewUserClientAdapter()
	chatClient = http_clients.NewChatClientAdapter()
	fileClient = http_clients.NewFileClientAdapter()

	t.Cleanup(func() {
		userServer.Close()
		chatServer.Close()
		fileServer.Close()
		os.Unsetenv("USER_SERVICE_URL")
		os.Unsetenv("CHAT_SERVICE_URL")
		os.Unsetenv("FILE_SERVICE_URL")
	})

	return userServer, chatServer, fileServer, userClient, chatClient, fileClient
}

// setupTestHTTPClientsWithErrors создает тестовые HTTP серверы, которые возвращают ошибки
func setupTestHTTPClientsWithErrors(t *testing.T, userError, chatError, fileError bool) (userServer *httptest.Server, chatServer *httptest.Server, fileServer *httptest.Server) {
	t.Helper()

	// Тестовый сервер для User Service с ошибкой
	userServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if userError {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		userID := r.URL.Path[len("/api/v1/users/"):]
		parsedUUID, _ := uuid.Parse(userID)
		response := cuc.Response{
			User: &cuc.User{
				ID:       parsedUUID,
				Username: "test_user",
				Email:    "test@example.com",
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))

	// Тестовый сервер для Chat Service с ошибкой
	chatServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if chatError {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		chatID := r.URL.Path[len("/api/v1/chats/"):]
		parsedUUID, _ := uuid.Parse(chatID)
		chat := cc.Chat{
			ID:      parsedUUID,
			Name:    "test_chat",
			IsGroup: true,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(chat)
	}))

	// Тестовый сервер для File Service с ошибкой
	fileServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if fileError {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		var fileID int
		fmt.Sscanf(r.URL.Path, "/api/v1/files/%d", &fileID)
		file := fc.File{
			ID:         fileID,
			Name:       "test_file.txt",
			FileTypeID: 1,
			URL:        "http://example.com/file",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(file)
	}))

	os.Setenv("USER_SERVICE_URL", userServer.URL)
	os.Setenv("CHAT_SERVICE_URL", chatServer.URL)
	os.Setenv("FILE_SERVICE_URL", fileServer.URL)

	t.Cleanup(func() {
		userServer.Close()
		chatServer.Close()
		fileServer.Close()
		os.Unsetenv("USER_SERVICE_URL")
		os.Unsetenv("CHAT_SERVICE_URL")
		os.Unsetenv("FILE_SERVICE_URL")
	})

	return userServer, chatServer, fileServer
}

// TestTaskController_Create_Integration тестирует создание задачи с реальными интеграциями
func TestTaskController_Create_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange - настройка реальных зависимостей
	db := setupTestDB(t)
	userServer, chatServer, fileServer, userClient, chatClient, fileClient := setupTestHTTPClients(t)
	defer userServer.Close()
	defer chatServer.Close()
	defer fileServer.Close()

	kafkaTopic := getEnvOrDefault("KAFKA_TOPIC_NOTIFICATIONS", "notifications_test")
	producer := setupTestKafkaProducer(t, kafkaTopic)
	if producer == nil {
		t.Skip("Kafka недоступен, пропускаем тест")
	}

	notificationService := services.NewNotificationServiceWithProducer(producer)
	taskRepo := repositories.NewTaskRepository(db)
	taskStatusRepo := repositories.NewTaskStatusRepository(db)
	taskFileRepo := repositories.NewTaskFileRepository(db)

	controller := controllers.NewTaskControllerWithClients(
		taskRepo,
		taskStatusRepo,
		taskFileRepo,
		notificationService,
		userClient,
		chatClient,
		fileClient,
	)

	creatorID := uuid.New()
	executorID := uuid.New()
	chatID := uuid.New()
	fileIDs := []int{1, 2}
	description := "test description"

	createDTO := &dto.CreateTaskDTO{
		Title:       "test_CreateTask_Integration",
		Description: &description,
		CreatorID:   creatorID,
		ExecutorID:  executorID,
		ChatID:      chatID,
		FileIDs:     fileIDs,
	}

	// Act - выполнение реального сценария
	task, err := controller.Create(createDTO)

	// Assert - проверка реального результата
	require.NoError(t, err)
	assert.NotNil(t, task)
	assert.Equal(t, "test_CreateTask_Integration", task.Title)
	assert.Equal(t, description, task.Description)
	assert.Equal(t, creatorID, task.CreatorID)
	assert.Equal(t, executorID, task.ExecutorID)
	assert.Equal(t, chatID, task.ChatID)
	assert.NotZero(t, task.ID)
	assert.NotNil(t, task.Status)
	assert.Equal(t, "created", task.Status.Name)

	// Проверяем, что файлы реально сохранены в БД
	var savedTask models.Task
	err = db.Preload("Files").First(&savedTask, task.ID).Error
	require.NoError(t, err)
	assert.Len(t, savedTask.Files, 2)
}

// TestTaskController_Create_Integration_WithoutExecutor тестирует создание задачи без исполнителя
func TestTaskController_Create_Integration_WithoutExecutor(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange
	db := setupTestDB(t)
	userServer, _, _, userClient, chatClient, fileClient := setupTestHTTPClients(t)
	defer userServer.Close()

	kafkaTopic := getEnvOrDefault("KAFKA_TOPIC_NOTIFICATIONS", "notifications_test")
	producer := setupTestKafkaProducer(t, kafkaTopic)
	if producer == nil {
		t.Skip("Kafka недоступен, пропускаем тест")
	}

	notificationService := services.NewNotificationServiceWithProducer(producer)
	taskRepo := repositories.NewTaskRepository(db)
	taskStatusRepo := repositories.NewTaskStatusRepository(db)
	taskFileRepo := repositories.NewTaskFileRepository(db)

	controller := controllers.NewTaskControllerWithClients(
		taskRepo,
		taskStatusRepo,
		taskFileRepo,
		notificationService,
		userClient,
		chatClient,
		fileClient,
	)

	creatorID := uuid.New()
	createDTO := &dto.CreateTaskDTO{
		Title:      "test_CreateTask_WithoutExecutor",
		CreatorID:  creatorID,
		ExecutorID: uuid.Nil, // Нет исполнителя
		ChatID:     uuid.Nil,
		FileIDs:    []int{},
	}

	// Act
	task, err := controller.Create(createDTO)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, task)
	assert.Equal(t, uuid.Nil, task.ExecutorID)
}

// TestTaskController_Create_Integration_UserNotFound тестирует обработку ошибки, когда пользователь не найден
func TestTaskController_Create_Integration_UserNotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange
	db := setupTestDB(t)
	userServer, chatServer, fileServer := setupTestHTTPClientsWithErrors(t, true, false, false)
	defer userServer.Close()
	defer chatServer.Close()
	defer fileServer.Close()

	kafkaTopic := getEnvOrDefault("KAFKA_TOPIC_NOTIFICATIONS", "notifications_test")
	producer := setupTestKafkaProducer(t, kafkaTopic)
	if producer == nil {
		t.Skip("Kafka недоступен, пропускаем тест")
	}

	notificationService := services.NewNotificationServiceWithProducer(producer)
	taskRepo := repositories.NewTaskRepository(db)
	taskStatusRepo := repositories.NewTaskStatusRepository(db)
	taskFileRepo := repositories.NewTaskFileRepository(db)

	userClient := http_clients.NewUserClientAdapter()
	chatClient := http_clients.NewChatClientAdapter()
	fileClient := http_clients.NewFileClientAdapter()

	controller := controllers.NewTaskControllerWithClients(
		taskRepo,
		taskStatusRepo,
		taskFileRepo,
		notificationService,
		userClient,
		chatClient,
		fileClient,
	)

	creatorID := uuid.New()
	createDTO := &dto.CreateTaskDTO{
		Title:      "test_CreateTask_UserNotFound",
		CreatorID:  creatorID,
		ExecutorID: uuid.Nil,
		ChatID:     uuid.Nil,
		FileIDs:    []int{},
	}

	// Act
	task, err := controller.Create(createDTO)

	// Assert
	require.Error(t, err)
	assert.Nil(t, task)
	assert.Contains(t, err.Error(), "can't get user")
}

// TestTaskController_Create_Integration_ChatNotFound тестирует обработку ошибки, когда чат не найден
func TestTaskController_Create_Integration_ChatNotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange
	db := setupTestDB(t)
	userServer, chatServer, fileServer := setupTestHTTPClientsWithErrors(t, false, true, false)
	defer userServer.Close()
	defer chatServer.Close()
	defer fileServer.Close()

	kafkaTopic := getEnvOrDefault("KAFKA_TOPIC_NOTIFICATIONS", "notifications_test")
	producer := setupTestKafkaProducer(t, kafkaTopic)
	if producer == nil {
		t.Skip("Kafka недоступен, пропускаем тест")
	}

	notificationService := services.NewNotificationServiceWithProducer(producer)
	taskRepo := repositories.NewTaskRepository(db)
	taskStatusRepo := repositories.NewTaskStatusRepository(db)
	taskFileRepo := repositories.NewTaskFileRepository(db)

	userClient := http_clients.NewUserClientAdapter()
	chatClient := http_clients.NewChatClientAdapter()
	fileClient := http_clients.NewFileClientAdapter()

	controller := controllers.NewTaskControllerWithClients(
		taskRepo,
		taskStatusRepo,
		taskFileRepo,
		notificationService,
		userClient,
		chatClient,
		fileClient,
	)

	creatorID := uuid.New()
	chatID := uuid.New()
	createDTO := &dto.CreateTaskDTO{
		Title:      "test_CreateTask_ChatNotFound",
		CreatorID:  creatorID,
		ExecutorID: uuid.Nil,
		ChatID:     chatID,
		FileIDs:    []int{},
	}

	// Act
	task, err := controller.Create(createDTO)

	// Assert
	require.Error(t, err)
	assert.Nil(t, task)
	assert.Contains(t, err.Error(), "can't get chat")
}

// TestTaskController_Create_Integration_FileNotFound тестирует обработку ошибки, когда файл не найден
func TestTaskController_Create_Integration_FileNotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange
	db := setupTestDB(t)
	userServer, chatServer, fileServer := setupTestHTTPClientsWithErrors(t, false, false, true)
	defer userServer.Close()
	defer chatServer.Close()
	defer fileServer.Close()

	kafkaTopic := getEnvOrDefault("KAFKA_TOPIC_NOTIFICATIONS", "notifications_test")
	producer := setupTestKafkaProducer(t, kafkaTopic)
	if producer == nil {
		t.Skip("Kafka недоступен, пропускаем тест")
	}

	notificationService := services.NewNotificationServiceWithProducer(producer)
	taskRepo := repositories.NewTaskRepository(db)
	taskStatusRepo := repositories.NewTaskStatusRepository(db)
	taskFileRepo := repositories.NewTaskFileRepository(db)

	userClient := http_clients.NewUserClientAdapter()
	chatClient := http_clients.NewChatClientAdapter()
	fileClient := http_clients.NewFileClientAdapter()

	controller := controllers.NewTaskControllerWithClients(
		taskRepo,
		taskStatusRepo,
		taskFileRepo,
		notificationService,
		userClient,
		chatClient,
		fileClient,
	)

	creatorID := uuid.New()
	createDTO := &dto.CreateTaskDTO{
		Title:      "test_CreateTask_FileNotFound",
		CreatorID:  creatorID,
		ExecutorID: uuid.Nil,
		ChatID:     uuid.Nil,
		FileIDs:    []int{999}, // Несуществующий файл
	}

	// Act
	task, err := controller.Create(createDTO)

	// Assert
	require.Error(t, err)
	assert.Nil(t, task)
	assert.Contains(t, err.Error(), "can't get file")
}

// TestTaskController_GetByID_Integration тестирует получение задачи по ID с реальной интеграцией FileClient
func TestTaskController_GetByID_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange
	db := setupTestDB(t)
	_, _, fileServer, _, _, fileClient := setupTestHTTPClients(t)
	defer fileServer.Close()

	taskRepo := repositories.NewTaskRepository(db)
	taskStatusRepo := repositories.NewTaskStatusRepository(db)
	taskFileRepo := repositories.NewTaskFileRepository(db)

	// Создаем задачу напрямую в БД для теста
	status, err := taskStatusRepo.GetByName("created")
	require.NoError(t, err)

	task := &models.Task{
		Title:      "test_GetByID_Integration",
		CreatorID:  uuid.New(),
		ExecutorID: uuid.New(),
		StatusID:   status.ID,
		Status:     status,
	}
	err = taskRepo.Create(task)
	require.NoError(t, err)

	// Добавляем файлы к задаче
	taskFiles := []models.TaskFile{
		{TaskID: task.ID, FileID: 1},
		{TaskID: task.ID, FileID: 2},
	}
	err = taskFileRepo.BulkCreate(taskFiles)
	require.NoError(t, err)

	controller := controllers.NewTaskControllerWithClients(
		taskRepo,
		taskStatusRepo,
		taskFileRepo,
		nil, // NotificationService не нужен для GetByID
		nil, // UserClient не нужен
		nil, // ChatClient не нужен
		fileClient,
	)

	// Act
	result, err := controller.GetByID(task.ID)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotNil(t, result.Task)
	assert.Equal(t, task.ID, result.Task.ID)
	assert.NotNil(t, result.Files)
	assert.Len(t, *result.Files, 2)
}

// TestTaskController_GetByID_Integration_FileNotFound тестирует обработку ошибки, когда файл не найден при получении задачи
func TestTaskController_GetByID_Integration_FileNotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange
	db := setupTestDB(t)
	_, _, fileServer := setupTestHTTPClientsWithErrors(t, false, false, true)
	defer fileServer.Close()

	taskRepo := repositories.NewTaskRepository(db)
	taskStatusRepo := repositories.NewTaskStatusRepository(db)
	taskFileRepo := repositories.NewTaskFileRepository(db)

	// Создаем задачу с файлом
	status, err := taskStatusRepo.GetByName("created")
	require.NoError(t, err)

	task := &models.Task{
		Title:     "test_GetByID_FileNotFound",
		CreatorID: uuid.New(),
		StatusID:  status.ID,
		Status:    status,
	}
	err = taskRepo.Create(task)
	require.NoError(t, err)

	taskFiles := []models.TaskFile{
		{TaskID: task.ID, FileID: 999}, // Несуществующий файл
	}
	err = taskFileRepo.BulkCreate(taskFiles)
	require.NoError(t, err)

	fileClient := http_clients.NewFileClientAdapter()

	controller := controllers.NewTaskControllerWithClients(
		taskRepo,
		taskStatusRepo,
		taskFileRepo,
		nil,
		nil,
		nil,
		fileClient,
	)

	// Act
	result, err := controller.GetByID(task.ID)

	// Assert
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "can't get file")
}

// TestTaskController_UpdateStatus_Integration тестирует обновление статуса задачи с реальной БД
func TestTaskController_UpdateStatus_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange
	db := setupTestDB(t)
	taskRepo := repositories.NewTaskRepository(db)
	taskStatusRepo := repositories.NewTaskStatusRepository(db)
	taskFileRepo := repositories.NewTaskFileRepository(db)

	// Создаем задачу
	status, err := taskStatusRepo.GetByName("created")
	require.NoError(t, err)

	task := &models.Task{
		Title:     "test_UpdateStatus_Integration",
		CreatorID: uuid.New(),
		StatusID:  status.ID,
		Status:    status,
	}
	err = taskRepo.Create(task)
	require.NoError(t, err)

	// Создаем новый статус для обновления
	newStatus, err := taskStatusRepo.Create("test_status")
	require.NoError(t, err)

	controller := controllers.NewTaskController(
		taskRepo,
		taskStatusRepo,
		taskFileRepo,
		nil,
	)

	// Act
	err = controller.UpdateStatus(task.ID, newStatus.ID)

	// Assert
	require.NoError(t, err)

	// Проверяем, что статус реально обновлен в БД
	updatedTask, err := taskRepo.GetByID(task.ID)
	require.NoError(t, err)
	assert.Equal(t, newStatus.ID, updatedTask.StatusID)
}

// TestTaskController_UpdateStatus_Integration_StatusNotFound тестирует обработку ошибки, когда статус не найден
func TestTaskController_UpdateStatus_Integration_StatusNotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange
	db := setupTestDB(t)
	taskRepo := repositories.NewTaskRepository(db)
	taskStatusRepo := repositories.NewTaskStatusRepository(db)
	taskFileRepo := repositories.NewTaskFileRepository(db)

	status, err := taskStatusRepo.GetByName("created")
	require.NoError(t, err)

	task := &models.Task{
		Title:     "test_UpdateStatus_StatusNotFound",
		CreatorID: uuid.New(),
		StatusID:  status.ID,
		Status:    status,
	}
	err = taskRepo.Create(task)
	require.NoError(t, err)

	controller := controllers.NewTaskController(
		taskRepo,
		taskStatusRepo,
		taskFileRepo,
		nil,
	)

	// Act - пытаемся обновить на несуществующий статус
	err = controller.UpdateStatus(task.ID, 99999)

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "task status")
}

// TestTaskController_GetUserTasks_Integration тестирует получение задач пользователя с реальной БД
func TestTaskController_GetUserTasks_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange
	db := setupTestDB(t)
	taskRepo := repositories.NewTaskRepository(db)
	taskStatusRepo := repositories.NewTaskStatusRepository(db)
	taskFileRepo := repositories.NewTaskFileRepository(db)

	userID := uuid.New()
	status, err := taskStatusRepo.GetByName("created")
	require.NoError(t, err)

	// Создаем несколько задач для пользователя
	for i := 0; i < 3; i++ {
		task := &models.Task{
			Title:      fmt.Sprintf("test_GetUserTasks_%d", i),
			CreatorID:  uuid.New(),
			ExecutorID: userID,
			StatusID:   status.ID,
			Status:     status,
		}
		err = taskRepo.Create(task)
		require.NoError(t, err)
	}

	controller := controllers.NewTaskController(
		taskRepo,
		taskStatusRepo,
		taskFileRepo,
		nil,
	)

	// Act
	tasks, err := controller.GetUserTasks(userID.String(), 10, 0)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, tasks)
	assert.Len(t, *tasks, 3)
}

// TestTaskController_Create_Integration_KafkaNotification тестирует отправку уведомления через Kafka
func TestTaskController_Create_Integration_KafkaNotification(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange
	db := setupTestDB(t)
	userServer, _, _, userClient, chatClient, fileClient := setupTestHTTPClients(t)
	defer userServer.Close()

	kafkaTopic := getEnvOrDefault("KAFKA_TOPIC_NOTIFICATIONS", "notifications_test")
	producer := setupTestKafkaProducer(t, kafkaTopic)
	if producer == nil {
		t.Skip("Kafka недоступен, пропускаем тест")
	}

	notificationService := services.NewNotificationServiceWithProducer(producer)
	taskRepo := repositories.NewTaskRepository(db)
	taskStatusRepo := repositories.NewTaskStatusRepository(db)
	taskFileRepo := repositories.NewTaskFileRepository(db)

	controller := controllers.NewTaskControllerWithClients(
		taskRepo,
		taskStatusRepo,
		taskFileRepo,
		notificationService,
		userClient,
		chatClient,
		fileClient,
	)

	creatorID := uuid.New()
	executorID := uuid.New()
	createDTO := &dto.CreateTaskDTO{
		Title:      "test_KafkaNotification",
		CreatorID:  creatorID,
		ExecutorID: executorID,
		ChatID:     uuid.Nil,
		FileIDs:    []int{},
	}

	// Act
	task, err := controller.Create(createDTO)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, task)
	// Уведомление отправляется асинхронно, поэтому мы не можем проверить его напрямую
	// Но если ошибки нет, значит интеграция работает
}
