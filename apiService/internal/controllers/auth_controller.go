package controllers

import (
	dto "apiService/internal/dto"
	"apiService/internal/http_clients"
	au "common/contracts/api-user"
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
	var avatarUploadWarning string

	uploadedFile, uploadErr := ctrl.fileClient.UploadFile(file)
	if uploadErr != nil {
		avatarUploadWarning = uploadErr.Error()
	} else {
		registerData.AvatarID = uploadedFile.ID
	}

	registerData.Username = userRequest.Username
	registerData.Email = userRequest.Email
	registerData.Password = userRequest.Password
	registerData.Description = userRequest.Description
	registerData.Gender = userRequest.Gender
	registerData.Age = userRequest.Age
	registerData.RoleID = userRequest.RoleID

	var response dto.RegisterUserResponseGateway
	response.Warning = &avatarUploadWarning
	// Register user
	user, userErr := ctrl.userClient.RegisterUser(registerData)
	if userErr != nil {
		userErrStr := userErr.Error()
		response.Error = &userErrStr
	}
	response.User = user

	return &response
}

func (ctrl *AuthController) Login(login *au.Login) (string, error) {
	return ctrl.userClient.Login(login)
}
