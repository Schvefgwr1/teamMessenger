package controllers

import (
	"errors"
	"testing"

	"chatService/internal/controllers"
	"chatService/internal/custom_errors"
	"chatService/internal/handlers/dto"
	"chatService/internal/models"
	cuc "common/contracts/user-contracts"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Тесты для MessageController.SendMessage

func TestMessageController_SendMessage_Success(t *testing.T) {
	t.Parallel()
	// Arrange
	mockMsgRepo := new(MockMessageRepository)
	mockChatRepo := new(MockChatRepository)
	mockChatUserRepo := new(MockChatUserRepository)
	mockFileClient := new(MockFileClient)
	mockUserClient := new(MockUserClient)

	chatID := uuid.New()
	senderID := uuid.New()
	fileIDs := []int{1, 2}
	createDTO := &dto.CreateMessageDTO{
		Content: "hello",
		FileIDs: fileIDs,
	}

	mockChatRepo.On("GetChatByID", chatID).Return(createTestChat(), nil)
	mockUserClient.On("GetUserByID", &senderID).Return(createTestUserResponse(), nil)
	mockFileClient.On("GetFileByID", 1).Return(createTestFile(), nil)
	mockFileClient.On("GetFileByID", 2).Return(createTestFile(), nil)

	mockMsgRepo.On("CreateMessage", mock.AnythingOfType("*models.Message")).Return(nil)
	mockMsgRepo.On("CreateMessageFile", mock.AnythingOfType("*models.MessageFile")).Return(nil)

	createdMsg := createTestMessage()
	createdMsg.ChatID = chatID
	createdMsg.SenderID = &senderID
	mockMsgRepo.On("GetMessageWithFile", mock.Anything).Return(createdMsg, nil)

	controller := controllers.NewMessageControllerWithClients(
		mockMsgRepo, mockChatRepo, mockChatUserRepo, mockFileClient, mockUserClient,
	)

	// Act
	result, err := controller.SendMessage(senderID, chatID, createDTO)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, chatID, result.ChatID)
	assert.Equal(t, senderID, *result.SenderID)

	mockChatRepo.AssertExpectations(t)
	mockUserClient.AssertExpectations(t)
	mockFileClient.AssertExpectations(t)
	mockMsgRepo.AssertExpectations(t)
}

func TestMessageController_SendMessage_WithMultipleFiles(t *testing.T) {
	t.Parallel()
	mockMsgRepo := new(MockMessageRepository)
	mockChatRepo := new(MockChatRepository)
	mockChatUserRepo := new(MockChatUserRepository)
	mockFileClient := new(MockFileClient)
	mockUserClient := new(MockUserClient)

	senderID := uuid.New()
	chatID := uuid.New()
	file1 := createTestFile()
	file1.ID = 1
	file2 := createTestFile()
	file2.ID = 2
	file3 := createTestFile()
	file3.ID = 3

	msg := createTestMessage()
	msg.ChatID = chatID
	msg.SenderID = &senderID
	msg.Files = []models.MessageFile{
		{MessageID: msg.ID, FileID: 1},
		{MessageID: msg.ID, FileID: 2},
		{MessageID: msg.ID, FileID: 3},
	}

	mockChatRepo.On("GetChatByID", chatID).Return(createTestChat(), nil)
	mockUserClient.On("GetUserByID", &senderID).Return(createTestUserResponse(), nil)
	mockFileClient.On("GetFileByID", 1).Return(file1, nil)
	mockFileClient.On("GetFileByID", 2).Return(file2, nil)
	mockFileClient.On("GetFileByID", 3).Return(file3, nil)
	mockMsgRepo.On("CreateMessage", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		m := args.Get(0).(*models.Message)
		m.ID = msg.ID
	})
	mockMsgRepo.On("CreateMessageFile", mock.MatchedBy(func(mf *models.MessageFile) bool {
		return mf.FileID == 1 || mf.FileID == 2 || mf.FileID == 3
	})).Return(nil).Times(3)
	mockMsgRepo.On("GetMessageWithFile", msg.ID).Return(msg, nil)

	controller := controllers.NewMessageControllerWithClients(
		mockMsgRepo, mockChatRepo, mockChatUserRepo, mockFileClient, mockUserClient,
	)

	result, err := controller.SendMessage(senderID, chatID, &dto.CreateMessageDTO{
		Content: "Test message",
		FileIDs: []int{1, 2, 3},
	})

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, msg.ID, result.ID)

	mockChatRepo.AssertExpectations(t)
	mockUserClient.AssertExpectations(t)
	mockFileClient.AssertExpectations(t)
	mockMsgRepo.AssertExpectations(t)
}

func TestMessageController_SendMessage_WithoutFiles(t *testing.T) {
	t.Parallel()
	mockMsgRepo := new(MockMessageRepository)
	mockChatRepo := new(MockChatRepository)
	mockChatUserRepo := new(MockChatUserRepository)
	mockFileClient := new(MockFileClient)
	mockUserClient := new(MockUserClient)

	senderID := uuid.New()
	chatID := uuid.New()
	msg := createTestMessage()
	msg.ChatID = chatID
	msg.SenderID = &senderID
	msg.Files = []models.MessageFile{}

	mockChatRepo.On("GetChatByID", chatID).Return(createTestChat(), nil)
	mockUserClient.On("GetUserByID", &senderID).Return(createTestUserResponse(), nil)
	mockMsgRepo.On("CreateMessage", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		m := args.Get(0).(*models.Message)
		m.ID = msg.ID
	})
	mockMsgRepo.On("GetMessageWithFile", msg.ID).Return(msg, nil)

	controller := controllers.NewMessageControllerWithClients(
		mockMsgRepo, mockChatRepo, mockChatUserRepo, mockFileClient, mockUserClient,
	)

	result, err := controller.SendMessage(senderID, chatID, &dto.CreateMessageDTO{
		Content: "Test message",
		FileIDs: []int{},
	})

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, msg.ID, result.ID)
	mockMsgRepo.AssertNotCalled(t, "CreateMessageFile", mock.Anything)

	mockChatRepo.AssertExpectations(t)
	mockUserClient.AssertExpectations(t)
	mockMsgRepo.AssertExpectations(t)
}

func TestMessageController_SendMessage_ChatNotFound(t *testing.T) {
	t.Parallel()
	// Arrange
	mockMsgRepo := new(MockMessageRepository)
	mockChatRepo := new(MockChatRepository)
	mockChatUserRepo := new(MockChatUserRepository)

	chatID := uuid.New()
	senderID := uuid.New()

	mockChatRepo.On("GetChatByID", chatID).Return(nil, errors.New("not found"))

	controller := controllers.NewMessageControllerWithClients(
		mockMsgRepo, mockChatRepo, mockChatUserRepo, nil, nil,
	)

	// Act
	result, err := controller.SendMessage(senderID, chatID, &dto.CreateMessageDTO{Content: "hi"})

	// Assert
	require.Error(t, err)
	assert.Nil(t, result)
	assert.True(t, errors.Is(err, custom_errors.ErrInvalidCredentials))

	mockChatRepo.AssertExpectations(t)
}

func TestMessageController_SendMessage_UserClientError(t *testing.T) {
	t.Parallel()
	mockMsgRepo := new(MockMessageRepository)
	mockChatRepo := new(MockChatRepository)
	mockChatUserRepo := new(MockChatUserRepository)
	mockUserClient := new(MockUserClient)

	chatID := uuid.New()
	senderID := uuid.New()

	mockChatRepo.On("GetChatByID", chatID).Return(createTestChat(), nil)
	mockUserClient.On("GetUserByID", &senderID).Return(nil, errors.New("http error"))

	controller := controllers.NewMessageControllerWithClients(
		mockMsgRepo, mockChatRepo, mockChatUserRepo, nil, mockUserClient,
	)

	result, err := controller.SendMessage(senderID, chatID, &dto.CreateMessageDTO{Content: "hi"})

	require.Error(t, err)
	assert.Nil(t, result)
	var userErr *custom_errors.UserClientError
	assert.True(t, errors.As(err, &userErr))

	mockUserClient.AssertExpectations(t)
	mockChatRepo.AssertExpectations(t)
}

func TestMessageController_SendMessage_FileHTTPError(t *testing.T) {
	t.Parallel()
	mockMsgRepo := new(MockMessageRepository)
	mockChatRepo := new(MockChatRepository)
	mockChatUserRepo := new(MockChatUserRepository)
	mockFileClient := new(MockFileClient)
	mockUserClient := new(MockUserClient)

	chatID := uuid.New()
	senderID := uuid.New()
	fileID := 1

	mockChatRepo.On("GetChatByID", chatID).Return(createTestChat(), nil)
	mockUserClient.On("GetUserByID", &senderID).Return(createTestUserResponse(), nil)
	mockFileClient.On("GetFileByID", fileID).Return(nil, errors.New("http err"))

	controller := controllers.NewMessageControllerWithClients(
		mockMsgRepo, mockChatRepo, mockChatUserRepo, mockFileClient, mockUserClient,
	)

	result, err := controller.SendMessage(senderID, chatID, &dto.CreateMessageDTO{
		Content: "text", FileIDs: []int{fileID},
	})

	require.Error(t, err)
	assert.Nil(t, result)
	var httpErr *custom_errors.GetFileHTTPError
	assert.True(t, errors.As(err, &httpErr))

	mockFileClient.AssertExpectations(t)
}

func TestMessageController_SendMessage_FileNotFound(t *testing.T) {
	t.Parallel()
	mockMsgRepo := new(MockMessageRepository)
	mockChatRepo := new(MockChatRepository)
	mockChatUserRepo := new(MockChatUserRepository)
	mockFileClient := new(MockFileClient)
	mockUserClient := new(MockUserClient)

	chatID := uuid.New()
	senderID := uuid.New()
	fileID := 1
	badFile := createTestFile()
	badFile.ID = 0

	mockChatRepo.On("GetChatByID", chatID).Return(createTestChat(), nil)
	mockUserClient.On("GetUserByID", &senderID).Return(createTestUserResponse(), nil)
	mockFileClient.On("GetFileByID", fileID).Return(badFile, nil)

	controller := controllers.NewMessageControllerWithClients(
		mockMsgRepo, mockChatRepo, mockChatUserRepo, mockFileClient, mockUserClient,
	)

	result, err := controller.SendMessage(senderID, chatID, &dto.CreateMessageDTO{
		Content: "text", FileIDs: []int{fileID},
	})

	require.Error(t, err)
	assert.Nil(t, result)
	var notFoundErr *custom_errors.FileNotFoundError
	assert.True(t, errors.As(err, &notFoundErr))

	mockFileClient.AssertExpectations(t)
}

func TestMessageController_SendMessage_CreateMessageError(t *testing.T) {
	t.Parallel()
	mockMsgRepo := new(MockMessageRepository)
	mockChatRepo := new(MockChatRepository)
	mockChatUserRepo := new(MockChatUserRepository)
	mockFileClient := new(MockFileClient)
	mockUserClient := new(MockUserClient)

	chatID := uuid.New()
	senderID := uuid.New()

	mockChatRepo.On("GetChatByID", chatID).Return(createTestChat(), nil)
	mockUserClient.On("GetUserByID", &senderID).Return(createTestUserResponse(), nil)
	mockMsgRepo.On("CreateMessage", mock.Anything).Return(errors.New("db"))

	controller := controllers.NewMessageControllerWithClients(
		mockMsgRepo, mockChatRepo, mockChatUserRepo, mockFileClient, mockUserClient,
	)

	result, err := controller.SendMessage(senderID, chatID, &dto.CreateMessageDTO{Content: "text", FileIDs: []int{}})

	require.Error(t, err)
	assert.Nil(t, result)
	var dbErr *custom_errors.DatabaseError
	assert.True(t, errors.As(err, &dbErr))
}

func TestMessageController_SendMessage_CreateMessageFileError(t *testing.T) {
	t.Parallel()
	mockMsgRepo := new(MockMessageRepository)
	mockChatRepo := new(MockChatRepository)
	mockChatUserRepo := new(MockChatUserRepository)
	mockFileClient := new(MockFileClient)
	mockUserClient := new(MockUserClient)

	chatID := uuid.New()
	senderID := uuid.New()
	fileID := 1

	mockChatRepo.On("GetChatByID", chatID).Return(createTestChat(), nil)
	mockUserClient.On("GetUserByID", &senderID).Return(createTestUserResponse(), nil)
	mockFileClient.On("GetFileByID", fileID).Return(createTestFile(), nil)
	mockMsgRepo.On("CreateMessage", mock.Anything).Return(nil)
	mockMsgRepo.On("CreateMessageFile", mock.Anything).Return(errors.New("db error"))

	controller := controllers.NewMessageControllerWithClients(
		mockMsgRepo, mockChatRepo, mockChatUserRepo, mockFileClient, mockUserClient,
	)

	result, err := controller.SendMessage(senderID, chatID, &dto.CreateMessageDTO{
		Content: "text", FileIDs: []int{fileID},
	})

	require.Error(t, err)
	assert.Nil(t, result)
	var dbErr *custom_errors.DatabaseError
	assert.True(t, errors.As(err, &dbErr))
}

func TestMessageController_SendMessage_GetMessageWithFileError(t *testing.T) {
	t.Parallel()
	mockMsgRepo := new(MockMessageRepository)
	mockChatRepo := new(MockChatRepository)
	mockChatUserRepo := new(MockChatUserRepository)
	mockFileClient := new(MockFileClient)
	mockUserClient := new(MockUserClient)

	chatID := uuid.New()
	senderID := uuid.New()

	mockChatRepo.On("GetChatByID", chatID).Return(createTestChat(), nil)
	mockUserClient.On("GetUserByID", &senderID).Return(createTestUserResponse(), nil)
	mockMsgRepo.On("CreateMessage", mock.Anything).Return(nil)
	mockMsgRepo.On("CreateMessageFile", mock.Anything).Return(nil)
	mockMsgRepo.On("GetMessageWithFile", mock.Anything).Return(nil, errors.New("db err"))

	controller := controllers.NewMessageControllerWithClients(
		mockMsgRepo, mockChatRepo, mockChatUserRepo, mockFileClient, mockUserClient,
	)

	result, err := controller.SendMessage(senderID, chatID, &dto.CreateMessageDTO{Content: "text", FileIDs: []int{}})

	require.Error(t, err)
	assert.Nil(t, result)
	var dbErr *custom_errors.DatabaseError
	assert.True(t, errors.As(err, &dbErr))
}

// Тесты для MessageController.GetChatMessages

func TestMessageController_GetChatMessages_Success(t *testing.T) {
	t.Parallel()
	mockMsgRepo := new(MockMessageRepository)
	mockChatRepo := new(MockChatRepository)
	mockChatUserRepo := new(MockChatUserRepository)
	mockFileClient := new(MockFileClient)

	chatID := uuid.New()
	msg := createTestMessage()
	msg.ChatID = chatID
	msg.Files = []models.MessageFile{{MessageID: msg.ID, FileID: 1}}

	mockChatRepo.On("GetChatByID", chatID).Return(createTestChat(), nil)
	mockMsgRepo.On("GetChatMessages", chatID, 0, 10).Return([]models.Message{*msg}, nil)
	mockFileClient.On("GetFileByID", 1).Return(createTestFile(), nil)

	controller := controllers.NewMessageControllerWithClients(
		mockMsgRepo, mockChatRepo, mockChatUserRepo, mockFileClient, nil,
	)

	result, err := controller.GetChatMessages(chatID, 0, 10)

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, *result, 1)
	assert.NotNil(t, *(*result)[0].Files)

	mockChatRepo.AssertExpectations(t)
	mockMsgRepo.AssertExpectations(t)
	mockFileClient.AssertExpectations(t)
}

func TestMessageController_GetChatMessages_MultipleMessagesWithFiles(t *testing.T) {
	t.Parallel()
	mockMsgRepo := new(MockMessageRepository)
	mockChatRepo := new(MockChatRepository)
	mockChatUserRepo := new(MockChatUserRepository)
	mockFileClient := new(MockFileClient)

	chatID := uuid.New()
	msg1 := createTestMessage()
	msg1.ChatID = chatID
	msg1.Files = []models.MessageFile{{MessageID: msg1.ID, FileID: 1}, {MessageID: msg1.ID, FileID: 2}}
	msg2 := createTestMessage()
	msg2.ChatID = chatID
	msg2.Files = []models.MessageFile{{MessageID: msg2.ID, FileID: 3}}

	mockChatRepo.On("GetChatByID", chatID).Return(createTestChat(), nil)
	mockMsgRepo.On("GetChatMessages", chatID, 0, 10).Return([]models.Message{*msg1, *msg2}, nil)
	mockFileClient.On("GetFileByID", 1).Return(createTestFile(), nil)
	mockFileClient.On("GetFileByID", 2).Return(createTestFile(), nil)
	mockFileClient.On("GetFileByID", 3).Return(createTestFile(), nil)

	controller := controllers.NewMessageControllerWithClients(
		mockMsgRepo, mockChatRepo, mockChatUserRepo, mockFileClient, nil,
	)

	result, err := controller.GetChatMessages(chatID, 0, 10)

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, *result, 2)
	assert.NotNil(t, *(*result)[0].Files)
	assert.Len(t, *(*result)[0].Files, 2)
	assert.NotNil(t, *(*result)[1].Files)
	assert.Len(t, *(*result)[1].Files, 1)

	mockChatRepo.AssertExpectations(t)
	mockMsgRepo.AssertExpectations(t)
	mockFileClient.AssertExpectations(t)
}

func TestMessageController_SendMessage_UserClientNilUser(t *testing.T) {
	t.Parallel()
	mockMsgRepo := new(MockMessageRepository)
	mockChatRepo := new(MockChatRepository)
	mockChatUserRepo := new(MockChatUserRepository)
	mockFileClient := new(MockFileClient)
	mockUserClient := new(MockUserClient)

	senderID := uuid.New()
	chatID := uuid.New()
	userResp := &cuc.Response{User: nil}

	mockChatRepo.On("GetChatByID", chatID).Return(createTestChat(), nil)
	mockUserClient.On("GetUserByID", &senderID).Return(userResp, nil)

	controller := controllers.NewMessageControllerWithClients(
		mockMsgRepo, mockChatRepo, mockChatUserRepo, mockFileClient, mockUserClient,
	)

	result, err := controller.SendMessage(senderID, chatID, &dto.CreateMessageDTO{Content: "test", FileIDs: []int{}})

	require.Error(t, err)
	assert.Nil(t, result)
	var userErr *custom_errors.UserClientError
	assert.True(t, errors.As(err, &userErr))

	mockChatRepo.AssertExpectations(t)
	mockUserClient.AssertExpectations(t)
}

func TestMessageController_GetChatMessages_EmptyMessages(t *testing.T) {
	t.Parallel()
	mockMsgRepo := new(MockMessageRepository)
	mockChatRepo := new(MockChatRepository)
	mockChatUserRepo := new(MockChatUserRepository)

	chatID := uuid.New()

	mockChatRepo.On("GetChatByID", chatID).Return(createTestChat(), nil)
	mockMsgRepo.On("GetChatMessages", chatID, 0, 10).Return([]models.Message{}, nil)

	controller := controllers.NewMessageControllerWithClients(
		mockMsgRepo, mockChatRepo, mockChatUserRepo, nil, nil,
	)

	result, err := controller.GetChatMessages(chatID, 0, 10)

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, *result, 0)

	mockChatRepo.AssertExpectations(t)
	mockMsgRepo.AssertExpectations(t)
}

func TestMessageController_GetChatMessages_ChatNotFound(t *testing.T) {
	t.Parallel()
	mockMsgRepo := new(MockMessageRepository)
	mockChatRepo := new(MockChatRepository)
	mockChatUserRepo := new(MockChatUserRepository)

	chatID := uuid.New()
	mockChatRepo.On("GetChatByID", chatID).Return(nil, errors.New("not found"))

	controller := controllers.NewMessageControllerWithClients(
		mockMsgRepo, mockChatRepo, mockChatUserRepo, nil, nil,
	)

	result, err := controller.GetChatMessages(chatID, 0, 10)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.True(t, errors.Is(err, custom_errors.ErrInvalidCredentials))

	mockChatRepo.AssertExpectations(t)
}

func TestMessageController_GetChatMessages_RepoError(t *testing.T) {
	t.Parallel()
	mockMsgRepo := new(MockMessageRepository)
	mockChatRepo := new(MockChatRepository)
	mockChatUserRepo := new(MockChatUserRepository)

	chatID := uuid.New()
	mockChatRepo.On("GetChatByID", chatID).Return(createTestChat(), nil)
	mockMsgRepo.On("GetChatMessages", chatID, 0, 10).Return(nil, errors.New("db"))

	controller := controllers.NewMessageControllerWithClients(
		mockMsgRepo, mockChatRepo, mockChatUserRepo, nil, nil,
	)

	result, err := controller.GetChatMessages(chatID, 0, 10)

	require.Error(t, err)
	assert.Nil(t, result)
	var dbErr *custom_errors.DatabaseError
	assert.True(t, errors.As(err, &dbErr))
}

func TestMessageController_GetChatMessages_FileHTTPError(t *testing.T) {
	t.Parallel()
	mockMsgRepo := new(MockMessageRepository)
	mockChatRepo := new(MockChatRepository)
	mockChatUserRepo := new(MockChatUserRepository)
	mockFileClient := new(MockFileClient)

	chatID := uuid.New()
	msg := createTestMessage()
	msg.ChatID = chatID
	msg.Files = []models.MessageFile{{MessageID: msg.ID, FileID: 1}}

	mockChatRepo.On("GetChatByID", chatID).Return(createTestChat(), nil)
	mockMsgRepo.On("GetChatMessages", chatID, 0, 10).Return([]models.Message{*msg}, nil)
	mockFileClient.On("GetFileByID", 1).Return(nil, errors.New("http"))

	controller := controllers.NewMessageControllerWithClients(
		mockMsgRepo, mockChatRepo, mockChatUserRepo, mockFileClient, nil,
	)

	result, err := controller.GetChatMessages(chatID, 0, 10)

	require.Error(t, err)
	assert.Nil(t, result)
	var httpErr *custom_errors.GetFileHTTPError
	assert.True(t, errors.As(err, &httpErr))
}

// Тесты для MessageController.SearchMessages

func TestMessageController_SearchMessages_Success(t *testing.T) {
	t.Parallel()
	mockMsgRepo := new(MockMessageRepository)
	mockChatRepo := new(MockChatRepository)
	mockChatUserRepo := new(MockChatUserRepository)

	userID := uuid.New()
	chatID := uuid.New()
	msg := createTestMessage()
	msg.ChatID = chatID

	mockChatRepo.On("GetChatByID", chatID).Return(createTestChat(), nil)
	mockChatUserRepo.On("GetChatUser", userID, chatID).Return(createTestChatUser(), nil)
	mockMsgRepo.On("SearchMessages", userID, chatID, "hi", 10, 0).
		Return([]models.Message{*msg}, int64(1), nil)

	controller := controllers.NewMessageControllerWithClients(
		mockMsgRepo, mockChatRepo, mockChatUserRepo, nil, nil,
	)

	result, err := controller.SearchMessages(userID, chatID, "hi", 10, 0)

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(1), *result.Total)
	assert.Len(t, *result.Messages, 1)

	mockChatRepo.AssertExpectations(t)
	mockChatUserRepo.AssertExpectations(t)
	mockMsgRepo.AssertExpectations(t)
}

func TestMessageController_SearchMessages_MultipleMessages(t *testing.T) {
	t.Parallel()
	mockMsgRepo := new(MockMessageRepository)
	mockChatRepo := new(MockChatRepository)
	mockChatUserRepo := new(MockChatUserRepository)

	userID := uuid.New()
	chatID := uuid.New()
	msg1 := createTestMessage()
	msg1.ChatID = chatID
	msg2 := createTestMessage()
	msg2.ChatID = chatID

	mockChatRepo.On("GetChatByID", chatID).Return(createTestChat(), nil)
	mockChatUserRepo.On("GetChatUser", userID, chatID).Return(createTestChatUser(), nil)
	mockMsgRepo.On("SearchMessages", userID, chatID, "test", 10, 0).
		Return([]models.Message{*msg1, *msg2}, int64(2), nil)

	controller := controllers.NewMessageControllerWithClients(
		mockMsgRepo, mockChatRepo, mockChatUserRepo, nil, nil,
	)

	result, err := controller.SearchMessages(userID, chatID, "test", 10, 0)

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(2), *result.Total)
	assert.Len(t, *result.Messages, 2)

	mockChatRepo.AssertExpectations(t)
	mockChatUserRepo.AssertExpectations(t)
	mockMsgRepo.AssertExpectations(t)
}

func TestMessageController_SearchMessages_GetChatUserError(t *testing.T) {
	t.Parallel()
	mockMsgRepo := new(MockMessageRepository)
	mockChatRepo := new(MockChatRepository)
	mockChatUserRepo := new(MockChatUserRepository)

	userID := uuid.New()
	chatID := uuid.New()

	mockChatRepo.On("GetChatByID", chatID).Return(createTestChat(), nil)
	mockChatUserRepo.On("GetChatUser", userID, chatID).Return(nil, errors.New("database error"))

	controller := controllers.NewMessageControllerWithClients(
		mockMsgRepo, mockChatRepo, mockChatUserRepo, nil, nil,
	)

	result, err := controller.SearchMessages(userID, chatID, "test", 10, 0)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "database error", err.Error())

	mockChatRepo.AssertExpectations(t)
	mockChatUserRepo.AssertExpectations(t)
}

func TestMessageController_SearchMessages_EmptyQuery(t *testing.T) {
	t.Parallel()
	mockMsgRepo := new(MockMessageRepository)
	mockChatRepo := new(MockChatRepository)
	mockChatUserRepo := new(MockChatUserRepository)

	userID := uuid.New()
	chatID := uuid.New()

	controller := controllers.NewMessageControllerWithClients(
		mockMsgRepo, mockChatRepo, mockChatUserRepo, nil, nil,
	)

	result, err := controller.SearchMessages(userID, chatID, "", 10, 0)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.True(t, errors.Is(err, custom_errors.ErrEmptyQuery))
}

func TestMessageController_SearchMessages_ChatNotFound(t *testing.T) {
	t.Parallel()
	mockMsgRepo := new(MockMessageRepository)
	mockChatRepo := new(MockChatRepository)
	mockChatUserRepo := new(MockChatUserRepository)

	userID := uuid.New()
	chatID := uuid.New()

	mockChatRepo.On("GetChatByID", chatID).Return(nil, errors.New("not found"))

	controller := controllers.NewMessageControllerWithClients(
		mockMsgRepo, mockChatRepo, mockChatUserRepo, nil, nil,
	)

	result, err := controller.SearchMessages(userID, chatID, "hi", 10, 0)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.True(t, errors.Is(err, custom_errors.ErrChatNotFound))

	mockChatRepo.AssertExpectations(t)
}

func TestMessageController_SearchMessages_Unauthorized(t *testing.T) {
	t.Parallel()
	mockMsgRepo := new(MockMessageRepository)
	mockChatRepo := new(MockChatRepository)
	mockChatUserRepo := new(MockChatUserRepository)

	userID := uuid.New()
	chatID := uuid.New()

	mockChatRepo.On("GetChatByID", chatID).Return(createTestChat(), nil)
	mockChatUserRepo.On("GetChatUser", userID, chatID).Return(nil, nil)

	controller := controllers.NewMessageControllerWithClients(
		mockMsgRepo, mockChatRepo, mockChatUserRepo, nil, nil,
	)

	result, err := controller.SearchMessages(userID, chatID, "hi", 10, 0)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.True(t, errors.Is(err, custom_errors.ErrUnauthorizedChat))

	mockChatUserRepo.AssertExpectations(t)
}

func TestMessageController_SearchMessages_RepoError(t *testing.T) {
	t.Parallel()
	mockMsgRepo := new(MockMessageRepository)
	mockChatRepo := new(MockChatRepository)
	mockChatUserRepo := new(MockChatUserRepository)

	userID := uuid.New()
	chatID := uuid.New()

	mockChatRepo.On("GetChatByID", chatID).Return(createTestChat(), nil)
	mockChatUserRepo.On("GetChatUser", userID, chatID).Return(createTestChatUser(), nil)
	mockMsgRepo.On("SearchMessages", userID, chatID, "hi", 10, 0).
		Return(nil, int64(0), errors.New("db"))

	controller := controllers.NewMessageControllerWithClients(
		mockMsgRepo, mockChatRepo, mockChatUserRepo, nil, nil,
	)

	result, err := controller.SearchMessages(userID, chatID, "hi", 10, 0)

	require.Error(t, err)
	assert.Nil(t, result)
	var dbErr *custom_errors.DatabaseError
	assert.True(t, errors.As(err, &dbErr))
}
