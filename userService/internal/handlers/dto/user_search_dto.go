package dto

import fc "common/contracts/file-contracts"

// UserSearchResult - результат поиска пользователя
type UserSearchResult struct {
	ID         string   `json:"id"`
	Username   string   `json:"username"`
	Email      string   `json:"email"`
	AvatarFile *fc.File `json:"avatarFile,omitempty"`
}

// UserSearchResponse - ответ на поиск пользователей
type UserSearchResponse struct {
	Users []UserSearchResult `json:"users"`
}
