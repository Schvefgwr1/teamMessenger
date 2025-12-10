package http_clients

import (
	cc "common/contracts/chat-contracts"
	fc "common/contracts/file-contracts"
	cuc "common/contracts/user-contracts"
	"github.com/google/uuid"
)

// UserClientInterface - интерфейс для HTTP клиента пользовательского сервиса для возможности мокирования
type UserClientInterface interface {
	GetUserByID(userID *uuid.UUID) (*cuc.Response, error)
}

// ChatClientInterface - интерфейс для HTTP клиента чат-сервиса для возможности мокирования
type ChatClientInterface interface {
	GetChatByID(chatID string) (*cc.Chat, error)
}

// FileClientInterface - интерфейс для HTTP клиента файлового сервиса для возможности мокирования
type FileClientInterface interface {
	GetFileByID(fileID int) (*fc.File, error)
}
