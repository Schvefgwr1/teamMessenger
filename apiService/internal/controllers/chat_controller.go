package controllers

import (
	"apiService/internal/dto"
	"apiService/internal/http_clients"
	"apiService/internal/services"
	ac "common/contracts/api-chat"
	"context"
	"errors"
	"github.com/google/uuid"
	"log"
	"time"
)

type ChatController struct {
	chatClient   http_clients.ChatClient
	fileClient   http_clients.FileClient
	cacheService *services.CacheService
}

func NewChatController(chatClient http_clients.ChatClient, fileClient http_clients.FileClient, cacheService *services.CacheService) *ChatController {
	return &ChatController{
		chatClient:   chatClient,
		fileClient:   fileClient,
		cacheService: cacheService,
	}
}

func (ctrl *ChatController) GetUserChats(userID uuid.UUID) ([]*ac.ChatResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Пытаемся получить из кеша
	var cachedChats []*ac.ChatResponse
	err := ctrl.cacheService.GetUserChatListCache(ctx, userID.String(), &cachedChats)
	if err == nil {
		return cachedChats, nil
	}

	chats, err := ctrl.chatClient.GetUserChats(userID)
	if err != nil {
		return nil, err
	}

	// Сохраняем в кеш
	if err := ctrl.cacheService.SetUserChatListCache(ctx, userID.String(), chats); err != nil {
		log.Printf("Failed to cache user chats for %s: %v", userID.String(), err)
	}

	return chats, nil
}

func (ctrl *ChatController) CreateChat(req *dto.CreateChatRequestGateway, ownerID uuid.UUID, userIDs []uuid.UUID) (*dto.CreateChatResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	createReq := &ac.CreateChatRequest{
		Name:        req.Name,
		Description: req.Description,
		OwnerID:     ownerID,
		UserIDs:     userIDs,
	}

	if req.Avatar != nil {
		uploadedFile, err := ctrl.fileClient.UploadFile(req.Avatar)
		if err != nil {
			return nil, err
		}
		createReq.AvatarFileID = uploadedFile.ID
	}

	serviceResp, err := ctrl.chatClient.CreateChat(createReq)
	if err != nil {
		return nil, errors.New("error of chat client")
	}

	// Инвалидируем кеш списков чатов для всех участников
	allUserIDs := append(userIDs, ownerID)
	for _, userID := range allUserIDs {
		if err := ctrl.cacheService.DeleteUserChatListCache(ctx, userID.String()); err != nil {
			log.Printf("Failed to invalidate chat list cache for user %s: %v", userID.String(), err)
		}
	}

	return &dto.CreateChatResponse{
		ID:           serviceResp.ChatID,
		Name:         createReq.Name,
		Description:  createReq.Description,
		OwnerID:      createReq.OwnerID,
		UserIDs:      createReq.UserIDs,
		AvatarFileID: createReq.AvatarFileID,
	}, nil
}

func (ctrl *ChatController) SendMessage(chatID uuid.UUID, senderID uuid.UUID, req *dto.SendMessageRequestGateway) (*ac.MessageResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	var fileIDs []int

	if len(req.Files) > 0 {
		fileIDs = make([]int, 0, len(req.Files))
		for _, file := range req.Files {
			uploadedFile, err := ctrl.fileClient.UploadFile(file)
			if err != nil {
				log.Printf("failed to upload file %s: %v\n", file.Filename, err)
				continue
			}
			if uploadedFile.ID != nil {
				fileIDs = append(fileIDs, *uploadedFile.ID)
			}
		}
	}

	createReq := &ac.CreateMessageRequest{
		Content: req.Content,
		FileIDs: fileIDs,
	}

	message, err := ctrl.chatClient.SendMessage(chatID, senderID, createReq)
	if err != nil {
		return nil, err
	}

	// Инвалидируем кеш сообщений для этого чата
	if err := ctrl.cacheService.DeleteChatMessagesCache(ctx, chatID.String()); err != nil {
		log.Printf("Failed to invalidate messages cache for chat %s: %v", chatID.String(), err)
	}

	return message, nil
}

func (ctrl *ChatController) GetChatMessages(chatID uuid.UUID, userID uuid.UUID, offset, limit int) ([]*ac.GetChatMessage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Кешируем только последние 20 сообщений (offset = 0, limit = 20)
	if offset == 0 && limit <= 20 {
		var cachedMessages []*ac.GetChatMessage
		err := ctrl.cacheService.GetChatMessagesCache(ctx, chatID.String(), &cachedMessages)
		if err == nil {
			log.Printf("Messages for chat %s found in cache", chatID.String())
			if limit < len(cachedMessages) {
				return cachedMessages[:limit], nil
			}
			return cachedMessages, nil
		}

		messages, err := ctrl.chatClient.GetChatMessages(chatID, userID, 0, 20) // Всегда запрашиваем последние 20
		if err != nil {
			return nil, err
		}

		// Сохраняем в кеш
		if err := ctrl.cacheService.SetChatMessagesCache(ctx, chatID.String(), messages); err != nil {
			log.Printf("Failed to cache messages for chat %s: %v", chatID.String(), err)
		}

		// Возвращаем запрошенное количество
		if limit < len(messages) {
			return messages[:limit], nil
		}
		return messages, nil
	}

	return ctrl.chatClient.GetChatMessages(chatID, userID, offset, limit)
}

func (ctrl *ChatController) SearchMessages(userID uuid.UUID, chatID uuid.UUID, query string, offset, limit int) (*ac.GetSearchResponse, error) {
	return ctrl.chatClient.SearchMessages(userID, chatID, query, offset, limit)
}
