package controllers

import (
	"apiService/internal/custom_errors"
	"apiService/internal/dto"
	"apiService/internal/http_clients"
	au "common/contracts/api-user"
	"fmt"
	"github.com/google/uuid"
	"mime/multipart"
)

type AuthController struct {
	fileClient http_clients.FileClient
	userClient http_clients.UserClient
}

func NewAuthController(fileClient http_clients.FileClient, userClient http_clients.UserClient) *AuthController {
	return &AuthController{fileClient: fileClient, userClient: userClient}
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

func (ctrl *AuthController) Login(login *au.Login) (string, uuid.UUID, error) {
	return ctrl.userClient.Login(login)
}
