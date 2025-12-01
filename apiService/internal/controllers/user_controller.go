package controllers

import (
	"apiService/internal/dto"
	"apiService/internal/http_clients"
	"apiService/internal/services"
	au "common/contracts/api-user"
	uc "common/contracts/user-contracts"
	"context"
	"errors"
	"log"
	"mime/multipart"
	"time"

	"github.com/google/uuid"
)

type UserController struct {
	fileClient   http_clients.FileClient
	userClient   http_clients.UserClient
	cacheService *services.CacheService
}

func NewUserController(fileClient http_clients.FileClient, userClient http_clients.UserClient, cacheService *services.CacheService) *UserController {
	return &UserController{
		fileClient:   fileClient,
		userClient:   userClient,
		cacheService: cacheService,
	}
}

func (ctrl *UserController) GetUser(id uuid.UUID) (*au.GetUserResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var cachedUser au.GetUserResponse
	err := ctrl.cacheService.GetUserCache(ctx, id.String(), &cachedUser)
	if err == nil {
		log.Printf("User %s found in cache", id.String())
		return &cachedUser, nil
	}

	user, err := ctrl.userClient.GetUserByID(id.String())
	if err != nil {
		return nil, err
	}

	// Сохраняем в кеш
	if err := ctrl.cacheService.SetUserCache(ctx, id.String(), user); err != nil {
		log.Printf("Failed to cache user %s: %v", id.String(), err)
	}

	return user, nil
}

func (ctrl *UserController) UpdateUser(id uuid.UUID, userRequest *dto.UpdateUserRequestGateway, file *multipart.FileHeader) (*au.UpdateUserResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	var updateData au.UpdateUserRequest

	// Map fields from Gateway DTO
	updateData.Username = userRequest.Username
	updateData.Description = userRequest.Description
	updateData.Gender = userRequest.Gender
	updateData.Age = userRequest.Age
	updateData.RoleID = userRequest.RoleID

	updateResponse, err := ctrl.userClient.UpdateUser(id.String(), &updateData)
	if err != nil {
		return nil, err
	}
	if updateResponse.Error != nil {
		return nil, errors.New(*updateResponse.Error)
	}

	if file != nil {
		uploadedFile, uploadErr := ctrl.fileClient.UploadFile(file)
		if uploadErr != nil {
			return nil, uploadErr
		} else {
			var updateAvatar = &au.UpdateUserRequest{AvatarFileID: uploadedFile.ID}
			updateAvResp, err := ctrl.userClient.UpdateUser(id.String(), updateAvatar)
			if err != nil {
				return nil, err
			}
			if updateAvResp.Error != nil {
				return nil, errors.New(*updateAvResp.Error)
			}
			updateResponse = updateAvResp
		}
	}

	// Инвалидируем кеш пользователя после обновления
	if err := ctrl.cacheService.DeleteUserCache(ctx, id.String()); err != nil {
		log.Printf("Failed to invalidate user cache for %s: %v", id.String(), err)
	}

	// Инвалидируем только список чатов этого пользователя (не всех!)
	if err := ctrl.cacheService.DeleteUserChatListCache(ctx, id.String()); err != nil {
		log.Printf("Failed to invalidate chat list cache for %s: %v", id.String(), err)
	}

	return updateResponse, nil
}

// GetAllPermissions - получить все разрешения с кешированием
func (ctrl *UserController) GetAllPermissions() ([]*uc.Permission, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Пытаемся получить из кеша
	cacheKey := "permissions:all"
	var cachedPermissions []*uc.Permission
	err := ctrl.cacheService.Get(ctx, cacheKey, &cachedPermissions)
	if err == nil {
		log.Printf("Permissions found in cache")
		return cachedPermissions, nil
	}

	// Получаем из сервиса
	permissions, err := ctrl.userClient.GetAllPermissions()
	if err != nil {
		return nil, err
	}

	// Сохраняем в кеш на 1 час (разрешения меняются редко)
	if err := ctrl.cacheService.Set(ctx, cacheKey, permissions, time.Hour); err != nil {
		log.Printf("Failed to cache permissions: %v", err)
	}

	return permissions, nil
}

// GetAllRoles - получить все роли с кешированием
func (ctrl *UserController) GetAllRoles() ([]*uc.Role, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Пытаемся получить из кеша
	cacheKey := "roles:all"
	var cachedRoles []*uc.Role
	err := ctrl.cacheService.Get(ctx, cacheKey, &cachedRoles)
	if err == nil {
		log.Printf("Roles found in cache")
		return cachedRoles, nil
	}

	// Получаем из сервиса
	roles, err := ctrl.userClient.GetAllRoles()
	if err != nil {
		return nil, err
	}

	// Сохраняем в кеш на 1 час (роли меняются редко)
	if err := ctrl.cacheService.Set(ctx, cacheKey, roles, time.Hour); err != nil {
		log.Printf("Failed to cache roles: %v", err)
	}

	return roles, nil
}

// CreateRole - создать новую роль с инвалидацией кеша
func (ctrl *UserController) CreateRole(req *au.CreateRoleRequest) (*uc.Role, error) {
	role, err := ctrl.userClient.CreateRole(req)
	if err != nil {
		return nil, err
	}

	// Инвалидация кеша списка ролей
	ctx := context.Background()
	cacheKey := "roles:all"
	_ = ctrl.cacheService.Delete(ctx, cacheKey)

	return role, nil
}

// GetUserProfileByID - получить профиль пользователя по ID с кешированием
func (ctrl *UserController) GetUserProfileByID(userID uuid.UUID) (*au.GetUserResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Пытаемся получить из кеша
	var cachedUser au.GetUserResponse
	err := ctrl.cacheService.GetUserCache(ctx, userID.String(), &cachedUser)
	if err == nil {
		log.Printf("User profile %s found in cache", userID.String())
		return &cachedUser, nil
	}

	// Получаем из сервиса
	userProfile, err := ctrl.userClient.GetUserByID(userID.String())
	if err != nil {
		return nil, err
	}

	// Сохраняем в кеш
	if err := ctrl.cacheService.SetUserCache(ctx, userID.String(), userProfile); err != nil {
		log.Printf("Failed to cache user profile %s: %v", userID.String(), err)
	}

	return userProfile, nil
}
