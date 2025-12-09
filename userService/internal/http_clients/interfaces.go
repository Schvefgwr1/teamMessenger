package http_clients

import (
	fc "common/contracts/file-contracts"
)

// FileClientInterface - интерфейс для HTTP клиента файлового сервиса для возможности мокирования
type FileClientInterface interface {
	GetFileByID(fileID int) (*fc.File, error)
}

// ChatClientInterface - интерфейс для HTTP клиента чат-сервиса для возможности мокирования
type ChatClientInterface interface {
	GetUserRoleInChat(chatID, userID, requesterID string) (*UserRoleInChatResponse, error)
}

// UserRoleInChatResponse - ответ с ролью пользователя в чате
type UserRoleInChatResponse struct {
	RoleName string `json:"roleName"`
}
