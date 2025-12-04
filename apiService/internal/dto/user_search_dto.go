package dto

import (
	fc "common/contracts/file-contracts"
	"encoding/json"
)

// UserSearchResult - результат поиска пользователя
type UserSearchResult struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`

	// AvatarFile - реальные данные (скрыто от Swagger)
	AvatarFile *fc.File `json:"-" swaggerignore:"true"`

	// AvatarFileSwagger - только для Swagger документации
	AvatarFileSwagger *FileSwagger `json:"avatarFile,omitempty" swaggertype:"object"`
}

// MarshalJSON кастомная сериализация
func (u UserSearchResult) MarshalJSON() ([]byte, error) {
	aux := struct {
		ID         string   `json:"id"`
		Username   string   `json:"username"`
		Email      string   `json:"email"`
		AvatarFile *fc.File `json:"avatarFile,omitempty"`
	}{
		ID:         u.ID,
		Username:   u.Username,
		Email:      u.Email,
		AvatarFile: u.AvatarFile,
	}
	return json.Marshal(aux)
}

// UnmarshalJSON кастомная десериализация
func (u *UserSearchResult) UnmarshalJSON(data []byte) error {
	aux := struct {
		ID         string   `json:"id"`
		Username   string   `json:"username"`
		Email      string   `json:"email"`
		AvatarFile *fc.File `json:"avatarFile,omitempty"`
	}{}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	u.ID = aux.ID
	u.Username = aux.Username
	u.Email = aux.Email
	u.AvatarFile = aux.AvatarFile

	return nil
}

// UserSearchResponse - ответ на поиск пользователей
type UserSearchResponse struct {
	Users []UserSearchResult `json:"users"`
}
