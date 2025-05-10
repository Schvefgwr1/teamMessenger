package controllers

import (
	httpClients "common/http_clients"
	"github.com/google/uuid"
	"userService/internal/custom_errors"
	"userService/internal/handlers/dto"
	"userService/internal/models"
	"userService/internal/repositories"
	"userService/internal/utils"
)

type AuthController struct {
	userRepo *repositories.UserRepository
	roleRepo *repositories.RoleRepository
}

func NewAuthController(userRepo *repositories.UserRepository, roleRepo *repositories.RoleRepository) *AuthController {
	return &AuthController{userRepo: userRepo, roleRepo: roleRepo}
}

func (c *AuthController) Register(req *dto.Register) (*models.User, error) {
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

func (c *AuthController) Login(req *dto.Login) (string, error) {
	user, err := c.userRepo.GetUserByUsername(req.Login)
	if err != nil {
		return "", custom_errors.ErrInvalidCredentials
	}

	if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		return "", custom_errors.ErrInvalidCredentials
	}

	var permNames []string
	for _, permission := range user.Role.Permissions {
		permNames = append(permNames, permission.Name)
	}

	token, err := utils.GenerateJWT(user.ID, permNames)
	if err != nil {
		return "", custom_errors.ErrTokenGeneration
	}

	return token, nil
}
