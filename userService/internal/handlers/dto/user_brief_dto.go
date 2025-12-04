package dto

import fc "common/contracts/file-contracts"

// UserBriefResponse - краткая информация о пользователе
type UserBriefResponse struct {
	Username     string   `json:"username"`
	Email        string   `json:"email"`
	Age          *int     `json:"age,omitempty"`
	Description  *string  `json:"description,omitempty"`
	AvatarFile   *fc.File `json:"avatarFile,omitempty"`
	ChatRoleName string   `json:"chatRoleName,omitempty"`
}
