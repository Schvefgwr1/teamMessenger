package controllers

import (
	"apiService/internal/dto"
	"apiService/internal/http_clients"
	ac "common/contracts/api-chat"
	"errors"
	"github.com/google/uuid"
	"log"
)

type ChatController struct {
	chatClient http_clients.ChatClient
	fileClient http_clients.FileClient
}

func NewChatController(chatClient http_clients.ChatClient, fileClient http_clients.FileClient) *ChatController {
	return &ChatController{
		chatClient: chatClient,
		fileClient: fileClient,
	}
}

func (ctrl *ChatController) GetUserChats(userID uuid.UUID) ([]*ac.ChatResponse, error) {
	return ctrl.chatClient.GetUserChats(userID)
}

func (ctrl *ChatController) CreateChat(req *dto.CreateChatRequestGateway, ownerID uuid.UUID, userIDs []uuid.UUID) (*dto.CreateChatResponse, error) {
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

	return ctrl.chatClient.SendMessage(chatID, senderID, createReq)
}

func (ctrl *ChatController) GetChatMessages(chatID uuid.UUID, userID uuid.UUID, offset, limit int) ([]*ac.GetChatMessage, error) {
	return ctrl.chatClient.GetChatMessages(chatID, userID, offset, limit)
}

func (ctrl *ChatController) SearchMessages(userID uuid.UUID, chatID uuid.UUID, query string, offset, limit int) (*ac.GetSearchResponse, error) {
	return ctrl.chatClient.SearchMessages(userID, chatID, query, offset, limit)
}
