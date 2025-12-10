package repositories

import (
	"fileService/internal/dto"
	"fileService/internal/models"
)

// FileRepositoryInterface определяет интерфейс для работы с файлами в БД
type FileRepositoryInterface interface {
	CreateFile(file *models.File) error
	GetFileByID(id int) (*models.File, error)
	GetFileByName(name string) (*models.File, error)
	UpdateFile(file *models.File) error
	GetFileNamesWithPagination(limit, offset int) (*[]dto.FileInformation, error)
}

// FileTypeRepositoryInterface определяет интерфейс для работы с типами файлов в БД
type FileTypeRepositoryInterface interface {
	CreateFileType(fileType *models.FileType) error
	GetFileTypeByID(id int) (*models.FileType, error)
	GetFileTypeByName(name string) (*models.FileType, error)
	DeleteFileTypeByID(id int) error
}
