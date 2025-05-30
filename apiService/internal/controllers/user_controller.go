package controllers

import (
	"apiService/internal/dto"
	"apiService/internal/http_clients"
	"apiService/internal/services"
	au "common/contracts/api-user"
	"context"
	"errors"
	"github.com/google/uuid"
	"log"
	"mime/multipart"
	"time"
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
	// Также инвалидируем списки чатов, где участвует пользователь
	if err := ctrl.cacheService.DeleteByPattern(ctx, "chat_list:*"); err != nil {
		log.Printf("Failed to invalidate chat list cache: %v", err)
	}

	return updateResponse, nil
}
