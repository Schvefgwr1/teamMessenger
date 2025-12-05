package dto

import (
	fc "common/contracts/file-contracts"
	"encoding/json"
	"taskService/internal/models"
)

type TaskResponse struct {
	Task *models.Task `json:"task"`

	// Files - реальные данные (скрыто от Swagger)
	Files *[]fc.File `json:"-" swaggerignore:"true"`

	// FilesSwagger - только для Swagger документации
	FilesSwagger *[]FileSwagger `json:"files,omitempty" swaggertype:"array,object"`
}

// MarshalJSON кастомная сериализация для правильной работы с Files
func (t TaskResponse) MarshalJSON() ([]byte, error) {
	// Создаем временную структуру для сериализации, используя только Files
	aux := struct {
		Task  *models.Task `json:"task"`
		Files *[]fc.File   `json:"files,omitempty"`
	}{
		Task:  t.Task,
		Files: t.Files,
	}
	return json.Marshal(aux)
}
