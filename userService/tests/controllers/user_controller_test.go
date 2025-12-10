package controllers

import (
	"errors"
	"testing"

	au "common/contracts/api-user"
	fc "common/contracts/file-contracts"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"userService/internal/controllers"
	"userService/internal/custom_errors"
	"userService/internal/http_clients"
	"userService/internal/models"
)

// Моки и вспомогательные функции вынесены в mocks.go

// Тесты для UserController.GetUserProfile

func TestUserController_GetUserProfile_Success(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepository)
	mockRoleRepo := new(MockRoleRepository)
	mockFileClient := new(MockFileClient)

	userID := uuid.New()
	expectedUser := createTestUser()
	expectedUser.ID = userID

	mockUserRepo.On("GetUserByID", userID).Return(expectedUser, nil)

	controller := controllers.NewUserControllerWithClients(mockUserRepo, mockRoleRepo, mockFileClient, nil)

	// Act
	user, file, err := controller.GetUserProfile(userID)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expectedUser, user)
	assert.Nil(t, file)

	mockUserRepo.AssertExpectations(t)
}

func TestUserController_GetUserProfile_WithAvatar_Success(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepository)
	mockRoleRepo := new(MockRoleRepository)
	mockFileClient := new(MockFileClient)

	userID := uuid.New()
	expectedUser := createTestUserWithAvatar()
	expectedUser.ID = userID
	expectedFile := createTestFile()

	mockUserRepo.On("GetUserByID", userID).Return(expectedUser, nil)
	mockFileClient.On("GetFileByID", *expectedUser.AvatarFileID).Return(expectedFile, nil)

	controller := controllers.NewUserControllerWithClients(mockUserRepo, mockRoleRepo, mockFileClient, nil)

	// Act
	user, file, err := controller.GetUserProfile(userID)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expectedUser, user)
	assert.Equal(t, expectedFile, file)

	mockUserRepo.AssertExpectations(t)
	mockFileClient.AssertExpectations(t)
}

func TestUserController_GetUserProfile_WithAvatar_FileError(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepository)
	mockRoleRepo := new(MockRoleRepository)
	mockFileClient := new(MockFileClient)

	userID := uuid.New()
	expectedUser := createTestUserWithAvatar()
	expectedUser.ID = userID
	fileError := errors.New("file not found")

	mockUserRepo.On("GetUserByID", userID).Return(expectedUser, nil)
	mockFileClient.On("GetFileByID", *expectedUser.AvatarFileID).Return(nil, fileError)

	controller := controllers.NewUserControllerWithClients(mockUserRepo, mockRoleRepo, mockFileClient, nil)

	// Act
	user, file, err := controller.GetUserProfile(userID)

	// Assert
	// Согласно коду контроллера, при ошибке загрузки файла возвращается user, nil file и ошибка
	// Но ошибка не критична - пользователь все равно возвращается
	require.Error(t, err) // Ошибка возвращается
	assert.Equal(t, expectedUser, user)
	assert.Nil(t, file)
	assert.Equal(t, fileError, err)

	mockUserRepo.AssertExpectations(t)
	mockFileClient.AssertExpectations(t)
}

func TestUserController_GetUserProfile_NotFound(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepository)
	mockRoleRepo := new(MockRoleRepository)
	mockFileClient := new(MockFileClient)

	userID := uuid.New()
	mockUserRepo.On("GetUserByID", userID).Return(nil, gorm.ErrRecordNotFound)

	controller := controllers.NewUserControllerWithClients(mockUserRepo, mockRoleRepo, mockFileClient, nil)

	// Act
	user, file, err := controller.GetUserProfile(userID)

	// Assert
	require.Error(t, err)
	assert.Nil(t, user)
	assert.Nil(t, file)
	assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))

	mockUserRepo.AssertExpectations(t)
}

// Тесты для UserController.UpdateUserProfile

func TestUserController_UpdateUserProfile_Success(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepository)
	mockRoleRepo := new(MockRoleRepository)
	mockFileClient := new(MockFileClient)

	userID := uuid.New()
	user := createTestUser()
	user.ID = userID
	newUsername := "newusername"

	req := &au.UpdateUserRequest{
		Username: stringPtr(newUsername),
	}

	mockUserRepo.On("GetUserByID", userID).Return(user, nil)
	mockUserRepo.On("GetUserByUsername", newUsername).Return(nil, gorm.ErrRecordNotFound)
	mockUserRepo.On("UpdateUser", mock.MatchedBy(func(u *models.User) bool {
		return u.Username == newUsername
	})).Return(nil)

	controller := controllers.NewUserControllerWithClients(mockUserRepo, mockRoleRepo, mockFileClient, nil)

	// Act
	err := controller.UpdateUserProfile(req, &userID)

	// Assert
	require.NoError(t, err)
	mockUserRepo.AssertExpectations(t)
}

func TestUserController_UpdateUserProfile_UserNotFound(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepository)
	mockRoleRepo := new(MockRoleRepository)
	mockFileClient := new(MockFileClient)

	userID := uuid.New()
	req := &au.UpdateUserRequest{
		Username: stringPtr("newusername"),
	}

	mockUserRepo.On("GetUserByID", userID).Return(nil, gorm.ErrRecordNotFound)

	controller := controllers.NewUserControllerWithClients(mockUserRepo, mockRoleRepo, mockFileClient, nil)

	// Act
	err := controller.UpdateUserProfile(req, &userID)

	// Assert
	require.Error(t, err)
	assert.True(t, errors.Is(err, custom_errors.ErrInvalidCredentials))

	mockUserRepo.AssertExpectations(t)
}

func TestUserController_UpdateUserProfile_UsernameConflict(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepository)
	mockRoleRepo := new(MockRoleRepository)
	mockFileClient := new(MockFileClient)

	userID := uuid.New()
	user := createTestUser()
	user.ID = userID
	newUsername := "existinguser"
	existingUser := createTestUser()
	existingUser.Username = newUsername

	req := &au.UpdateUserRequest{
		Username: stringPtr(newUsername),
	}

	mockUserRepo.On("GetUserByID", userID).Return(user, nil)
	mockUserRepo.On("GetUserByUsername", newUsername).Return(existingUser, nil)

	controller := controllers.NewUserControllerWithClients(mockUserRepo, mockRoleRepo, mockFileClient, nil)

	// Act
	err := controller.UpdateUserProfile(req, &userID)

	// Assert
	require.Error(t, err)
	var usernameConflictErr *custom_errors.UserUsernameConflictError
	assert.True(t, errors.As(err, &usernameConflictErr))

	mockUserRepo.AssertExpectations(t)
}

func TestUserController_UpdateUserProfile_InvalidRole(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepository)
	mockRoleRepo := new(MockRoleRepository)
	mockFileClient := new(MockFileClient)

	userID := uuid.New()
	user := createTestUser()
	user.ID = userID
	roleID := 999

	req := &au.UpdateUserRequest{
		RoleID: intPtr(roleID),
	}

	mockUserRepo.On("GetUserByID", userID).Return(user, nil)
	mockRoleRepo.On("GetRoleByID", roleID).Return(nil, gorm.ErrRecordNotFound)

	controller := controllers.NewUserControllerWithClients(mockUserRepo, mockRoleRepo, mockFileClient, nil)

	// Act
	err := controller.UpdateUserProfile(req, &userID)

	// Assert
	require.Error(t, err)
	var roleNotFoundErr *custom_errors.RoleNotFoundError
	assert.True(t, errors.As(err, &roleNotFoundErr))

	mockUserRepo.AssertExpectations(t)
	mockRoleRepo.AssertExpectations(t)
}

func TestUserController_UpdateUserProfile_ValidRole(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepository)
	mockRoleRepo := new(MockRoleRepository)
	mockFileClient := new(MockFileClient)

	userID := uuid.New()
	user := createTestUser()
	user.ID = userID
	roleID := 2
	role := createTestRole()
	roleIDVal := 2
	role.ID = &roleIDVal

	req := &au.UpdateUserRequest{
		RoleID: intPtr(roleID),
	}

	mockUserRepo.On("GetUserByID", userID).Return(user, nil)
	mockRoleRepo.On("GetRoleByID", roleID).Return(role, nil)
	mockUserRepo.On("UpdateUser", mock.MatchedBy(func(u *models.User) bool {
		return u.RoleID == roleID
	})).Return(nil)

	controller := controllers.NewUserControllerWithClients(mockUserRepo, mockRoleRepo, mockFileClient, nil)

	// Act
	err := controller.UpdateUserProfile(req, &userID)

	// Assert
	require.NoError(t, err)
	mockUserRepo.AssertExpectations(t)
	mockRoleRepo.AssertExpectations(t)
}

func TestUserController_UpdateUserProfile_InvalidAvatarFile(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepository)
	mockRoleRepo := new(MockRoleRepository)
	mockFileClient := new(MockFileClient)

	userID := uuid.New()
	user := createTestUser()
	user.ID = userID
	avatarID := 999
	fileError := errors.New("file not found")

	req := &au.UpdateUserRequest{
		AvatarFileID: intPtr(avatarID),
	}

	mockUserRepo.On("GetUserByID", userID).Return(user, nil)
	mockFileClient.On("GetFileByID", avatarID).Return(nil, fileError)

	controller := controllers.NewUserControllerWithClients(mockUserRepo, mockRoleRepo, mockFileClient, nil)

	// Act
	err := controller.UpdateUserProfile(req, &userID)

	// Assert
	require.Error(t, err)
	var getFileHTTPErr *custom_errors.GetFileHTTPError
	assert.True(t, errors.As(err, &getFileHTTPErr))

	mockUserRepo.AssertExpectations(t)
	mockFileClient.AssertExpectations(t)
}

func TestUserController_UpdateUserProfile_FileNotFound(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepository)
	mockRoleRepo := new(MockRoleRepository)
	mockFileClient := new(MockFileClient)

	userID := uuid.New()
	user := createTestUser()
	user.ID = userID
	avatarID := 999
	file := &fc.File{
		ID: 0, // Невалидный ID
	}

	req := &au.UpdateUserRequest{
		AvatarFileID: intPtr(avatarID),
	}

	mockUserRepo.On("GetUserByID", userID).Return(user, nil)
	mockFileClient.On("GetFileByID", avatarID).Return(file, nil)

	controller := controllers.NewUserControllerWithClients(mockUserRepo, mockRoleRepo, mockFileClient, nil)

	// Act
	err := controller.UpdateUserProfile(req, &userID)

	// Assert
	require.Error(t, err)
	var fileNotFoundErr *custom_errors.FileNotFoundError
	assert.True(t, errors.As(err, &fileNotFoundErr))

	mockUserRepo.AssertExpectations(t)
	mockFileClient.AssertExpectations(t)
}

func TestUserController_UpdateUserProfile_ValidAvatarFile(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepository)
	mockRoleRepo := new(MockRoleRepository)
	mockFileClient := new(MockFileClient)

	userID := uuid.New()
	user := createTestUser()
	user.ID = userID
	avatarID := 1
	file := createTestFile()

	req := &au.UpdateUserRequest{
		AvatarFileID: intPtr(avatarID),
	}

	mockUserRepo.On("GetUserByID", userID).Return(user, nil)
	mockFileClient.On("GetFileByID", avatarID).Return(file, nil)
	mockUserRepo.On("UpdateUser", mock.MatchedBy(func(u *models.User) bool {
		return u.AvatarFileID != nil && *u.AvatarFileID == file.ID
	})).Return(nil)

	controller := controllers.NewUserControllerWithClients(mockUserRepo, mockRoleRepo, mockFileClient, nil)

	// Act
	err := controller.UpdateUserProfile(req, &userID)

	// Assert
	require.NoError(t, err)
	mockUserRepo.AssertExpectations(t)
	mockFileClient.AssertExpectations(t)
}

func TestUserController_UpdateUserProfile_UpdateDescription(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepository)
	mockRoleRepo := new(MockRoleRepository)
	mockFileClient := new(MockFileClient)

	userID := uuid.New()
	user := createTestUser()
	user.ID = userID
	description := "New description"

	req := &au.UpdateUserRequest{
		Description: &description,
	}

	mockUserRepo.On("GetUserByID", userID).Return(user, nil)
	mockUserRepo.On("UpdateUser", mock.MatchedBy(func(u *models.User) bool {
		return u.Description != nil && *u.Description == description
	})).Return(nil)

	controller := controllers.NewUserControllerWithClients(mockUserRepo, mockRoleRepo, mockFileClient, nil)

	// Act
	err := controller.UpdateUserProfile(req, &userID)

	// Assert
	require.NoError(t, err)
	mockUserRepo.AssertExpectations(t)
}

// Тесты для UserController.GetUserBrief

func TestUserController_GetUserBrief_Success(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepository)
	mockRoleRepo := new(MockRoleRepository)
	mockFileClient := new(MockFileClient)
	mockChatClient := new(MockChatClient)

	userID := uuid.New()
	user := createTestUser()
	user.ID = userID
	chatID := "chat-123"
	requesterID := "requester-123"

	mockUserRepo.On("GetUserByID", userID).Return(user, nil)
	// GetUserRoleInChat может быть вызван, но ошибка игнорируется, поэтому настраиваем мок с возможностью ошибки
	mockChatClient.On("GetUserRoleInChat", chatID, userID.String(), requesterID).Return(nil, errors.New("role not found")).Maybe()

	controller := controllers.NewUserControllerWithClients(mockUserRepo, mockRoleRepo, mockFileClient, mockChatClient)

	// Act
	result, err := controller.GetUserBrief(userID, chatID, requesterID)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, user.Username, result.Username)
	assert.Equal(t, user.Email, result.Email)

	mockUserRepo.AssertExpectations(t)
}

func TestUserController_GetUserBrief_WithAvatar(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepository)
	mockRoleRepo := new(MockRoleRepository)
	mockFileClient := new(MockFileClient)
	mockChatClient := new(MockChatClient)

	userID := uuid.New()
	user := createTestUserWithAvatar()
	user.ID = userID
	chatID := "chat-123"
	requesterID := "requester-123"
	file := createTestFile()

	mockUserRepo.On("GetUserByID", userID).Return(user, nil)
	mockFileClient.On("GetFileByID", *user.AvatarFileID).Return(file, nil)
	// GetUserRoleInChat может быть вызван, но ошибка игнорируется, поэтому настраиваем мок с возможностью ошибки
	mockChatClient.On("GetUserRoleInChat", chatID, userID.String(), requesterID).Return(nil, errors.New("role not found")).Maybe()

	controller := controllers.NewUserControllerWithClients(mockUserRepo, mockRoleRepo, mockFileClient, mockChatClient)

	// Act
	result, err := controller.GetUserBrief(userID, chatID, requesterID)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, file, result.AvatarFile)

	mockUserRepo.AssertExpectations(t)
	mockFileClient.AssertExpectations(t)
}

func TestUserController_GetUserBrief_WithChatRole(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepository)
	mockRoleRepo := new(MockRoleRepository)
	mockFileClient := new(MockFileClient)
	mockChatClient := new(MockChatClient)

	userID := uuid.New()
	user := createTestUser()
	user.ID = userID
	chatID := "chat-123"
	requesterID := "requester-123"
	roleResp := &http_clients.UserRoleInChatResponse{
		RoleName: "admin",
	}

	mockUserRepo.On("GetUserByID", userID).Return(user, nil)
	mockChatClient.On("GetUserRoleInChat", chatID, userID.String(), requesterID).Return(roleResp, nil)

	controller := controllers.NewUserControllerWithClients(mockUserRepo, mockRoleRepo, mockFileClient, mockChatClient)

	// Act
	result, err := controller.GetUserBrief(userID, chatID, requesterID)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, roleResp.RoleName, result.ChatRoleName)

	mockUserRepo.AssertExpectations(t)
	mockChatClient.AssertExpectations(t)
}

func TestUserController_GetUserBrief_UserNotFound(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepository)
	mockRoleRepo := new(MockRoleRepository)
	mockFileClient := new(MockFileClient)
	mockChatClient := new(MockChatClient)

	userID := uuid.New()
	chatID := "chat-123"
	requesterID := "requester-123"

	mockUserRepo.On("GetUserByID", userID).Return(nil, gorm.ErrRecordNotFound)

	controller := controllers.NewUserControllerWithClients(mockUserRepo, mockRoleRepo, mockFileClient, mockChatClient)

	// Act
	result, err := controller.GetUserBrief(userID, chatID, requesterID)

	// Assert
	require.Error(t, err)
	assert.Nil(t, result)
	assert.True(t, errors.Is(err, custom_errors.ErrInvalidCredentials))

	mockUserRepo.AssertExpectations(t)
}

func TestUserController_GetUserBrief_EmptyChatID(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepository)
	mockRoleRepo := new(MockRoleRepository)
	mockFileClient := new(MockFileClient)
	mockChatClient := new(MockChatClient)

	userID := uuid.New()
	user := createTestUser()
	user.ID = userID
	chatID := ""
	requesterID := "requester-123"

	mockUserRepo.On("GetUserByID", userID).Return(user, nil)

	controller := controllers.NewUserControllerWithClients(mockUserRepo, mockRoleRepo, mockFileClient, mockChatClient)

	// Act
	result, err := controller.GetUserBrief(userID, chatID, requesterID)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Empty(t, result.ChatRoleName)

	mockUserRepo.AssertExpectations(t)
	mockChatClient.AssertNotCalled(t, "GetUserRoleInChat", mock.Anything, mock.Anything, mock.Anything)
}

// Тесты для UserController.SearchUsers

func TestUserController_SearchUsers_Success(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepository)
	mockRoleRepo := new(MockRoleRepository)
	mockFileClient := new(MockFileClient)

	query := "test"
	limit := 10
	users := []*models.User{
		createTestUser(),
		createTestUser(),
	}

	mockUserRepo.On("SearchUsers", query, limit).Return(users, nil)

	controller := controllers.NewUserControllerWithClients(mockUserRepo, mockRoleRepo, mockFileClient, nil)

	// Act
	result, err := controller.SearchUsers(query, limit)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Users, 2)

	mockUserRepo.AssertExpectations(t)
}

func TestUserController_SearchUsers_WithAvatars(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepository)
	mockRoleRepo := new(MockRoleRepository)
	mockFileClient := new(MockFileClient)

	query := "test"
	limit := 10
	user1 := createTestUserWithAvatar()
	user2 := createTestUser()
	users := []*models.User{user1, user2}
	file := createTestFile()

	mockUserRepo.On("SearchUsers", query, limit).Return(users, nil)
	mockFileClient.On("GetFileByID", *user1.AvatarFileID).Return(file, nil)

	controller := controllers.NewUserControllerWithClients(mockUserRepo, mockRoleRepo, mockFileClient, nil)

	// Act
	result, err := controller.SearchUsers(query, limit)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Users, 2)
	assert.Equal(t, file, result.Users[0].AvatarFile)
	assert.Nil(t, result.Users[1].AvatarFile)

	mockUserRepo.AssertExpectations(t)
	mockFileClient.AssertExpectations(t)
}

func TestUserController_SearchUsers_LimitNormalization(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepository)
	mockRoleRepo := new(MockRoleRepository)
	mockFileClient := new(MockFileClient)

	query := "test"
	limit := 0 // Должен быть нормализован до 10
	users := []*models.User{}

	mockUserRepo.On("SearchUsers", query, 10).Return(users, nil)

	controller := controllers.NewUserControllerWithClients(mockUserRepo, mockRoleRepo, mockFileClient, nil)

	// Act
	result, err := controller.SearchUsers(query, limit)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)

	mockUserRepo.AssertExpectations(t)
}

func TestUserController_SearchUsers_LimitMax(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepository)
	mockRoleRepo := new(MockRoleRepository)
	mockFileClient := new(MockFileClient)

	query := "test"
	limit := 25 // Должен быть нормализован до 10
	users := []*models.User{}

	mockUserRepo.On("SearchUsers", query, 10).Return(users, nil)

	controller := controllers.NewUserControllerWithClients(mockUserRepo, mockRoleRepo, mockFileClient, nil)

	// Act
	result, err := controller.SearchUsers(query, limit)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)

	mockUserRepo.AssertExpectations(t)
}

func TestUserController_SearchUsers_RepositoryError(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepository)
	mockRoleRepo := new(MockRoleRepository)
	mockFileClient := new(MockFileClient)

	query := "test"
	limit := 10
	repoError := errors.New("database error")

	mockUserRepo.On("SearchUsers", query, limit).Return(nil, repoError)

	controller := controllers.NewUserControllerWithClients(mockUserRepo, mockRoleRepo, mockFileClient, nil)

	// Act
	result, err := controller.SearchUsers(query, limit)

	// Assert
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, repoError, err)

	mockUserRepo.AssertExpectations(t)
}

// Тесты для UserController.UpdateUserRole

func TestUserController_UpdateUserRole_Success(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepository)
	mockRoleRepo := new(MockRoleRepository)
	mockFileClient := new(MockFileClient)

	userID := uuid.New()
	user := createTestUser()
	user.ID = userID
	roleID := 2
	role := createTestRole()
	roleIDVal := 2
	role.ID = &roleIDVal

	mockUserRepo.On("GetUserByID", userID).Return(user, nil)
	mockRoleRepo.On("GetRoleByID", roleID).Return(role, nil)
	mockUserRepo.On("UpdateUser", mock.MatchedBy(func(u *models.User) bool {
		return u.RoleID == roleID
	})).Return(nil)

	controller := controllers.NewUserControllerWithClients(mockUserRepo, mockRoleRepo, mockFileClient, nil)

	// Act
	err := controller.UpdateUserRole(userID, roleID)

	// Assert
	require.NoError(t, err)
	mockUserRepo.AssertExpectations(t)
	mockRoleRepo.AssertExpectations(t)
}

func TestUserController_UpdateUserRole_UserNotFound(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepository)
	mockRoleRepo := new(MockRoleRepository)
	mockFileClient := new(MockFileClient)

	userID := uuid.New()
	roleID := 2

	mockUserRepo.On("GetUserByID", userID).Return(nil, gorm.ErrRecordNotFound)

	controller := controllers.NewUserControllerWithClients(mockUserRepo, mockRoleRepo, mockFileClient, nil)

	// Act
	err := controller.UpdateUserRole(userID, roleID)

	// Assert
	require.Error(t, err)
	assert.True(t, errors.Is(err, custom_errors.ErrInvalidCredentials))

	mockUserRepo.AssertExpectations(t)
	mockRoleRepo.AssertNotCalled(t, "GetRoleByID", mock.Anything)
}

func TestUserController_UpdateUserRole_RoleNotFound(t *testing.T) {
	// Arrange
	mockUserRepo := new(MockUserRepository)
	mockRoleRepo := new(MockRoleRepository)
	mockFileClient := new(MockFileClient)

	userID := uuid.New()
	user := createTestUser()
	user.ID = userID
	roleID := 999

	mockUserRepo.On("GetUserByID", userID).Return(user, nil)
	mockRoleRepo.On("GetRoleByID", roleID).Return(nil, gorm.ErrRecordNotFound)

	controller := controllers.NewUserControllerWithClients(mockUserRepo, mockRoleRepo, mockFileClient, nil)

	// Act
	err := controller.UpdateUserRole(userID, roleID)

	// Assert
	require.Error(t, err)
	var roleNotFoundErr *custom_errors.RoleNotFoundError
	assert.True(t, errors.As(err, &roleNotFoundErr))

	mockUserRepo.AssertExpectations(t)
	mockRoleRepo.AssertExpectations(t)
}
