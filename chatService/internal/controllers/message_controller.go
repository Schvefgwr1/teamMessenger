package controllers

import (
	"chatService/internal/custom_errors"
	"chatService/internal/handlers/dto"
	"chatService/internal/models"
	"chatService/internal/repositories"
	ac "common/contracts/api-chat"
	fc "common/contracts/file-contracts"
	httpClients "common/http_clients"
	"github.com/google/uuid"
	"time"
)

type MessageController struct {
	MessageRepo  repositories.MessageRepository
	ChatRepo     repositories.ChatRepository
	ChatUserRepo repositories.ChatUserRepository
}

func NewMessageController(messageRepo repositories.MessageRepository, chatRepo repositories.ChatRepository, chatUserRepo repositories.ChatUserRepository) *MessageController {
	return &MessageController{messageRepo, chatRepo, chatUserRepo}
}

func (c *MessageController) SendMessage(senderID, chatID uuid.UUID, dto *dto.CreateMessageDTO) (*models.Message, error) {
	_, err := c.ChatRepo.GetChatByID(chatID)
	if err != nil {
		return nil, custom_errors.ErrInvalidCredentials
	}

	userResp, err := httpClients.GetUserByID(&senderID)
	if err != nil {
		return nil, custom_errors.NewUserClientError(err.Error())
	}
	if userResp.User == nil {
		return nil, custom_errors.NewUserClientError("sender not found")
	}

	for _, fileID := range dto.FileIDs {
		file, err := httpClients.GetFileByID(fileID)
		if err != nil {
			return nil, custom_errors.NewGetFileHTTPError(fileID, err.Error())
		}
		if file.ID <= 0 {
			return nil, custom_errors.NewFileNotFoundError(fileID)
		}
	}

	msg := &models.Message{
		ID:        uuid.New(),
		ChatID:    chatID,
		SenderID:  &senderID,
		Content:   dto.Content,
		CreatedAt: time.Now(),
	}

	if err := c.MessageRepo.CreateMessage(msg); err != nil {
		return nil, custom_errors.NewDatabaseError(err.Error())
	}

	for _, fileID := range dto.FileIDs {
		messageFile := &models.MessageFile{
			MessageID: msg.ID,
			FileID:    fileID,
		}
		if err := c.MessageRepo.CreateMessageFile(messageFile); err != nil {
			return nil, custom_errors.NewDatabaseError(err.Error())
		}
	}

	newMsg, errMsg := c.MessageRepo.GetMessageWithFile(msg.ID)
	if errMsg != nil {
		return nil, custom_errors.NewDatabaseError(errMsg.Error())
	}
	return newMsg, nil
}

func (c *MessageController) GetChatMessages(chatID uuid.UUID, offset, limit int) (*[]dto.GetChatMessage, error) {
	_, err := c.ChatRepo.GetChatByID(chatID)
	if err != nil {
		return nil, custom_errors.ErrInvalidCredentials
	}

	messages, err := c.MessageRepo.GetChatMessages(chatID, offset, limit)
	if err != nil {
		return nil, custom_errors.NewDatabaseError(err.Error())
	}
	var messagesResponse []dto.GetChatMessage
	for _, message := range messages {
		var files []*fc.File
		for _, file := range message.Files {
			fileHTTP, err := httpClients.GetFileByID(file.FileID)
			if err != nil {
				return nil, custom_errors.NewGetFileHTTPError(file.FileID, err.Error())
			}
			files = append(files, fileHTTP)
		}
		messagesResponse = append(messagesResponse, dto.GetChatMessage{
			ID:        message.ID,
			ChatID:    message.ChatID,
			SenderID:  message.SenderID,
			Content:   message.Content,
			UpdatedAt: message.UpdatedAt,
			CreatedAt: message.CreatedAt,
			Files:     &files,
		})
	}
	return &messagesResponse, nil
}

func (c *MessageController) SearchMessages(userID, chatID uuid.UUID, query string, limit, offset int) (*ac.GetSearchResponse, error) {
	if query == "" {
		return nil, custom_errors.ErrEmptyQuery
	}

	_, err := c.ChatRepo.GetChatByID(chatID)
	if err != nil {
		return nil, custom_errors.ErrChatNotFound
	}

	user, err := c.ChatUserRepo.GetChatUser(userID, chatID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, custom_errors.ErrUnauthorizedChat
	}

	messages, total, err := c.MessageRepo.SearchMessages(userID, chatID, query, limit, offset)
	if err != nil {
		return nil, custom_errors.NewDatabaseError(err.Error())
	}

	var messageResponse []ac.GetChatMessage
	for _, message := range messages {
		messageResponse = append(messageResponse, ac.GetChatMessage{
			ID:        message.ID,
			ChatID:    message.ChatID,
			SenderID:  message.SenderID,
			Content:   message.Content,
			UpdatedAt: message.UpdatedAt,
			CreatedAt: message.CreatedAt,
			Files:     nil,
		})
	}

	return &ac.GetSearchResponse{Messages: &messageResponse, Total: &total}, nil
}
