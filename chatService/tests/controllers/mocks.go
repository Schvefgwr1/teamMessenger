package controllers

import (
	"chatService/internal/models"
	fc "common/contracts/file-contracts"
	cuc "common/contracts/user-contracts"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"time"
)

// MockChatRepository - мок для ChatRepository
type MockChatRepository struct {
	mock.Mock
}

func (m *MockChatRepository) GetUserChats(userID uuid.UUID) ([]models.Chat, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Chat), args.Error(1)
}

func (m *MockChatRepository) CreateChat(chat *models.Chat) error {
	args := m.Called(chat)
	return args.Error(0)
}

func (m *MockChatRepository) UpdateChat(chat *models.Chat) error {
	args := m.Called(chat)
	return args.Error(0)
}

func (m *MockChatRepository) DeleteChat(chatID uuid.UUID) error {
	args := m.Called(chatID)
	return args.Error(0)
}

func (m *MockChatRepository) GetChatByID(chatID uuid.UUID) (*models.Chat, error) {
	args := m.Called(chatID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Chat), args.Error(1)
}

// MockChatUserRepository - мок для ChatUserRepository
type MockChatUserRepository struct {
	mock.Mock
}

func (m *MockChatUserRepository) AddUserToChat(chatUser *models.ChatUser) error {
	args := m.Called(chatUser)
	return args.Error(0)
}

func (m *MockChatUserRepository) ChangeUserRole(chatID, userID uuid.UUID, roleID int) error {
	args := m.Called(chatID, userID, roleID)
	return args.Error(0)
}

func (m *MockChatUserRepository) GetUserRole(chatID, userID uuid.UUID) (*models.ChatRole, error) {
	args := m.Called(chatID, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ChatRole), args.Error(1)
}

func (m *MockChatUserRepository) GetChatUserWithRoleAndPermissions(userID, chatID uuid.UUID) (*models.ChatUser, error) {
	args := m.Called(userID, chatID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ChatUser), args.Error(1)
}

func (m *MockChatUserRepository) GetChatUser(userID, chatID uuid.UUID) (*models.ChatUser, error) {
	args := m.Called(userID, chatID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ChatUser), args.Error(1)
}

func (m *MockChatUserRepository) GetChatUsers(chatID uuid.UUID) ([]models.ChatUser, error) {
	args := m.Called(chatID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.ChatUser), args.Error(1)
}

func (m *MockChatUserRepository) RemoveUserFromChat(chatID, userID uuid.UUID) error {
	args := m.Called(chatID, userID)
	return args.Error(0)
}

func (m *MockChatUserRepository) DeleteChatUsersByChatID(chatID uuid.UUID) error {
	args := m.Called(chatID)
	return args.Error(0)
}

// MockChatRoleRepository - мок для ChatRoleRepository
type MockChatRoleRepository struct {
	mock.Mock
}

func (m *MockChatRoleRepository) GetRoleByID(roleID int) (*models.ChatRole, error) {
	args := m.Called(roleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ChatRole), args.Error(1)
}

func (m *MockChatRoleRepository) GetRoleByName(roleName string) (*models.ChatRole, error) {
	args := m.Called(roleName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ChatRole), args.Error(1)
}

func (m *MockChatRoleRepository) GetAllRoles() ([]models.ChatRole, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.ChatRole), args.Error(1)
}

func (m *MockChatRoleRepository) CreateRole(role *models.ChatRole, permissionIDs []int) error {
	args := m.Called(role, permissionIDs)
	return args.Error(0)
}

func (m *MockChatRoleRepository) DeleteRole(roleID int) error {
	args := m.Called(roleID)
	return args.Error(0)
}

func (m *MockChatRoleRepository) UpdateRolePermissions(roleID int, permissionIDs []int) error {
	args := m.Called(roleID, permissionIDs)
	return args.Error(0)
}

// MockMessageRepository - мок для MessageRepository
type MockMessageRepository struct {
	mock.Mock
}

func (m *MockMessageRepository) CreateMessage(message *models.Message) error {
	args := m.Called(message)
	return args.Error(0)
}

func (m *MockMessageRepository) CreateMessageFile(msgFile *models.MessageFile) error {
	args := m.Called(msgFile)
	return args.Error(0)
}

func (m *MockMessageRepository) GetMessageWithFile(msgID uuid.UUID) (*models.Message, error) {
	args := m.Called(msgID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Message), args.Error(1)
}

func (m *MockMessageRepository) GetChatMessages(chatID uuid.UUID, offset, limit int) ([]models.Message, error) {
	args := m.Called(chatID, offset, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Message), args.Error(1)
}

func (m *MockMessageRepository) SearchMessages(userID, chatID uuid.UUID, text string, limit, offset int) ([]models.Message, int64, error) {
	args := m.Called(userID, chatID, text, limit, offset)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]models.Message), args.Get(1).(int64), args.Error(2)
}

// MockNotificationService - мок для NotificationService
type MockNotificationService struct {
	mock.Mock
}

func (m *MockNotificationService) SendChatCreatedNotification(
	chatID uuid.UUID,
	chatName string,
	creatorName string,
	isGroup bool,
	description string,
	userEmail string,
) error {
	args := m.Called(chatID, chatName, creatorName, isGroup, description, userEmail)
	return args.Error(0)
}

func (m *MockNotificationService) Close() error {
	args := m.Called()
	return args.Error(0)
}

// Вспомогательные функции для создания тестовых данных
func createTestChat() *models.Chat {
	chatID := uuid.New()
	description := "Test chat description"
	return &models.Chat{
		ID:          chatID,
		Name:        "Test Chat",
		IsGroup:     true,
		Description: &description,
		CreatedAt:   time.Now(),
	}
}

func createTestChatWithoutAvatar() *models.Chat {
	chatID := uuid.New()
	return &models.Chat{
		ID:        chatID,
		Name:      "Test Chat",
		IsGroup:   false,
		CreatedAt: time.Now(),
	}
}

func createTestChatRole() *models.ChatRole {
	return &models.ChatRole{
		ID:          1,
		Name:        "owner",
		Permissions: []models.ChatPermission{},
	}
}

func createTestChatRoleWithID(id int, name string) *models.ChatRole {
	return &models.ChatRole{
		ID:          id,
		Name:        name,
		Permissions: []models.ChatPermission{},
	}
}

func createTestChatUser() *models.ChatUser {
	return &models.ChatUser{
		ChatID: uuid.New(),
		UserID: uuid.New(),
		RoleID: 1,
		Role:   *createTestChatRole(),
	}
}

func createTestChatUserWithRole(chatID, userID uuid.UUID, roleID int, roleName string) *models.ChatUser {
	return &models.ChatUser{
		ChatID: chatID,
		UserID: userID,
		RoleID: roleID,
		Role: models.ChatRole{
			ID:          roleID,
			Name:        roleName,
			Permissions: []models.ChatPermission{},
		},
	}
}

func createTestMessage() *models.Message {
	return &models.Message{
		ID:        uuid.New(),
		ChatID:    uuid.New(),
		SenderID:  func() *uuid.UUID { id := uuid.New(); return &id }(),
		Content:   "Test message",
		CreatedAt: time.Now(),
		Files:     []models.MessageFile{},
	}
}

func createTestFile() *fc.File {
	return &fc.File{
		ID:   1,
		Name: "test.jpg",
		URL:  "http://example.com/test.jpg",
	}
}

func createTestUserResponse() *cuc.Response {
	userID := uuid.New()
	return &cuc.Response{
		User: &cuc.User{
			ID:       userID,
			Username: "testuser",
			Email:    "test@example.com",
		},
	}
}

func createTestUserResponseWithEmail(email string) *cuc.Response {
	userID := uuid.New()
	return &cuc.Response{
		User: &cuc.User{
			ID:       userID,
			Username: "testuser",
			Email:    email,
		},
	}
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

// MockUserClient - мок для UserClientInterface
type MockUserClient struct {
	mock.Mock
}

func (m *MockUserClient) GetUserByID(userID *uuid.UUID) (*cuc.Response, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*cuc.Response), args.Error(1)
}

func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}

// MockChatPermissionRepository - мок для ChatPermissionRepository
type MockChatPermissionRepository struct {
	mock.Mock
}

func (m *MockChatPermissionRepository) GetAllPermissions() ([]models.ChatPermission, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.ChatPermission), args.Error(1)
}

func (m *MockChatPermissionRepository) GetPermissionByID(id int) (*models.ChatPermission, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ChatPermission), args.Error(1)
}

func (m *MockChatPermissionRepository) GetPermissionByName(name string) (*models.ChatPermission, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ChatPermission), args.Error(1)
}

func (m *MockChatPermissionRepository) CreatePermission(permission *models.ChatPermission) error {
	args := m.Called(permission)
	return args.Error(0)
}

func (m *MockChatPermissionRepository) DeletePermission(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

// Вспомогательные функции для создания тестовых данных для permissions
func createTestChatPermission() *models.ChatPermission {
	return &models.ChatPermission{
		ID:   1,
		Name: "send_message",
	}
}

func createTestChatPermissionWithID(id int, name string) *models.ChatPermission {
	return &models.ChatPermission{
		ID:   id,
		Name: name,
	}
}

func createTestChatRoleWithPermissions(roleID int, roleName string, permissions []models.ChatPermission) *models.ChatRole {
	return &models.ChatRole{
		ID:          roleID,
		Name:        roleName,
		Permissions: permissions,
	}
}
