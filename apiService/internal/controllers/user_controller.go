package controllers

import (
	"apiService/internal/http_clients"
	au "common/contracts/api-user"
	"github.com/google/uuid"
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
