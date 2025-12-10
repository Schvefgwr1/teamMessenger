package controllers

import (
	"apiService/internal/controllers"
	"apiService/internal/dto"
	"apiService/internal/services"
	"context"
	"errors"
	"testing"

	ac "common/contracts/api-chat"
	af "common/contracts/api-file"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"mime/multipart"
)

// Тесты для ChatController.GetUserChats

func TestChatController_GetUserChats_Success_FromCache(t *testing.T) {
	// Arrange
	mockChatClient := new(MockChatClient)
	mockFileClient := new(MockFileClient)
	redisClient := setupTestRedis(t)
	defer redisClient.Close()
	cacheService := services.NewCacheService(redisClient)

	controller := controllers.NewChatController(mockChatClient, mockFileClient, cacheService)

	userID := uuid.New()
	cachedChats := []*ac.ChatResponse{
		{ID: uuid.New(), Name: "Test Chat"},
	}

	// Сохраняем чаты в кеш
	cacheService.SetUserChatListCache(context.Background(), userID.String(), cachedChats)

	// Act
	result, err := controller.GetUserChats(userID)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 1)

	mockChatClient.AssertNotCalled(t, "GetUserChats", mock.Anything)
}

func TestChatController_GetUserChats_Success_FromService(t *testing.T) {
	// Arrange
	mockChatClient := new(MockChatClient)
	mockFileClient := new(MockFileClient)
	redisClient := setupTestRedis(t)
	defer redisClient.Close()
	cacheService := services.NewCacheService(redisClient)

	controller := controllers.NewChatController(mockChatClient, mockFileClient, cacheService)

	userID := uuid.New()
	expectedChats := []*ac.ChatResponse{
		{ID: uuid.New(), Name: "Test Chat"},
	}

	mockChatClient.On("GetUserChats", userID).Return(expectedChats, nil)

	// Act
	result, err := controller.GetUserChats(userID)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 1)

	mockChatClient.AssertExpectations(t)
}

func TestChatController_GetUserChats_ServiceError(t *testing.T) {
	// Arrange
	mockChatClient := new(MockChatClient)
	mockFileClient := new(MockFileClient)
	redisClient := setupTestRedis(t)
	defer redisClient.Close()
	cacheService := services.NewCacheService(redisClient)

	controller := controllers.NewChatController(mockChatClient, mockFileClient, cacheService)

	userID := uuid.New()
	serviceError := errors.New("service error")

	mockChatClient.On("GetUserChats", userID).Return(nil, serviceError)

	// Act
	result, err := controller.GetUserChats(userID)

	// Assert
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, serviceError, err)

	mockChatClient.AssertExpectations(t)
}

// Тесты для ChatController.CreateChat

func TestChatController_CreateChat_Success_WithoutAvatar(t *testing.T) {
	// Arrange
	mockChatClient := new(MockChatClient)
	mockFileClient := new(MockFileClient)
	redisClient := setupTestRedis(t)
	defer redisClient.Close()
	cacheService := services.NewCacheService(redisClient)

	controller := controllers.NewChatController(mockChatClient, mockFileClient, cacheService)

	ownerID := uuid.New()
	userIDs := []uuid.UUID{uuid.New()}
	req := &dto.CreateChatRequestGateway{
		Name:        "Test Chat",
		Description: stringPtr("Test Description"),
		Avatar:      nil,
	}

	chatID := uuid.New()
	serviceResp := &ac.CreateChatServiceResponse{
		ChatID: chatID,
	}

	mockChatClient.On("CreateChat", mock.Anything).Return(serviceResp, nil)

	// Act
	result, err := controller.CreateChat(req, ownerID, userIDs)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, chatID, result.ID)
	assert.Equal(t, req.Name, result.Name)

	mockChatClient.AssertExpectations(t)
}

func TestChatController_CreateChat_Success_WithAvatar(t *testing.T) {
	// Arrange
	mockChatClient := new(MockChatClient)
	mockFileClient := new(MockFileClient)
	redisClient := setupTestRedis(t)
	defer redisClient.Close()
	cacheService := services.NewCacheService(redisClient)

	controller := controllers.NewChatController(mockChatClient, mockFileClient, cacheService)

	ownerID := uuid.New()
	userIDs := []uuid.UUID{uuid.New()}
	req := &dto.CreateChatRequestGateway{
		Name:        "Test Chat",
		Description: stringPtr("Test Description"),
		Avatar:      &multipart.FileHeader{Filename: "avatar.jpg"},
	}

	uploadedFile := &af.FileUploadResponse{
		ID: intPtr(1),
	}

	chatID := uuid.New()
	serviceResp := &ac.CreateChatServiceResponse{
		ChatID: chatID,
	}

	mockFileClient.On("UploadFile", req.Avatar).Return(uploadedFile, nil)
	mockChatClient.On("CreateChat", mock.Anything).Return(serviceResp, nil)

	// Act
	result, err := controller.CreateChat(req, ownerID, userIDs)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, chatID, result.ID)
	assert.Equal(t, req.Name, result.Name)
	assert.NotNil(t, result.AvatarFileID)

	mockFileClient.AssertExpectations(t)
	mockChatClient.AssertExpectations(t)
}

func TestChatController_CreateChat_FileUploadError(t *testing.T) {
	// Arrange
	mockChatClient := new(MockChatClient)
	mockFileClient := new(MockFileClient)
	redisClient := setupTestRedis(t)
	defer redisClient.Close()
	cacheService := services.NewCacheService(redisClient)

	controller := controllers.NewChatController(mockChatClient, mockFileClient, cacheService)

	ownerID := uuid.New()
	userIDs := []uuid.UUID{uuid.New()}
	req := &dto.CreateChatRequestGateway{
		Name:        "Test Chat",
		Description: stringPtr("Test Description"),
		Avatar:      &multipart.FileHeader{Filename: "avatar.jpg"},
	}

	uploadError := errors.New("upload error")

	mockFileClient.On("UploadFile", req.Avatar).Return(nil, uploadError)

	// Act
	result, err := controller.CreateChat(req, ownerID, userIDs)

	// Assert
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, uploadError, err)

	mockFileClient.AssertExpectations(t)
	mockChatClient.AssertNotCalled(t, "CreateChat", mock.Anything)
}

func TestChatController_CreateChat_ServiceError(t *testing.T) {
	// Arrange
	mockChatClient := new(MockChatClient)
	mockFileClient := new(MockFileClient)
	redisClient := setupTestRedis(t)
	defer redisClient.Close()
	cacheService := services.NewCacheService(redisClient)

	controller := controllers.NewChatController(mockChatClient, mockFileClient, cacheService)

	ownerID := uuid.New()
	userIDs := []uuid.UUID{uuid.New()}
	req := &dto.CreateChatRequestGateway{
		Name:        "Test Chat",
		Description: stringPtr("Test Description"),
		Avatar:      nil,
	}

	serviceError := errors.New("service error")

	mockChatClient.On("CreateChat", mock.Anything).Return(nil, serviceError)

	// Act
	result, err := controller.CreateChat(req, ownerID, userIDs)

	// Assert
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "error of chat client")

	mockChatClient.AssertExpectations(t)
}

// Тесты для ChatController.SendMessage

func TestChatController_SendMessage_Success(t *testing.T) {
	// Arrange
	mockChatClient := new(MockChatClient)
	mockFileClient := new(MockFileClient)
	redisClient := setupTestRedis(t)
	defer redisClient.Close()
	cacheService := services.NewCacheService(redisClient)

	controller := controllers.NewChatController(mockChatClient, mockFileClient, cacheService)

	chatID := uuid.New()
	senderID := uuid.New()
	req := &dto.SendMessageRequestGateway{
		Content: "Test message",
		Files:   nil,
	}

	expectedMessage := &ac.MessageResponse{
		ID:       uuid.New(),
		ChatID:   chatID,
		SenderID: senderID,
		Content:  "Test message",
	}

	mockChatClient.On("SendMessage", chatID, senderID, mock.Anything).Return(expectedMessage, nil)

	// Act
	result, err := controller.SendMessage(chatID, senderID, req)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedMessage.Content, result.Content)

	mockChatClient.AssertExpectations(t)
}

func TestChatController_SendMessage_Success_WithFiles(t *testing.T) {
	// Arrange
	mockChatClient := new(MockChatClient)
	mockFileClient := new(MockFileClient)
	redisClient := setupTestRedis(t)
	defer redisClient.Close()
	cacheService := services.NewCacheService(redisClient)

	controller := controllers.NewChatController(mockChatClient, mockFileClient, cacheService)

	chatID := uuid.New()
	senderID := uuid.New()
	fileHeader1 := &multipart.FileHeader{Filename: "file1.txt"}
	fileHeader2 := &multipart.FileHeader{Filename: "file2.txt"}
	req := &dto.SendMessageRequestGateway{
		Content: "Test message with files",
		Files:   []*multipart.FileHeader{fileHeader1, fileHeader2},
	}

	uploadedFile1 := &af.FileUploadResponse{ID: intPtr(1)}
	uploadedFile2 := &af.FileUploadResponse{ID: intPtr(2)}

	expectedMessage := &ac.MessageResponse{
		ID:       uuid.New(),
		ChatID:   chatID,
		SenderID: senderID,
		Content:  "Test message with files",
	}

	mockFileClient.On("UploadFile", fileHeader1).Return(uploadedFile1, nil)
	mockFileClient.On("UploadFile", fileHeader2).Return(uploadedFile2, nil)
	mockChatClient.On("SendMessage", chatID, senderID, mock.Anything).Return(expectedMessage, nil)

	// Act
	result, err := controller.SendMessage(chatID, senderID, req)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedMessage.Content, result.Content)

	mockFileClient.AssertExpectations(t)
	mockChatClient.AssertExpectations(t)
}

func TestChatController_SendMessage_PartialFileUploadError(t *testing.T) {
	// Arrange
	mockChatClient := new(MockChatClient)
	mockFileClient := new(MockFileClient)
	redisClient := setupTestRedis(t)
	defer redisClient.Close()
	cacheService := services.NewCacheService(redisClient)

	controller := controllers.NewChatController(mockChatClient, mockFileClient, cacheService)

	chatID := uuid.New()
	senderID := uuid.New()
	fileHeader1 := &multipart.FileHeader{Filename: "file1.txt"}
	fileHeader2 := &multipart.FileHeader{Filename: "file2.txt"}
	req := &dto.SendMessageRequestGateway{
		Content: "Test message with files",
		Files:   []*multipart.FileHeader{fileHeader1, fileHeader2},
	}

	uploadedFile1 := &af.FileUploadResponse{ID: intPtr(1)}
	uploadError := errors.New("upload error")

	expectedMessage := &ac.MessageResponse{
		ID:       uuid.New(),
		ChatID:   chatID,
		SenderID: senderID,
		Content:  "Test message with files",
	}

	mockFileClient.On("UploadFile", fileHeader1).Return(uploadedFile1, nil)
	mockFileClient.On("UploadFile", fileHeader2).Return(nil, uploadError)
	mockChatClient.On("SendMessage", chatID, senderID, mock.Anything).Return(expectedMessage, nil)

	// Act
	result, err := controller.SendMessage(chatID, senderID, req)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	// Сообщение должно быть отправлено даже если один файл не загрузился

	mockFileClient.AssertExpectations(t)
	mockChatClient.AssertExpectations(t)
}

func TestChatController_SendMessage_ServiceError(t *testing.T) {
	// Arrange
	mockChatClient := new(MockChatClient)
	mockFileClient := new(MockFileClient)
	redisClient := setupTestRedis(t)
	defer redisClient.Close()
	cacheService := services.NewCacheService(redisClient)

	controller := controllers.NewChatController(mockChatClient, mockFileClient, cacheService)

	chatID := uuid.New()
	senderID := uuid.New()
	req := &dto.SendMessageRequestGateway{
		Content: "Test message",
		Files:   nil,
	}

	serviceError := errors.New("service error")

	mockChatClient.On("SendMessage", chatID, senderID, mock.Anything).Return(nil, serviceError)

	// Act
	result, err := controller.SendMessage(chatID, senderID, req)

	// Assert
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, serviceError, err)

	mockChatClient.AssertExpectations(t)
}

// Тесты для ChatController.GetChatMessages

func TestChatController_GetChatMessages_Success_FromCache(t *testing.T) {
	// Arrange
	mockChatClient := new(MockChatClient)
	mockFileClient := new(MockFileClient)
	redisClient := setupTestRedis(t)
	defer redisClient.Close()
	cacheService := services.NewCacheService(redisClient)

	controller := controllers.NewChatController(mockChatClient, mockFileClient, cacheService)

	chatID := uuid.New()
	userID := uuid.New()
	cachedMessages := []*ac.GetChatMessage{
		{ID: uuid.New(), Content: "Test message"},
	}

	// Сохраняем сообщения в кеш
	cacheService.SetChatMessagesCache(context.Background(), chatID.String(), cachedMessages)

	// Act
	result, err := controller.GetChatMessages(chatID, userID, 0, 10)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 1)

	mockChatClient.AssertNotCalled(t, "GetChatMessages", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

func TestChatController_GetChatMessages_Success_FromService(t *testing.T) {
	// Arrange
	mockChatClient := new(MockChatClient)
	mockFileClient := new(MockFileClient)
	redisClient := setupTestRedis(t)
	defer redisClient.Close()
	cacheService := services.NewCacheService(redisClient)

	controller := controllers.NewChatController(mockChatClient, mockFileClient, cacheService)

	chatID := uuid.New()
	userID := uuid.New()
	expectedMessages := []*ac.GetChatMessage{
		{ID: uuid.New(), Content: "Test message"},
	}

	mockChatClient.On("GetChatMessages", chatID, userID, 0, 20).Return(expectedMessages, nil)

	// Act
	result, err := controller.GetChatMessages(chatID, userID, 0, 10)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 1)

	mockChatClient.AssertExpectations(t)
}

func TestChatController_GetChatMessages_Success_FromCache_LimitLessThanCached(t *testing.T) {
	// Arrange
	mockChatClient := new(MockChatClient)
	mockFileClient := new(MockFileClient)
	redisClient := setupTestRedis(t)
	defer redisClient.Close()
	cacheService := services.NewCacheService(redisClient)

	controller := controllers.NewChatController(mockChatClient, mockFileClient, cacheService)

	chatID := uuid.New()
	userID := uuid.New()
	cachedMessages := []*ac.GetChatMessage{
		{ID: uuid.New(), Content: "Message 1"},
		{ID: uuid.New(), Content: "Message 2"},
		{ID: uuid.New(), Content: "Message 3"},
	}

	// Сохраняем сообщения в кеш
	cacheService.SetChatMessagesCache(context.Background(), chatID.String(), cachedMessages)

	// Act - запрашиваем только 2 сообщения из 3
	result, err := controller.GetChatMessages(chatID, userID, 0, 2)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	assert.Equal(t, cachedMessages[0].Content, result[0].Content)

	mockChatClient.AssertNotCalled(t, "GetChatMessages", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

func TestChatController_GetChatMessages_Success_FromService_LimitLessThanMessages(t *testing.T) {
	// Arrange
	mockChatClient := new(MockChatClient)
	mockFileClient := new(MockFileClient)
	redisClient := setupTestRedis(t)
	defer redisClient.Close()
	cacheService := services.NewCacheService(redisClient)

	controller := controllers.NewChatController(mockChatClient, mockFileClient, cacheService)

	chatID := uuid.New()
	userID := uuid.New()
	expectedMessages := []*ac.GetChatMessage{
		{ID: uuid.New(), Content: "Message 1"},
		{ID: uuid.New(), Content: "Message 2"},
		{ID: uuid.New(), Content: "Message 3"},
	}

	mockChatClient.On("GetChatMessages", chatID, userID, 0, 20).Return(expectedMessages, nil)

	// Act - запрашиваем только 2 сообщения из 3
	result, err := controller.GetChatMessages(chatID, userID, 0, 2)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	assert.Equal(t, expectedMessages[0].Content, result[0].Content)

	mockChatClient.AssertExpectations(t)
}

func TestChatController_GetChatMessages_Success_WithOffset(t *testing.T) {
	// Arrange
	mockChatClient := new(MockChatClient)
	mockFileClient := new(MockFileClient)
	redisClient := setupTestRedis(t)
	defer redisClient.Close()
	cacheService := services.NewCacheService(redisClient)

	controller := controllers.NewChatController(mockChatClient, mockFileClient, cacheService)

	chatID := uuid.New()
	userID := uuid.New()
	expectedMessages := []*ac.GetChatMessage{
		{ID: uuid.New(), Content: "Message 3"},
	}

	// При offset > 0 кеш не используется
	mockChatClient.On("GetChatMessages", chatID, userID, 10, 20).Return(expectedMessages, nil)

	// Act
	result, err := controller.GetChatMessages(chatID, userID, 10, 20)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 1)

	mockChatClient.AssertExpectations(t)
}

func TestChatController_GetChatMessages_Success_LimitGreaterThan20(t *testing.T) {
	// Arrange
	mockChatClient := new(MockChatClient)
	mockFileClient := new(MockFileClient)
	redisClient := setupTestRedis(t)
	defer redisClient.Close()
	cacheService := services.NewCacheService(redisClient)

	controller := controllers.NewChatController(mockChatClient, mockFileClient, cacheService)

	chatID := uuid.New()
	userID := uuid.New()
	expectedMessages := []*ac.GetChatMessage{
		{ID: uuid.New(), Content: "Message 1"},
	}

	// При limit > 20 кеш не используется
	mockChatClient.On("GetChatMessages", chatID, userID, 0, 30).Return(expectedMessages, nil)

	// Act
	result, err := controller.GetChatMessages(chatID, userID, 0, 30)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)

	mockChatClient.AssertExpectations(t)
}

// Тесты для ChatController.GetMyRoleInChat

func TestChatController_GetMyRoleInChat_Success_FromCache(t *testing.T) {
	// Arrange
	mockChatClient := new(MockChatClient)
	mockFileClient := new(MockFileClient)
	redisClient := setupTestRedis(t)
	defer redisClient.Close()
	cacheService := services.NewCacheService(redisClient)

	controller := controllers.NewChatController(mockChatClient, mockFileClient, cacheService)

	chatID := uuid.New()
	userID := uuid.New()
	cachedRole := &ac.MyRoleResponse{
		RoleID:   1,
		RoleName: "admin",
	}

	// Сохраняем роль в кеш
	cacheService.SetChatUserRoleCache(context.Background(), chatID.String(), userID.String(), cachedRole)

	// Act
	result, err := controller.GetMyRoleInChat(chatID, userID)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, cachedRole.RoleID, result.RoleID)

	mockChatClient.AssertNotCalled(t, "GetMyRoleInChat", mock.Anything, mock.Anything)
}

func TestChatController_GetMyRoleInChat_Success_FromService(t *testing.T) {
	// Arrange
	mockChatClient := new(MockChatClient)
	mockFileClient := new(MockFileClient)
	redisClient := setupTestRedis(t)
	defer redisClient.Close()
	cacheService := services.NewCacheService(redisClient)

	controller := controllers.NewChatController(mockChatClient, mockFileClient, cacheService)

	chatID := uuid.New()
	userID := uuid.New()
	expectedRole := &ac.MyRoleResponse{
		RoleID:   1,
		RoleName: "admin",
	}

	mockChatClient.On("GetMyRoleInChat", chatID, userID).Return(expectedRole, nil)

	// Act
	result, err := controller.GetMyRoleInChat(chatID, userID)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedRole.RoleID, result.RoleID)

	mockChatClient.AssertExpectations(t)
}

// Тесты для ChatController.DeleteChat

func TestChatController_DeleteChat_Success(t *testing.T) {
	// Arrange
	mockChatClient := new(MockChatClient)
	mockFileClient := new(MockFileClient)
	redisClient := setupTestRedis(t)
	defer redisClient.Close()
	cacheService := services.NewCacheService(redisClient)

	controller := controllers.NewChatController(mockChatClient, mockFileClient, cacheService)

	chatID := uuid.New()
	userID := uuid.New()

	mockChatClient.On("DeleteChat", chatID, userID).Return(nil)

	// Act
	err := controller.DeleteChat(chatID, userID)

	// Assert
	require.NoError(t, err)

	mockChatClient.AssertExpectations(t)
}

// Тесты для ChatController.GetChatMembers

func TestChatController_GetChatMembers_Success(t *testing.T) {
	// Arrange
	mockChatClient := new(MockChatClient)
	mockFileClient := new(MockFileClient)
	redisClient := setupTestRedis(t)
	defer redisClient.Close()
	cacheService := services.NewCacheService(redisClient)

	controller := controllers.NewChatController(mockChatClient, mockFileClient, cacheService)

	chatID := uuid.New()
	expectedMembers := []*ac.ChatMember{
		{UserID: uuid.New().String(), RoleID: 1, RoleName: "user"},
	}

	mockChatClient.On("GetChatMembers", chatID).Return(expectedMembers, nil)

	// Act
	result, err := controller.GetChatMembers(chatID)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 1)

	mockChatClient.AssertExpectations(t)
}

// Тесты для ChatController.SearchMessages

func TestChatController_SearchMessages_Success_FromCache(t *testing.T) {
	// Arrange
	mockChatClient := new(MockChatClient)
	mockFileClient := new(MockFileClient)
	redisClient := setupTestRedis(t)
	defer redisClient.Close()
	cacheService := services.NewCacheService(redisClient)

	controller := controllers.NewChatController(mockChatClient, mockFileClient, cacheService)

	userID := uuid.New()
	chatID := uuid.New()
	query := "test"
	queryHash := "74657374" // hex hash of "test"

	cachedResult := &ac.GetSearchResponse{
		Messages: &[]ac.GetChatMessage{
			{ID: uuid.New(), Content: "test message"},
		},
		Total: func() *int64 { v := int64(1); return &v }(),
	}

	// Сохраняем результат поиска в кеш
	cacheService.SetSearchCache(context.Background(), chatID.String(), queryHash, cachedResult)

	// Act
	result, err := controller.SearchMessages(userID, chatID, query, 0, 10)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotNil(t, result.Messages)

	mockChatClient.AssertNotCalled(t, "SearchMessages", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

func TestChatController_SearchMessages_Success_FromService(t *testing.T) {
	// Arrange
	mockChatClient := new(MockChatClient)
	mockFileClient := new(MockFileClient)
	redisClient := setupTestRedis(t)
	defer redisClient.Close()
	cacheService := services.NewCacheService(redisClient)

	controller := controllers.NewChatController(mockChatClient, mockFileClient, cacheService)

	userID := uuid.New()
	chatID := uuid.New()
	query := "test"
	expectedResult := &ac.GetSearchResponse{
		Messages: &[]ac.GetChatMessage{
			{ID: uuid.New(), Content: "test message"},
		},
		Total: func() *int64 { v := int64(1); return &v }(),
	}

	mockChatClient.On("SearchMessages", userID, chatID, query, 0, 10).Return(expectedResult, nil)

	// Act
	result, err := controller.SearchMessages(userID, chatID, query, 0, 10)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotNil(t, result.Messages)

	mockChatClient.AssertExpectations(t)
}

func TestChatController_SearchMessages_ServiceError(t *testing.T) {
	// Arrange
	mockChatClient := new(MockChatClient)
	mockFileClient := new(MockFileClient)
	redisClient := setupTestRedis(t)
	defer redisClient.Close()
	cacheService := services.NewCacheService(redisClient)

	controller := controllers.NewChatController(mockChatClient, mockFileClient, cacheService)

	userID := uuid.New()
	chatID := uuid.New()
	query := "test"
	serviceError := errors.New("service error")

	mockChatClient.On("SearchMessages", userID, chatID, query, 0, 10).Return(nil, serviceError)

	// Act
	result, err := controller.SearchMessages(userID, chatID, query, 0, 10)

	// Assert
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, serviceError, err)

	mockChatClient.AssertExpectations(t)
}

func TestChatController_SearchMessages_WithOffset(t *testing.T) {
	// Arrange
	mockChatClient := new(MockChatClient)
	mockFileClient := new(MockFileClient)
	redisClient := setupTestRedis(t)
	defer redisClient.Close()
	cacheService := services.NewCacheService(redisClient)

	controller := controllers.NewChatController(mockChatClient, mockFileClient, cacheService)

	userID := uuid.New()
	chatID := uuid.New()
	query := "test"
	expectedResult := &ac.GetSearchResponse{
		Messages: &[]ac.GetChatMessage{
			{ID: uuid.New(), Content: "test message"},
		},
		Total: func() *int64 { v := int64(1); return &v }(),
	}

	// При offset > 0 кеш не используется
	mockChatClient.On("SearchMessages", userID, chatID, query, 10, 20).Return(expectedResult, nil)

	// Act
	result, err := controller.SearchMessages(userID, chatID, query, 10, 20)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)

	mockChatClient.AssertExpectations(t)
}

// Тесты для ChatController.UpdateChat

func TestChatController_UpdateChat_Success_WithoutAvatar(t *testing.T) {
	// Arrange
	mockChatClient := new(MockChatClient)
	mockFileClient := new(MockFileClient)
	redisClient := setupTestRedis(t)
	defer redisClient.Close()
	cacheService := services.NewCacheService(redisClient)

	controller := controllers.NewChatController(mockChatClient, mockFileClient, cacheService)

	chatID := uuid.New()
	userID := uuid.New()
	req := &dto.UpdateChatRequestGateway{
		Name:        stringPtr("Updated Chat"),
		Description: stringPtr("Updated Description"),
		Avatar:      nil,
	}

	updateReq := &ac.UpdateChatRequest{
		Name:         stringPtr("Updated Chat"),
		Description:  stringPtr("Updated Description"),
		AvatarFileID: nil,
	}

	expectedResponse := &ac.UpdateChatResponse{
		Chat: ac.ChatResponse{
			ID:   chatID,
			Name: "Updated Chat",
		},
	}

	mockChatClient.On("UpdateChat", chatID, mock.Anything, userID).Return(expectedResponse, nil)

	// Act
	result, err := controller.UpdateChat(chatID, req, updateReq, userID)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedResponse.Chat.Name, result.Chat.Name)

	mockChatClient.AssertExpectations(t)
}

func TestChatController_UpdateChat_Success_WithAvatar(t *testing.T) {
	// Arrange
	mockChatClient := new(MockChatClient)
	mockFileClient := new(MockFileClient)
	redisClient := setupTestRedis(t)
	defer redisClient.Close()
	cacheService := services.NewCacheService(redisClient)

	controller := controllers.NewChatController(mockChatClient, mockFileClient, cacheService)

	chatID := uuid.New()
	userID := uuid.New()
	req := &dto.UpdateChatRequestGateway{
		Name:        stringPtr("Updated Chat"),
		Description: stringPtr("Updated Description"),
		Avatar:      &multipart.FileHeader{Filename: "avatar.jpg"},
	}

	updateReq := &ac.UpdateChatRequest{
		Name:         stringPtr("Updated Chat"),
		Description:  stringPtr("Updated Description"),
		AvatarFileID: nil,
	}

	uploadedFile := &af.FileUploadResponse{
		ID: intPtr(1),
	}

	expectedResponse := &ac.UpdateChatResponse{
		Chat: ac.ChatResponse{
			ID:   chatID,
			Name: "Updated Chat",
		},
	}

	mockFileClient.On("UploadFile", req.Avatar).Return(uploadedFile, nil)
	mockChatClient.On("UpdateChat", chatID, mock.Anything, userID).Return(expectedResponse, nil)

	// Act
	result, err := controller.UpdateChat(chatID, req, updateReq, userID)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)

	mockFileClient.AssertExpectations(t)
	mockChatClient.AssertExpectations(t)
}

func TestChatController_UpdateChat_FileUploadError(t *testing.T) {
	// Arrange
	mockChatClient := new(MockChatClient)
	mockFileClient := new(MockFileClient)
	redisClient := setupTestRedis(t)
	defer redisClient.Close()
	cacheService := services.NewCacheService(redisClient)

	controller := controllers.NewChatController(mockChatClient, mockFileClient, cacheService)

	chatID := uuid.New()
	userID := uuid.New()
	req := &dto.UpdateChatRequestGateway{
		Name:        stringPtr("Updated Chat"),
		Description: stringPtr("Updated Description"),
		Avatar:      &multipart.FileHeader{Filename: "avatar.jpg"},
	}

	updateReq := &ac.UpdateChatRequest{
		Name:         stringPtr("Updated Chat"),
		Description:  stringPtr("Updated Description"),
		AvatarFileID: nil,
	}

	uploadError := errors.New("upload error")

	mockFileClient.On("UploadFile", req.Avatar).Return(nil, uploadError)

	// Act
	result, err := controller.UpdateChat(chatID, req, updateReq, userID)

	// Assert
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, uploadError, err)

	mockFileClient.AssertExpectations(t)
	mockChatClient.AssertNotCalled(t, "UpdateChat", mock.Anything, mock.Anything, mock.Anything)
}

func TestChatController_UpdateChat_ServiceError(t *testing.T) {
	// Arrange
	mockChatClient := new(MockChatClient)
	mockFileClient := new(MockFileClient)
	redisClient := setupTestRedis(t)
	defer redisClient.Close()
	cacheService := services.NewCacheService(redisClient)

	controller := controllers.NewChatController(mockChatClient, mockFileClient, cacheService)

	chatID := uuid.New()
	userID := uuid.New()
	req := &dto.UpdateChatRequestGateway{
		Name:        stringPtr("Updated Chat"),
		Description: stringPtr("Updated Description"),
		Avatar:      nil,
	}

	updateReq := &ac.UpdateChatRequest{
		Name:         stringPtr("Updated Chat"),
		Description:  stringPtr("Updated Description"),
		AvatarFileID: nil,
	}

	serviceError := errors.New("service error")

	mockChatClient.On("UpdateChat", chatID, mock.Anything, userID).Return(nil, serviceError)

	// Act
	result, err := controller.UpdateChat(chatID, req, updateReq, userID)

	// Assert
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, serviceError, err)

	mockChatClient.AssertExpectations(t)
}

// Тесты для ChatController.BanUser

func TestChatController_BanUser_Success(t *testing.T) {
	// Arrange
	mockChatClient := new(MockChatClient)
	mockFileClient := new(MockFileClient)
	redisClient := setupTestRedis(t)
	defer redisClient.Close()
	cacheService := services.NewCacheService(redisClient)

	controller := controllers.NewChatController(mockChatClient, mockFileClient, cacheService)

	chatID := uuid.New()
	userID := uuid.New()
	ownerID := uuid.New()

	mockChatClient.On("BanUser", chatID, userID, ownerID).Return(nil)

	// Act
	err := controller.BanUser(chatID, userID, ownerID)

	// Assert
	require.NoError(t, err)

	mockChatClient.AssertExpectations(t)
}

func TestChatController_BanUser_ServiceError(t *testing.T) {
	// Arrange
	mockChatClient := new(MockChatClient)
	mockFileClient := new(MockFileClient)
	redisClient := setupTestRedis(t)
	defer redisClient.Close()
	cacheService := services.NewCacheService(redisClient)

	controller := controllers.NewChatController(mockChatClient, mockFileClient, cacheService)

	chatID := uuid.New()
	userID := uuid.New()
	ownerID := uuid.New()
	serviceError := errors.New("service error")

	mockChatClient.On("BanUser", chatID, userID, ownerID).Return(serviceError)

	// Act
	err := controller.BanUser(chatID, userID, ownerID)

	// Assert
	require.Error(t, err)
	assert.Equal(t, serviceError, err)

	mockChatClient.AssertExpectations(t)
}

// Тесты для ChatController.ChangeUserRole

func TestChatController_ChangeUserRole_Success(t *testing.T) {
	// Arrange
	mockChatClient := new(MockChatClient)
	mockFileClient := new(MockFileClient)
	redisClient := setupTestRedis(t)
	defer redisClient.Close()
	cacheService := services.NewCacheService(redisClient)

	controller := controllers.NewChatController(mockChatClient, mockFileClient, cacheService)

	chatID := uuid.New()
	ownerID := uuid.New()
	userID := uuid.New()
	changeRoleReq := &ac.ChangeRoleRequest{
		UserID: userID,
		RoleID: 2,
	}

	mockChatClient.On("ChangeUserRole", chatID, ownerID, changeRoleReq).Return(nil)

	// Act
	err := controller.ChangeUserRole(chatID, ownerID, changeRoleReq)

	// Assert
	require.NoError(t, err)

	mockChatClient.AssertExpectations(t)
}

func TestChatController_ChangeUserRole_ServiceError(t *testing.T) {
	// Arrange
	mockChatClient := new(MockChatClient)
	mockFileClient := new(MockFileClient)
	redisClient := setupTestRedis(t)
	defer redisClient.Close()
	cacheService := services.NewCacheService(redisClient)

	controller := controllers.NewChatController(mockChatClient, mockFileClient, cacheService)

	chatID := uuid.New()
	ownerID := uuid.New()
	userID := uuid.New()
	changeRoleReq := &ac.ChangeRoleRequest{
		UserID: userID,
		RoleID: 2,
	}

	serviceError := errors.New("service error")

	mockChatClient.On("ChangeUserRole", chatID, ownerID, changeRoleReq).Return(serviceError)

	// Act
	err := controller.ChangeUserRole(chatID, ownerID, changeRoleReq)

	// Assert
	require.Error(t, err)
	assert.Equal(t, serviceError, err)

	mockChatClient.AssertExpectations(t)
}
