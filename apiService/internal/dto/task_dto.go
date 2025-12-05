package dto

import (
	"errors"
	"mime/multipart"

	"github.com/google/uuid"
)

// CreateTaskRequestGateway - запрос на создание задачи через API Gateway
type CreateTaskRequestGateway struct {
	Title       string                  `form:"title" binding:"required"`
	Description *string                 `form:"description"`
	ExecutorID  string                  `form:"executor_id" binding:"required"`
	ChatID      *string                 `form:"chat_id"`
	Files       []*multipart.FileHeader `form:"files"`
}

// ParseUUIDs парсит строковые UUID в структуру CreateTaskRequestGateway
func (req *CreateTaskRequestGateway) ParseUUIDs() (*uuid.UUID, *uuid.UUID, error) {
	var executorID *uuid.UUID
	var chatID *uuid.UUID

	if req.ExecutorID != "" {
		parsed, err := uuid.Parse(req.ExecutorID)
		if err != nil {
			return nil, nil, err
		}
		executorID = &parsed
	} else {
		return nil, nil, errors.New("executor_id is required")
	}

	if req.ChatID != nil {
		parsed, err := uuid.Parse(*req.ChatID)
		if err != nil {
			return nil, nil, err
		}
		chatID = &parsed
	}

	return executorID, chatID, nil
}

// CreateStatusRequestGateway - запрос на создание статуса задачи через API Gateway
type CreateStatusRequestGateway struct {
	Name string `json:"name" binding:"required"`
}
