package controllers

import (
	"apiService/internal/dto"
	ac "common/contracts/api-chat"
	at "common/contracts/api-task"
	au "common/contracts/api-user"
	uc "common/contracts/user-contracts"
	"context"
	"github.com/google/uuid"
	"mime/multipart"
)

// ChatControllerInterface - интерфейс для ChatController
type ChatControllerInterface interface {
	GetUserChats(userID uuid.UUID) ([]*ac.ChatResponse, error)
	CreateChat(req *dto.CreateChatRequestGateway, ownerID uuid.UUID, userIDs []uuid.UUID) (*dto.CreateChatResponse, error)
	SendMessage(chatID uuid.UUID, senderID uuid.UUID, req *dto.SendMessageRequestGateway) (*ac.MessageResponse, error)
	GetChatMessages(chatID uuid.UUID, userID uuid.UUID, offset, limit int) ([]*ac.GetChatMessage, error)
	SearchMessages(userID uuid.UUID, chatID uuid.UUID, query string, offset, limit int) (*ac.GetSearchResponse, error)
	UpdateChat(chatID uuid.UUID, req *dto.UpdateChatRequestGateway, updateReq *ac.UpdateChatRequest, userID uuid.UUID) (*ac.UpdateChatResponse, error)
	DeleteChat(chatID, userID uuid.UUID) error
	BanUser(chatID, userID, ownerID uuid.UUID) error
	ChangeUserRole(chatID, ownerID uuid.UUID, changeRoleReq *ac.ChangeRoleRequest) error
	GetMyRoleInChat(chatID, userID uuid.UUID) (*ac.MyRoleResponse, error)
	GetChatMembers(chatID uuid.UUID) ([]*ac.ChatMember, error)
}

// UserControllerInterface - интерфейс для UserController
type UserControllerInterface interface {
	GetUser(id uuid.UUID) (*au.GetUserResponse, error)
	UpdateUser(id uuid.UUID, userRequest *dto.UpdateUserRequestGateway, file *multipart.FileHeader) (*au.UpdateUserResponse, error)
	GetAllPermissions() ([]*uc.Permission, error)
	GetAllRoles() ([]*uc.Role, error)
	CreateRole(req *au.CreateRoleRequest) (*uc.Role, error)
	UpdateUserRole(userID uuid.UUID, roleID int) error
	UpdateRolePermissions(roleID int, permissionIDs []int) error
	DeleteRole(roleID int) error
	GetUserProfileByID(userID uuid.UUID) (*au.GetUserResponse, error)
	GetUserBrief(userID uuid.UUID, chatID string, requesterID uuid.UUID) (*dto.UserBriefResponse, error)
	SearchUsers(query string, limit int) (*dto.UserSearchResponse, error)
}

// AuthControllerInterface - интерфейс для AuthController
type AuthControllerInterface interface {
	Register(req *dto.RegisterUserRequestGateway, file *multipart.FileHeader) *dto.RegisterUserResponseGateway
	Login(ctx context.Context, loginData *au.Login) (string, uuid.UUID, error)
	Logout(ctx context.Context, userID uuid.UUID, token string) error
}

// TaskControllerInterface - интерфейс для TaskController
type TaskControllerInterface interface {
	CreateTask(req *dto.CreateTaskRequestGateway, creatorID uuid.UUID) (*at.TaskResponse, error)
	UpdateTaskStatus(taskID, statusID int) error
	GetTaskByID(taskID int) (*at.TaskServiceResponse, error)
	GetUserTasks(userID string, limit, offset int) (*[]at.TaskToList, error)
	GetAllStatuses() ([]at.TaskStatus, error)
	CreateStatus(statusName string) (*at.TaskStatus, error)
	GetStatusByID(statusID int) (*at.TaskStatus, error)
	DeleteStatus(statusID int) error
}

// ChatRolePermissionControllerInterface - интерфейс для ChatRolePermissionController
type ChatRolePermissionControllerInterface interface {
	GetAllRoles() ([]dto.ChatRoleResponseGateway, error)
	GetRoleByID(roleID int) (*dto.ChatRoleResponseGateway, error)
	CreateRole(req *dto.CreateChatRoleRequestGateway) (*dto.ChatRoleResponseGateway, error)
	DeleteRole(roleID int) error
	UpdateRolePermissions(roleID int, req *dto.UpdateChatRolePermissionsRequestGateway) (*dto.ChatRoleResponseGateway, error)
	GetAllPermissions() ([]dto.ChatPermissionResponseGateway, error)
	CreatePermission(req *dto.CreateChatPermissionRequestGateway) (*dto.ChatPermissionResponseGateway, error)
	DeletePermission(permissionID int) error
}
