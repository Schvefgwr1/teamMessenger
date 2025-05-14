package api_user

import (
	fc "common/contracts/file-contracts"
	uc "common/contracts/user-contracts"
)

type GetUserResponse struct {
	File *fc.File `json:"file"`
	User *uc.User `json:"user"`
}
