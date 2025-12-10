package handlers

import (
	"apiService/internal/dto"
	ac "common/contracts/api-chat"
	at "common/contracts/api-task"
	au "common/contracts/api-user"
	uc "common/contracts/user-contracts"
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"mime/multipart"
)

// MockUserController - мок для UserController
type MockUserController struct {
	mock.Mock
}

func (m *MockUserController) GetUser(userID uuid.UUID) (*au.GetUserResponse, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*au.GetUserResponse), args.Error(1)
}

func (m *MockUserController) UpdateUser(userID uuid.UUID, req *dto.UpdateUserRequestGateway, file *multipart.FileHeader) (*au.UpdateUserResponse, error) {
	args := m.Called(userID, req, file)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*au.UpdateUserResponse), args.Error(1)
}

func (m *MockUserController) GetAllPermissions() ([]*uc.Permission, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*uc.Permission), args.Error(1)
}

func (m *MockUserController) GetAllRoles() ([]*uc.Role, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*uc.Role), args.Error(1)
}

func (m *MockUserController) CreateRole(req *au.CreateRoleRequest) (*uc.Role, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*uc.Role), args.Error(1)
}

func (m *MockUserController) UpdateUserRole(userID uuid.UUID, roleID int) error {
	args := m.Called(userID, roleID)
	return args.Error(0)
}

func (m *MockUserController) UpdateRolePermissions(roleID int, permissionIDs []int) error {
	args := m.Called(roleID, permissionIDs)
	return args.Error(0)
}

func (m *MockUserController) DeleteRole(roleID int) error {
	args := m.Called(roleID)
	return args.Error(0)
}

func (m *MockUserController) GetUserProfileByID(userID uuid.UUID) (*au.GetUserResponse, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*au.GetUserResponse), args.Error(1)
}

func (m *MockUserController) GetUserBrief(userID uuid.UUID, chatID string, requesterID uuid.UUID) (*dto.UserBriefResponse, error) {
	args := m.Called(userID, chatID, requesterID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.UserBriefResponse), args.Error(1)
}

func (m *MockUserController) SearchUsers(query string, limit int) (*dto.UserSearchResponse, error) {
	args := m.Called(query, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.UserSearchResponse), args.Error(1)
}

// MockAuthController - мок для AuthController
type MockAuthController struct {
	mock.Mock
}

func (m *MockAuthController) Register(req *dto.RegisterUserRequestGateway, file *multipart.FileHeader) *dto.RegisterUserResponseGateway {
	args := m.Called(req, file)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*dto.RegisterUserResponseGateway)
}

func (m *MockAuthController) Login(ctx context.Context, loginData *au.Login) (string, uuid.UUID, error) {
	args := m.Called(ctx, loginData)
	return args.String(0), args.Get(1).(uuid.UUID), args.Error(2)
}

func (m *MockAuthController) Logout(ctx context.Context, userID uuid.UUID, token string) error {
	args := m.Called(ctx, userID, token)
	return args.Error(0)
}

// MockChatController - мок для ChatController
type MockChatController struct {
	mock.Mock
}

func (m *MockChatController) GetUserChats(userID uuid.UUID) ([]*ac.ChatResponse, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*ac.ChatResponse), args.Error(1)
}

func (m *MockChatController) CreateChat(req *dto.CreateChatRequestGateway, ownerID uuid.UUID, userIDs []uuid.UUID) (*dto.CreateChatResponse, error) {
	args := m.Called(req, ownerID, userIDs)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.CreateChatResponse), args.Error(1)
}

func (m *MockChatController) SendMessage(chatID uuid.UUID, senderID uuid.UUID, req *dto.SendMessageRequestGateway) (*ac.MessageResponse, error) {
	args := m.Called(chatID, senderID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ac.MessageResponse), args.Error(1)
}

func (m *MockChatController) GetChatMessages(chatID uuid.UUID, userID uuid.UUID, offset, limit int) ([]*ac.GetChatMessage, error) {
	args := m.Called(chatID, userID, offset, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*ac.GetChatMessage), args.Error(1)
}

func (m *MockChatController) SearchMessages(userID uuid.UUID, chatID uuid.UUID, query string, offset, limit int) (*ac.GetSearchResponse, error) {
	args := m.Called(userID, chatID, query, offset, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ac.GetSearchResponse), args.Error(1)
}

func (m *MockChatController) UpdateChat(chatID uuid.UUID, req *dto.UpdateChatRequestGateway, updateReq *ac.UpdateChatRequest, userID uuid.UUID) (*ac.UpdateChatResponse, error) {
	args := m.Called(chatID, req, updateReq, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ac.UpdateChatResponse), args.Error(1)
}

func (m *MockChatController) DeleteChat(chatID, userID uuid.UUID) error {
	args := m.Called(chatID, userID)
	return args.Error(0)
}

func (m *MockChatController) BanUser(chatID, userID, ownerID uuid.UUID) error {
	args := m.Called(chatID, userID, ownerID)
	return args.Error(0)
}

func (m *MockChatController) ChangeUserRole(chatID, ownerID uuid.UUID, changeRoleReq *ac.ChangeRoleRequest) error {
	args := m.Called(chatID, ownerID, changeRoleReq)
	return args.Error(0)
}

func (m *MockChatController) GetMyRoleInChat(chatID, userID uuid.UUID) (*ac.MyRoleResponse, error) {
	args := m.Called(chatID, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ac.MyRoleResponse), args.Error(1)
}

func (m *MockChatController) GetChatMembers(chatID uuid.UUID) ([]*ac.ChatMember, error) {
	args := m.Called(chatID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*ac.ChatMember), args.Error(1)
}

// MockTaskController - мок для TaskController
type MockTaskController struct {
	mock.Mock
}

func (m *MockTaskController) CreateTask(req *dto.CreateTaskRequestGateway, creatorID uuid.UUID) (*at.TaskResponse, error) {
	args := m.Called(req, creatorID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*at.TaskResponse), args.Error(1)
}

func (m *MockTaskController) UpdateTaskStatus(taskID, statusID int) error {
	args := m.Called(taskID, statusID)
	return args.Error(0)
}

func (m *MockTaskController) GetTaskByID(taskID int) (*at.TaskServiceResponse, error) {
	args := m.Called(taskID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*at.TaskServiceResponse), args.Error(1)
}

func (m *MockTaskController) GetUserTasks(userID string, limit, offset int) (*[]at.TaskToList, error) {
	args := m.Called(userID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*[]at.TaskToList), args.Error(1)
}

func (m *MockTaskController) GetAllStatuses() ([]at.TaskStatus, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]at.TaskStatus), args.Error(1)
}

func (m *MockTaskController) CreateStatus(statusName string) (*at.TaskStatus, error) {
	args := m.Called(statusName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*at.TaskStatus), args.Error(1)
}

func (m *MockTaskController) GetStatusByID(statusID int) (*at.TaskStatus, error) {
	args := m.Called(statusID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*at.TaskStatus), args.Error(1)
}

func (m *MockTaskController) DeleteStatus(statusID int) error {
	args := m.Called(statusID)
	return args.Error(0)
}

// MockChatRolePermissionController - мок для ChatRolePermissionController
type MockChatRolePermissionController struct {
	mock.Mock
}

func (m *MockChatRolePermissionController) GetAllRoles() ([]dto.ChatRoleResponseGateway, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]dto.ChatRoleResponseGateway), args.Error(1)
}

func (m *MockChatRolePermissionController) GetRoleByID(roleID int) (*dto.ChatRoleResponseGateway, error) {
	args := m.Called(roleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.ChatRoleResponseGateway), args.Error(1)
}

func (m *MockChatRolePermissionController) CreateRole(req *dto.CreateChatRoleRequestGateway) (*dto.ChatRoleResponseGateway, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.ChatRoleResponseGateway), args.Error(1)
}

func (m *MockChatRolePermissionController) DeleteRole(roleID int) error {
	args := m.Called(roleID)
	return args.Error(0)
}

func (m *MockChatRolePermissionController) UpdateRolePermissions(roleID int, req *dto.UpdateChatRolePermissionsRequestGateway) (*dto.ChatRoleResponseGateway, error) {
	args := m.Called(roleID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.ChatRoleResponseGateway), args.Error(1)
}

func (m *MockChatRolePermissionController) GetAllPermissions() ([]dto.ChatPermissionResponseGateway, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]dto.ChatPermissionResponseGateway), args.Error(1)
}

func (m *MockChatRolePermissionController) CreatePermission(req *dto.CreateChatPermissionRequestGateway) (*dto.ChatPermissionResponseGateway, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.ChatPermissionResponseGateway), args.Error(1)
}

func (m *MockChatRolePermissionController) DeletePermission(permissionID int) error {
	args := m.Called(permissionID)
	return args.Error(0)
}
