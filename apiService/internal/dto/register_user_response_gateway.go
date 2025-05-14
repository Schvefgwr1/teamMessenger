package dto

import au "common/contracts/api-user"

type RegisterUserResponseGateway struct {
	User    *au.RegisterUserResponse `json:"user,omitempty"`
	Warning *string                  `json:"warning,omitempty"`
	Error   *string                  `json:"error,omitempty"`
}
