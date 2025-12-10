package controllers

import (
	"fileService/internal/models"
	"fileService/internal/repositories"
)

type FileTypeController struct {
	repo repositories.FileTypeRepositoryInterface
}

func NewFileTypeController(repo repositories.FileTypeRepositoryInterface) *FileTypeController {
	return &FileTypeController{repo: repo}
}

// CreateFileType создает новый тип файла
func (c *FileTypeController) CreateFileType(name string) (*models.FileType, error) {
	fileType := &models.FileType{Name: name}
	err := c.repo.CreateFileType(fileType)
	return fileType, err
}

// GetFileTypeByID получает тип файла по ID
func (c *FileTypeController) GetFileTypeByID(id int) (*models.FileType, error) {
	return c.repo.GetFileTypeByID(id)
}

// GetFileTypeByName получает тип файла по названию
func (c *FileTypeController) GetFileTypeByName(name string) (*models.FileType, error) {
	return c.repo.GetFileTypeByName("application/" + name)
}

// DeleteFileTypeByID удаляет тип файла по ID
func (c *FileTypeController) DeleteFileTypeByID(id int) error {
	return c.repo.DeleteFileTypeByID(id)
}
