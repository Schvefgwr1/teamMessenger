package controllers

import (
	fc "common/contracts/file-contracts"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"userService/internal/http_clients"
	"userService/internal/models"
)

// MockUserRepository - мок для UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) CreateUser(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) GetUserByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetUserByUsername(username string) (*models.User, error) {
	args := m.Called(username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetUserByID(id uuid.UUID) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) UpdateUser(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) SearchUsers(query string, limit int) ([]*models.User, error) {
	args := m.Called(query, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.User), args.Error(1)
}

// MockRoleRepository - мок для RoleRepository
type MockRoleRepository struct {
	mock.Mock
}

func (m *MockRoleRepository) GetAllRoles() ([]models.Role, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Role), args.Error(1)
}

func (m *MockRoleRepository) GetRoleByID(id int) (*models.Role, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Role), args.Error(1)
}

func (m *MockRoleRepository) CreateRole(role *models.Role) error {
	args := m.Called(role)
	return args.Error(0)
}

func (m *MockRoleRepository) DeleteRole(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockRoleRepository) UpdateRolePermissions(roleID int, permissionIDs []int) error {
	args := m.Called(roleID, permissionIDs)
	return args.Error(0)
}

// MockPermissionRepository - мок для PermissionRepository
type MockPermissionRepository struct {
	mock.Mock
}

func (m *MockPermissionRepository) GetAllPermissions() ([]models.Permission, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Permission), args.Error(1)
}

func (m *MockPermissionRepository) GetPermissionById(id int) (*models.Permission, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Permission), args.Error(1)
}

// MockFileClient - мок для FileClientInterface
type MockFileClient struct {
	mock.Mock
}

func (m *MockFileClient) GetFileByID(fileID int) (*fc.File, error) {
	args := m.Called(fileID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*fc.File), args.Error(1)
}

// MockChatClient - мок для ChatClientInterface
type MockChatClient struct {
	mock.Mock
}

func (m *MockChatClient) GetUserRoleInChat(chatID, userID, requesterID string) (*http_clients.UserRoleInChatResponse, error) {
	args := m.Called(chatID, userID, requesterID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*http_clients.UserRoleInChatResponse), args.Error(1)
}

// MockNotificationService - мок для NotificationServiceInterface
type MockNotificationService struct {
	mock.Mock
}

func (m *MockNotificationService) SendLoginNotification(
	userID uuid.UUID,
	username string,
	email string,
	ipAddress string,
	userAgent string,
) error {
	args := m.Called(userID, username, email, ipAddress, userAgent)
	return args.Error(0)
}

func (m *MockNotificationService) Close() error {
	args := m.Called()
	return args.Error(0)
}

// Вспомогательные функции для создания тестовых данных
func createTestUser() *models.User {
	userID := uuid.New()
	username := "testuser"
	email := "test@example.com"
	roleID := 1
	return &models.User{
		ID:       userID,
		Username: username,
		Email:    email,
		RoleID:   roleID,
	}
}

func createTestUserWithAvatar() *models.User {
	user := createTestUser()
	avatarID := 1
	user.AvatarFileID = &avatarID
	return user
}

func createTestRole() *models.Role {
	roleID := 1
	roleName := "user"
	return &models.Role{
		ID:   &roleID,
		Name: roleName,
	}
}

func createTestRoleWithID(id int) *models.Role {
	return &models.Role{
		ID:   &id,
		Name: "user",
	}
}

func createTestFile() *fc.File {
	return &fc.File{
		ID:  1,
		URL: "http://example.com/avatar.jpg",
	}
}

func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}
