package controllers

import (
	au "common/contracts/api-user"
	fc "common/contracts/file-contracts"
	httpClients "common/http_clients"
	"github.com/google/uuid"
	"log"
	"userService/internal/custom_errors"
	"userService/internal/handlers/dto"
	"userService/internal/models"
	"userService/internal/repositories"
)

type UserController struct {
	userRepo *repositories.UserRepository
	roleRepo *repositories.RoleRepository
}

func NewUserController(userRepo *repositories.UserRepository, roleRepo *repositories.RoleRepository) *UserController {
	return &UserController{userRepo: userRepo, roleRepo: roleRepo}
}

func (c *UserController) GetUserProfile(id uuid.UUID) (*models.User, *fc.File, error) {
	user, err := c.userRepo.GetUserByID(id)
	if err != nil {
		return nil, nil, err
	}

	if user.AvatarFileID != nil {
		file, err := httpClients.GetFileByID(*user.AvatarFileID)
		if err != nil {
			return user, nil, err
		} else {
			return user, file, nil
		}
	}
	return user, nil, nil
}

func (c *UserController) UpdateUserProfile(req *au.UpdateUserRequest, userId *uuid.UUID) error {
	var user *models.User
	user, err := c.userRepo.GetUserByID(*userId)
	if err != nil {
		return custom_errors.ErrInvalidCredentials
	}

	if req.Username != nil {
		existingUser, err := c.userRepo.GetUserByUsername(*req.Username)
		if err == nil && existingUser != nil {
			return custom_errors.NewUserUsernameConflictError(*req.Username)
		} else {
			user.Username = *req.Username
		}
	}

	if req.RoleID != nil {
		existingRole, err := c.roleRepo.GetRoleByID(*req.RoleID)
		if err != nil || existingRole == nil {
			return custom_errors.NewRoleNotFoundError(*req.RoleID)
		} else {
			user.RoleID = *req.RoleID
		}
	}
	if req.AvatarFileID != nil {
		file, errHTTP := httpClients.GetFileByID(*req.AvatarFileID)
		if errHTTP != nil {
			return custom_errors.NewGetFileHTTPError(*req.AvatarFileID, errHTTP.Error())
		}
		if file.ID <= 0 {
			return custom_errors.NewFileNotFoundError(*req.AvatarFileID)
		}

		user.AvatarFileID = &file.ID
	}

	if req.Description != nil {
		user.Description = req.Description
	}
	if req.Gender != nil {
		user.Gender = req.Gender
	}
	if req.Age != nil {
		user.Age = req.Age
	}

	return c.userRepo.UpdateUser(user)
}

// GetUserBrief возвращает краткую информацию о пользователе с ролью в чате
func (c *UserController) GetUserBrief(userID uuid.UUID, chatID string, requesterID string) (*dto.UserBriefResponse, error) {
	user, err := c.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, custom_errors.ErrInvalidCredentials
	}

	response := &dto.UserBriefResponse{
		Username:    user.Username,
		Email:       user.Email,
		Age:         user.Age,
		Description: user.Description,
	}

	// Загружаем аватар если есть
	if user.AvatarFileID != nil {
		log.Printf("[GetUserBrief] Loading avatar for user %s, avatarFileID: %d", userID.String(), *user.AvatarFileID)
		file, err := httpClients.GetFileByID(*user.AvatarFileID)
		if err != nil {
			log.Printf("[GetUserBrief] Error loading avatar: %v", err)
		} else if file != nil {
			log.Printf("[GetUserBrief] Avatar loaded successfully: ID=%d, URL=%s", file.ID, file.URL)
			response.AvatarFile = file
		}
	} else {
		log.Printf("[GetUserBrief] User %s has no avatar", userID.String())
	}

	// Получаем роль в чате если передан chatID
	if chatID != "" && requesterID != "" {
		log.Printf("[GetUserBrief] Getting chat role for user %s in chat %s, requester: %s", userID.String(), chatID, requesterID)
		roleResp, err := httpClients.GetUserRoleInChat(chatID, userID.String(), requesterID)
		if err != nil {
			log.Printf("[GetUserBrief] Error getting chat role: %v", err)
		} else if roleResp != nil {
			log.Printf("[GetUserBrief] Chat role received: %s", roleResp.RoleName)
			response.ChatRoleName = roleResp.RoleName
		}
	} else {
		log.Printf("[GetUserBrief] chatID or requesterID is empty, skipping chat role. chatID=%s, requesterID=%s", chatID, requesterID)
	}

	return response, nil
}

// SearchUsers ищет пользователей по имени или email
func (c *UserController) SearchUsers(query string, limit int) (*dto.UserSearchResponse, error) {
	if limit <= 0 || limit > 20 {
		limit = 10
	}

	users, err := c.userRepo.SearchUsers(query, limit)
	if err != nil {
		return nil, err
	}

	results := make([]dto.UserSearchResult, 0, len(users))
	for _, user := range users {
		result := dto.UserSearchResult{
			ID:       user.ID.String(),
			Username: user.Username,
			Email:    user.Email,
		}

		// Загружаем аватар если есть
		if user.AvatarFileID != nil {
			file, err := httpClients.GetFileByID(*user.AvatarFileID)
			if err == nil && file != nil {
				result.AvatarFile = file
			}
		}

		results = append(results, result)
	}

	return &dto.UserSearchResponse{Users: results}, nil
}
