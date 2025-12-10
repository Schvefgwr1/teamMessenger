package http_clients

import (
	cc "common/contracts/chat-contracts"
	fc "common/contracts/file-contracts"
	cuc "common/contracts/user-contracts"
	commonHttpClients "common/http_clients"
	"github.com/google/uuid"
)

// UserClientAdapter - адаптер для common/http_clients.GetUserByID
type UserClientAdapter struct{}

func NewUserClientAdapter() UserClientInterface {
	return &UserClientAdapter{}
}

func (a *UserClientAdapter) GetUserByID(userID *uuid.UUID) (*cuc.Response, error) {
	return commonHttpClients.GetUserByID(userID)
}

// ChatClientAdapter - адаптер для common/http_clients.GetChatByID
type ChatClientAdapter struct{}

func NewChatClientAdapter() ChatClientInterface {
	return &ChatClientAdapter{}
}

func (a *ChatClientAdapter) GetChatByID(chatID string) (*cc.Chat, error) {
	return commonHttpClients.GetChatByID(chatID)
}

// FileClientAdapter - адаптер для common/http_clients.GetFileByID
type FileClientAdapter struct{}

func NewFileClientAdapter() FileClientInterface {
	return &FileClientAdapter{}
}

func (a *FileClientAdapter) GetFileByID(fileID int) (*fc.File, error) {
	return commonHttpClients.GetFileByID(fileID)
}
