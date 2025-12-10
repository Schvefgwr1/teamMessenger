package http_clients

import (
	fc "common/contracts/file-contracts"
	cuc "common/contracts/user-contracts"
	commonHttpClients "common/http_clients"
	"github.com/google/uuid"
)

// FileClientInterface - интерфейс для HTTP клиента файлового сервиса для возможности мокирования
type FileClientInterface interface {
	GetFileByID(fileID int) (*fc.File, error)
}

// UserClientInterface - интерфейс для HTTP клиента пользовательского сервиса для возможности мокирования
type UserClientInterface interface {
	GetUserByID(userID *uuid.UUID) (*cuc.Response, error)
}

// FileClientAdapter - адаптер для обертки функций из common/http_clients
type FileClientAdapter struct{}

func NewFileClientAdapter() FileClientInterface {
	return &FileClientAdapter{}
}

func (f *FileClientAdapter) GetFileByID(fileID int) (*fc.File, error) {
	return commonHttpClients.GetFileByID(fileID)
}

// UserClientAdapter - адаптер для обертки функций из common/http_clients
type UserClientAdapter struct{}

func NewUserClientAdapter() UserClientInterface {
	return &UserClientAdapter{}
}

func (u *UserClientAdapter) GetUserByID(userID *uuid.UUID) (*cuc.Response, error) {
	return commonHttpClients.GetUserByID(userID)
}
