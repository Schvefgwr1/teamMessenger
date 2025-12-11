//go:build integration
// +build integration

package integration

import (
	"fmt"
	"testing"

	"chatService/internal/controllers"
	"chatService/internal/handlers/dto"
	"chatService/internal/http_clients"
	"chatService/internal/models"
	"chatService/internal/repositories"
	ac "common/contracts/api-chat"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestMessageController_SendMessage_Integration тестирует отправку сообщения с реальными интеграциями
func TestMessageController_SendMessage_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange
	db := setupTestDB(t)
	userServer, fileServer, userClient, fileClient := setupTestHTTPClients(t)
	defer userServer.Close()
	defer fileServer.Close()

	messageRepo := repositories.NewMessageRepository(db)
	chatRepo := repositories.NewChatRepository(db)
	chatUserRepo := repositories.NewChatUserRepository(db)

	// Создаем чат
	chat := &models.Chat{
		ID:      uuid.New(),
		Name:    "test_SendMessage_Integration",
		IsGroup: true,
	}
	err := chatRepo.CreateChat(chat)
	require.NoError(t, err)

	// Добавляем пользователя в чат (требуется для внешнего ключа)
	chatRoleRepo := repositories.NewChatRoleRepository(db)
	mainRole, err := chatRoleRepo.GetRoleByName("main")
	require.NoError(t, err)

	senderID := uuid.New()
	chatUser := &models.ChatUser{
		ChatID: chat.ID,
		UserID: senderID,
		RoleID: mainRole.ID,
	}
	err = chatUserRepo.AddUserToChat(chatUser)
	require.NoError(t, err)

	controller := controllers.NewMessageControllerWithClients(
		messageRepo,
		chatRepo,
		chatUserRepo,
		fileClient,
		userClient,
	)
	content := "test_SendMessage_Integration message content"
	fileIDs := []int{1, 2}

	createDTO := &dto.CreateMessageDTO{
		Content: content,
		FileIDs: fileIDs,
	}

	// Act
	message, err := controller.SendMessage(senderID, chat.ID, createDTO)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, message)
	assert.Equal(t, content, message.Content)
	assert.Equal(t, chat.ID, message.ChatID)
	assert.Equal(t, senderID, *message.SenderID)
	assert.NotZero(t, message.ID)

	// Проверяем, что сообщение реально сохранено в БД
	var savedMessage models.Message
	err = db.Where("id = ?", message.ID).First(&savedMessage).Error
	require.NoError(t, err)
	assert.Equal(t, content, savedMessage.Content)
}

// TestMessageController_SendMessage_Integration_WithoutFiles тестирует отправку сообщения без файлов
func TestMessageController_SendMessage_Integration_WithoutFiles(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange
	db := setupTestDB(t)
	userServer, _, userClient, fileClient := setupTestHTTPClients(t)
	defer userServer.Close()

	messageRepo := repositories.NewMessageRepository(db)
	chatRepo := repositories.NewChatRepository(db)
	chatUserRepo := repositories.NewChatUserRepository(db)

	// Создаем чат
	chat := &models.Chat{
		ID:      uuid.New(),
		Name:    "test_SendMessage_WithoutFiles",
		IsGroup: true,
	}
	err := chatRepo.CreateChat(chat)
	require.NoError(t, err)

	// Добавляем пользователя в чат (требуется для внешнего ключа)
	chatRoleRepo := repositories.NewChatRoleRepository(db)
	mainRole, err := chatRoleRepo.GetRoleByName("main")
	require.NoError(t, err)

	senderID := uuid.New()
	chatUser := &models.ChatUser{
		ChatID: chat.ID,
		UserID: senderID,
		RoleID: mainRole.ID,
	}
	err = chatUserRepo.AddUserToChat(chatUser)
	require.NoError(t, err)

	controller := controllers.NewMessageControllerWithClients(
		messageRepo,
		chatRepo,
		chatUserRepo,
		fileClient,
		userClient,
	)
	content := "test message without files"

	createDTO := &dto.CreateMessageDTO{
		Content: content,
		FileIDs: []int{}, // Нет файлов
	}

	// Act
	message, err := controller.SendMessage(senderID, chat.ID, createDTO)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, message)
	assert.Equal(t, content, message.Content)
	assert.Empty(t, message.Files)
}

// TestMessageController_SendMessage_Integration_ChatNotFound тестирует обработку ошибки, когда чат не найден
func TestMessageController_SendMessage_Integration_ChatNotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange
	db := setupTestDB(t)
	messageRepo := repositories.NewMessageRepository(db)
	chatRepo := repositories.NewChatRepository(db)
	chatUserRepo := repositories.NewChatUserRepository(db)

	controller := controllers.NewMessageController(
		messageRepo,
		chatRepo,
		chatUserRepo,
	)

	nonExistentChatID := uuid.New()
	senderID := uuid.New()

	createDTO := &dto.CreateMessageDTO{
		Content: "test message",
		FileIDs: []int{},
	}

	// Act
	message, err := controller.SendMessage(senderID, nonExistentChatID, createDTO)

	// Assert
	require.Error(t, err)
	assert.Nil(t, message)
}

// TestMessageController_SendMessage_Integration_UserNotFound тестирует обработку ошибки, когда пользователь не найден
func TestMessageController_SendMessage_Integration_UserNotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange
	db := setupTestDB(t)
	userServer, fileServer := setupTestHTTPClientsWithErrors(t, true, false)
	defer userServer.Close()
	defer fileServer.Close()

	messageRepo := repositories.NewMessageRepository(db)
	chatRepo := repositories.NewChatRepository(db)
	chatUserRepo := repositories.NewChatUserRepository(db)

	// Создаем чат
	chat := &models.Chat{
		ID:      uuid.New(),
		Name:    "test_SendMessage_UserNotFound",
		IsGroup: true,
	}
	err := chatRepo.CreateChat(chat)
	require.NoError(t, err)

	userClient := http_clients.NewUserClientAdapter()
	fileClient := http_clients.NewFileClientAdapter()

	controller := controllers.NewMessageControllerWithClients(
		messageRepo,
		chatRepo,
		chatUserRepo,
		fileClient,
		userClient,
	)

	nonExistentSenderID := uuid.New()

	createDTO := &dto.CreateMessageDTO{
		Content: "test message",
		FileIDs: []int{},
	}

	// Act
	message, err := controller.SendMessage(nonExistentSenderID, chat.ID, createDTO)

	// Assert
	require.Error(t, err)
	assert.Nil(t, message)
}

// TestMessageController_SendMessage_Integration_FileNotFound тестирует обработку ошибки, когда файл не найден
func TestMessageController_SendMessage_Integration_FileNotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange
	db := setupTestDB(t)
	userServer, fileServer := setupTestHTTPClientsWithErrors(t, false, true)
	defer userServer.Close()
	defer fileServer.Close()

	messageRepo := repositories.NewMessageRepository(db)
	chatRepo := repositories.NewChatRepository(db)
	chatUserRepo := repositories.NewChatUserRepository(db)

	// Создаем чат
	chat := &models.Chat{
		ID:      uuid.New(),
		Name:    "test_SendMessage_FileNotFound",
		IsGroup: true,
	}
	err := chatRepo.CreateChat(chat)
	require.NoError(t, err)

	userClient := http_clients.NewUserClientAdapter()
	fileClient := http_clients.NewFileClientAdapter()

	controller := controllers.NewMessageControllerWithClients(
		messageRepo,
		chatRepo,
		chatUserRepo,
		fileClient,
		userClient,
	)

	senderID := uuid.New()

	createDTO := &dto.CreateMessageDTO{
		Content: "test message",
		FileIDs: []int{999}, // Несуществующий файл
	}

	// Act
	message, err := controller.SendMessage(senderID, chat.ID, createDTO)

	// Assert
	require.Error(t, err)
	assert.Nil(t, message)
}

// TestMessageController_GetChatMessages_Integration тестирует получение сообщений чата с реальной интеграцией FileClient
func TestMessageController_GetChatMessages_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange
	db := setupTestDB(t)
	_, fileServer, _, fileClient := setupTestHTTPClients(t)
	defer fileServer.Close()

	messageRepo := repositories.NewMessageRepository(db)
	chatRepo := repositories.NewChatRepository(db)
	chatUserRepo := repositories.NewChatUserRepository(db)

	// Создаем чат
	chat := &models.Chat{
		ID:      uuid.New(),
		Name:    "test_GetChatMessages_Integration",
		IsGroup: true,
	}
	err := chatRepo.CreateChat(chat)
	require.NoError(t, err)

	// Добавляем пользователя в чат (требуется для внешнего ключа)
	chatRoleRepo := repositories.NewChatRoleRepository(db)
	mainRole, err := chatRoleRepo.GetRoleByName("main")
	require.NoError(t, err)

	senderID := uuid.New()
	chatUser := &models.ChatUser{
		ChatID: chat.ID,
		UserID: senderID,
		RoleID: mainRole.ID,
	}
	err = chatUserRepo.AddUserToChat(chatUser)
	require.NoError(t, err)

	// Создаем несколько сообщений напрямую в БД
	for i := 0; i < 3; i++ {
		message := &models.Message{
			ID:       uuid.New(),
			ChatID:   chat.ID,
			SenderID: &senderID,
			Content:  fmt.Sprintf("test_GetChatMessages_%d", i),
		}
		err = messageRepo.CreateMessage(message)
		require.NoError(t, err)

		// Добавляем файл к сообщению
		messageFile := &models.MessageFile{
			MessageID: message.ID,
			FileID:    1,
		}
		err = messageRepo.CreateMessageFile(messageFile)
		require.NoError(t, err)
	}

	controller := controllers.NewMessageControllerWithClients(
		messageRepo,
		chatRepo,
		chatUserRepo,
		fileClient,
		http_clients.NewUserClientAdapter(),
	)

	// Act
	messages, err := controller.GetChatMessages(chat.ID, 0, 10)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, messages)
	assert.Len(t, *messages, 3)

	for _, msg := range *messages {
		assert.Equal(t, chat.ID, msg.ChatID)
		assert.NotNil(t, msg.Files)
		assert.Len(t, *msg.Files, 1)
	}
}

// TestMessageController_GetChatMessages_Integration_WithPagination тестирует пагинацию сообщений
func TestMessageController_GetChatMessages_Integration_WithPagination(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange
	db := setupTestDB(t)
	_, fileServer, _, fileClient := setupTestHTTPClients(t)
	defer fileServer.Close()

	messageRepo := repositories.NewMessageRepository(db)
	chatRepo := repositories.NewChatRepository(db)
	chatUserRepo := repositories.NewChatUserRepository(db)

	// Создаем чат
	chat := &models.Chat{
		ID:      uuid.New(),
		Name:    "test_GetChatMessages_WithPagination",
		IsGroup: true,
	}
	err := chatRepo.CreateChat(chat)
	require.NoError(t, err)

	// Добавляем пользователя в чат (требуется для внешнего ключа)
	chatRoleRepo := repositories.NewChatRoleRepository(db)
	mainRole, err := chatRoleRepo.GetRoleByName("main")
	require.NoError(t, err)

	senderID := uuid.New()
	chatUser := &models.ChatUser{
		ChatID: chat.ID,
		UserID: senderID,
		RoleID: mainRole.ID,
	}
	err = chatUserRepo.AddUserToChat(chatUser)
	require.NoError(t, err)

	// Создаем 5 сообщений
	for i := 0; i < 5; i++ {
		message := &models.Message{
			ID:       uuid.New(),
			ChatID:   chat.ID,
			SenderID: &senderID,
			Content:  fmt.Sprintf("test message %d", i),
		}
		err = messageRepo.CreateMessage(message)
		require.NoError(t, err)
	}

	controller := controllers.NewMessageControllerWithClients(
		messageRepo,
		chatRepo,
		chatUserRepo,
		fileClient,
		http_clients.NewUserClientAdapter(),
	)

	// Act - получаем первые 2 сообщения
	messages, err := controller.GetChatMessages(chat.ID, 0, 2)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, messages)
	assert.Len(t, *messages, 2)
}

// TestMessageController_GetChatMessages_Integration_ChatNotFound тестирует обработку ошибки, когда чат не найден
func TestMessageController_GetChatMessages_Integration_ChatNotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange
	db := setupTestDB(t)
	messageRepo := repositories.NewMessageRepository(db)
	chatRepo := repositories.NewChatRepository(db)
	chatUserRepo := repositories.NewChatUserRepository(db)

	controller := controllers.NewMessageController(
		messageRepo,
		chatRepo,
		chatUserRepo,
	)

	nonExistentChatID := uuid.New()

	// Act
	messages, err := controller.GetChatMessages(nonExistentChatID, 0, 10)

	// Assert
	require.Error(t, err)
	assert.Nil(t, messages)
}

// TestMessageController_SearchMessages_Integration тестирует поиск сообщений с реальной БД
func TestMessageController_SearchMessages_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange
	db := setupTestDB(t)
	messageRepo := repositories.NewMessageRepository(db)
	chatRepo := repositories.NewChatRepository(db)
	chatUserRepo := repositories.NewChatUserRepository(db)

	// Создаем чат
	chat := &models.Chat{
		ID:      uuid.New(),
		Name:    "test_SearchMessages_Integration",
		IsGroup: true,
	}
	err := chatRepo.CreateChat(chat)
	require.NoError(t, err)

	// Добавляем пользователя в чат
	mainRole, err := repositories.NewChatRoleRepository(db).GetRoleByName("main")
	require.NoError(t, err)

	userID := uuid.New()
	chatUser := &models.ChatUser{
		ChatID: chat.ID,
		UserID: userID,
		RoleID: mainRole.ID,
	}
	err = chatUserRepo.AddUserToChat(chatUser)
	require.NoError(t, err)

	// Добавляем отправителя в чат (требуется для внешнего ключа)
	senderID := uuid.New()
	senderChatUser := &models.ChatUser{
		ChatID: chat.ID,
		UserID: senderID,
		RoleID: mainRole.ID,
	}
	err = chatUserRepo.AddUserToChat(senderChatUser)
	require.NoError(t, err)

	// Создаем сообщения с разным содержимым
	messages := []*models.Message{
		{
			ID:       uuid.New(),
			ChatID:   chat.ID,
			SenderID: &senderID,
			Content:  "test_SearchMessages hello world",
		},
		{
			ID:       uuid.New(),
			ChatID:   chat.ID,
			SenderID: &senderID,
			Content:  "test_SearchMessages hello there",
		},
		{
			ID:       uuid.New(),
			ChatID:   chat.ID,
			SenderID: &senderID,
			Content:  "test_SearchMessages goodbye",
		},
	}

	for _, msg := range messages {
		err = messageRepo.CreateMessage(msg)
		require.NoError(t, err)
	}

	controller := controllers.NewMessageController(
		messageRepo,
		chatRepo,
		chatUserRepo,
	)

	// Act - ищем сообщения со словом "hello"
	var searchResult *ac.GetSearchResponse
	searchResult, err = controller.SearchMessages(userID, chat.ID, "hello", 10, 0)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, searchResult)
	assert.NotNil(t, searchResult.Messages)
	assert.Len(t, *searchResult.Messages, 2) // Должно найти 2 сообщения со словом "hello"
	assert.NotNil(t, searchResult.Total)
	assert.Equal(t, int64(2), *searchResult.Total)
}

// TestMessageController_SearchMessages_Integration_EmptyQuery тестирует обработку пустого запроса
func TestMessageController_SearchMessages_Integration_EmptyQuery(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange
	db := setupTestDB(t)
	messageRepo := repositories.NewMessageRepository(db)
	chatRepo := repositories.NewChatRepository(db)
	chatUserRepo := repositories.NewChatUserRepository(db)

	controller := controllers.NewMessageController(
		messageRepo,
		chatRepo,
		chatUserRepo,
	)

	userID := uuid.New()
	chatID := uuid.New()

	// Act
	searchResult, err := controller.SearchMessages(userID, chatID, "", 10, 0)

	// Assert
	require.Error(t, err)
	assert.Nil(t, searchResult)
}

// TestMessageController_SearchMessages_Integration_ChatNotFound тестирует обработку ошибки, когда чат не найден
func TestMessageController_SearchMessages_Integration_ChatNotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange
	db := setupTestDB(t)
	messageRepo := repositories.NewMessageRepository(db)
	chatRepo := repositories.NewChatRepository(db)
	chatUserRepo := repositories.NewChatUserRepository(db)

	controller := controllers.NewMessageController(
		messageRepo,
		chatRepo,
		chatUserRepo,
	)

	userID := uuid.New()
	nonExistentChatID := uuid.New()

	// Act
	searchResult, err := controller.SearchMessages(userID, nonExistentChatID, "test", 10, 0)

	// Assert
	require.Error(t, err)
	assert.Nil(t, searchResult)
}

// TestMessageController_SearchMessages_Integration_Unauthorized тестирует обработку ошибки, когда пользователь не в чате
func TestMessageController_SearchMessages_Integration_Unauthorized(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange
	db := setupTestDB(t)
	messageRepo := repositories.NewMessageRepository(db)
	chatRepo := repositories.NewChatRepository(db)
	chatUserRepo := repositories.NewChatUserRepository(db)

	// Создаем чат
	chat := &models.Chat{
		ID:      uuid.New(),
		Name:    "test_SearchMessages_Unauthorized",
		IsGroup: true,
	}
	err := chatRepo.CreateChat(chat)
	require.NoError(t, err)

	controller := controllers.NewMessageController(
		messageRepo,
		chatRepo,
		chatUserRepo,
	)

	nonExistentUserID := uuid.New()

	// Act
	searchResult, err := controller.SearchMessages(nonExistentUserID, chat.ID, "test", 10, 0)

	// Assert
	require.Error(t, err)
	assert.Nil(t, searchResult)
}
