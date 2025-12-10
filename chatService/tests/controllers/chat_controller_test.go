package controllers

import (
	"errors"
	"testing"

	"chatService/internal/controllers"
	"chatService/internal/custom_errors"
	"chatService/internal/handlers/dto"
	"chatService/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Тесты для ChatController.GetChatByID

func TestChatController_GetChatByID_Success(t *testing.T) {
	// Arrange
	mockChatRepo := new(MockChatRepository)
	mockChatUserRepo := new(MockChatUserRepository)
	mockChatRoleRepo := new(MockChatRoleRepository)
	mockNotificationService := new(MockNotificationService)
	mockFileClient := new(MockFileClient)
	mockUserClient := new(MockUserClient)

	chatID := uuid.New()
	expectedChat := createTestChat()
	expectedChat.ID = chatID

	mockChatRepo.On("GetChatByID", chatID).Return(expectedChat, nil)

	controller := controllers.NewChatControllerWithClients(
		mockChatRepo,
		mockChatUserRepo,
		mockChatRoleRepo,
		mockNotificationService,
		mockFileClient,
		mockUserClient,
	)

	// Act
	result, err := controller.GetChatByID(chatID)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedChat.ID, result.ID)
	assert.Equal(t, expectedChat.Name, result.Name)

	mockChatRepo.AssertExpectations(t)
}

func TestChatController_GetChatByID_WithAvatar_Success(t *testing.T) {
	// Arrange
	mockChatRepo := new(MockChatRepository)
	mockChatUserRepo := new(MockChatUserRepository)
	mockChatRoleRepo := new(MockChatRoleRepository)
	mockNotificationService := new(MockNotificationService)
	mockFileClient := new(MockFileClient)
	mockUserClient := new(MockUserClient)

	chatID := uuid.New()
	expectedChat := createTestChat()
	expectedChat.ID = chatID
	avatarFileID := 1
	expectedChat.AvatarFileID = &avatarFileID
	expectedFile := createTestFile()

	mockChatRepo.On("GetChatByID", chatID).Return(expectedChat, nil)
	mockFileClient.On("GetFileByID", avatarFileID).Return(expectedFile, nil)

	controller := controllers.NewChatControllerWithClients(
		mockChatRepo,
		mockChatUserRepo,
		mockChatRoleRepo,
		mockNotificationService,
		mockFileClient,
		mockUserClient,
	)

	// Act
	result, err := controller.GetChatByID(chatID)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedChat.ID, result.ID)
	assert.Equal(t, expectedFile, result.AvatarFile)

	mockChatRepo.AssertExpectations(t)
	mockFileClient.AssertExpectations(t)
}

func TestChatController_GetChatByID_WithAvatar_FileClientError(t *testing.T) {
	// Arrange
	mockChatRepo := new(MockChatRepository)
	mockChatUserRepo := new(MockChatUserRepository)
	mockChatRoleRepo := new(MockChatRoleRepository)
	mockNotificationService := new(MockNotificationService)
	mockFileClient := new(MockFileClient)
	mockUserClient := new(MockUserClient)

	chatID := uuid.New()
	expectedChat := createTestChat()
	expectedChat.ID = chatID
	avatarFileID := 1
	expectedChat.AvatarFileID = &avatarFileID

	mockChatRepo.On("GetChatByID", chatID).Return(expectedChat, nil)
	mockFileClient.On("GetFileByID", avatarFileID).Return(nil, errors.New("file service error"))

	controller := controllers.NewChatControllerWithClients(
		mockChatRepo,
		mockChatUserRepo,
		mockChatRoleRepo,
		mockNotificationService,
		mockFileClient,
		mockUserClient,
	)

	// Act
	result, err := controller.GetChatByID(chatID)

	// Assert
	require.NoError(t, err) // Ошибка FileClient не должна прерывать процесс
	assert.NotNil(t, result)
	assert.Equal(t, expectedChat.ID, result.ID)
	assert.Nil(t, result.AvatarFile) // Аватар не должен быть установлен при ошибке

	mockChatRepo.AssertExpectations(t)
	mockFileClient.AssertExpectations(t)
}

func TestChatController_GetChatByID_NotFound(t *testing.T) {
	// Arrange
	mockChatRepo := new(MockChatRepository)
	mockChatUserRepo := new(MockChatUserRepository)
	mockChatRoleRepo := new(MockChatRoleRepository)
	mockNotificationService := new(MockNotificationService)
	mockFileClient := new(MockFileClient)
	mockUserClient := new(MockUserClient)

	chatID := uuid.New()
	mockChatRepo.On("GetChatByID", chatID).Return(nil, errors.New("record not found"))

	controller := controllers.NewChatControllerWithClients(
		mockChatRepo,
		mockChatUserRepo,
		mockChatRoleRepo,
		mockNotificationService,
		mockFileClient,
		mockUserClient,
	)

	// Act
	result, err := controller.GetChatByID(chatID)

	// Assert
	require.Error(t, err)
	assert.Nil(t, result)
	assert.True(t, errors.Is(err, custom_errors.ErrChatNotFound))

	mockChatRepo.AssertExpectations(t)
}

func TestChatController_GetChatByID_DatabaseError(t *testing.T) {
	// Arrange
	mockChatRepo := new(MockChatRepository)
	mockChatUserRepo := new(MockChatUserRepository)
	mockChatRoleRepo := new(MockChatRoleRepository)
	mockNotificationService := new(MockNotificationService)
	mockFileClient := new(MockFileClient)
	mockUserClient := new(MockUserClient)

	chatID := uuid.New()
	dbError := errors.New("database connection failed")
	mockChatRepo.On("GetChatByID", chatID).Return(nil, dbError)

	controller := controllers.NewChatControllerWithClients(
		mockChatRepo,
		mockChatUserRepo,
		mockChatRoleRepo,
		mockNotificationService,
		mockFileClient,
		mockUserClient,
	)

	// Act
	result, err := controller.GetChatByID(chatID)

	// Assert
	require.Error(t, err)
	assert.Nil(t, result)
	var dbErr *custom_errors.DatabaseError
	assert.True(t, errors.As(err, &dbErr))

	mockChatRepo.AssertExpectations(t)
}

// Тесты для ChatController.GetUserChats

func TestChatController_GetUserChats_Success(t *testing.T) {
	// Arrange
	mockChatRepo := new(MockChatRepository)
	mockChatUserRepo := new(MockChatUserRepository)
	mockChatRoleRepo := new(MockChatRoleRepository)
	mockNotificationService := new(MockNotificationService)
	mockFileClient := new(MockFileClient)
	mockUserClient := new(MockUserClient)

	userID := uuid.New()
	chats := []models.Chat{
		*createTestChat(),
		*createTestChatWithoutAvatar(),
	}

	mockChatRepo.On("GetUserChats", userID).Return(chats, nil)

	controller := controllers.NewChatControllerWithClients(
		mockChatRepo,
		mockChatUserRepo,
		mockChatRoleRepo,
		mockNotificationService,
		mockFileClient,
		mockUserClient,
	)

	// Act
	result, err := controller.GetUserChats(userID)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, *result, 2)

	mockChatRepo.AssertExpectations(t)
}

func TestChatController_GetUserChats_WithAvatars(t *testing.T) {
	// Arrange
	mockChatRepo := new(MockChatRepository)
	mockChatUserRepo := new(MockChatUserRepository)
	mockChatRoleRepo := new(MockChatRoleRepository)
	mockNotificationService := new(MockNotificationService)
	mockFileClient := new(MockFileClient)
	mockUserClient := new(MockUserClient)

	userID := uuid.New()
	avatarFileID := 1
	chat1 := createTestChat()
	chat1.AvatarFileID = &avatarFileID
	chat2 := createTestChatWithoutAvatar()
	chats := []models.Chat{*chat1, *chat2}
	expectedFile := createTestFile()

	mockChatRepo.On("GetUserChats", userID).Return(chats, nil)
	mockFileClient.On("GetFileByID", avatarFileID).Return(expectedFile, nil)

	controller := controllers.NewChatControllerWithClients(
		mockChatRepo,
		mockChatUserRepo,
		mockChatRoleRepo,
		mockNotificationService,
		mockFileClient,
		mockUserClient,
	)

	// Act
	result, err := controller.GetUserChats(userID)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, *result, 2)
	assert.Equal(t, expectedFile, (*result)[0].AvatarFile)
	assert.Nil(t, (*result)[1].AvatarFile)

	mockChatRepo.AssertExpectations(t)
	mockFileClient.AssertExpectations(t)
}

func TestChatController_GetUserChats_RepositoryError(t *testing.T) {
	// Arrange
	mockChatRepo := new(MockChatRepository)
	mockChatUserRepo := new(MockChatUserRepository)
	mockChatRoleRepo := new(MockChatRoleRepository)
	mockNotificationService := new(MockNotificationService)
	mockFileClient := new(MockFileClient)
	mockUserClient := new(MockUserClient)

	userID := uuid.New()
	repoError := errors.New("database error")
	mockChatRepo.On("GetUserChats", userID).Return(nil, repoError)

	controller := controllers.NewChatControllerWithClients(
		mockChatRepo,
		mockChatUserRepo,
		mockChatRoleRepo,
		mockNotificationService,
		mockFileClient,
		mockUserClient,
	)

	// Act
	result, err := controller.GetUserChats(userID)

	// Assert
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, repoError, err)

	mockChatRepo.AssertExpectations(t)
}

// Тесты для ChatController.CreateChat

func TestChatController_CreateChat_Success(t *testing.T) {
	// Arrange
	mockChatRepo := new(MockChatRepository)
	mockChatUserRepo := new(MockChatUserRepository)
	mockChatRoleRepo := new(MockChatRoleRepository)
	mockNotificationService := new(MockNotificationService)
	mockFileClient := new(MockFileClient)
	mockUserClient := new(MockUserClient)

	ownerID := uuid.New()
	userID1 := uuid.New()
	userID2 := uuid.New()
	ownerRole := createTestChatRoleWithID(1, "owner")
	mainRole := createTestChatRoleWithID(2, "main")
	ownerUser := createTestUserResponse()
	ownerUser.User.ID = ownerID
	user1 := createTestUserResponseWithEmail("user1@example.com")
	user1.User.ID = userID1
	user2 := createTestUserResponseWithEmail("user2@example.com")
	user2.User.ID = userID2

	createChatDTO := &dto.CreateChatDTO{
		Name:        "Test Chat",
		Description: stringPtr("Test Description"),
		OwnerID:     ownerID,
		UserIDs:     []uuid.UUID{userID1, userID2},
	}

	mockChatRoleRepo.On("GetRoleByName", "owner").Return(ownerRole, nil)
	mockChatRoleRepo.On("GetRoleByName", "main").Return(mainRole, nil)
	mockUserClient.On("GetUserByID", &ownerID).Return(ownerUser, nil)
	mockUserClient.On("GetUserByID", &userID1).Return(user1, nil)
	mockUserClient.On("GetUserByID", &userID2).Return(user2, nil)
	mockChatRepo.On("CreateChat", mock.AnythingOfType("*models.Chat")).Return(nil).Run(func(args mock.Arguments) {
		chat := args.Get(0).(*models.Chat)
		chat.ID = uuid.New()
	})
	mockChatUserRepo.On("AddUserToChat", mock.MatchedBy(func(cu *models.ChatUser) bool {
		return cu.RoleID == ownerRole.ID && cu.UserID == ownerID
	})).Return(nil)
	mockChatUserRepo.On("AddUserToChat", mock.MatchedBy(func(cu *models.ChatUser) bool {
		return cu.RoleID == mainRole.ID && cu.UserID == userID1
	})).Return(nil)
	mockChatUserRepo.On("AddUserToChat", mock.MatchedBy(func(cu *models.ChatUser) bool {
		return cu.RoleID == mainRole.ID && cu.UserID == userID2
	})).Return(nil)
	mockNotificationService.On("SendChatCreatedNotification",
		mock.AnythingOfType("uuid.UUID"),
		createChatDTO.Name,
		ownerUser.User.Username,
		true,
		*createChatDTO.Description,
		user1.User.Email,
	).Return(nil)
	mockNotificationService.On("SendChatCreatedNotification",
		mock.AnythingOfType("uuid.UUID"),
		createChatDTO.Name,
		ownerUser.User.Username,
		true,
		*createChatDTO.Description,
		user2.User.Email,
	).Return(nil)

	controller := controllers.NewChatControllerWithClients(
		mockChatRepo,
		mockChatUserRepo,
		mockChatRoleRepo,
		mockNotificationService,
		mockFileClient,
		mockUserClient,
	)

	// Act
	result, err := controller.CreateChat(createChatDTO)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)

	mockChatRoleRepo.AssertExpectations(t)
	mockUserClient.AssertExpectations(t)
	mockChatRepo.AssertExpectations(t)
	mockChatUserRepo.AssertExpectations(t)
	mockNotificationService.AssertExpectations(t)
}
