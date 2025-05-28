package controllers

import (
	"apiService/internal/dto"
	"apiService/internal/http_clients"
	au "common/contracts/api-user"
	"errors"
	"github.com/google/uuid"
	"mime/multipart"
)

type UserController struct {
	fileClient http_clients.FileClient
	userClient http_clients.UserClient
}

func NewUserController(fileClient http_clients.FileClient, userClient http_clients.UserClient) *UserController {
	return &UserController{fileClient: fileClient, userClient: userClient}
}

func (ctrl *UserController) GetUser(id uuid.UUID) (*au.GetUserResponse, error) {
	return ctrl.userClient.GetUserByID(id.String())
}

func (ctrl *UserController) UpdateUser(id uuid.UUID, userRequest *dto.UpdateUserRequestGateway, file *multipart.FileHeader) (*au.UpdateUserResponse, error) {
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
			return updateAvResp, nil
		}
	}
	return updateResponse, nil
}
