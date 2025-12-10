package controllers

import (
	"apiService/internal/dto"
	"apiService/internal/services"
	ac "common/contracts/api-chat"
	af "common/contracts/api-file"
	at "common/contracts/api-task"
	au "common/contracts/api-user"
	uc "common/contracts/user-contracts"
	"context"
	"crypto/rsa"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"mime/multipart"
	"time"
)

// MockUserClient - мок для UserClient
type MockUserClient struct {
	mock.Mock
}

func (m *MockUserClient) RegisterUser(data au.RegisterUserRequest) (*au.RegisterUserResponse, error) {
	args := m.Called(data)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*au.RegisterUserResponse), args.Error(1)
}

func (m *MockUserClient) Login(body *au.Login) (string, uuid.UUID, error) {
	args := m.Called(body)
	return args.String(0), args.Get(1).(uuid.UUID), args.Error(2)
}

func (m *MockUserClient) GetUserByID(s string) (*au.GetUserResponse, error) {
	args := m.Called(s)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*au.GetUserResponse), args.Error(1)
}

func (m *MockUserClient) UpdateUser(userID string, req *au.UpdateUserRequest) (*au.UpdateUserResponse, error) {
	args := m.Called(userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*au.UpdateUserResponse), args.Error(1)
}

func (m *MockUserClient) GetPublicKey() (*rsa.PublicKey, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*rsa.PublicKey), args.Error(1)
}

func (m *MockUserClient) GetAllPermissions() ([]*uc.Permission, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*uc.Permission), args.Error(1)
}

func (m *MockUserClient) GetAllRoles() ([]*uc.Role, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*uc.Role), args.Error(1)
}

func (m *MockUserClient) CreateRole(req *au.CreateRoleRequest) (*uc.Role, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*uc.Role), args.Error(1)
}

func (m *MockUserClient) UpdateUserRole(userID string, roleID int) error {
	args := m.Called(userID, roleID)
	return args.Error(0)
}

func (m *MockUserClient) UpdateRolePermissions(roleID int, permissionIDs []int) error {
	args := m.Called(roleID, permissionIDs)
	return args.Error(0)
}

func (m *MockUserClient) DeleteRole(roleID int) error {
	args := m.Called(roleID)
	return args.Error(0)
}

func (m *MockUserClient) GetUserBrief(userID, chatID, requesterID string) (*dto.UserBriefResponse, error) {
	args := m.Called(userID, chatID, requesterID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.UserBriefResponse), args.Error(1)
}

func (m *MockUserClient) SearchUsers(query string, limit int) (*dto.UserSearchResponse, error) {
	args := m.Called(query, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.UserSearchResponse), args.Error(1)
}

// MockFileClient - мок для FileClient
type MockFileClient struct {
	mock.Mock
}

func (m *MockFileClient) UploadFile(file *multipart.FileHeader) (*af.FileUploadResponse, error) {
	args := m.Called(file)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*af.FileUploadResponse), args.Error(1)
}

// MockCacheService - мок для CacheService
type MockCacheService struct {
	mock.Mock
}

func (m *MockCacheService) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	args := m.Called(ctx, key, value, ttl)
	return args.Error(0)
}

func (m *MockCacheService) Get(ctx context.Context, key string, dest interface{}) error {
	args := m.Called(ctx, key, dest)
	return args.Error(0)
}

func (m *MockCacheService) Delete(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

func (m *MockCacheService) DeleteByPattern(ctx context.Context, pattern string) error {
	args := m.Called(ctx, pattern)
	return args.Error(0)
}

func (m *MockCacheService) Exists(ctx context.Context, key string) (bool, error) {
	args := m.Called(ctx, key)
	return args.Bool(0), args.Error(1)
}

func (m *MockCacheService) SetTTL(ctx context.Context, key string, ttl time.Duration) error {
	args := m.Called(ctx, key, ttl)
	return args.Error(0)
}

func (m *MockCacheService) GetTTL(ctx context.Context, key string) (time.Duration, error) {
	args := m.Called(ctx, key)
	return args.Get(0).(time.Duration), args.Error(1)
}

func (m *MockCacheService) SetUserCache(ctx context.Context, userID string, userData interface{}) error {
	args := m.Called(ctx, userID, userData)
	return args.Error(0)
}

func (m *MockCacheService) GetUserCache(ctx context.Context, userID string, dest interface{}) error {
	args := m.Called(ctx, userID, dest)
	return args.Error(0)
}

func (m *MockCacheService) DeleteUserCache(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockCacheService) SetChatMessagesCache(ctx context.Context, chatID string, messages interface{}) error {
	args := m.Called(ctx, chatID, messages)
	return args.Error(0)
}

func (m *MockCacheService) GetChatMessagesCache(ctx context.Context, chatID string, dest interface{}) error {
	args := m.Called(ctx, chatID, dest)
	return args.Error(0)
}

func (m *MockCacheService) DeleteChatMessagesCache(ctx context.Context, chatID string) error {
	args := m.Called(ctx, chatID)
	return args.Error(0)
}

func (m *MockCacheService) SetUserChatListCache(ctx context.Context, userID string, chats interface{}) error {
	args := m.Called(ctx, userID, chats)
	return args.Error(0)
}

func (m *MockCacheService) GetUserChatListCache(ctx context.Context, userID string, dest interface{}) error {
	args := m.Called(ctx, userID, dest)
	return args.Error(0)
}

func (m *MockCacheService) DeleteUserChatListCache(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockCacheService) ChatInfoCacheKey(chatID string) string {
	args := m.Called(chatID)
	return args.String(0)
}

func (m *MockCacheService) SetChatInfoCache(ctx context.Context, chatID string, chatInfo interface{}) error {
	args := m.Called(ctx, chatID, chatInfo)
	return args.Error(0)
}

func (m *MockCacheService) GetChatInfoCache(ctx context.Context, chatID string, dest interface{}) error {
	args := m.Called(ctx, chatID, dest)
	return args.Error(0)
}

func (m *MockCacheService) DeleteChatInfoCache(ctx context.Context, chatID string) error {
	args := m.Called(ctx, chatID)
	return args.Error(0)
}

func (m *MockCacheService) ChatMembersCacheKey(chatID string) string {
	args := m.Called(chatID)
	return args.String(0)
}

func (m *MockCacheService) DeleteChatMembersCache(ctx context.Context, chatID string) error {
	args := m.Called(ctx, chatID)
	return args.Error(0)
}

func (m *MockCacheService) ChatUserRoleCacheKey(chatID, userID string) string {
	args := m.Called(chatID, userID)
	return args.String(0)
}

func (m *MockCacheService) DeleteChatUserRoleCache(ctx context.Context, chatID, userID string) error {
	args := m.Called(ctx, chatID, userID)
	return args.Error(0)
}

func (m *MockCacheService) GetChatUserRoleCache(ctx context.Context, chatID, userID string, dest interface{}) error {
	args := m.Called(ctx, chatID, userID, dest)
	return args.Error(0)
}

func (m *MockCacheService) SetChatUserRoleCache(ctx context.Context, chatID, userID string, role interface{}) error {
	args := m.Called(ctx, chatID, userID, role)
	return args.Error(0)
}

func (m *MockCacheService) TaskCacheKey(taskID int) string {
	args := m.Called(taskID)
	return args.String(0)
}

func (m *MockCacheService) UserTasksCacheKey(userID string) string {
	args := m.Called(userID)
	return args.String(0)
}

func (m *MockCacheService) SetTaskCache(ctx context.Context, taskID int, task interface{}) error {
	args := m.Called(ctx, taskID, task)
	return args.Error(0)
}

func (m *MockCacheService) GetTaskCache(ctx context.Context, taskID int, dest interface{}) error {
	args := m.Called(ctx, taskID, dest)
	return args.Error(0)
}

func (m *MockCacheService) DeleteTaskCache(ctx context.Context, taskID int) error {
	args := m.Called(ctx, taskID)
	return args.Error(0)
}

func (m *MockCacheService) SetUserTasksCache(ctx context.Context, userID string, tasks interface{}) error {
	args := m.Called(ctx, userID, tasks)
	return args.Error(0)
}

func (m *MockCacheService) GetUserTasksCache(ctx context.Context, userID string, dest interface{}) error {
	args := m.Called(ctx, userID, dest)
	return args.Error(0)
}

func (m *MockCacheService) DeleteUserTasksCache(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockCacheService) SetChatRolesCache(ctx context.Context, roles interface{}) error {
	args := m.Called(ctx, roles)
	return args.Error(0)
}

func (m *MockCacheService) GetChatRolesCache(ctx context.Context, dest interface{}) error {
	args := m.Called(ctx, dest)
	return args.Error(0)
}

func (m *MockCacheService) DeleteChatRolesCache(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockCacheService) SetChatPermissionsCache(ctx context.Context, permissions interface{}) error {
	args := m.Called(ctx, permissions)
	return args.Error(0)
}

func (m *MockCacheService) GetChatPermissionsCache(ctx context.Context, dest interface{}) error {
	args := m.Called(ctx, dest)
	return args.Error(0)
}

func (m *MockCacheService) DeleteChatPermissionsCache(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockCacheService) SearchCacheKey(chatID, queryHash string) string {
	args := m.Called(chatID, queryHash)
	return args.String(0)
}

func (m *MockCacheService) SetSearchCache(ctx context.Context, chatID, queryHash string, result interface{}) error {
	args := m.Called(ctx, chatID, queryHash, result)
	return args.Error(0)
}

func (m *MockCacheService) GetSearchCache(ctx context.Context, chatID, queryHash string, dest interface{}) error {
	args := m.Called(ctx, chatID, queryHash, dest)
	return args.Error(0)
}

func (m *MockCacheService) DeleteSearchCacheByChat(ctx context.Context, chatID string) error {
	args := m.Called(ctx, chatID)
	return args.Error(0)
}

func (m *MockCacheService) UserCacheKey(userID string) string {
	args := m.Called(userID)
	return args.String(0)
}

func (m *MockCacheService) ChatMessagesCacheKey(chatID string) string {
	args := m.Called(chatID)
	return args.String(0)
}

func (m *MockCacheService) UserChatListCacheKey(userID string) string {
	args := m.Called(userID)
	return args.String(0)
}

// MockSessionService - мок для SessionService
type MockSessionService struct {
	mock.Mock
}

func (m *MockSessionService) CreateSession(ctx context.Context, userID uuid.UUID, token string, expiresAt time.Time) error {
	args := m.Called(ctx, userID, token, expiresAt)
	return args.Error(0)
}

func (m *MockSessionService) GetSession(ctx context.Context, userID uuid.UUID, token string) (*services.JWTSession, error) {
	args := m.Called(ctx, userID, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*services.JWTSession), args.Error(1)
}

func (m *MockSessionService) RevokeSession(ctx context.Context, userID uuid.UUID, token string) error {
	args := m.Called(ctx, userID, token)
	return args.Error(0)
}

func (m *MockSessionService) RevokeAllUserSessions(ctx context.Context, userID uuid.UUID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockSessionService) IsSessionValid(ctx context.Context, userID uuid.UUID, token string) (bool, error) {
	args := m.Called(ctx, userID, token)
	return args.Bool(0), args.Error(1)
}

// Вспомогательные функции для создания тестовых данных
func createTestUserResponse() *au.GetUserResponse {
	userID := uuid.New()
	return &au.GetUserResponse{
		User: &uc.User{
			ID:       userID,
			Username: "testuser",
			Email:    "test@example.com",
		},
		File: nil,
	}
}

func createTestUpdateUserResponse() *au.UpdateUserResponse {
	return &au.UpdateUserResponse{
		Error:   nil,
		Message: nil,
	}
}

func createTestRole() *uc.Role {
	return &uc.Role{
		ID:   1,
		Name: "user",
	}
}

func createTestPermission() *uc.Permission {
	return &uc.Permission{
		ID:   1,
		Name: "read",
	}
}

func stringPtr(s string) *string {
	return &s
}

// MockChatClient - мок для ChatClient
type MockChatClient struct {
	mock.Mock
}

func (m *MockChatClient) GetUserChats(userID uuid.UUID) ([]*ac.ChatResponse, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*ac.ChatResponse), args.Error(1)
}

func (m *MockChatClient) CreateChat(req *ac.CreateChatRequest) (*ac.CreateChatServiceResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ac.CreateChatServiceResponse), args.Error(1)
}

func (m *MockChatClient) SendMessage(chatID uuid.UUID, senderID uuid.UUID, req *ac.CreateMessageRequest) (*ac.MessageResponse, error) {
	args := m.Called(chatID, senderID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ac.MessageResponse), args.Error(1)
}

func (m *MockChatClient) GetChatMessages(chatID uuid.UUID, userID uuid.UUID, offset, limit int) ([]*ac.GetChatMessage, error) {
	args := m.Called(chatID, userID, offset, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*ac.GetChatMessage), args.Error(1)
}

func (m *MockChatClient) SearchMessages(userID uuid.UUID, chatID uuid.UUID, query string, offset, limit int) (*ac.GetSearchResponse, error) {
	args := m.Called(userID, chatID, query, offset, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ac.GetSearchResponse), args.Error(1)
}

func (m *MockChatClient) UpdateChat(chatID uuid.UUID, updateReq *ac.UpdateChatRequest, userID uuid.UUID) (*ac.UpdateChatResponse, error) {
	args := m.Called(chatID, updateReq, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ac.UpdateChatResponse), args.Error(1)
}

func (m *MockChatClient) DeleteChat(chatID, userID uuid.UUID) error {
	args := m.Called(chatID, userID)
	return args.Error(0)
}

func (m *MockChatClient) BanUser(chatID, userID, ownerID uuid.UUID) error {
	args := m.Called(chatID, userID, ownerID)
	return args.Error(0)
}

func (m *MockChatClient) ChangeUserRole(chatID, ownerID uuid.UUID, changeRoleReq *ac.ChangeRoleRequest) error {
	args := m.Called(chatID, ownerID, changeRoleReq)
	return args.Error(0)
}

func (m *MockChatClient) GetMyRoleInChat(chatID, userID uuid.UUID) (*ac.MyRoleResponse, error) {
	args := m.Called(chatID, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ac.MyRoleResponse), args.Error(1)
}

func (m *MockChatClient) GetChatMembers(chatID uuid.UUID) ([]*ac.ChatMember, error) {
	args := m.Called(chatID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*ac.ChatMember), args.Error(1)
}

// MockTaskClient - мок для TaskClient
type MockTaskClient struct {
	mock.Mock
}

func (m *MockTaskClient) CreateTask(req *at.CreateTaskRequest) (*at.TaskResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*at.TaskResponse), args.Error(1)
}

func (m *MockTaskClient) UpdateTaskStatus(taskID, statusID int) error {
	args := m.Called(taskID, statusID)
	return args.Error(0)
}

func (m *MockTaskClient) GetTaskByID(taskID int) (*at.TaskServiceResponse, error) {
	args := m.Called(taskID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*at.TaskServiceResponse), args.Error(1)
}

func (m *MockTaskClient) GetUserTasks(userID string, limit, offset int) (*[]at.TaskToList, error) {
	args := m.Called(userID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*[]at.TaskToList), args.Error(1)
}

func (m *MockTaskClient) GetAllStatuses() ([]at.TaskStatus, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]at.TaskStatus), args.Error(1)
}

func (m *MockTaskClient) CreateStatus(req *at.CreateStatusRequest) (*at.TaskStatus, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*at.TaskStatus), args.Error(1)
}

func (m *MockTaskClient) GetStatusByID(statusID int) (*at.TaskStatus, error) {
	args := m.Called(statusID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*at.TaskStatus), args.Error(1)
}

func (m *MockTaskClient) DeleteStatus(statusID int) error {
	args := m.Called(statusID)
	return args.Error(0)
}

// MockChatRolePermissionClient - мок для ChatRolePermissionClient
type MockChatRolePermissionClient struct {
	mock.Mock
}

func (m *MockChatRolePermissionClient) GetAllRoles() ([]ac.RoleResponse, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]ac.RoleResponse), args.Error(1)
}

func (m *MockChatRolePermissionClient) GetRoleByID(roleID int) (*ac.RoleResponse, error) {
	args := m.Called(roleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ac.RoleResponse), args.Error(1)
}

func (m *MockChatRolePermissionClient) CreateRole(req *ac.CreateRoleRequest) (*ac.RoleResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ac.RoleResponse), args.Error(1)
}

func (m *MockChatRolePermissionClient) DeleteRole(roleID int) error {
	args := m.Called(roleID)
	return args.Error(0)
}

func (m *MockChatRolePermissionClient) UpdateRolePermissions(roleID int, req *ac.UpdateRolePermissionsRequest) (*ac.RoleResponse, error) {
	args := m.Called(roleID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ac.RoleResponse), args.Error(1)
}

func (m *MockChatRolePermissionClient) GetAllPermissions() ([]ac.PermissionResponse, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]ac.PermissionResponse), args.Error(1)
}

func (m *MockChatRolePermissionClient) CreatePermission(req *ac.CreatePermissionRequest) (*ac.PermissionResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ac.PermissionResponse), args.Error(1)
}

func (m *MockChatRolePermissionClient) DeletePermission(permissionID int) error {
	args := m.Called(permissionID)
	return args.Error(0)
}
