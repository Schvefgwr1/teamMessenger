package controllers

import (
	"chatService/internal/custom_errors"
	"chatService/internal/handlers/dto"
	"chatService/internal/models"
	"chatService/internal/repositories"
	httpClients "common/http_clients"
	"github.com/google/uuid"
)

type ChatController struct {
	ChatRepo     repositories.ChatRepository
	ChatUserRepo repositories.ChatUserRepository
	ChatRoleRepo repositories.ChatRoleRepository
}

func NewChatController(
	chatRepo repositories.ChatRepository,
	chatUserRepo repositories.ChatUserRepository,
	chatRoleRepo repositories.ChatRoleRepository,
) *ChatController {
	return &ChatController{chatRepo, chatUserRepo, chatRoleRepo}
}

func (c *ChatController) GetUserChats(userID uuid.UUID) (*[]dto.ChatResponse, error) {
	chats, err := c.ChatRepo.GetUserChats(userID)
	if err != nil {
		return nil, err
	}

	var result []dto.ChatResponse
	for _, chat := range chats {
		result = append(result, dto.ChatResponse{
			ID:           chat.ID,
			Name:         chat.Name,
			IsGroup:      chat.IsGroup,
			Description:  chat.Description,
			AvatarFileID: chat.AvatarFileID,
			CreatedAt:    chat.CreatedAt,
		})
	}
	return &result, nil
}

func (c *ChatController) CreateChat(dto *dto.CreateChatDTO) (*uuid.UUID, error) {
	ownerRole, err := c.ChatRoleRepo.GetRoleByName("owner")
	if err != nil {
		return nil, custom_errors.ErrInvalidCredentials
	}
	mainRole, err := c.ChatRoleRepo.GetRoleByName("main")
	if err != nil {
		return nil, custom_errors.ErrInvalidCredentials
	}

	newChat := &models.Chat{
		ID:          uuid.New(),
		Name:        dto.Name,
		Description: dto.Description,
	}

	if dto.AvatarFileID != nil {
		file, errHTTP := httpClients.GetFileByID(*dto.AvatarFileID)
		if errHTTP != nil {
			return nil, custom_errors.NewGetFileHTTPError(*dto.AvatarFileID, errHTTP.Error())
		}
		if file.ID <= 0 {
			return nil, custom_errors.NewFileNotFoundError(*dto.AvatarFileID)
		}

		newChat.AvatarFileID = &file.ID
	}
	if dto.Description != nil {
		newChat.Description = dto.Description
	}
	if len(dto.UserIDs) <= 1 {
		newChat.IsGroup = false
	} else {
		newChat.IsGroup = true
	}

	userOwnerClientResponse, err := httpClients.GetUserByID(&dto.OwnerID)
	if err != nil {
		return nil, custom_errors.NewUserClientError(err.Error())
	}
	if userOwnerClientResponse.User == nil {
		return nil, custom_errors.NewUserClientError("nil user")
	}

	err = c.ChatRepo.CreateChat(newChat)
	if err != nil {
		return nil, custom_errors.NewDatabaseError(err.Error())
	}

	chatUser := &models.ChatUser{
		ChatID: newChat.ID,
		UserID: userOwnerClientResponse.User.ID,
		RoleID: ownerRole.ID,
	}
	err = c.ChatUserRepo.AddUserToChat(chatUser)
	if err != nil {
		return nil, custom_errors.NewDatabaseError(err.Error())
	}
	for _, userID := range dto.UserIDs {
		userClientResponse, err := httpClients.GetUserByID(&userID)
		if err != nil {
			return nil, custom_errors.NewUserClientError(err.Error())
		}
		if userClientResponse.User == nil {
			return nil, custom_errors.NewUserClientError("nil user")
		}
		newChatUser := &models.ChatUser{
			ChatID: newChat.ID,
			UserID: userClientResponse.User.ID,
			RoleID: mainRole.ID,
		}
		err = c.ChatUserRepo.AddUserToChat(newChatUser)
		if err != nil {
			return nil, custom_errors.NewDatabaseError(err.Error())
		}
	}
	return &newChat.ID, nil
}

func (c *ChatController) UpdateChat(chatID uuid.UUID, updateChatDTO *dto.UpdateChatDTO) (*dto.UpdateChatResponse, error) {
	chat, err := c.ChatRepo.GetChatByID(chatID)
	if err != nil {
		return nil, custom_errors.NewDatabaseError("chat not found: " + err.Error())
	}

	if updateChatDTO.Name != nil {
		chat.Name = *updateChatDTO.Name
	}
	if updateChatDTO.Description != nil {
		chat.Description = updateChatDTO.Description
	}
	if updateChatDTO.AvatarFileID != nil {
		file, errHTTP := httpClients.GetFileByID(*updateChatDTO.AvatarFileID)
		if errHTTP != nil {
			return nil, custom_errors.NewGetFileHTTPError(*updateChatDTO.AvatarFileID, errHTTP.Error())
		}
		if file.ID <= 0 {
			return nil, custom_errors.NewFileNotFoundError(*updateChatDTO.AvatarFileID)
		}
		chat.AvatarFileID = &file.ID
	}

	if err := c.ChatRepo.UpdateChat(chat); err != nil {
		return nil, custom_errors.NewDatabaseError(err.Error())
	}

	mainRole, err := c.ChatRoleRepo.GetRoleByName("main")
	if err != nil {
		return nil, custom_errors.ErrInvalidCredentials
	}

	var updateUsers []dto.UpdateUser
	for _, userID := range updateChatDTO.AddUserIDs {
		userResp, err := httpClients.GetUserByID(&userID)
		if err != nil {
			return nil, custom_errors.NewUserClientError(err.Error())
		}
		if userResp.User == nil {
			return nil, custom_errors.NewUserClientError("nil user")
		}
		newChatUser := &models.ChatUser{
			ChatID: chatID,
			UserID: userResp.User.ID,
			RoleID: mainRole.ID,
		}
		if err := c.ChatUserRepo.AddUserToChat(newChatUser); err != nil {
			return nil, custom_errors.NewDatabaseError(err.Error())
		}
		updateUsers = append(updateUsers, dto.UpdateUser{UserID: newChatUser.UserID, State: "created"})
	}

	for _, userID := range updateChatDTO.RemoveUserIDs {
		if err := c.ChatUserRepo.RemoveUserFromChat(chatID, userID); err != nil {
			return nil, custom_errors.NewDatabaseError(err.Error())
		}
		updateUsers = append(updateUsers, dto.UpdateUser{UserID: userID, State: "deleted"})
	}

	return &dto.UpdateChatResponse{Chat: dto.ChatResponse{
		ID:           chat.ID,
		Name:         chat.Name,
		IsGroup:      chat.IsGroup,
		Description:  chat.Description,
		AvatarFileID: chat.AvatarFileID,
		CreatedAt:    chat.CreatedAt,
	}, UpdateUsers: updateUsers}, nil
}

func (c *ChatController) DeleteChat(chatID uuid.UUID) error {
	if err := c.ChatUserRepo.DeleteChatUsersByChatID(chatID); err != nil {
		return custom_errors.NewDatabaseError(err.Error())
	}

	if err := c.ChatRepo.DeleteChat(chatID); err != nil {
		return custom_errors.NewDatabaseError(err.Error())
	}

	return nil
}

func (c *ChatController) ChangeUserRole(chatID, userID uuid.UUID, roleID int) error {
	chatUser, err := c.ChatUserRepo.GetChatUser(userID, chatID)
	if chatUser == nil || err != nil {
		return custom_errors.ErrInvalidCredentials
	}
	if _, err := c.ChatRoleRepo.GetRoleByID(roleID); err != nil {
		return custom_errors.ErrInvalidCredentials
	}
	return c.ChatUserRepo.ChangeUserRole(chatID, userID, roleID)
}

func (c *ChatController) BanUser(chatID, userID uuid.UUID) error {
	chatUser, err := c.ChatUserRepo.GetChatUser(userID, chatID)
	if chatUser == nil || err != nil {
		return custom_errors.ErrInvalidCredentials
	}
	bannedRole, err := c.ChatRoleRepo.GetRoleByName("banned")
	if err != nil {
		return custom_errors.ErrInternalServerError
	}
	return c.ChatUserRepo.ChangeUserRole(chatID, userID, bannedRole.ID)

}
