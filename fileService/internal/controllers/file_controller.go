package controllers

import (
	"common/config"
	fc "common/contracts/file-contracts"
	"context"
	"errors"
	"fileService/internal/dto"
	"fileService/internal/models"
	"fileService/internal/repositories"
	"fmt"
	"github.com/minio/minio-go/v7"
	"golang.org/x/xerrors"
	"log"
	"mime/multipart"
	"strconv"
	"strings"
	"time"
)

type FileController struct {
	repo        *repositories.FileRepository
	repoType    *repositories.FileTypeRepository
	minioClient *minio.Client
	minioConfig *config.MinIO
}

func NewFileController(repo *repositories.FileRepository, repoType *repositories.FileTypeRepository, minioClient *minio.Client, minioConfig *config.MinIO) *FileController {
	return &FileController{
		repo:        repo,
		repoType:    repoType,
		minioClient: minioClient,
		minioConfig: minioConfig,
	}
}

// UploadFile загружает файл в MinIO и сохраняет метаданные в БД
func (c *FileController) UploadFile(fileHeader *multipart.FileHeader) (*models.File, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bucketName := c.minioConfig.Bucket
	objectName := strconv.Itoa(time.Now().Nanosecond()) + "_" + fileHeader.Filename

	// Проверяем, поддерживается ли тип файла
	fileType, err := c.repoType.GetFileTypeByName(fileHeader.Header.Get("Content-Type"))
	if err != nil || fileType == nil {
		return nil, errors.New(fmt.Sprintf(
			"Unsupported file type: %s",
			fileHeader.Header.Get("Content-Type")),
		)
	}

	// Проверяем, существует ли файл в БД
	existingFile, err := c.repo.GetFileByName(objectName)
	if err == nil && existingFile != nil {
		return nil, fmt.Errorf("file %s already exists in database", objectName)
	}

	// Проверяем, существует ли файл в MinIO
	_, err = c.minioClient.StatObject(context.Background(), bucketName, objectName, minio.StatObjectOptions{})
	if err == nil {
		return nil, fmt.Errorf("file %s already exists in MinIO", objectName)
	}

	_, err = c.minioClient.PutObject(context.Background(), bucketName, objectName, file, fileHeader.Size, minio.PutObjectOptions{
		ContentType: fileType.Name,
	})
	if err != nil {
		return nil, err
	}

	fileURL := fmt.Sprintf("http://%s/%s/%s", c.getExternalHost(), bucketName, objectName)

	newFile := &models.File{
		Name:       objectName,
		FileTypeID: fileType.ID,
		URL:        fileURL,
	}

	if err := c.repo.CreateFile(newFile); err != nil {
		return nil, err
	}

	return newFile, nil
}

// GetFile получает актуальный URL файла из MinIO
func (c *FileController) GetFile(id int) (*fc.File, error) {
	file, err := c.repo.GetFileByID(id)
	if err != nil {
		return nil, err
	}

	// Проверяем, существует ли файл в MinIO
	_, err = c.minioClient.StatObject(context.Background(), c.minioConfig.Bucket, file.Name, minio.StatObjectOptions{})
	if err != nil {
		if minio.ToErrorResponse(err).Code == "NoSuchKey" {
			return nil, errors.New(fmt.Sprintf("Don't have file: %s in filesource", file.Name))
		}
		return nil, err
	}

	actualURL := fmt.Sprintf("http://%s/%s/%s", c.getExternalHost(), c.minioConfig.Bucket, file.Name)
	if file.URL != actualURL {
		file.URL = actualURL
		errorRepo := c.repo.UpdateFile(file)
		if errorRepo != nil {
			return nil, errorRepo
		}
	}

	fileTypeContract := &fc.FileType{
		ID:   file.FileTypeID,
		Name: file.FileType.Name,
	}
	fileContract := &fc.File{
		ID:         file.ID,
		Name:       file.Name,
		FileTypeID: file.FileTypeID,
		URL:        file.URL,
		CreatedAt:  file.CreatedAt,
		FileType:   *fileTypeContract,
	}

	return fileContract, nil
}

// RenameFile изменяет имя файла в MinIO и обновляет БД
func (c *FileController) RenameFile(id int, newName string) (*models.File, error) {
	file, err := c.repo.GetFileByID(id)
	if err != nil {
		return nil, err
	}

	oldObjectName := file.Name
	newObjectName := newName

	var newObjectType string
	parts := strings.Split(file.FileType.Name, "/")
	if len(parts) == 2 {
		newObjectType = parts[1]
	} else {
		log.Printf("unsupported FileType in database: %d", file.FileTypeID)
		return nil, errors.New("internal server error with file types")
	}

	if _, err = c.repo.GetFileByName(newObjectName + "." + newObjectType); err == nil {
		return file, xerrors.Errorf("file with new name: %s already exists", newObjectName+"."+newObjectType)
	}

	// Копируем объект в MinIO с новым именем
	srcOpts := minio.CopySrcOptions{
		Bucket: c.minioConfig.Bucket,
		Object: oldObjectName,
	}
	destOpts := minio.CopyDestOptions{
		Bucket: c.minioConfig.Bucket,
		Object: newObjectName + "." + newObjectType,
	}
	_, err = c.minioClient.CopyObject(context.Background(), destOpts, srcOpts)
	if err != nil {
		return nil, err
	}

	// Удаляем старый файл в MinIO
	err = c.minioClient.RemoveObject(context.Background(), c.minioConfig.Bucket, oldObjectName, minio.RemoveObjectOptions{})
	if err != nil {
		return nil, err
	}

	// Обновляем запись в БД
	file.Name = newObjectName + "." + newObjectType
	file.URL = fmt.Sprintf("http://%s/%s/%s", c.getExternalHost(), c.minioConfig.Bucket, file.Name)
	err = c.repo.UpdateFile(file)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (c *FileController) GetFileNamesWithPagination(limit, offset int) (*[]dto.FileInformation, error) {
	return c.repo.GetFileNamesWithPagination(limit, offset)
}

// getExternalHost возвращает внешний хост MinIO или обычный хост если внешний не задан
func (c *FileController) getExternalHost() string {
	if c.minioConfig.ExternalHost != "" {
		return c.minioConfig.ExternalHost
	}
	return c.minioConfig.Host
}
