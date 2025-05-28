package controllers

import (
	au "common/contracts/api-user"
	fc "common/contracts/file-contracts"
	httpClients "common/http_clients"
	"github.com/google/uuid"
	"userService/internal/custom_errors"
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
