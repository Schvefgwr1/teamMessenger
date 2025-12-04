package controllers

import (
	"apiService/internal/dto"
	"apiService/internal/http_clients"
	"apiService/internal/services"
	ac "common/contracts/api-chat"
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
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

	// Инвалидируем кеш поиска для этого чата
	if err := ctrl.cacheService.DeleteSearchCacheByChat(ctx, chatID.String()); err != nil {
		log.Printf("Failed to invalidate search cache for chat %s: %v", chatID.String(), err)
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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Кешируем только первую страницу поиска (offset = 0)
	if offset == 0 && limit <= 20 {
		// Создаём хеш из query для ключа кеша
		queryHash := fmt.Sprintf("%x", query) // Простой хеш для короткого ключа

		var cachedResult ac.GetSearchResponse
		err := ctrl.cacheService.GetSearchCache(ctx, chatID.String(), queryHash, &cachedResult)
		if err == nil {
			log.Printf("Search results for chat %s found in cache", chatID.String())
			return &cachedResult, nil
		}

		// Получаем из сервиса
		result, err := ctrl.chatClient.SearchMessages(userID, chatID, query, offset, limit)
		if err != nil {
			return nil, err
		}

		// Сохраняем в кеш
		if result != nil {
			if err := ctrl.cacheService.SetSearchCache(ctx, chatID.String(), queryHash, result); err != nil {
				log.Printf("Failed to cache search results for chat %s: %v", chatID.String(), err)
			}
		}

		return result, nil
	}

	return ctrl.chatClient.SearchMessages(userID, chatID, query, offset, limit)
}

// UpdateChat - обновление чата с инвалидацией кеша
func (ctrl *ChatController) UpdateChat(chatID uuid.UUID, req *dto.UpdateChatRequestGateway, updateReq *ac.UpdateChatRequest, userID uuid.UUID) (*ac.UpdateChatResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Если передан новый аватар, загружаем его
	if req.Avatar != nil {
		uploadedFile, err := ctrl.fileClient.UploadFile(req.Avatar)
		if err != nil {
			return nil, err
		}
		updateReq.AvatarFileID = uploadedFile.ID
	}

	result, err := ctrl.chatClient.UpdateChat(chatID, updateReq, userID)
	if err != nil {
		return nil, err
	}

	// Инвалидация кеша информации о чате
	_ = ctrl.cacheService.DeleteChatInfoCache(ctx, chatID.String())

	// Инвалидация кеша участников чата
	_ = ctrl.cacheService.DeleteChatMembersCache(ctx, chatID.String())

	// Инвалидация кеша списков чатов для всех затронутых пользователей
	if updateReq.AddUserIDs != nil {
		for _, uid := range updateReq.AddUserIDs {
			_ = ctrl.cacheService.DeleteUserChatListCache(ctx, uid.String())
		}
	}
	if updateReq.RemoveUserIDs != nil {
		for _, uid := range updateReq.RemoveUserIDs {
			_ = ctrl.cacheService.DeleteUserChatListCache(ctx, uid.String())
		}
	}

	return result, nil
}

// DeleteChat - удаление чата с инвалидацией кеша
func (ctrl *ChatController) DeleteChat(chatID, userID uuid.UUID) error {
	err := ctrl.chatClient.DeleteChat(chatID, userID)
	if err != nil {
		return err
	}

	ctx := context.Background()

	// Инвалидация кеша информации о чате
	_ = ctrl.cacheService.DeleteChatInfoCache(ctx, chatID.String())

	// Инвалидация кеша участников чата
	_ = ctrl.cacheService.DeleteChatMembersCache(ctx, chatID.String())

	// Инвалидация кеша сообщений чата
	_ = ctrl.cacheService.DeleteChatMessagesCache(ctx, chatID.String())

	return nil
}

// BanUser - блокировка пользователя в чате с инвалидацией кеша
func (ctrl *ChatController) BanUser(chatID, userID, ownerID uuid.UUID) error {
	err := ctrl.chatClient.BanUser(chatID, userID, ownerID)
	if err != nil {
		return err
	}

	ctx := context.Background()

	// Инвалидация кеша участников чата
	_ = ctrl.cacheService.DeleteChatMembersCache(ctx, chatID.String())

	// Инвалидация кеша роли пользователя в чате
	_ = ctrl.cacheService.DeleteChatUserRoleCache(ctx, chatID.String(), userID.String())

	// Инвалидация списка чатов пользователя
	_ = ctrl.cacheService.DeleteUserChatListCache(ctx, userID.String())

	return nil
}

// ChangeUserRole - изменение роли пользователя в чате с инвалидацией кеша
func (ctrl *ChatController) ChangeUserRole(chatID, ownerID uuid.UUID, changeRoleReq *ac.ChangeRoleRequest) error {
	err := ctrl.chatClient.ChangeUserRole(chatID, ownerID, changeRoleReq)
	if err != nil {
		return err
	}

	ctx := context.Background()

	// Инвалидация кеша роли пользователя в чате
	_ = ctrl.cacheService.DeleteChatUserRoleCache(ctx, chatID.String(), changeRoleReq.UserID.String())

	// Инвалидация кеша участников чата (роли могут влиять на отображение)
	_ = ctrl.cacheService.DeleteChatMembersCache(ctx, chatID.String())

	return nil
}

// GetMyRoleInChat - получение своей роли в чате с permissions
func (ctrl *ChatController) GetMyRoleInChat(chatID, userID uuid.UUID) (*ac.MyRoleResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Пытаемся получить из кеша
	var cachedRole ac.MyRoleResponse
	err := ctrl.cacheService.GetChatUserRoleCache(ctx, chatID.String(), userID.String(), &cachedRole)
	if err == nil {
		return &cachedRole, nil
	}

	// Получаем из сервиса
	role, err := ctrl.chatClient.GetMyRoleInChat(chatID, userID)
	if err != nil {
		return nil, err
	}

	// Сохраняем в кеш
	if err := ctrl.cacheService.SetChatUserRoleCache(ctx, chatID.String(), userID.String(), role); err != nil {
		log.Printf("Failed to cache user role for chat %s user %s: %v", chatID.String(), userID.String(), err)
	}

	return role, nil
}

// GetChatMembers - получение списка участников чата
func (ctrl *ChatController) GetChatMembers(chatID uuid.UUID) ([]*ac.ChatMember, error) {
	return ctrl.chatClient.GetChatMembers(chatID)
}
