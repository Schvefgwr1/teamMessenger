package http_clients

import (
	fc "common/contracts/file-contracts"
	commonHttpClients "common/http_clients"
)

// FileClientAdapter - адаптер для common/http_clients.GetFileByID
type FileClientAdapter struct{}

func NewFileClientAdapter() FileClientInterface {
	return &FileClientAdapter{}
}

func (a *FileClientAdapter) GetFileByID(fileID int) (*fc.File, error) {
	return commonHttpClients.GetFileByID(fileID)
}

// ChatClientAdapter - адаптер для common/http_clients.GetUserRoleInChat
type ChatClientAdapter struct{}

func NewChatClientAdapter() ChatClientInterface {
	return &ChatClientAdapter{}
}

func (a *ChatClientAdapter) GetUserRoleInChat(chatID, userID, requesterID string) (*UserRoleInChatResponse, error) {
	commonResp, err := commonHttpClients.GetUserRoleInChat(chatID, userID, requesterID)
	if err != nil {
		return nil, err
	}
	if commonResp == nil {
		return nil, nil
	}
	return &UserRoleInChatResponse{
		RoleName: commonResp.RoleName,
	}, nil
}
