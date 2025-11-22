package dto

import (
	fc "common/contracts/file-contracts"
	"taskService/internal/models"
)

type TaskResponse struct {
	Task *models.Task `json:"task"`

	// Files Swagger override
	Files *[]fc.File `json:"files" swaggertype:"array,object" swaggerignore:"true"`

	// Эта часть будет видна Swagger
	FilesSwagger *[]FileSwagger `json:"files"`
}
