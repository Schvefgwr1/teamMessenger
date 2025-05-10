package services

import (
	"chatService/internal/repositories"
	"github.com/google/uuid"
)

type ChatPermissionService struct {
	chatUserRepo repositories.ChatUserRepository
}

func NewChatPermissionService(chatUserRepo repositories.ChatUserRepository) *ChatPermissionService {
	return &ChatPermissionService{chatUserRepo: chatUserRepo}
}

func (s *ChatPermissionService) HasPermission(userID, chatID uuid.UUID, permissionName string) (bool, error) {
	chatUser, err := s.chatUserRepo.GetChatUserWithRoleAndPermissions(userID, chatID)
	if err != nil {
		return false, err
	}

	for _, permission := range chatUser.Role.Permissions {
		if permission.Name == permissionName {
			return true, nil
		}
	}

	return false, nil
}
