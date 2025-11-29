package controllers

import (
	"apiService/internal/custom_errors"
	"apiService/internal/dto"
	"apiService/internal/http_clients"
	"apiService/internal/services"
	au "common/contracts/api-user"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"mime/multipart"
)

type AuthController struct {
	fileClient     http_clients.FileClient
	userClient     http_clients.UserClient
	sessionService *services.SessionService
}

func NewAuthController(fileClient http_clients.FileClient, userClient http_clients.UserClient, sessionService *services.SessionService) *AuthController {
	return &AuthController{
		fileClient:     fileClient,
		userClient:     userClient,
		sessionService: sessionService,
	}
}

func (ctrl *AuthController) Register(userRequest *dto.RegisterUserRequestGateway, file *multipart.FileHeader) *dto.RegisterUserResponseGateway {
	var registerData au.RegisterUserRequest

	registerData.Username = userRequest.Username
	registerData.Email = userRequest.Email
	registerData.Password = userRequest.Password
	registerData.Description = userRequest.Description
	registerData.Gender = userRequest.Gender
	registerData.Age = userRequest.Age
	registerData.RoleID = userRequest.RoleID

	var response dto.RegisterUserResponseGateway
	user, userErr := ctrl.userClient.RegisterUser(registerData)
	if userErr != nil {
		userErrStr := userErr.Error()
		response.Error = &userErrStr
		return &response
	}
	if user == nil {
		response.Error = &custom_errors.ErrNilUserInClient
		return &response
	}
	response.User = user

	var avatarUploadWarning string
	if file != nil {
		uploadedFile, uploadErr := ctrl.fileClient.UploadFile(file)
		if uploadErr != nil {
			avatarUploadWarning = uploadErr.Error()
		} else {
			updateUserResp, err := ctrl.userClient.UpdateUser(
				user.ID.String(),
				&au.UpdateUserRequest{AvatarFileID: uploadedFile.ID},
			)
			if updateUserResp == nil {
				response.Error = &custom_errors.ErrNilUserInClient
				return &response
			}
			if err != nil {
				avatarUploadWarning = fmt.Sprintf("Error of add avatar in user service: %s %s", err.Error(), updateUserResp.Error)
			}
			response.User.AvatarFileID = uploadedFile.ID
		}
	}

	response.Warning = &avatarUploadWarning
	return &response
}

func (ctrl *AuthController) Login(ctx context.Context, login *au.Login) (string, uuid.UUID, error) {
	token, userID, err := ctrl.userClient.Login(login)
	if err != nil {
		return "", uuid.Nil, err
	}

	// Создаем сессию в Redis если есть sessionService
	if ctrl.sessionService != nil && token != "" && userID != uuid.Nil {
		// Отзываем все старые сессии пользователя перед созданием новой
		if err := ctrl.sessionService.RevokeAllUserSessions(ctx, userID); err != nil {
			// Логируем ошибку, но продолжаем создание новой сессии
			fmt.Printf("Failed to revoke old sessions for user %s: %v\n", userID.String(), err)
		}

		expiresAt := time.Now().Add(24 * time.Hour)
		if err := ctrl.sessionService.CreateSession(ctx, userID, token, expiresAt); err != nil {
			// Логируем ошибку, но не прерываем процесс логина
			fmt.Printf("Failed to create session in Redis: %v\n", err)
		}
	}

	return token, userID, nil
}

func (ctrl *AuthController) Logout(ctx context.Context, userID uuid.UUID, token string) error {
	// Если sessionService отсутствует, считаем что выход успешен
	if ctrl.sessionService == nil {
		return nil
	}

	return ctrl.sessionService.RevokeSession(ctx, userID, token)
}
