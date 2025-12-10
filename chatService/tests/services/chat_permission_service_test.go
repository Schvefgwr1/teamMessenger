package services

import (
	"errors"
	"testing"

	"chatService/internal/models"
	"chatService/internal/services"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockChatUserRepository для тестов ChatPermissionService
type MockChatUserRepositoryForPermissionService struct {
	mock.Mock
}

func (m *MockChatUserRepositoryForPermissionService) AddUserToChat(chatUser *models.ChatUser) error {
	args := m.Called(chatUser)
	return args.Error(0)
}

func (m *MockChatUserRepositoryForPermissionService) ChangeUserRole(chatID, userID uuid.UUID, roleID int) error {
	args := m.Called(chatID, userID, roleID)
	return args.Error(0)
}

func (m *MockChatUserRepositoryForPermissionService) GetUserRole(chatID, userID uuid.UUID) (*models.ChatRole, error) {
	args := m.Called(chatID, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ChatRole), args.Error(1)
}

func (m *MockChatUserRepositoryForPermissionService) GetChatUserWithRoleAndPermissions(userID, chatID uuid.UUID) (*models.ChatUser, error) {
	args := m.Called(userID, chatID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ChatUser), args.Error(1)
}

func (m *MockChatUserRepositoryForPermissionService) GetChatUser(userID, chatID uuid.UUID) (*models.ChatUser, error) {
	args := m.Called(userID, chatID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ChatUser), args.Error(1)
}

func (m *MockChatUserRepositoryForPermissionService) GetChatUsers(chatID uuid.UUID) ([]models.ChatUser, error) {
	args := m.Called(chatID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.ChatUser), args.Error(1)
}

func (m *MockChatUserRepositoryForPermissionService) RemoveUserFromChat(chatID, userID uuid.UUID) error {
	args := m.Called(chatID, userID)
	return args.Error(0)
}

func (m *MockChatUserRepositoryForPermissionService) DeleteChatUsersByChatID(chatID uuid.UUID) error {
	args := m.Called(chatID)
	return args.Error(0)
}

func TestChatPermissionService_HasPermission_WithPermission(t *testing.T) {
	// Arrange
	mockRepo := new(MockChatUserRepositoryForPermissionService)
	service := services.NewChatPermissionService(mockRepo)

	userID := uuid.New()
	chatID := uuid.New()
	permissionName := "send_message"

	chatUser := &models.ChatUser{
		UserID: userID,
		ChatID: chatID,
		Role: models.ChatRole{
			ID:   1,
			Name: "main",
			Permissions: []models.ChatPermission{
				{ID: 1, Name: "send_message"},
				{ID: 2, Name: "delete_message"},
			},
		},
	}

	mockRepo.On("GetChatUserWithRoleAndPermissions", userID, chatID).Return(chatUser, nil)

	// Act
	hasPermission, err := service.HasPermission(userID, chatID, permissionName)

	// Assert
	require.NoError(t, err)
	assert.True(t, hasPermission)
	mockRepo.AssertExpectations(t)
}

func TestChatPermissionService_HasPermission_WithoutPermission(t *testing.T) {
	// Arrange
	mockRepo := new(MockChatUserRepositoryForPermissionService)
	service := services.NewChatPermissionService(mockRepo)

	userID := uuid.New()
	chatID := uuid.New()
	permissionName := "admin_permission"

	chatUser := &models.ChatUser{
		UserID: userID,
		ChatID: chatID,
		Role: models.ChatRole{
			ID:   1,
			Name: "main",
			Permissions: []models.ChatPermission{
				{ID: 1, Name: "send_message"},
				{ID: 2, Name: "delete_message"},
			},
		},
	}

	mockRepo.On("GetChatUserWithRoleAndPermissions", userID, chatID).Return(chatUser, nil)

	// Act
	hasPermission, err := service.HasPermission(userID, chatID, permissionName)

	// Assert
	require.NoError(t, err)
	assert.False(t, hasPermission)
	mockRepo.AssertExpectations(t)
}

func TestChatPermissionService_HasPermission_EmptyPermissions(t *testing.T) {
	// Arrange
	mockRepo := new(MockChatUserRepositoryForPermissionService)
	service := services.NewChatPermissionService(mockRepo)

	userID := uuid.New()
	chatID := uuid.New()
	permissionName := "send_message"

	chatUser := &models.ChatUser{
		UserID: userID,
		ChatID: chatID,
		Role: models.ChatRole{
			ID:          1,
			Name:        "main",
			Permissions: []models.ChatPermission{},
		},
	}

	mockRepo.On("GetChatUserWithRoleAndPermissions", userID, chatID).Return(chatUser, nil)

	// Act
	hasPermission, err := service.HasPermission(userID, chatID, permissionName)

	// Assert
	require.NoError(t, err)
	assert.False(t, hasPermission)
	mockRepo.AssertExpectations(t)
}

func TestChatPermissionService_HasPermission_UserNotInChat(t *testing.T) {
	// Arrange
	mockRepo := new(MockChatUserRepositoryForPermissionService)
	service := services.NewChatPermissionService(mockRepo)

	userID := uuid.New()
	chatID := uuid.New()
	permissionName := "send_message"

	mockRepo.On("GetChatUserWithRoleAndPermissions", userID, chatID).Return(nil, errors.New("user not in chat"))

	// Act
	hasPermission, err := service.HasPermission(userID, chatID, permissionName)

	// Assert
	require.Error(t, err)
	assert.False(t, hasPermission)
	mockRepo.AssertExpectations(t)
}
