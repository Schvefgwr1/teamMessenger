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

// Тесты для ChatController.CreateChat (дополнительные сценарии)

func TestChatController_CreateChat_WithoutUsers_Success(t *testing.T) {
	// Arrange
	mockChatRepo := new(MockChatRepository)
	mockChatUserRepo := new(MockChatUserRepository)
	mockChatRoleRepo := new(MockChatRoleRepository)
	mockNotificationService := new(MockNotificationService)
	mockFileClient := new(MockFileClient)
	mockUserClient := new(MockUserClient)

	ownerID := uuid.New()
	ownerRole := createTestChatRoleWithID(1, "owner")
	ownerUser := createTestUserResponse()
	ownerUser.User.ID = ownerID

	createChatDTO := &dto.CreateChatDTO{
		Name:    "Test Chat",
		OwnerID: ownerID,
		UserIDs: []uuid.UUID{},
	}

	mockChatRoleRepo.On("GetRoleByName", "owner").Return(ownerRole, nil)
	mockChatRoleRepo.On("GetRoleByName", "main").Return(createTestChatRoleWithID(2, "main"), nil)
	mockUserClient.On("GetUserByID", &ownerID).Return(ownerUser, nil)
	mockChatRepo.On("CreateChat", mock.MatchedBy(func(chat *models.Chat) bool {
		return chat.IsGroup == false
	})).Return(nil).Run(func(args mock.Arguments) {
		chat := args.Get(0).(*models.Chat)
		chat.ID = uuid.New()
	})
	mockChatUserRepo.On("AddUserToChat", mock.MatchedBy(func(cu *models.ChatUser) bool {
		return cu.RoleID == ownerRole.ID && cu.UserID == ownerID
	})).Return(nil)

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
	mockNotificationService.AssertNotCalled(t, "SendChatCreatedNotification", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

func TestChatController_CreateChat_OwnerRoleNotFound(t *testing.T) {
	// Arrange
	mockChatRepo := new(MockChatRepository)
	mockChatUserRepo := new(MockChatUserRepository)
	mockChatRoleRepo := new(MockChatRoleRepository)
	mockNotificationService := new(MockNotificationService)
	mockFileClient := new(MockFileClient)
	mockUserClient := new(MockUserClient)

	ownerID := uuid.New()
	createChatDTO := &dto.CreateChatDTO{
		Name:    "Test Chat",
		OwnerID: ownerID,
		UserIDs: []uuid.UUID{},
	}

	mockChatRoleRepo.On("GetRoleByName", "owner").Return(nil, errors.New("role not found"))

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
	require.Error(t, err)
	assert.Nil(t, result)
	assert.True(t, errors.Is(err, custom_errors.ErrInvalidCredentials))

	mockChatRoleRepo.AssertExpectations(t)
	mockUserClient.AssertNotCalled(t, "GetUserByID", mock.Anything)
}

func TestChatController_CreateChat_WithAvatar_Success(t *testing.T) {
	// Arrange
	mockChatRepo := new(MockChatRepository)
	mockChatUserRepo := new(MockChatUserRepository)
	mockChatRoleRepo := new(MockChatRoleRepository)
	mockNotificationService := new(MockNotificationService)
	mockFileClient := new(MockFileClient)
	mockUserClient := new(MockUserClient)

	ownerID := uuid.New()
	avatarFileID := 1
	ownerRole := createTestChatRoleWithID(1, "owner")
	mainRole := createTestChatRoleWithID(2, "main")
	ownerUser := createTestUserResponse()
	ownerUser.User.ID = ownerID
	expectedFile := createTestFile()

	createChatDTO := &dto.CreateChatDTO{
		Name:         "Test Chat",
		AvatarFileID: &avatarFileID,
		OwnerID:      ownerID,
		UserIDs:      []uuid.UUID{},
	}

	mockChatRoleRepo.On("GetRoleByName", "owner").Return(ownerRole, nil)
	mockChatRoleRepo.On("GetRoleByName", "main").Return(mainRole, nil)
	mockFileClient.On("GetFileByID", avatarFileID).Return(expectedFile, nil)
	mockUserClient.On("GetUserByID", &ownerID).Return(ownerUser, nil)
	mockChatRepo.On("CreateChat", mock.MatchedBy(func(chat *models.Chat) bool {
		return chat.AvatarFileID != nil && *chat.AvatarFileID == expectedFile.ID
	})).Return(nil).Run(func(args mock.Arguments) {
		chat := args.Get(0).(*models.Chat)
		chat.ID = uuid.New()
	})
	mockChatUserRepo.On("AddUserToChat", mock.AnythingOfType("*models.ChatUser")).Return(nil)

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
	mockFileClient.AssertExpectations(t)
	mockUserClient.AssertExpectations(t)
	mockChatRepo.AssertExpectations(t)
}

func TestChatController_CreateChat_InvalidAvatarFile(t *testing.T) {
	// Arrange
	mockChatRepo := new(MockChatRepository)
	mockChatUserRepo := new(MockChatUserRepository)
	mockChatRoleRepo := new(MockChatRoleRepository)
	mockNotificationService := new(MockNotificationService)
	mockFileClient := new(MockFileClient)
	mockUserClient := new(MockUserClient)

	ownerID := uuid.New()
	avatarFileID := 1
	ownerRole := createTestChatRoleWithID(1, "owner")
	mainRole := createTestChatRoleWithID(2, "main")
	fileError := errors.New("file not found")

	createChatDTO := &dto.CreateChatDTO{
		Name:         "Test Chat",
		AvatarFileID: &avatarFileID,
		OwnerID:      ownerID,
		UserIDs:      []uuid.UUID{},
	}

	mockChatRoleRepo.On("GetRoleByName", "owner").Return(ownerRole, nil)
	mockChatRoleRepo.On("GetRoleByName", "main").Return(mainRole, nil)
	mockFileClient.On("GetFileByID", avatarFileID).Return(nil, fileError)

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
	require.Error(t, err)
	assert.Nil(t, result)
	var getFileErr *custom_errors.GetFileHTTPError
	assert.True(t, errors.As(err, &getFileErr))

	mockChatRoleRepo.AssertExpectations(t)
	mockFileClient.AssertExpectations(t)
	mockChatRepo.AssertNotCalled(t, "CreateChat", mock.Anything)
}

// Тесты для ChatController.UpdateChat

func TestChatController_UpdateChat_Success(t *testing.T) {
	// Arrange
	mockChatRepo := new(MockChatRepository)
	mockChatUserRepo := new(MockChatUserRepository)
	mockChatRoleRepo := new(MockChatRoleRepository)
	mockNotificationService := new(MockNotificationService)
	mockFileClient := new(MockFileClient)
	mockUserClient := new(MockUserClient)

	chatID := uuid.New()
	chat := createTestChat()
	chat.ID = chatID
	newName := "Updated Chat Name"
	mainRole := createTestChatRoleWithID(2, "main")

	updateChatDTO := &dto.UpdateChatDTO{
		Name: &newName,
	}

	mockChatRepo.On("GetChatByID", chatID).Return(chat, nil)
	mockChatRepo.On("UpdateChat", mock.MatchedBy(func(c *models.Chat) bool {
		return c.Name == newName
	})).Return(nil)
	mockChatRoleRepo.On("GetRoleByName", "main").Return(mainRole, nil)

	controller := controllers.NewChatControllerWithClients(
		mockChatRepo,
		mockChatUserRepo,
		mockChatRoleRepo,
		mockNotificationService,
		mockFileClient,
		mockUserClient,
	)

	// Act
	result, err := controller.UpdateChat(chatID, updateChatDTO)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, newName, result.Chat.Name)

	mockChatRepo.AssertExpectations(t)
	mockChatRoleRepo.AssertExpectations(t)
}

func TestChatController_UpdateChat_ChatNotFound(t *testing.T) {
	// Arrange
	mockChatRepo := new(MockChatRepository)
	mockChatUserRepo := new(MockChatUserRepository)
	mockChatRoleRepo := new(MockChatRoleRepository)
	mockNotificationService := new(MockNotificationService)
	mockFileClient := new(MockFileClient)
	mockUserClient := new(MockUserClient)

	chatID := uuid.New()
	updateChatDTO := &dto.UpdateChatDTO{}

	mockChatRepo.On("GetChatByID", chatID).Return(nil, errors.New("chat not found"))

	controller := controllers.NewChatControllerWithClients(
		mockChatRepo,
		mockChatUserRepo,
		mockChatRoleRepo,
		mockNotificationService,
		mockFileClient,
		mockUserClient,
	)

	// Act
	result, err := controller.UpdateChat(chatID, updateChatDTO)

	// Assert
	require.Error(t, err)
	assert.Nil(t, result)
	var dbErr *custom_errors.DatabaseError
	assert.True(t, errors.As(err, &dbErr))

	mockChatRepo.AssertExpectations(t)
}

func TestChatController_UpdateChat_AddUser_Success(t *testing.T) {
	// Arrange
	mockChatRepo := new(MockChatRepository)
	mockChatUserRepo := new(MockChatUserRepository)
	mockChatRoleRepo := new(MockChatRoleRepository)
	mockNotificationService := new(MockNotificationService)
	mockFileClient := new(MockFileClient)
	mockUserClient := new(MockUserClient)

	chatID := uuid.New()
	userID := uuid.New()
	chat := createTestChat()
	chat.ID = chatID
	mainRole := createTestChatRoleWithID(2, "main")
	user := createTestUserResponseWithEmail("user@example.com")
	user.User.ID = userID

	updateChatDTO := &dto.UpdateChatDTO{
		AddUserIDs: []uuid.UUID{userID},
	}

	mockChatRepo.On("GetChatByID", chatID).Return(chat, nil)
	mockChatRepo.On("UpdateChat", chat).Return(nil)
	mockChatRoleRepo.On("GetRoleByName", "main").Return(mainRole, nil)
	mockUserClient.On("GetUserByID", &userID).Return(user, nil)
	mockChatUserRepo.On("AddUserToChat", mock.MatchedBy(func(cu *models.ChatUser) bool {
		return cu.ChatID == chatID && cu.UserID == userID && cu.RoleID == mainRole.ID
	})).Return(nil)
	mockNotificationService.On("SendChatCreatedNotification",
		chatID, chat.Name, "Администратор", chat.IsGroup, mock.Anything, user.User.Email,
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
	result, err := controller.UpdateChat(chatID, updateChatDTO)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.UpdateUsers, 1)
	assert.Equal(t, userID, result.UpdateUsers[0].UserID)
	assert.Equal(t, "created", result.UpdateUsers[0].State)

	mockChatRepo.AssertExpectations(t)
	mockChatRoleRepo.AssertExpectations(t)
	mockUserClient.AssertExpectations(t)
	mockChatUserRepo.AssertExpectations(t)
	mockNotificationService.AssertExpectations(t)
}

func TestChatController_UpdateChat_RemoveUser_Success(t *testing.T) {
	// Arrange
	mockChatRepo := new(MockChatRepository)
	mockChatUserRepo := new(MockChatUserRepository)
	mockChatRoleRepo := new(MockChatRoleRepository)
	mockNotificationService := new(MockNotificationService)
	mockFileClient := new(MockFileClient)
	mockUserClient := new(MockUserClient)

	chatID := uuid.New()
	userID := uuid.New()
	chat := createTestChat()
	chat.ID = chatID
	mainRole := createTestChatRoleWithID(2, "main")

	updateChatDTO := &dto.UpdateChatDTO{
		RemoveUserIDs: []uuid.UUID{userID},
	}

	mockChatRepo.On("GetChatByID", chatID).Return(chat, nil)
	mockChatRepo.On("UpdateChat", chat).Return(nil)
	mockChatRoleRepo.On("GetRoleByName", "main").Return(mainRole, nil)
	mockChatUserRepo.On("RemoveUserFromChat", chatID, userID).Return(nil)

	controller := controllers.NewChatControllerWithClients(
		mockChatRepo,
		mockChatUserRepo,
		mockChatRoleRepo,
		mockNotificationService,
		mockFileClient,
		mockUserClient,
	)

	// Act
	result, err := controller.UpdateChat(chatID, updateChatDTO)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.UpdateUsers, 1)
	assert.Equal(t, userID, result.UpdateUsers[0].UserID)
	assert.Equal(t, "deleted", result.UpdateUsers[0].State)

	mockChatRepo.AssertExpectations(t)
	mockChatRoleRepo.AssertExpectations(t)
	mockChatUserRepo.AssertExpectations(t)
}

func TestChatController_UpdateChat_UpdateDescriptionOnly(t *testing.T) {
	// Arrange
	mockChatRepo := new(MockChatRepository)
	mockChatUserRepo := new(MockChatUserRepository)
	mockChatRoleRepo := new(MockChatRoleRepository)
	mockNotificationService := new(MockNotificationService)
	mockFileClient := new(MockFileClient)
	mockUserClient := new(MockUserClient)

	chatID := uuid.New()
	chat := createTestChat()
	chat.ID = chatID
	newDescription := "Updated description"
	mainRole := createTestChatRoleWithID(2, "main")

	updateChatDTO := &dto.UpdateChatDTO{
		Description: &newDescription,
	}

	mockChatRepo.On("GetChatByID", chatID).Return(chat, nil)
	mockChatRepo.On("UpdateChat", mock.MatchedBy(func(c *models.Chat) bool {
		return c.Description != nil && *c.Description == newDescription
	})).Return(nil)
	mockChatRoleRepo.On("GetRoleByName", "main").Return(mainRole, nil)

	controller := controllers.NewChatControllerWithClients(
		mockChatRepo,
		mockChatUserRepo,
		mockChatRoleRepo,
		mockNotificationService,
		mockFileClient,
		mockUserClient,
	)

	// Act
	result, err := controller.UpdateChat(chatID, updateChatDTO)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, newDescription, *result.Chat.Description)

	mockChatRepo.AssertExpectations(t)
	mockChatRoleRepo.AssertExpectations(t)
}

func TestChatController_UpdateChat_NotificationServiceError(t *testing.T) {
	// Arrange
	mockChatRepo := new(MockChatRepository)
	mockChatUserRepo := new(MockChatUserRepository)
	mockChatRoleRepo := new(MockChatRoleRepository)
	mockNotificationService := new(MockNotificationService)
	mockFileClient := new(MockFileClient)
	mockUserClient := new(MockUserClient)

	chatID := uuid.New()
	userID := uuid.New()
	chat := createTestChat()
	chat.ID = chatID
	mainRole := createTestChatRoleWithID(2, "main")
	user := createTestUserResponseWithEmail("user@example.com")
	user.User.ID = userID

	updateChatDTO := &dto.UpdateChatDTO{
		AddUserIDs: []uuid.UUID{userID},
	}

	mockChatRepo.On("GetChatByID", chatID).Return(chat, nil)
	mockChatRepo.On("UpdateChat", chat).Return(nil)
	mockChatRoleRepo.On("GetRoleByName", "main").Return(mainRole, nil)
	mockUserClient.On("GetUserByID", &userID).Return(user, nil)
	mockChatUserRepo.On("AddUserToChat", mock.Anything).Return(nil)
	mockNotificationService.On("SendChatCreatedNotification", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("notification error"))

	controller := controllers.NewChatControllerWithClients(
		mockChatRepo,
		mockChatUserRepo,
		mockChatRoleRepo,
		mockNotificationService,
		mockFileClient,
		mockUserClient,
	)

	// Act
	result, err := controller.UpdateChat(chatID, updateChatDTO)

	// Assert
	require.NoError(t, err) // Ошибка NotificationService не должна прерывать процесс
	assert.NotNil(t, result)
	assert.Len(t, result.UpdateUsers, 1)

	mockChatRepo.AssertExpectations(t)
	mockChatRoleRepo.AssertExpectations(t)
	mockUserClient.AssertExpectations(t)
	mockChatUserRepo.AssertExpectations(t)
	mockNotificationService.AssertExpectations(t)
}

func TestChatController_UpdateChat_MultipleUsers(t *testing.T) {
	// Arrange
	mockChatRepo := new(MockChatRepository)
	mockChatUserRepo := new(MockChatUserRepository)
	mockChatRoleRepo := new(MockChatRoleRepository)
	mockNotificationService := new(MockNotificationService)
	mockFileClient := new(MockFileClient)
	mockUserClient := new(MockUserClient)

	chatID := uuid.New()
	userID1 := uuid.New()
	userID2 := uuid.New()
	chat := createTestChat()
	chat.ID = chatID
	mainRole := createTestChatRoleWithID(2, "main")
	user1 := createTestUserResponseWithEmail("user1@example.com")
	user1.User.ID = userID1
	user2 := createTestUserResponseWithEmail("user2@example.com")
	user2.User.ID = userID2

	updateChatDTO := &dto.UpdateChatDTO{
		AddUserIDs:    []uuid.UUID{userID1, userID2},
		RemoveUserIDs: []uuid.UUID{userID1},
	}

	mockChatRepo.On("GetChatByID", chatID).Return(chat, nil)
	mockChatRepo.On("UpdateChat", chat).Return(nil)
	mockChatRoleRepo.On("GetRoleByName", "main").Return(mainRole, nil)
	mockUserClient.On("GetUserByID", &userID1).Return(user1, nil)
	mockUserClient.On("GetUserByID", &userID2).Return(user2, nil)
	mockChatUserRepo.On("AddUserToChat", mock.MatchedBy(func(cu *models.ChatUser) bool {
		return cu.UserID == userID1 || cu.UserID == userID2
	})).Return(nil).Times(2)
	mockChatUserRepo.On("RemoveUserFromChat", chatID, userID1).Return(nil)
	mockNotificationService.On("SendChatCreatedNotification", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, user1.User.Email).Return(nil)
	mockNotificationService.On("SendChatCreatedNotification", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, user2.User.Email).Return(nil)

	controller := controllers.NewChatControllerWithClients(
		mockChatRepo,
		mockChatUserRepo,
		mockChatRoleRepo,
		mockNotificationService,
		mockFileClient,
		mockUserClient,
	)

	// Act
	result, err := controller.UpdateChat(chatID, updateChatDTO)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.UpdateUsers, 3) // 2 добавленных + 1 удаленный

	mockChatRepo.AssertExpectations(t)
	mockChatRoleRepo.AssertExpectations(t)
	mockUserClient.AssertExpectations(t)
	mockChatUserRepo.AssertExpectations(t)
	mockNotificationService.AssertExpectations(t)
}

// Тесты для ChatController.DeleteChat

func TestChatController_DeleteChat_Success(t *testing.T) {
	// Arrange
	mockChatRepo := new(MockChatRepository)
	mockChatUserRepo := new(MockChatUserRepository)
	mockChatRoleRepo := new(MockChatRoleRepository)
	mockNotificationService := new(MockNotificationService)
	mockFileClient := new(MockFileClient)
	mockUserClient := new(MockUserClient)

	chatID := uuid.New()

	mockChatUserRepo.On("DeleteChatUsersByChatID", chatID).Return(nil)
	mockChatRepo.On("DeleteChat", chatID).Return(nil)

	controller := controllers.NewChatControllerWithClients(
		mockChatRepo,
		mockChatUserRepo,
		mockChatRoleRepo,
		mockNotificationService,
		mockFileClient,
		mockUserClient,
	)

	// Act
	err := controller.DeleteChat(chatID)

	// Assert
	require.NoError(t, err)

	mockChatUserRepo.AssertExpectations(t)
	mockChatRepo.AssertExpectations(t)
}

func TestChatController_DeleteChat_DeleteUsersError(t *testing.T) {
	// Arrange
	mockChatRepo := new(MockChatRepository)
	mockChatUserRepo := new(MockChatUserRepository)
	mockChatRoleRepo := new(MockChatRoleRepository)
	mockNotificationService := new(MockNotificationService)
	mockFileClient := new(MockFileClient)
	mockUserClient := new(MockUserClient)

	chatID := uuid.New()
	deleteError := errors.New("failed to delete users")

	mockChatUserRepo.On("DeleteChatUsersByChatID", chatID).Return(deleteError)

	controller := controllers.NewChatControllerWithClients(
		mockChatRepo,
		mockChatUserRepo,
		mockChatRoleRepo,
		mockNotificationService,
		mockFileClient,
		mockUserClient,
	)

	// Act
	err := controller.DeleteChat(chatID)

	// Assert
	require.Error(t, err)
	var dbErr *custom_errors.DatabaseError
	assert.True(t, errors.As(err, &dbErr))

	mockChatUserRepo.AssertExpectations(t)
	mockChatRepo.AssertNotCalled(t, "DeleteChat", mock.Anything)
}

// Тесты для ChatController.ChangeUserRole

func TestChatController_ChangeUserRole_Success(t *testing.T) {
	// Arrange
	mockChatRepo := new(MockChatRepository)
	mockChatUserRepo := new(MockChatUserRepository)
	mockChatRoleRepo := new(MockChatRoleRepository)
	mockNotificationService := new(MockNotificationService)
	mockFileClient := new(MockFileClient)
	mockUserClient := new(MockUserClient)

	chatID := uuid.New()
	userID := uuid.New()
	roleID := 2
	chatUser := createTestChatUserWithRole(chatID, userID, 1, "main")
	newRole := createTestChatRoleWithID(roleID, "admin")

	mockChatUserRepo.On("GetChatUser", userID, chatID).Return(chatUser, nil)
	mockChatRoleRepo.On("GetRoleByID", roleID).Return(newRole, nil)
	mockChatUserRepo.On("ChangeUserRole", chatID, userID, roleID).Return(nil)

	controller := controllers.NewChatControllerWithClients(
		mockChatRepo,
		mockChatUserRepo,
		mockChatRoleRepo,
		mockNotificationService,
		mockFileClient,
		mockUserClient,
	)

	// Act
	err := controller.ChangeUserRole(chatID, userID, roleID)

	// Assert
	require.NoError(t, err)

	mockChatUserRepo.AssertExpectations(t)
	mockChatRoleRepo.AssertExpectations(t)
}

func TestChatController_ChangeUserRole_UserNotInChat(t *testing.T) {
	// Arrange
	mockChatRepo := new(MockChatRepository)
	mockChatUserRepo := new(MockChatUserRepository)
	mockChatRoleRepo := new(MockChatRoleRepository)
	mockNotificationService := new(MockNotificationService)
	mockFileClient := new(MockFileClient)
	mockUserClient := new(MockUserClient)

	chatID := uuid.New()
	userID := uuid.New()
	roleID := 2

	mockChatUserRepo.On("GetChatUser", userID, chatID).Return(nil, errors.New("user not found"))

	controller := controllers.NewChatControllerWithClients(
		mockChatRepo,
		mockChatUserRepo,
		mockChatRoleRepo,
		mockNotificationService,
		mockFileClient,
		mockUserClient,
	)

	// Act
	err := controller.ChangeUserRole(chatID, userID, roleID)

	// Assert
	require.Error(t, err)
	assert.True(t, errors.Is(err, custom_errors.ErrInvalidCredentials))

	mockChatUserRepo.AssertExpectations(t)
	mockChatRoleRepo.AssertNotCalled(t, "GetRoleByID", mock.Anything)
}

func TestChatController_ChangeUserRole_RoleNotFound(t *testing.T) {
	// Arrange
	mockChatRepo := new(MockChatRepository)
	mockChatUserRepo := new(MockChatUserRepository)
	mockChatRoleRepo := new(MockChatRoleRepository)
	mockNotificationService := new(MockNotificationService)
	mockFileClient := new(MockFileClient)
	mockUserClient := new(MockUserClient)

	chatID := uuid.New()
	userID := uuid.New()
	roleID := 999
	chatUser := createTestChatUserWithRole(chatID, userID, 1, "main")

	mockChatUserRepo.On("GetChatUser", userID, chatID).Return(chatUser, nil)
	mockChatRoleRepo.On("GetRoleByID", roleID).Return(nil, errors.New("role not found"))

	controller := controllers.NewChatControllerWithClients(
		mockChatRepo,
		mockChatUserRepo,
		mockChatRoleRepo,
		mockNotificationService,
		mockFileClient,
		mockUserClient,
	)

	// Act
	err := controller.ChangeUserRole(chatID, userID, roleID)

	// Assert
	require.Error(t, err)
	assert.True(t, errors.Is(err, custom_errors.ErrInvalidCredentials))

	mockChatUserRepo.AssertExpectations(t)
	mockChatRoleRepo.AssertExpectations(t)
}

func TestChatController_ChangeUserRole_ChangeRoleError(t *testing.T) {
	// Arrange
	mockChatRepo := new(MockChatRepository)
	mockChatUserRepo := new(MockChatUserRepository)
	mockChatRoleRepo := new(MockChatRoleRepository)
	mockNotificationService := new(MockNotificationService)
	mockFileClient := new(MockFileClient)
	mockUserClient := new(MockUserClient)

	chatID := uuid.New()
	userID := uuid.New()
	roleID := 2
	chatUser := createTestChatUserWithRole(chatID, userID, 1, "main")
	newRole := createTestChatRoleWithID(roleID, "admin")

	mockChatUserRepo.On("GetChatUser", userID, chatID).Return(chatUser, nil)
	mockChatRoleRepo.On("GetRoleByID", roleID).Return(newRole, nil)
	mockChatUserRepo.On("ChangeUserRole", chatID, userID, roleID).Return(errors.New("database error"))

	controller := controllers.NewChatControllerWithClients(
		mockChatRepo,
		mockChatUserRepo,
		mockChatRoleRepo,
		mockNotificationService,
		mockFileClient,
		mockUserClient,
	)

	// Act
	err := controller.ChangeUserRole(chatID, userID, roleID)

	// Assert
	require.Error(t, err)
	assert.Equal(t, "database error", err.Error())

	mockChatUserRepo.AssertExpectations(t)
	mockChatRoleRepo.AssertExpectations(t)
}

func TestChatController_ChangeUserRole_ChatUserNil(t *testing.T) {
	// Arrange
	mockChatRepo := new(MockChatRepository)
	mockChatUserRepo := new(MockChatUserRepository)
	mockChatRoleRepo := new(MockChatRoleRepository)
	mockNotificationService := new(MockNotificationService)
	mockFileClient := new(MockFileClient)
	mockUserClient := new(MockUserClient)

	chatID := uuid.New()
	userID := uuid.New()
	roleID := 2

	mockChatUserRepo.On("GetChatUser", userID, chatID).Return(nil, nil)

	controller := controllers.NewChatControllerWithClients(
		mockChatRepo,
		mockChatUserRepo,
		mockChatRoleRepo,
		mockNotificationService,
		mockFileClient,
		mockUserClient,
	)

	// Act
	err := controller.ChangeUserRole(chatID, userID, roleID)

	// Assert
	require.Error(t, err)
	assert.True(t, errors.Is(err, custom_errors.ErrInvalidCredentials))

	mockChatUserRepo.AssertExpectations(t)
	mockChatRoleRepo.AssertNotCalled(t, "GetRoleByID", mock.Anything)
}

// Тесты для ChatController.BanUser

func TestChatController_BanUser_Success(t *testing.T) {
	// Arrange
	mockChatRepo := new(MockChatRepository)
	mockChatUserRepo := new(MockChatUserRepository)
	mockChatRoleRepo := new(MockChatRoleRepository)
	mockNotificationService := new(MockNotificationService)
	mockFileClient := new(MockFileClient)
	mockUserClient := new(MockUserClient)

	chatID := uuid.New()
	userID := uuid.New()
	chatUser := createTestChatUserWithRole(chatID, userID, 1, "main")
	bannedRole := createTestChatRoleWithID(3, "banned")

	mockChatUserRepo.On("GetChatUser", userID, chatID).Return(chatUser, nil)
	mockChatRoleRepo.On("GetRoleByName", "banned").Return(bannedRole, nil)
	mockChatUserRepo.On("ChangeUserRole", chatID, userID, bannedRole.ID).Return(nil)

	controller := controllers.NewChatControllerWithClients(
		mockChatRepo,
		mockChatUserRepo,
		mockChatRoleRepo,
		mockNotificationService,
		mockFileClient,
		mockUserClient,
	)

	// Act
	err := controller.BanUser(chatID, userID)

	// Assert
	require.NoError(t, err)

	mockChatUserRepo.AssertExpectations(t)
	mockChatRoleRepo.AssertExpectations(t)
}

func TestChatController_BanUser_UserNotInChat(t *testing.T) {
	// Arrange
	mockChatRepo := new(MockChatRepository)
	mockChatUserRepo := new(MockChatUserRepository)
	mockChatRoleRepo := new(MockChatRoleRepository)
	mockNotificationService := new(MockNotificationService)
	mockFileClient := new(MockFileClient)
	mockUserClient := new(MockUserClient)

	chatID := uuid.New()
	userID := uuid.New()

	mockChatUserRepo.On("GetChatUser", userID, chatID).Return(nil, errors.New("user not found"))

	controller := controllers.NewChatControllerWithClients(
		mockChatRepo,
		mockChatUserRepo,
		mockChatRoleRepo,
		mockNotificationService,
		mockFileClient,
		mockUserClient,
	)

	// Act
	err := controller.BanUser(chatID, userID)

	// Assert
	require.Error(t, err)
	assert.True(t, errors.Is(err, custom_errors.ErrInvalidCredentials))

	mockChatUserRepo.AssertExpectations(t)
	mockChatRoleRepo.AssertNotCalled(t, "GetRoleByName", mock.Anything)
}

func TestChatController_BanUser_BannedRoleNotFound(t *testing.T) {
	// Arrange
	mockChatRepo := new(MockChatRepository)
	mockChatUserRepo := new(MockChatUserRepository)
	mockChatRoleRepo := new(MockChatRoleRepository)
	mockNotificationService := new(MockNotificationService)
	mockFileClient := new(MockFileClient)
	mockUserClient := new(MockUserClient)

	chatID := uuid.New()
	userID := uuid.New()
	chatUser := createTestChatUserWithRole(chatID, userID, 1, "main")

	mockChatUserRepo.On("GetChatUser", userID, chatID).Return(chatUser, nil)
	mockChatRoleRepo.On("GetRoleByName", "banned").Return(nil, errors.New("role not found"))

	controller := controllers.NewChatControllerWithClients(
		mockChatRepo,
		mockChatUserRepo,
		mockChatRoleRepo,
		mockNotificationService,
		mockFileClient,
		mockUserClient,
	)

	// Act
	err := controller.BanUser(chatID, userID)

	// Assert
	require.Error(t, err)
	assert.True(t, errors.Is(err, custom_errors.ErrInternalServerError))

	mockChatUserRepo.AssertExpectations(t)
	mockChatRoleRepo.AssertExpectations(t)
}

func TestChatController_BanUser_ChangeRoleError(t *testing.T) {
	// Arrange
	mockChatRepo := new(MockChatRepository)
	mockChatUserRepo := new(MockChatUserRepository)
	mockChatRoleRepo := new(MockChatRoleRepository)
	mockNotificationService := new(MockNotificationService)
	mockFileClient := new(MockFileClient)
	mockUserClient := new(MockUserClient)

	chatID := uuid.New()
	userID := uuid.New()
	chatUser := createTestChatUserWithRole(chatID, userID, 1, "main")
	bannedRole := createTestChatRoleWithID(3, "banned")

	mockChatUserRepo.On("GetChatUser", userID, chatID).Return(chatUser, nil)
	mockChatRoleRepo.On("GetRoleByName", "banned").Return(bannedRole, nil)
	mockChatUserRepo.On("ChangeUserRole", chatID, userID, bannedRole.ID).Return(errors.New("database error"))

	controller := controllers.NewChatControllerWithClients(
		mockChatRepo,
		mockChatUserRepo,
		mockChatRoleRepo,
		mockNotificationService,
		mockFileClient,
		mockUserClient,
	)

	// Act
	err := controller.BanUser(chatID, userID)

	// Assert
	require.Error(t, err)
	assert.Equal(t, "database error", err.Error())

	mockChatUserRepo.AssertExpectations(t)
	mockChatRoleRepo.AssertExpectations(t)
}

func TestChatController_BanUser_ChatUserNil(t *testing.T) {
	// Arrange
	mockChatRepo := new(MockChatRepository)
	mockChatUserRepo := new(MockChatUserRepository)
	mockChatRoleRepo := new(MockChatRoleRepository)
	mockNotificationService := new(MockNotificationService)
	mockFileClient := new(MockFileClient)
	mockUserClient := new(MockUserClient)

	chatID := uuid.New()
	userID := uuid.New()

	mockChatUserRepo.On("GetChatUser", userID, chatID).Return(nil, nil)

	controller := controllers.NewChatControllerWithClients(
		mockChatRepo,
		mockChatUserRepo,
		mockChatRoleRepo,
		mockNotificationService,
		mockFileClient,
		mockUserClient,
	)

	// Act
	err := controller.BanUser(chatID, userID)

	// Assert
	require.Error(t, err)
	assert.True(t, errors.Is(err, custom_errors.ErrInvalidCredentials))

	mockChatUserRepo.AssertExpectations(t)
	mockChatRoleRepo.AssertNotCalled(t, "GetRoleByName", mock.Anything)
}

// Тесты для ChatController.GetUserRoleInChat

func TestChatController_GetUserRoleInChat_Success(t *testing.T) {
	// Arrange
	mockChatRepo := new(MockChatRepository)
	mockChatUserRepo := new(MockChatUserRepository)
	mockChatRoleRepo := new(MockChatRoleRepository)
	mockNotificationService := new(MockNotificationService)
	mockFileClient := new(MockFileClient)
	mockUserClient := new(MockUserClient)

	chatID := uuid.New()
	userID := uuid.New()
	requesterID := uuid.New()
	requesterChatUser := createTestChatUserWithRole(chatID, requesterID, 1, "main")
	userRole := createTestChatRoleWithID(2, "admin")

	mockChatUserRepo.On("GetChatUser", requesterID, chatID).Return(requesterChatUser, nil)
	mockChatUserRepo.On("GetUserRole", chatID, userID).Return(userRole, nil)

	controller := controllers.NewChatControllerWithClients(
		mockChatRepo,
		mockChatUserRepo,
		mockChatRoleRepo,
		mockNotificationService,
		mockFileClient,
		mockUserClient,
	)

	// Act
	roleName, err := controller.GetUserRoleInChat(chatID, userID, requesterID)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, userRole.Name, roleName)

	mockChatUserRepo.AssertExpectations(t)
}

func TestChatController_GetUserRoleInChat_RequesterNotInChat(t *testing.T) {
	// Arrange
	mockChatRepo := new(MockChatRepository)
	mockChatUserRepo := new(MockChatUserRepository)
	mockChatRoleRepo := new(MockChatRoleRepository)
	mockNotificationService := new(MockNotificationService)
	mockFileClient := new(MockFileClient)
	mockUserClient := new(MockUserClient)

	chatID := uuid.New()
	userID := uuid.New()
	requesterID := uuid.New()

	mockChatUserRepo.On("GetChatUser", requesterID, chatID).Return(nil, errors.New("user not found"))

	controller := controllers.NewChatControllerWithClients(
		mockChatRepo,
		mockChatUserRepo,
		mockChatRoleRepo,
		mockNotificationService,
		mockFileClient,
		mockUserClient,
	)

	// Act
	roleName, err := controller.GetUserRoleInChat(chatID, userID, requesterID)

	// Assert
	require.Error(t, err)
	assert.Empty(t, roleName)
	assert.True(t, errors.Is(err, custom_errors.ErrUnauthorizedChat))

	mockChatUserRepo.AssertExpectations(t)
}

func TestChatController_GetUserRoleInChat_UserNotInChat(t *testing.T) {
	// Arrange
	mockChatRepo := new(MockChatRepository)
	mockChatUserRepo := new(MockChatUserRepository)
	mockChatRoleRepo := new(MockChatRoleRepository)
	mockNotificationService := new(MockNotificationService)
	mockFileClient := new(MockFileClient)
	mockUserClient := new(MockUserClient)

	chatID := uuid.New()
	userID := uuid.New()
	requesterID := uuid.New()
	requesterChatUser := createTestChatUserWithRole(chatID, requesterID, 1, "main")

	mockChatUserRepo.On("GetChatUser", requesterID, chatID).Return(requesterChatUser, nil)
	mockChatUserRepo.On("GetUserRole", chatID, userID).Return(nil, errors.New("user not in chat"))

	controller := controllers.NewChatControllerWithClients(
		mockChatRepo,
		mockChatUserRepo,
		mockChatRoleRepo,
		mockNotificationService,
		mockFileClient,
		mockUserClient,
	)

	// Act
	roleName, err := controller.GetUserRoleInChat(chatID, userID, requesterID)

	// Assert
	require.Error(t, err)
	assert.Empty(t, roleName)
	assert.True(t, errors.Is(err, custom_errors.ErrUserNotInChat))

	mockChatUserRepo.AssertExpectations(t)
}

func TestChatController_GetUserRoleInChat_GetUserRoleError(t *testing.T) {
	// Arrange
	mockChatRepo := new(MockChatRepository)
	mockChatUserRepo := new(MockChatUserRepository)
	mockChatRoleRepo := new(MockChatRoleRepository)
	mockNotificationService := new(MockNotificationService)
	mockFileClient := new(MockFileClient)
	mockUserClient := new(MockUserClient)

	chatID := uuid.New()
	userID := uuid.New()
	requesterID := uuid.New()
	requesterChatUser := createTestChatUserWithRole(chatID, requesterID, 1, "main")

	mockChatUserRepo.On("GetChatUser", requesterID, chatID).Return(requesterChatUser, nil)
	mockChatUserRepo.On("GetUserRole", chatID, userID).Return(nil, errors.New("user not in chat"))

	controller := controllers.NewChatControllerWithClients(
		mockChatRepo,
		mockChatUserRepo,
		mockChatRoleRepo,
		mockNotificationService,
		mockFileClient,
		mockUserClient,
	)

	// Act
	roleName, err := controller.GetUserRoleInChat(chatID, userID, requesterID)

	// Assert
	require.Error(t, err)
	assert.Equal(t, "", roleName)
	assert.True(t, errors.Is(err, custom_errors.ErrUserNotInChat))

	mockChatUserRepo.AssertExpectations(t)
}

// Тесты для ChatController.GetMyRoleWithPermissions

func TestChatController_GetMyRoleWithPermissions_Success(t *testing.T) {
	// Arrange
	mockChatRepo := new(MockChatRepository)
	mockChatUserRepo := new(MockChatUserRepository)
	mockChatRoleRepo := new(MockChatRoleRepository)
	mockNotificationService := new(MockNotificationService)
	mockFileClient := new(MockFileClient)
	mockUserClient := new(MockUserClient)

	chatID := uuid.New()
	userID := uuid.New()
	chatUser := createTestChatUserWithRole(chatID, userID, 1, "main")
	chatUser.Role.Permissions = []models.ChatPermission{
		{ID: 1, Name: "send_message"},
		{ID: 2, Name: "delete_message"},
	}

	mockChatUserRepo.On("GetChatUserWithRoleAndPermissions", userID, chatID).Return(chatUser, nil)

	controller := controllers.NewChatControllerWithClients(
		mockChatRepo,
		mockChatUserRepo,
		mockChatRoleRepo,
		mockNotificationService,
		mockFileClient,
		mockUserClient,
	)

	// Act
	role, err := controller.GetMyRoleWithPermissions(chatID, userID)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, role)
	assert.Equal(t, chatUser.Role.ID, role.ID)
	assert.Equal(t, chatUser.Role.Name, role.Name)
	assert.Len(t, role.Permissions, 2)

	mockChatUserRepo.AssertExpectations(t)
}

func TestChatController_GetMyRoleWithPermissions_UserNotInChat(t *testing.T) {
	// Arrange
	mockChatRepo := new(MockChatRepository)
	mockChatUserRepo := new(MockChatUserRepository)
	mockChatRoleRepo := new(MockChatRoleRepository)
	mockNotificationService := new(MockNotificationService)
	mockFileClient := new(MockFileClient)
	mockUserClient := new(MockUserClient)

	chatID := uuid.New()
	userID := uuid.New()

	mockChatUserRepo.On("GetChatUserWithRoleAndPermissions", userID, chatID).Return(nil, errors.New("user not found"))

	controller := controllers.NewChatControllerWithClients(
		mockChatRepo,
		mockChatUserRepo,
		mockChatRoleRepo,
		mockNotificationService,
		mockFileClient,
		mockUserClient,
	)

	// Act
	role, err := controller.GetMyRoleWithPermissions(chatID, userID)

	// Assert
	require.Error(t, err)
	assert.Nil(t, role)
	assert.True(t, errors.Is(err, custom_errors.ErrUserNotInChat))

	mockChatUserRepo.AssertExpectations(t)
}

// Тесты для ChatController.GetChatMembers

func TestChatController_GetChatMembers_Success(t *testing.T) {
	// Arrange
	mockChatRepo := new(MockChatRepository)
	mockChatUserRepo := new(MockChatUserRepository)
	mockChatRoleRepo := new(MockChatRoleRepository)
	mockNotificationService := new(MockNotificationService)
	mockFileClient := new(MockFileClient)
	mockUserClient := new(MockUserClient)

	chatID := uuid.New()
	chatUsers := []models.ChatUser{
		*createTestChatUserWithRole(chatID, uuid.New(), 1, "owner"),
		*createTestChatUserWithRole(chatID, uuid.New(), 2, "main"),
	}

	mockChatUserRepo.On("GetChatUsers", chatID).Return(chatUsers, nil)

	controller := controllers.NewChatControllerWithClients(
		mockChatRepo,
		mockChatUserRepo,
		mockChatRoleRepo,
		mockNotificationService,
		mockFileClient,
		mockUserClient,
	)

	// Act
	result, err := controller.GetChatMembers(chatID)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 2)

	mockChatUserRepo.AssertExpectations(t)
}

func TestChatController_GetChatMembers_RepositoryError(t *testing.T) {
	// Arrange
	mockChatRepo := new(MockChatRepository)
	mockChatUserRepo := new(MockChatUserRepository)
	mockChatRoleRepo := new(MockChatRoleRepository)
	mockNotificationService := new(MockNotificationService)
	mockFileClient := new(MockFileClient)
	mockUserClient := new(MockUserClient)

	chatID := uuid.New()
	repoError := errors.New("database error")

	mockChatUserRepo.On("GetChatUsers", chatID).Return(nil, repoError)

	controller := controllers.NewChatControllerWithClients(
		mockChatRepo,
		mockChatUserRepo,
		mockChatRoleRepo,
		mockNotificationService,
		mockFileClient,
		mockUserClient,
	)

	// Act
	result, err := controller.GetChatMembers(chatID)

	// Assert
	require.Error(t, err)
	assert.Nil(t, result)
	var dbErr *custom_errors.DatabaseError
	assert.True(t, errors.As(err, &dbErr))

	mockChatUserRepo.AssertExpectations(t)
}

func TestChatController_CreateChat_WithMultipleUsers_WithNotifications(t *testing.T) {
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
	ownerUser := createTestUserResponseWithEmail("owner@example.com")
	ownerUser.User.ID = ownerID
	ownerUser.User.Username = "owneruser"
	user1 := createTestUserResponseWithEmail("user1@example.com")
	user1.User.ID = userID1
	user2 := createTestUserResponseWithEmail("user2@example.com")
	user2.User.ID = userID2

	createChatDTO := &dto.CreateChatDTO{
		Name:        "Group Chat",
		OwnerID:     ownerID,
		UserIDs:     []uuid.UUID{userID1, userID2},
		Description: stringPtr("Test description"),
	}

	mockChatRoleRepo.On("GetRoleByName", "owner").Return(ownerRole, nil)
	mockChatRoleRepo.On("GetRoleByName", "main").Return(mainRole, nil)
	mockUserClient.On("GetUserByID", &ownerID).Return(ownerUser, nil)
	mockUserClient.On("GetUserByID", &userID1).Return(user1, nil)
	mockUserClient.On("GetUserByID", &userID2).Return(user2, nil)
	mockChatRepo.On("CreateChat", mock.MatchedBy(func(chat *models.Chat) bool {
		return chat.IsGroup == true && chat.Name == "Group Chat"
	})).Return(nil).Run(func(args mock.Arguments) {
		chat := args.Get(0).(*models.Chat)
		chat.ID = uuid.New()
	})
	mockChatUserRepo.On("AddUserToChat", mock.MatchedBy(func(cu *models.ChatUser) bool {
		return cu.RoleID == ownerRole.ID && cu.UserID == ownerID
	})).Return(nil)
	mockChatUserRepo.On("AddUserToChat", mock.MatchedBy(func(cu *models.ChatUser) bool {
		return cu.RoleID == mainRole.ID && (cu.UserID == userID1 || cu.UserID == userID2)
	})).Return(nil).Times(2)
	mockNotificationService.On("SendChatCreatedNotification", mock.Anything, "Group Chat", "owneruser", true, "Test description", "user1@example.com").Return(nil)
	mockNotificationService.On("SendChatCreatedNotification", mock.Anything, "Group Chat", "owneruser", true, "Test description", "user2@example.com").Return(nil)

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

func TestChatController_CreateChat_NotificationError_DoesNotFail(t *testing.T) {
	// Arrange
	mockChatRepo := new(MockChatRepository)
	mockChatUserRepo := new(MockChatUserRepository)
	mockChatRoleRepo := new(MockChatRoleRepository)
	mockNotificationService := new(MockNotificationService)
	mockFileClient := new(MockFileClient)
	mockUserClient := new(MockUserClient)

	ownerID := uuid.New()
	userID1 := uuid.New()
	ownerRole := createTestChatRoleWithID(1, "owner")
	mainRole := createTestChatRoleWithID(2, "main")
	ownerUser := createTestUserResponseWithEmail("owner@example.com")
	ownerUser.User.ID = ownerID
	user1 := createTestUserResponseWithEmail("user1@example.com")
	user1.User.ID = userID1

	createChatDTO := &dto.CreateChatDTO{
		Name:    "Test Chat",
		OwnerID: ownerID,
		UserIDs: []uuid.UUID{userID1},
	}

	mockChatRoleRepo.On("GetRoleByName", "owner").Return(ownerRole, nil)
	mockChatRoleRepo.On("GetRoleByName", "main").Return(mainRole, nil)
	mockUserClient.On("GetUserByID", &ownerID).Return(ownerUser, nil)
	mockUserClient.On("GetUserByID", &userID1).Return(user1, nil)
	mockChatRepo.On("CreateChat", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		chat := args.Get(0).(*models.Chat)
		chat.ID = uuid.New()
	})
	mockChatUserRepo.On("AddUserToChat", mock.Anything).Return(nil).Times(2)
	mockNotificationService.On("SendChatCreatedNotification", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("notification error"))

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
	require.NoError(t, err) // Ошибка уведомления не должна прерывать процесс
	assert.NotNil(t, result)

	mockNotificationService.AssertExpectations(t)
}

func TestChatController_CreateChat_UserWithoutEmail_NoNotification(t *testing.T) {
	// Arrange
	mockChatRepo := new(MockChatRepository)
	mockChatUserRepo := new(MockChatUserRepository)
	mockChatRoleRepo := new(MockChatRoleRepository)
	mockNotificationService := new(MockNotificationService)
	mockFileClient := new(MockFileClient)
	mockUserClient := new(MockUserClient)

	ownerID := uuid.New()
	userID1 := uuid.New()
	ownerRole := createTestChatRoleWithID(1, "owner")
	mainRole := createTestChatRoleWithID(2, "main")
	ownerUser := createTestUserResponse()
	ownerUser.User.ID = ownerID
	user1 := createTestUserResponse()
	user1.User.ID = userID1
	user1.User.Email = "" // Пользователь без email

	createChatDTO := &dto.CreateChatDTO{
		Name:    "Test Chat",
		OwnerID: ownerID,
		UserIDs: []uuid.UUID{userID1},
	}

	mockChatRoleRepo.On("GetRoleByName", "owner").Return(ownerRole, nil)
	mockChatRoleRepo.On("GetRoleByName", "main").Return(mainRole, nil)
	mockUserClient.On("GetUserByID", &ownerID).Return(ownerUser, nil)
	mockUserClient.On("GetUserByID", &userID1).Return(user1, nil)
	mockChatRepo.On("CreateChat", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		chat := args.Get(0).(*models.Chat)
		chat.ID = uuid.New()
	})
	mockChatUserRepo.On("AddUserToChat", mock.Anything).Return(nil).Times(2)

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
	mockNotificationService.AssertNotCalled(t, "SendChatCreatedNotification", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

func TestChatController_CreateChat_AddUserToChatError(t *testing.T) {
	// Arrange
	mockChatRepo := new(MockChatRepository)
	mockChatUserRepo := new(MockChatUserRepository)
	mockChatRoleRepo := new(MockChatRoleRepository)
	mockNotificationService := new(MockNotificationService)
	mockFileClient := new(MockFileClient)
	mockUserClient := new(MockUserClient)

	ownerID := uuid.New()
	userID1 := uuid.New()
	ownerRole := createTestChatRoleWithID(1, "owner")
	mainRole := createTestChatRoleWithID(2, "main")
	ownerUser := createTestUserResponse()
	ownerUser.User.ID = ownerID
	user1 := createTestUserResponse()
	user1.User.ID = userID1

	createChatDTO := &dto.CreateChatDTO{
		Name:    "Test Chat",
		OwnerID: ownerID,
		UserIDs: []uuid.UUID{userID1},
	}

	mockChatRoleRepo.On("GetRoleByName", "owner").Return(ownerRole, nil)
	mockChatRoleRepo.On("GetRoleByName", "main").Return(mainRole, nil)
	mockUserClient.On("GetUserByID", &ownerID).Return(ownerUser, nil)
	mockUserClient.On("GetUserByID", &userID1).Return(user1, nil)
	mockChatRepo.On("CreateChat", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		chat := args.Get(0).(*models.Chat)
		chat.ID = uuid.New()
	})
	mockChatUserRepo.On("AddUserToChat", mock.MatchedBy(func(cu *models.ChatUser) bool {
		return cu.RoleID == ownerRole.ID
	})).Return(nil)
	mockChatUserRepo.On("AddUserToChat", mock.MatchedBy(func(cu *models.ChatUser) bool {
		return cu.RoleID == mainRole.ID
	})).Return(errors.New("database error"))

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
	require.Error(t, err)
	assert.Nil(t, result)
	var dbErr *custom_errors.DatabaseError
	assert.True(t, errors.As(err, &dbErr))
}

func TestChatController_GetUserChats_EmptyList(t *testing.T) {
	// Arrange
	mockChatRepo := new(MockChatRepository)
	mockChatUserRepo := new(MockChatUserRepository)
	mockChatRoleRepo := new(MockChatRoleRepository)
	mockNotificationService := new(MockNotificationService)
	mockFileClient := new(MockFileClient)
	mockUserClient := new(MockUserClient)

	userID := uuid.New()

	mockChatRepo.On("GetUserChats", userID).Return([]models.Chat{}, nil)

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
	assert.Len(t, *result, 0)

	mockChatRepo.AssertExpectations(t)
}
