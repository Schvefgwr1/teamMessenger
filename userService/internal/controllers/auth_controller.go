package controllers

import (
	au "common/contracts/api-user"
	httpClients "common/http_clients"
	"github.com/google/uuid"
	"log"
	"userService/internal/custom_errors"
	"userService/internal/models"
	"userService/internal/repositories"
	"userService/internal/services"
	"userService/internal/utils"
)

type AuthController struct {
	userRepo            *repositories.UserRepository
	roleRepo            *repositories.RoleRepository
	notificationService *services.NotificationService
}

func NewAuthController(
	userRepo *repositories.UserRepository,
	roleRepo *repositories.RoleRepository,
	notificationService *services.NotificationService,
) *AuthController {
	return &AuthController{
		userRepo:            userRepo,
		roleRepo:            roleRepo,
		notificationService: notificationService,
	}
}

func (c *AuthController) Register(req *au.RegisterUserRequest) (*models.User, error) {
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	existingUser, err := c.userRepo.GetUserByEmail(req.Email)
	if err == nil && existingUser != nil {
		return nil, custom_errors.NewUserEmailConflictError(req.Email)
	}

	existingUser, err = c.userRepo.GetUserByUsername(req.Username)
	if err == nil && existingUser != nil {
		return nil, custom_errors.NewUserUsernameConflictError(req.Username)
	}

	existingRole, err := c.roleRepo.GetRoleByID(req.RoleID)
	if err != nil || existingRole == nil {
		return nil, custom_errors.NewRoleNotFoundError(req.RoleID)
	}

	user := &models.User{
		ID:           uuid.New(),
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		Gender:       &req.Gender,
		Age:          &req.Age,
		RoleID:       *existingRole.ID,
		Role:         *existingRole,
	}

	if req.Description != nil {
		user.Description = req.Description
	}

	if req.AvatarID != nil {
		file, errHTTP := httpClients.GetFileByID(*req.AvatarID)
		if errHTTP != nil {
			return nil, custom_errors.NewGetFileHTTPError(*req.AvatarID, errHTTP.Error())
		}
		if file.ID <= 0 {
			return nil, custom_errors.NewFileNotFoundError(*req.AvatarID)
		}

		user.AvatarFileID = &file.ID
	}
	if err := c.userRepo.CreateUser(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (c *AuthController) Login(req *au.Login, ipAddress, userAgent string) (string, uuid.UUID, error) {
	user, err := c.userRepo.GetUserByUsername(req.Login)
	if err != nil {
		return "", uuid.Nil, custom_errors.ErrInvalidCredentials
	}

	if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		return "", uuid.Nil, custom_errors.ErrInvalidCredentials
	}

	var permNames []string
	for _, permission := range user.Role.Permissions {
		permNames = append(permNames, permission.Name)
	}

	token, err := utils.GenerateJWT(user.ID, permNames)
	if err != nil {
		return "", uuid.Nil, custom_errors.ErrTokenGeneration
	}

	// Отправляем уведомление о входе в систему
	if c.notificationService != nil && user.Email != "" {
		if err := c.notificationService.SendLoginNotification(
			user.ID,
			user.Username,
			user.Email,
			ipAddress,
			userAgent,
		); err != nil {
			// Логируем ошибку, но не прерываем процесс входа
			log.Printf("Failed to send login notification for user %s: %v", user.Username, err)
		}
	}

	return token, user.ID, nil
}
