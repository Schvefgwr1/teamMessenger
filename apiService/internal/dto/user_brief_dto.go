package dto

import (
	fc "common/contracts/file-contracts"
	"encoding/json"
)

// FileSwagger - структура файла для Swagger документации
type FileSwagger struct {
	ID         int         `json:"id"`
	Name       string      `json:"name"`
	FileTypeID int         `json:"file_type_id"`
	URL        string      `json:"url"`
	CreatedAt  string      `json:"created_at"`
	FileType   interface{} `json:"file_type,omitempty"`
}

// UserBriefResponse - краткая информация о пользователе
type UserBriefResponse struct {
	Username    string  `json:"username"`
	Email       string  `json:"email"`
	Age         *int    `json:"age,omitempty"`
	Description *string `json:"description,omitempty"`

	// AvatarFile - реальные данные (скрыто от Swagger)
	AvatarFile *fc.File `json:"-" swaggerignore:"true"`

	// AvatarFileSwagger - только для Swagger документации
	AvatarFileSwagger *FileSwagger `json:"avatarFile,omitempty" swaggertype:"object"`

	ChatRoleName string `json:"chatRoleName,omitempty"`
}

// MarshalJSON кастомная сериализация для правильной работы с AvatarFile
func (u UserBriefResponse) MarshalJSON() ([]byte, error) {
	// Создаем временную структуру для сериализации, используя только AvatarFile
	aux := struct {
		Username     string   `json:"username"`
		Email        string   `json:"email"`
		Age          *int     `json:"age,omitempty"`
		Description  *string  `json:"description,omitempty"`
		AvatarFile   *fc.File `json:"avatarFile,omitempty"`
		ChatRoleName string   `json:"chatRoleName,omitempty"`
	}{
		Username:     u.Username,
		Email:        u.Email,
		Age:          u.Age,
		Description:  u.Description,
		AvatarFile:   u.AvatarFile,
		ChatRoleName: u.ChatRoleName,
	}
	return json.Marshal(aux)
}

// UnmarshalJSON кастомная десериализация для правильной работы с AvatarFile
func (u *UserBriefResponse) UnmarshalJSON(data []byte) error {
	// Используем временную структуру для десериализации
	aux := struct {
		Username     string   `json:"username"`
		Email        string   `json:"email"`
		Age          *int     `json:"age,omitempty"`
		Description  *string  `json:"description,omitempty"`
		AvatarFile   *fc.File `json:"avatarFile,omitempty"`
		ChatRoleName string   `json:"chatRoleName,omitempty"`
	}{}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	u.Username = aux.Username
	u.Email = aux.Email
	u.Age = aux.Age
	u.Description = aux.Description
	u.AvatarFile = aux.AvatarFile
	u.ChatRoleName = aux.ChatRoleName

	return nil
}
