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

	"chatService/internal/controllers"
	"chatService/internal/handlers/dto"
	"chatService/internal/http_clients"
	"chatService/internal/models"
	"chatService/internal/repositories"
	"chatService/internal/services"
	fc "common/contracts/file-contracts"
	cuc "common/contracts/user-contracts"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupTestHTTPClients создает тестовые HTTP серверы для внешних сервисов
func setupTestHTTPClients(t *testing.T) (userServer *httptest.Server, fileServer *httptest.Server, userClient http_clients.UserClientInterface, fileClient http_clients.FileClientInterface) {
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
	os.Setenv("FILE_SERVICE_URL", fileServer.URL)

	// Создаем адаптеры клиентов
	userClient = http_clients.NewUserClientAdapter()
	fileClient = http_clients.NewFileClientAdapter()

	t.Cleanup(func() {
		userServer.Close()
		fileServer.Close()
		os.Unsetenv("USER_SERVICE_URL")
		os.Unsetenv("FILE_SERVICE_URL")
	})

	return userServer, fileServer, userClient, fileClient
}

// setupTestHTTPClientsWithErrors создает тестовые HTTP серверы, которые возвращают ошибки
func setupTestHTTPClientsWithErrors(t *testing.T, userError, fileError bool) (userServer *httptest.Server, fileServer *httptest.Server) {
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
	os.Setenv("FILE_SERVICE_URL", fileServer.URL)

	t.Cleanup(func() {
		userServer.Close()
		fileServer.Close()
		os.Unsetenv("USER_SERVICE_URL")
		os.Unsetenv("FILE_SERVICE_URL")
	})

	return userServer, fileServer
}

// TestChatController_GetChatByID_Integration тестирует получение чата по ID с реальной интеграцией FileClient
func TestChatController_GetChatByID_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange
	db := setupTestDB(t)
	_, fileServer, _, fileClient := setupTestHTTPClients(t)
	defer fileServer.Close()

	chatRepo := repositories.NewChatRepository(db)
	chatUserRepo := repositories.NewChatUserRepository(db)
	chatRoleRepo := repositories.NewChatRoleRepository(db)

	// Создаем чат напрямую в БД для теста
	chat := &models.Chat{
		ID:      uuid.New(),
		Name:    "test_GetChatByID_Integration",
		IsGroup: true,
	}
	err := chatRepo.CreateChat(chat)
	require.NoError(t, err)

	controller := controllers.NewChatControllerWithClients(
		chatRepo,
		chatUserRepo,
		chatRoleRepo,
		nil, // NotificationService не нужен для GetChatByID
		fileClient,
		http_clients.NewUserClientAdapter(),
	)

	// Act
	result, err := controller.GetChatByID(chat.ID)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, chat.ID, result.ID)
	assert.Equal(t, "test_GetChatByID_Integration", result.Name)
	assert.Equal(t, true, result.IsGroup)
}

// TestChatController_GetChatByID_Integration_NotFound тестирует обработку ошибки, когда чат не найден
func TestChatController_GetChatByID_Integration_NotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange
	db := setupTestDB(t)
	chatRepo := repositories.NewChatRepository(db)
	chatUserRepo := repositories.NewChatUserRepository(db)
	chatRoleRepo := repositories.NewChatRoleRepository(db)

	controller := controllers.NewChatController(
		chatRepo,
		chatUserRepo,
		chatRoleRepo,
		nil,
	)

	nonExistentChatID := uuid.New()

	// Act
	result, err := controller.GetChatByID(nonExistentChatID)

	// Assert
	require.Error(t, err)
	assert.Nil(t, result)
}

// TestChatController_GetUserChats_Integration тестирует получение чатов пользователя с реальной БД
func TestChatController_GetUserChats_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange
	db := setupTestDB(t)
	_, fileServer, _, fileClient := setupTestHTTPClients(t)
	defer fileServer.Close()

	chatRepo := repositories.NewChatRepository(db)
	chatUserRepo := repositories.NewChatUserRepository(db)
	chatRoleRepo := repositories.NewChatRoleRepository(db)

	userID := uuid.New()

	// Получаем роль main для добавления пользователя
	mainRole, err := chatRoleRepo.GetRoleByName("main")
	require.NoError(t, err)

	// Создаем несколько чатов для пользователя
	for i := 0; i < 3; i++ {
		chat := &models.Chat{
			ID:      uuid.New(),
			Name:    fmt.Sprintf("test_GetUserChats_%d", i),
			IsGroup: true,
		}
		err = chatRepo.CreateChat(chat)
		require.NoError(t, err)

		chatUser := &models.ChatUser{
			ChatID: chat.ID,
			UserID: userID,
			RoleID: mainRole.ID,
		}
		err = chatUserRepo.AddUserToChat(chatUser)
		require.NoError(t, err)
	}

	controller := controllers.NewChatControllerWithClients(
		chatRepo,
		chatUserRepo,
		chatRoleRepo,
		nil,
		fileClient,
		http_clients.NewUserClientAdapter(),
	)

	// Act
	chats, err := controller.GetUserChats(userID)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, chats)
	assert.Len(t, *chats, 3)
}

// TestChatController_CreateChat_Integration тестирует создание чата с реальными интеграциями
func TestChatController_CreateChat_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange
	db := setupTestDB(t)
	userServer, fileServer, userClient, fileClient := setupTestHTTPClients(t)
	defer userServer.Close()
	defer fileServer.Close()

	kafkaTopic := getEnvOrDefault("KAFKA_TOPIC_NOTIFICATIONS", "notifications_test")
	producer := setupTestKafkaProducer(t, kafkaTopic)
	if producer == nil {
		t.Skip("Kafka недоступен, пропускаем тест")
	}

	notificationService := services.NewNotificationServiceWithProducer(producer)
	chatRepo := repositories.NewChatRepository(db)
	chatUserRepo := repositories.NewChatUserRepository(db)
	chatRoleRepo := repositories.NewChatRoleRepository(db)

	controller := controllers.NewChatControllerWithClients(
		chatRepo,
		chatUserRepo,
		chatRoleRepo,
		notificationService,
		fileClient,
		userClient,
	)

	ownerID := uuid.New()
	userID1 := uuid.New()
	userID2 := uuid.New()
	description := "test description"

	createDTO := &dto.CreateChatDTO{
		Name:        "test_CreateChat_Integration",
		Description: &description,
		OwnerID:     ownerID,
		UserIDs:     []uuid.UUID{userID1, userID2},
	}

	// Act
	chatID, err := controller.CreateChat(createDTO)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, chatID)

	// Проверяем, что чат реально создан в БД
	chat, err := chatRepo.GetChatByID(*chatID)
	require.NoError(t, err)
	assert.Equal(t, "test_CreateChat_Integration", chat.Name)
	assert.Equal(t, true, chat.IsGroup) // Должен быть групповым, т.к. больше 1 пользователя
}

// TestChatController_CreateChat_Integration_WithoutUsers тестирует создание чата без дополнительных пользователей
func TestChatController_CreateChat_Integration_WithoutUsers(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange
	db := setupTestDB(t)
	userServer, _, userClient, fileClient := setupTestHTTPClients(t)
	defer userServer.Close()

	kafkaTopic := getEnvOrDefault("KAFKA_TOPIC_NOTIFICATIONS", "notifications_test")
	producer := setupTestKafkaProducer(t, kafkaTopic)
	if producer == nil {
		t.Skip("Kafka недоступен, пропускаем тест")
	}

	notificationService := services.NewNotificationServiceWithProducer(producer)
	chatRepo := repositories.NewChatRepository(db)
	chatUserRepo := repositories.NewChatUserRepository(db)
	chatRoleRepo := repositories.NewChatRoleRepository(db)

	controller := controllers.NewChatControllerWithClients(
		chatRepo,
		chatUserRepo,
		chatRoleRepo,
		notificationService,
		fileClient,
		userClient,
	)

	ownerID := uuid.New()

	createDTO := &dto.CreateChatDTO{
		Name:    "test_CreateChat_WithoutUsers",
		OwnerID: ownerID,
		UserIDs: []uuid.UUID{}, // Нет дополнительных пользователей
	}

	// Act
	chatID, err := controller.CreateChat(createDTO)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, chatID)

	// Проверяем, что чат создан как личный (не групповой)
	chat, err := chatRepo.GetChatByID(*chatID)
	require.NoError(t, err)
	assert.Equal(t, false, chat.IsGroup)
}

// TestChatController_CreateChat_Integration_UserNotFound тестирует обработку ошибки, когда пользователь не найден
func TestChatController_CreateChat_Integration_UserNotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange
	db := setupTestDB(t)
	userServer, fileServer := setupTestHTTPClientsWithErrors(t, true, false)
	defer userServer.Close()
	defer fileServer.Close()

	kafkaTopic := getEnvOrDefault("KAFKA_TOPIC_NOTIFICATIONS", "notifications_test")
	producer := setupTestKafkaProducer(t, kafkaTopic)
	if producer == nil {
		t.Skip("Kafka недоступен, пропускаем тест")
	}

	notificationService := services.NewNotificationServiceWithProducer(producer)
	chatRepo := repositories.NewChatRepository(db)
	chatUserRepo := repositories.NewChatUserRepository(db)
	chatRoleRepo := repositories.NewChatRoleRepository(db)

	userClient := http_clients.NewUserClientAdapter()
	fileClient := http_clients.NewFileClientAdapter()

	controller := controllers.NewChatControllerWithClients(
		chatRepo,
		chatUserRepo,
		chatRoleRepo,
		notificationService,
		fileClient,
		userClient,
	)

	ownerID := uuid.New()

	createDTO := &dto.CreateChatDTO{
		Name:    "test_CreateChat_UserNotFound",
		OwnerID: ownerID,
		UserIDs: []uuid.UUID{},
	}

	// Act
	chatID, err := controller.CreateChat(createDTO)

	// Assert
	require.Error(t, err)
	assert.Nil(t, chatID)
}

// TestChatController_CreateChat_Integration_FileNotFound тестирует обработку ошибки, когда файл не найден
func TestChatController_CreateChat_Integration_FileNotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange
	db := setupTestDB(t)
	userServer, fileServer := setupTestHTTPClientsWithErrors(t, false, true)
	defer userServer.Close()
	defer fileServer.Close()

	kafkaTopic := getEnvOrDefault("KAFKA_TOPIC_NOTIFICATIONS", "notifications_test")
	producer := setupTestKafkaProducer(t, kafkaTopic)
	if producer == nil {
		t.Skip("Kafka недоступен, пропускаем тест")
	}

	notificationService := services.NewNotificationServiceWithProducer(producer)
	chatRepo := repositories.NewChatRepository(db)
	chatUserRepo := repositories.NewChatUserRepository(db)
	chatRoleRepo := repositories.NewChatRoleRepository(db)

	userClient := http_clients.NewUserClientAdapter()
	fileClient := http_clients.NewFileClientAdapter()

	controller := controllers.NewChatControllerWithClients(
		chatRepo,
		chatUserRepo,
		chatRoleRepo,
		notificationService,
		fileClient,
		userClient,
	)

	ownerID := uuid.New()
	avatarFileID := 999 // Несуществующий файл

	createDTO := &dto.CreateChatDTO{
		Name:         "test_CreateChat_FileNotFound",
		OwnerID:      ownerID,
		UserIDs:      []uuid.UUID{},
		AvatarFileID: &avatarFileID,
	}

	// Act
	chatID, err := controller.CreateChat(createDTO)

	// Assert
	require.Error(t, err)
	assert.Nil(t, chatID)
}

// TestChatController_UpdateChat_Integration тестирует обновление чата с реальной БД
func TestChatController_UpdateChat_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange
	db := setupTestDB(t)
	userServer, fileServer, userClient, fileClient := setupTestHTTPClients(t)
	defer userServer.Close()
	defer fileServer.Close()

	kafkaTopic := getEnvOrDefault("KAFKA_TOPIC_NOTIFICATIONS", "notifications_test")
	producer := setupTestKafkaProducer(t, kafkaTopic)
	if producer == nil {
		t.Skip("Kafka недоступен, пропускаем тест")
	}

	notificationService := services.NewNotificationServiceWithProducer(producer)
	chatRepo := repositories.NewChatRepository(db)
	chatUserRepo := repositories.NewChatUserRepository(db)
	chatRoleRepo := repositories.NewChatRoleRepository(db)

	// Создаем чат
	chat := &models.Chat{
		ID:      uuid.New(),
		Name:    "test_UpdateChat_Integration_Old",
		IsGroup: true,
	}
	err := chatRepo.CreateChat(chat)
	require.NoError(t, err)

	controller := controllers.NewChatControllerWithClients(
		chatRepo,
		chatUserRepo,
		chatRoleRepo,
		notificationService,
		fileClient,
		userClient,
	)

	newName := "test_UpdateChat_Integration_New"
	newDescription := "updated description"

	updateDTO := &dto.UpdateChatDTO{
		Name:          &newName,
		Description:   &newDescription,
		AddUserIDs:    []uuid.UUID{},
		RemoveUserIDs: []uuid.UUID{},
	}

	// Act
	result, err := controller.UpdateChat(chat.ID, updateDTO)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, newName, result.Chat.Name)
	assert.Equal(t, newDescription, *result.Chat.Description)

	// Проверяем, что чат реально обновлен в БД
	updatedChat, err := chatRepo.GetChatByID(chat.ID)
	require.NoError(t, err)
	assert.Equal(t, newName, updatedChat.Name)
}

// TestChatController_UpdateChat_Integration_AddUsers тестирует добавление пользователей в чат
func TestChatController_UpdateChat_Integration_AddUsers(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange
	db := setupTestDB(t)
	userServer, fileServer, userClient, fileClient := setupTestHTTPClients(t)
	defer userServer.Close()
	defer fileServer.Close()

	kafkaTopic := getEnvOrDefault("KAFKA_TOPIC_NOTIFICATIONS", "notifications_test")
	producer := setupTestKafkaProducer(t, kafkaTopic)
	if producer == nil {
		t.Skip("Kafka недоступен, пропускаем тест")
	}

	notificationService := services.NewNotificationServiceWithProducer(producer)
	chatRepo := repositories.NewChatRepository(db)
	chatUserRepo := repositories.NewChatUserRepository(db)
	chatRoleRepo := repositories.NewChatRoleRepository(db)

	// Создаем чат
	chat := &models.Chat{
		ID:      uuid.New(),
		Name:    "test_UpdateChat_AddUsers",
		IsGroup: true,
	}
	err := chatRepo.CreateChat(chat)
	require.NoError(t, err)

	controller := controllers.NewChatControllerWithClients(
		chatRepo,
		chatUserRepo,
		chatRoleRepo,
		notificationService,
		fileClient,
		userClient,
	)

	newUserID := uuid.New()

	updateDTO := &dto.UpdateChatDTO{
		AddUserIDs:    []uuid.UUID{newUserID},
		RemoveUserIDs: []uuid.UUID{},
	}

	// Act
	result, err := controller.UpdateChat(chat.ID, updateDTO)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.UpdateUsers, 1)
	assert.Equal(t, newUserID, result.UpdateUsers[0].UserID)
	assert.Equal(t, "created", result.UpdateUsers[0].State)

	// Проверяем, что пользователь реально добавлен в БД
	chatUsers, err := chatUserRepo.GetChatUsers(chat.ID)
	require.NoError(t, err)
	assert.Len(t, chatUsers, 1)
}

// TestChatController_UpdateChat_Integration_RemoveUsers тестирует удаление пользователей из чата
func TestChatController_UpdateChat_Integration_RemoveUsers(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange
	db := setupTestDB(t)
	userServer, fileServer, userClient, fileClient := setupTestHTTPClients(t)
	defer userServer.Close()
	defer fileServer.Close()

	kafkaTopic := getEnvOrDefault("KAFKA_TOPIC_NOTIFICATIONS", "notifications_test")
	producer := setupTestKafkaProducer(t, kafkaTopic)
	if producer == nil {
		t.Skip("Kafka недоступен, пропускаем тест")
	}

	notificationService := services.NewNotificationServiceWithProducer(producer)
	chatRepo := repositories.NewChatRepository(db)
	chatUserRepo := repositories.NewChatUserRepository(db)
	chatRoleRepo := repositories.NewChatRoleRepository(db)

	// Создаем чат
	chat := &models.Chat{
		ID:      uuid.New(),
		Name:    "test_UpdateChat_RemoveUsers",
		IsGroup: true,
	}
	err := chatRepo.CreateChat(chat)
	require.NoError(t, err)

	// Добавляем пользователя в чат
	mainRole, err := chatRoleRepo.GetRoleByName("main")
	require.NoError(t, err)

	userID := uuid.New()
	chatUser := &models.ChatUser{
		ChatID: chat.ID,
		UserID: userID,
		RoleID: mainRole.ID,
	}
	err = chatUserRepo.AddUserToChat(chatUser)
	require.NoError(t, err)

	controller := controllers.NewChatControllerWithClients(
		chatRepo,
		chatUserRepo,
		chatRoleRepo,
		notificationService,
		fileClient,
		userClient,
	)

	updateDTO := &dto.UpdateChatDTO{
		AddUserIDs:    []uuid.UUID{},
		RemoveUserIDs: []uuid.UUID{userID},
	}

	// Act
	result, err := controller.UpdateChat(chat.ID, updateDTO)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.UpdateUsers, 1)
	assert.Equal(t, userID, result.UpdateUsers[0].UserID)
	assert.Equal(t, "deleted", result.UpdateUsers[0].State)

	// Проверяем, что пользователь реально удален из БД
	chatUsers, err := chatUserRepo.GetChatUsers(chat.ID)
	require.NoError(t, err)
	assert.Len(t, chatUsers, 0)
}

// TestChatController_UpdateChat_Integration_ChatNotFound тестирует обработку ошибки, когда чат не найден
func TestChatController_UpdateChat_Integration_ChatNotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange
	db := setupTestDB(t)
	chatRepo := repositories.NewChatRepository(db)
	chatUserRepo := repositories.NewChatUserRepository(db)
	chatRoleRepo := repositories.NewChatRoleRepository(db)

	controller := controllers.NewChatController(
		chatRepo,
		chatUserRepo,
		chatRoleRepo,
		nil,
	)

	nonExistentChatID := uuid.New()
	newName := "new name"

	updateDTO := &dto.UpdateChatDTO{
		Name: &newName,
	}

	// Act
	result, err := controller.UpdateChat(nonExistentChatID, updateDTO)

	// Assert
	require.Error(t, err)
	assert.Nil(t, result)
}

// TestChatController_DeleteChat_Integration тестирует удаление чата с реальной БД
func TestChatController_DeleteChat_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange
	db := setupTestDB(t)
	chatRepo := repositories.NewChatRepository(db)
	chatUserRepo := repositories.NewChatUserRepository(db)
	chatRoleRepo := repositories.NewChatRoleRepository(db)

	// Создаем чат
	chat := &models.Chat{
		ID:      uuid.New(),
		Name:    "test_DeleteChat_Integration",
		IsGroup: true,
	}
	err := chatRepo.CreateChat(chat)
	require.NoError(t, err)

	// Добавляем пользователя в чат
	mainRole, err := chatRoleRepo.GetRoleByName("main")
	require.NoError(t, err)

	chatUser := &models.ChatUser{
		ChatID: chat.ID,
		UserID: uuid.New(),
		RoleID: mainRole.ID,
	}
	err = chatUserRepo.AddUserToChat(chatUser)
	require.NoError(t, err)

	controller := controllers.NewChatController(
		chatRepo,
		chatUserRepo,
		chatRoleRepo,
		nil,
	)

	// Act
	err = controller.DeleteChat(chat.ID)

	// Assert
	require.NoError(t, err)

	// Проверяем, что чат реально удален из БД
	_, err = chatRepo.GetChatByID(chat.ID)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "record not found")
}

// TestChatController_ChangeUserRole_Integration тестирует изменение роли пользователя с реальной БД
func TestChatController_ChangeUserRole_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange
	db := setupTestDB(t)
	chatRepo := repositories.NewChatRepository(db)
	chatUserRepo := repositories.NewChatUserRepository(db)
	chatRoleRepo := repositories.NewChatRoleRepository(db)

	// Создаем чат
	chat := &models.Chat{
		ID:      uuid.New(),
		Name:    "test_ChangeUserRole_Integration",
		IsGroup: true,
	}
	err := chatRepo.CreateChat(chat)
	require.NoError(t, err)

	// Получаем роли
	mainRole, err := chatRoleRepo.GetRoleByName("main")
	require.NoError(t, err)
	ownerRole, err := chatRoleRepo.GetRoleByName("owner")
	require.NoError(t, err)

	// Добавляем пользователя в чат с ролью main
	userID := uuid.New()
	chatUser := &models.ChatUser{
		ChatID: chat.ID,
		UserID: userID,
		RoleID: mainRole.ID,
	}
	err = chatUserRepo.AddUserToChat(chatUser)
	require.NoError(t, err)

	controller := controllers.NewChatController(
		chatRepo,
		chatUserRepo,
		chatRoleRepo,
		nil,
	)

	// Act
	err = controller.ChangeUserRole(chat.ID, userID, ownerRole.ID)

	// Assert
	require.NoError(t, err)

	// Проверяем, что роль реально изменена в БД
	updatedChatUser, err := chatUserRepo.GetChatUser(userID, chat.ID)
	require.NoError(t, err)
	assert.Equal(t, ownerRole.ID, updatedChatUser.RoleID)
}

// TestChatController_ChangeUserRole_Integration_UserNotInChat тестирует обработку ошибки, когда пользователь не в чате
func TestChatController_ChangeUserRole_Integration_UserNotInChat(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange
	db := setupTestDB(t)
	chatRepo := repositories.NewChatRepository(db)
	chatUserRepo := repositories.NewChatUserRepository(db)
	chatRoleRepo := repositories.NewChatRoleRepository(db)

	// Создаем чат
	chat := &models.Chat{
		ID:      uuid.New(),
		Name:    "test_ChangeUserRole_UserNotInChat",
		IsGroup: true,
	}
	err := chatRepo.CreateChat(chat)
	require.NoError(t, err)

	ownerRole, err := chatRoleRepo.GetRoleByName("owner")
	require.NoError(t, err)

	controller := controllers.NewChatController(
		chatRepo,
		chatUserRepo,
		chatRoleRepo,
		nil,
	)

	nonExistentUserID := uuid.New()

	// Act
	err = controller.ChangeUserRole(chat.ID, nonExistentUserID, ownerRole.ID)

	// Assert
	require.Error(t, err)
}

// TestChatController_BanUser_Integration тестирует бан пользователя с реальной БД
func TestChatController_BanUser_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange
	db := setupTestDB(t)
	chatRepo := repositories.NewChatRepository(db)
	chatUserRepo := repositories.NewChatUserRepository(db)
	chatRoleRepo := repositories.NewChatRoleRepository(db)

	// Создаем чат
	chat := &models.Chat{
		ID:      uuid.New(),
		Name:    "test_BanUser_Integration",
		IsGroup: true,
	}
	err := chatRepo.CreateChat(chat)
	require.NoError(t, err)

	// Получаем роль main
	mainRole, err := chatRoleRepo.GetRoleByName("main")
	require.NoError(t, err)

	// Добавляем пользователя в чат
	userID := uuid.New()
	chatUser := &models.ChatUser{
		ChatID: chat.ID,
		UserID: userID,
		RoleID: mainRole.ID,
	}
	err = chatUserRepo.AddUserToChat(chatUser)
	require.NoError(t, err)

	controller := controllers.NewChatController(
		chatRepo,
		chatUserRepo,
		chatRoleRepo,
		nil,
	)

	// Act
	err = controller.BanUser(chat.ID, userID)

	// Assert
	require.NoError(t, err)

	// Проверяем, что роль реально изменена на banned в БД
	bannedRole, err := chatRoleRepo.GetRoleByName("banned")
	require.NoError(t, err)

	updatedChatUser, err := chatUserRepo.GetChatUser(userID, chat.ID)
	require.NoError(t, err)
	assert.Equal(t, bannedRole.ID, updatedChatUser.RoleID)
}

// TestChatController_GetUserRoleInChat_Integration тестирует получение роли пользователя в чате
func TestChatController_GetUserRoleInChat_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange
	db := setupTestDB(t)
	chatRepo := repositories.NewChatRepository(db)
	chatUserRepo := repositories.NewChatUserRepository(db)
	chatRoleRepo := repositories.NewChatRoleRepository(db)

	// Создаем чат
	chat := &models.Chat{
		ID:      uuid.New(),
		Name:    "test_GetUserRoleInChat_Integration",
		IsGroup: true,
	}
	err := chatRepo.CreateChat(chat)
	require.NoError(t, err)

	// Получаем роль main
	mainRole, err := chatRoleRepo.GetRoleByName("main")
	require.NoError(t, err)

	// Добавляем двух пользователей в чат
	requesterID := uuid.New()
	userID := uuid.New()

	requesterChatUser := &models.ChatUser{
		ChatID: chat.ID,
		UserID: requesterID,
		RoleID: mainRole.ID,
	}
	err = chatUserRepo.AddUserToChat(requesterChatUser)
	require.NoError(t, err)

	userChatUser := &models.ChatUser{
		ChatID: chat.ID,
		UserID: userID,
		RoleID: mainRole.ID,
	}
	err = chatUserRepo.AddUserToChat(userChatUser)
	require.NoError(t, err)

	controller := controllers.NewChatController(
		chatRepo,
		chatUserRepo,
		chatRoleRepo,
		nil,
	)

	// Act
	roleName, err := controller.GetUserRoleInChat(chat.ID, userID, requesterID)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "main", roleName)
}

// TestChatController_GetUserRoleInChat_Integration_Unauthorized тестирует обработку ошибки, когда запрашивающий не в чате
func TestChatController_GetUserRoleInChat_Integration_Unauthorized(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange
	db := setupTestDB(t)
	chatRepo := repositories.NewChatRepository(db)
	chatUserRepo := repositories.NewChatUserRepository(db)
	chatRoleRepo := repositories.NewChatRoleRepository(db)

	// Создаем чат
	chat := &models.Chat{
		ID:      uuid.New(),
		Name:    "test_GetUserRoleInChat_Unauthorized",
		IsGroup: true,
	}
	err := chatRepo.CreateChat(chat)
	require.NoError(t, err)

	controller := controllers.NewChatController(
		chatRepo,
		chatUserRepo,
		chatRoleRepo,
		nil,
	)

	requesterID := uuid.New()
	userID := uuid.New()

	// Act
	roleName, err := controller.GetUserRoleInChat(chat.ID, userID, requesterID)

	// Assert
	require.Error(t, err)
	assert.Empty(t, roleName)
}

// TestChatController_GetMyRoleWithPermissions_Integration тестирует получение роли с правами доступа
func TestChatController_GetMyRoleWithPermissions_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange
	db := setupTestDB(t)
	chatRepo := repositories.NewChatRepository(db)
	chatUserRepo := repositories.NewChatUserRepository(db)
	chatRoleRepo := repositories.NewChatRoleRepository(db)

	// Создаем чат
	chat := &models.Chat{
		ID:      uuid.New(),
		Name:    "test_GetMyRoleWithPermissions_Integration",
		IsGroup: true,
	}
	err := chatRepo.CreateChat(chat)
	require.NoError(t, err)

	// Получаем роль owner (которая имеет все permissions)
	ownerRole, err := chatRoleRepo.GetRoleByName("owner")
	require.NoError(t, err)

	// Добавляем пользователя в чат с ролью owner
	userID := uuid.New()
	chatUser := &models.ChatUser{
		ChatID: chat.ID,
		UserID: userID,
		RoleID: ownerRole.ID,
	}
	err = chatUserRepo.AddUserToChat(chatUser)
	require.NoError(t, err)

	controller := controllers.NewChatController(
		chatRepo,
		chatUserRepo,
		chatRoleRepo,
		nil,
	)

	// Act
	role, err := controller.GetMyRoleWithPermissions(chat.ID, userID)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, role)
	assert.Equal(t, "owner", role.Name)
	assert.NotEmpty(t, role.Permissions)
}

// TestChatController_GetChatMembers_Integration тестирует получение участников чата
func TestChatController_GetChatMembers_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange
	db := setupTestDB(t)
	chatRepo := repositories.NewChatRepository(db)
	chatUserRepo := repositories.NewChatUserRepository(db)
	chatRoleRepo := repositories.NewChatRoleRepository(db)

	// Создаем чат
	chat := &models.Chat{
		ID:      uuid.New(),
		Name:    "test_GetChatMembers_Integration",
		IsGroup: true,
	}
	err := chatRepo.CreateChat(chat)
	require.NoError(t, err)

	// Получаем роль main
	mainRole, err := chatRoleRepo.GetRoleByName("main")
	require.NoError(t, err)

	// Добавляем нескольких пользователей в чат
	for i := 0; i < 3; i++ {
		chatUser := &models.ChatUser{
			ChatID: chat.ID,
			UserID: uuid.New(),
			RoleID: mainRole.ID,
		}
		err = chatUserRepo.AddUserToChat(chatUser)
		require.NoError(t, err)
	}

	controller := controllers.NewChatController(
		chatRepo,
		chatUserRepo,
		chatRoleRepo,
		nil,
	)

	// Act
	members, err := controller.GetChatMembers(chat.ID)

	// Assert
	require.NoError(t, err)
	assert.Len(t, members, 3)
	for _, member := range members {
		assert.Equal(t, chat.ID, member.ChatID)
		assert.NotNil(t, member.Role)
	}
}
