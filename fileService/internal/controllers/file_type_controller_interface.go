package controllers

import (
	"fileService/internal/models"
)

// FileTypeControllerInterface определяет интерфейс для FileTypeController
type FileTypeControllerInterface interface {
	CreateFileType(name string) (*models.FileType, error)
	GetFileTypeByID(id int) (*models.FileType, error)
	GetFileTypeByName(name string) (*models.FileType, error)
	DeleteFileTypeByID(id int) error
}
