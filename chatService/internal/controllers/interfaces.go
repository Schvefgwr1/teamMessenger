package controllers

import (
	"chatService/internal/handlers/dto"
	"chatService/internal/models"
	ac "common/contracts/api-chat"
	"github.com/google/uuid"
)

// ChatControllerInterface - интерфейс для ChatController (для мокирования в тестах)
type ChatControllerInterface interface {
	ChangeUserRole(chatID, userID uuid.UUID, roleID int) error
	GetUserChats(userID uuid.UUID) (*[]dto.ChatResponse, error)
	CreateChat(dto *dto.CreateChatDTO) (*uuid.UUID, error)
	UpdateChat(chatID uuid.UUID, updateChatDTO *dto.UpdateChatDTO) (*dto.UpdateChatResponse, error)
	DeleteChat(chatID uuid.UUID) error
	BanUser(chatID, userID uuid.UUID) error
	GetUserRoleInChat(chatID, userID, requesterID uuid.UUID) (string, error)
	GetMyRoleWithPermissions(chatID, userID uuid.UUID) (*models.ChatRole, error)
	GetChatByID(chatID uuid.UUID) (*dto.ChatResponse, error)
	GetChatMembers(chatID uuid.UUID) ([]models.ChatUser, error)
}

// MessageControllerInterface - интерфейс для MessageController (для мокирования в тестах)
type MessageControllerInterface interface {
	SendMessage(senderID, chatID uuid.UUID, dto *dto.CreateMessageDTO) (*models.Message, error)
	GetChatMessages(chatID uuid.UUID, offset, limit int) (*[]dto.GetChatMessage, error)
	SearchMessages(userID, chatID uuid.UUID, query string, limit, offset int) (*ac.GetSearchResponse, error)
}
