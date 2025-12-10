package controllers

import (
	fc "common/contracts/file-contracts"
	"fileService/internal/dto"
	"fileService/internal/models"
	"mime/multipart"
)

// FileControllerInterface определяет интерфейс для FileController
type FileControllerInterface interface {
	UploadFile(fileHeader *multipart.FileHeader) (*models.File, error)
	GetFile(id int) (*fc.File, error)
	RenameFile(id int, newName string) (*models.File, error)
	GetFileNamesWithPagination(limit, offset int) (*[]dto.FileInformation, error)
}
