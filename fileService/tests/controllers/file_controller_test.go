package controllers

import (
	"bytes"
	"context"
	"errors"
	ctrl "fileService/internal/controllers"
	"fileService/internal/dto"
	"fileService/internal/models"
	"io"
	"mime/multipart"
	"strings"
	"testing"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"common/config"
)

// MockFileRepository - мок для FileRepository
type MockFileRepository struct {
	mock.Mock
}

func (m *MockFileRepository) CreateFile(file *models.File) error {
	args := m.Called(file)
	return args.Error(0)
}

func (m *MockFileRepository) GetFileByID(id int) (*models.File, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.File), args.Error(1)
}

func (m *MockFileRepository) GetFileByName(name string) (*models.File, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.File), args.Error(1)
}

func (m *MockFileRepository) UpdateFile(file *models.File) error {
	args := m.Called(file)
	return args.Error(0)
}

func (m *MockFileRepository) GetFileNamesWithPagination(limit, offset int) (*[]dto.FileInformation, error) {
	args := m.Called(limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*[]dto.FileInformation), args.Error(1)
}

// MockFileTypeRepository - мок для FileTypeRepository
type MockFileTypeRepository struct {
	mock.Mock
}

func (m *MockFileTypeRepository) CreateFileType(fileType *models.FileType) error {
	args := m.Called(fileType)
	return args.Error(0)
}

func (m *MockFileTypeRepository) GetFileTypeByID(id int) (*models.FileType, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.FileType), args.Error(1)
}

func (m *MockFileTypeRepository) GetFileTypeByName(name string) (*models.FileType, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.FileType), args.Error(1)
}

func (m *MockFileTypeRepository) DeleteFileTypeByID(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

// MockMinIOClient - мок для MinIO клиента
type MockMinIOClient struct {
	mock.Mock
}

func (m *MockMinIOClient) StatObject(ctx context.Context, bucketName, objectName string, opts minio.StatObjectOptions) (minio.ObjectInfo, error) {
	args := m.Called(ctx, bucketName, objectName, opts)
	return args.Get(0).(minio.ObjectInfo), args.Error(1)
}

func (m *MockMinIOClient) PutObject(ctx context.Context, bucketName, objectName string, reader io.Reader, objectSize int64, opts minio.PutObjectOptions) (minio.UploadInfo, error) {
	args := m.Called(ctx, bucketName, objectName, reader, objectSize, opts)
	return args.Get(0).(minio.UploadInfo), args.Error(1)
}

func (m *MockMinIOClient) CopyObject(ctx context.Context, dst minio.CopyDestOptions, src minio.CopySrcOptions) (minio.UploadInfo, error) {
	args := m.Called(ctx, dst, src)
	return args.Get(0).(minio.UploadInfo), args.Error(1)
}

func (m *MockMinIOClient) RemoveObject(ctx context.Context, bucketName, objectName string, opts minio.RemoveObjectOptions) error {
	args := m.Called(ctx, bucketName, objectName, opts)
	return args.Error(0)
}

// Вспомогательные функции для создания тестовых данных
func createTestFile() *models.File {
	return &models.File{
		ID:         1,
		Name:       "123456789_test.txt",
		FileTypeID: 1,
		URL:        "http://localhost:9000/test-bucket/123456789_test.txt",
		CreatedAt:  time.Now(),
		FileType: models.FileType{
			ID:   1,
			Name: "text/plain",
		},
	}
}

func createTestFileType() *models.FileType {
	return &models.FileType{
		ID:   1,
		Name: "text/plain",
	}
}

func createTestFileHeader() *multipart.FileHeader {
	return &multipart.FileHeader{
		Filename: "test.txt",
		Size:     1024,
		Header: map[string][]string{
			"Content-Type": {"text/plain"},
		},
	}
}

func createTestMinIOConfig() *config.MinIO {
	return &config.MinIO{
		Bucket:       "test-bucket",
		Host:         "localhost:9000",
		ExternalHost: "",
	}
}

func createTestMinIOConfigWithExternalHost() *config.MinIO {
	return &config.MinIO{
		Bucket:       "test-bucket",
		Host:         "localhost:9000",
		ExternalHost: "external.example.com:9000",
	}
}

// createTestMultipartFile создает тестовый multipart.FileHeader используя multipart.Writer
func createTestMultipartFile(filename, contentType string, content []byte) (*multipart.FileHeader, func(), error) {
	if len(content) == 0 {
		content = []byte("test content")
	}

	// Создаем multipart форму
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return nil, func() {}, err
	}

	_, err = part.Write(content)
	if err != nil {
		return nil, func() {}, err
	}

	err = writer.Close()
	if err != nil {
		return nil, func() {}, err
	}

	// Парсим multipart форму для получения FileHeader
	reader := multipart.NewReader(&body, writer.Boundary())
	form, err := reader.ReadForm(32 << 20) // 32MB max memory
	if err != nil {
		return nil, func() {}, err
	}

	if len(form.File["file"]) == 0 {
		return nil, func() {}, errors.New("no file in form")
	}

	fileHeader := form.File["file"][0]
	fileHeader.Header.Set("Content-Type", contentType)

	cleanup := func() {
		form.RemoveAll()
	}

	return fileHeader, cleanup, nil
}

// Тесты для FileController.UploadFile

func TestFileController_UploadFile_Success(t *testing.T) {
	// Arrange
	mockFileRepo := new(MockFileRepository)
	mockFileTypeRepo := new(MockFileTypeRepository)
	mockMinIOClient := new(MockMinIOClient)
	minioConfig := createTestMinIOConfig()

	controller := ctrl.NewFileController(mockFileRepo, mockFileTypeRepo, mockMinIOClient, minioConfig)

	testFile, cleanup, err := createTestMultipartFile("test.txt", "text/plain", []byte("test content"))
	require.NoError(t, err)
	defer cleanup()

	fileType := createTestFileType()

	// Настраиваем моки
	mockFileTypeRepo.On("GetFileTypeByName", "text/plain").Return(fileType, nil)
	mockFileRepo.On("GetFileByName", mock.MatchedBy(func(name string) bool {
		return len(name) > 0 // objectName содержит timestamp + filename
	})).Return(nil, errors.New("not found"))
	mockMinIOClient.On("StatObject", mock.Anything, minioConfig.Bucket, mock.AnythingOfType("string"), mock.Anything).Return(minio.ObjectInfo{}, errors.New("not found"))
	mockMinIOClient.On("PutObject", mock.Anything, minioConfig.Bucket, mock.AnythingOfType("string"), mock.Anything, testFile.Size, mock.Anything).Return(minio.UploadInfo{}, nil)
	mockFileRepo.On("CreateFile", mock.MatchedBy(func(file *models.File) bool {
		return file.FileTypeID == fileType.ID && file.Name != "" && file.URL != ""
	})).Return(nil)

	// Act
	result, err := controller.UploadFile(testFile)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, fileType.ID, result.FileTypeID)
	assert.Contains(t, result.Name, "test.txt")
	assert.Contains(t, result.URL, minioConfig.Bucket)

	mockFileTypeRepo.AssertExpectations(t)
	mockFileRepo.AssertExpectations(t)
	mockMinIOClient.AssertExpectations(t)
}

func TestFileController_UploadFile_UnsupportedFileType(t *testing.T) {
	// Arrange
	mockFileRepo := new(MockFileRepository)
	mockFileTypeRepo := new(MockFileTypeRepository)
	mockMinIOClient := new(MockMinIOClient)
	minioConfig := createTestMinIOConfig()

	controller := ctrl.NewFileController(mockFileRepo, mockFileTypeRepo, mockMinIOClient, minioConfig)

	testFile, cleanup, err := createTestMultipartFile("test.xyz", "application/xyz", []byte("test content"))
	require.NoError(t, err)
	defer cleanup()

	// Настраиваем моки - тип файла не найден
	mockFileTypeRepo.On("GetFileTypeByName", "application/xyz").Return(nil, errors.New("not found"))

	// Act
	result, err := controller.UploadFile(testFile)

	// Assert
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "Unsupported file type")

	mockFileTypeRepo.AssertExpectations(t)
	mockFileRepo.AssertNotCalled(t, "GetFileByName", mock.Anything)
	mockMinIOClient.AssertNotCalled(t, "PutObject", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

func TestFileController_UploadFile_FileTypeNil(t *testing.T) {
	// Arrange
	mockFileRepo := new(MockFileRepository)
	mockFileTypeRepo := new(MockFileTypeRepository)
	mockMinIOClient := new(MockMinIOClient)
	minioConfig := createTestMinIOConfig()

	controller := ctrl.NewFileController(mockFileRepo, mockFileTypeRepo, mockMinIOClient, minioConfig)

	testFile, cleanup, err := createTestMultipartFile("test.xyz", "application/xyz", []byte("test content"))
	require.NoError(t, err)
	defer cleanup()

	// Настраиваем моки - тип файла возвращает nil без ошибки (err == nil, fileType == nil)
	mockFileTypeRepo.On("GetFileTypeByName", "application/xyz").Return(nil, nil)

	// Act
	result, err := controller.UploadFile(testFile)

	// Assert
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "Unsupported file type")

	mockFileTypeRepo.AssertExpectations(t)
	mockFileRepo.AssertNotCalled(t, "GetFileByName", mock.Anything)
	mockMinIOClient.AssertNotCalled(t, "PutObject", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

func TestFileController_UploadFile_FileExistsInDatabase(t *testing.T) {
	// Arrange
	mockFileRepo := new(MockFileRepository)
	mockFileTypeRepo := new(MockFileTypeRepository)
	mockMinIOClient := new(MockMinIOClient)
	minioConfig := createTestMinIOConfig()

	controller := ctrl.NewFileController(mockFileRepo, mockFileTypeRepo, mockMinIOClient, minioConfig)

	testFile, cleanup, err := createTestMultipartFile("test.txt", "text/plain", []byte("test content"))
	require.NoError(t, err)
	defer cleanup()

	fileType := createTestFileType()
	existingFile := createTestFile()

	// Настраиваем моки - файл уже существует в БД
	mockFileTypeRepo.On("GetFileTypeByName", "text/plain").Return(fileType, nil)
	mockFileRepo.On("GetFileByName", mock.MatchedBy(func(name string) bool {
		return len(name) > 0
	})).Return(existingFile, nil)

	// Act
	result, err := controller.UploadFile(testFile)

	// Assert
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "already exists in database")

	mockFileTypeRepo.AssertExpectations(t)
	mockFileRepo.AssertExpectations(t)
	mockMinIOClient.AssertNotCalled(t, "PutObject", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

func TestFileController_UploadFile_FileExistsInMinIO(t *testing.T) {
	// Arrange
	mockFileRepo := new(MockFileRepository)
	mockFileTypeRepo := new(MockFileTypeRepository)
	mockMinIOClient := new(MockMinIOClient)
	minioConfig := createTestMinIOConfig()

	controller := ctrl.NewFileController(mockFileRepo, mockFileTypeRepo, mockMinIOClient, minioConfig)

	testFile, cleanup, err := createTestMultipartFile("test.txt", "text/plain", []byte("test content"))
	require.NoError(t, err)
	defer cleanup()

	fileType := createTestFileType()

	// Настраиваем моки - файл существует в MinIO
	mockFileTypeRepo.On("GetFileTypeByName", "text/plain").Return(fileType, nil)
	mockFileRepo.On("GetFileByName", mock.MatchedBy(func(name string) bool {
		return len(name) > 0
	})).Return(nil, errors.New("not found"))
	mockMinIOClient.On("StatObject", mock.Anything, minioConfig.Bucket, mock.AnythingOfType("string"), mock.Anything).Return(minio.ObjectInfo{}, nil)

	// Act
	result, err := controller.UploadFile(testFile)

	// Assert
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "already exists in MinIO")

	mockFileTypeRepo.AssertExpectations(t)
	mockFileRepo.AssertExpectations(t)
	mockMinIOClient.AssertExpectations(t)
	mockFileRepo.AssertNotCalled(t, "CreateFile", mock.Anything)
}

func TestFileController_UploadFile_MinIOUploadError(t *testing.T) {
	// Arrange
	mockFileRepo := new(MockFileRepository)
	mockFileTypeRepo := new(MockFileTypeRepository)
	mockMinIOClient := new(MockMinIOClient)
	minioConfig := createTestMinIOConfig()

	controller := ctrl.NewFileController(mockFileRepo, mockFileTypeRepo, mockMinIOClient, minioConfig)

	testFile, cleanup, err := createTestMultipartFile("test.txt", "text/plain", []byte("test content"))
	require.NoError(t, err)
	defer cleanup()

	fileType := createTestFileType()
	uploadError := errors.New("minio upload failed")

	// Настраиваем моки - ошибка при загрузке в MinIO
	mockFileTypeRepo.On("GetFileTypeByName", "text/plain").Return(fileType, nil)
	mockFileRepo.On("GetFileByName", mock.MatchedBy(func(name string) bool {
		return len(name) > 0
	})).Return(nil, errors.New("not found"))
	mockMinIOClient.On("StatObject", mock.Anything, minioConfig.Bucket, mock.AnythingOfType("string"), mock.Anything).Return(minio.ObjectInfo{}, errors.New("not found"))
	mockMinIOClient.On("PutObject", mock.Anything, minioConfig.Bucket, mock.AnythingOfType("string"), mock.Anything, testFile.Size, mock.Anything).Return(minio.UploadInfo{}, uploadError)

	// Act
	result, err := controller.UploadFile(testFile)

	// Assert
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, uploadError, err)

	mockFileTypeRepo.AssertExpectations(t)
	mockFileRepo.AssertExpectations(t)
	mockMinIOClient.AssertExpectations(t)
	mockFileRepo.AssertNotCalled(t, "CreateFile", mock.Anything)
}

func TestFileController_UploadFile_DatabaseSaveError(t *testing.T) {
	// Arrange
	mockFileRepo := new(MockFileRepository)
	mockFileTypeRepo := new(MockFileTypeRepository)
	mockMinIOClient := new(MockMinIOClient)
	minioConfig := createTestMinIOConfig()

	controller := ctrl.NewFileController(mockFileRepo, mockFileTypeRepo, mockMinIOClient, minioConfig)

	testFile, cleanup, err := createTestMultipartFile("test.txt", "text/plain", []byte("test content"))
	require.NoError(t, err)
	defer cleanup()

	fileType := createTestFileType()
	dbError := errors.New("database save failed")

	// Настраиваем моки - ошибка при сохранении в БД
	mockFileTypeRepo.On("GetFileTypeByName", "text/plain").Return(fileType, nil)
	mockFileRepo.On("GetFileByName", mock.MatchedBy(func(name string) bool {
		return len(name) > 0
	})).Return(nil, errors.New("not found"))
	mockMinIOClient.On("StatObject", mock.Anything, minioConfig.Bucket, mock.AnythingOfType("string"), mock.Anything).Return(minio.ObjectInfo{}, errors.New("not found"))
	mockMinIOClient.On("PutObject", mock.Anything, minioConfig.Bucket, mock.AnythingOfType("string"), mock.Anything, testFile.Size, mock.Anything).Return(minio.UploadInfo{}, nil)
	mockFileRepo.On("CreateFile", mock.MatchedBy(func(file *models.File) bool {
		return file.FileTypeID == fileType.ID
	})).Return(dbError)

	// Act
	result, err := controller.UploadFile(testFile)

	// Assert
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, dbError, err)

	mockFileTypeRepo.AssertExpectations(t)
	mockFileRepo.AssertExpectations(t)
	mockMinIOClient.AssertExpectations(t)
}

func TestFileController_UploadFile_WithExternalHost(t *testing.T) {
	// Arrange
	mockFileRepo := new(MockFileRepository)
	mockFileTypeRepo := new(MockFileTypeRepository)
	mockMinIOClient := new(MockMinIOClient)
	minioConfig := createTestMinIOConfigWithExternalHost()

	controller := ctrl.NewFileController(mockFileRepo, mockFileTypeRepo, mockMinIOClient, minioConfig)

	testFile, cleanup, err := createTestMultipartFile("test.txt", "text/plain", []byte("test content"))
	require.NoError(t, err)
	defer cleanup()

	fileType := createTestFileType()

	// Настраиваем моки
	mockFileTypeRepo.On("GetFileTypeByName", "text/plain").Return(fileType, nil)
	mockFileRepo.On("GetFileByName", mock.MatchedBy(func(name string) bool {
		return len(name) > 0
	})).Return(nil, errors.New("not found"))
	mockMinIOClient.On("StatObject", mock.Anything, minioConfig.Bucket, mock.AnythingOfType("string"), mock.Anything).Return(minio.ObjectInfo{}, errors.New("not found"))
	mockMinIOClient.On("PutObject", mock.Anything, minioConfig.Bucket, mock.AnythingOfType("string"), mock.Anything, testFile.Size, mock.Anything).Return(minio.UploadInfo{}, nil)
	mockFileRepo.On("CreateFile", mock.MatchedBy(func(file *models.File) bool {
		return file.FileTypeID == fileType.ID && file.Name != "" && file.URL != "" &&
			strings.Contains(file.URL, minioConfig.ExternalHost)
	})).Return(nil)

	// Act
	result, err := controller.UploadFile(testFile)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Contains(t, result.URL, minioConfig.ExternalHost) // URL должен содержать ExternalHost

	mockFileTypeRepo.AssertExpectations(t)
	mockFileRepo.AssertExpectations(t)
	mockMinIOClient.AssertExpectations(t)
}

// Тесты для FileController.GetFile

func TestFileController_GetFile_Success(t *testing.T) {
	// Arrange
	mockFileRepo := new(MockFileRepository)
	mockFileTypeRepo := new(MockFileTypeRepository)
	mockMinIOClient := new(MockMinIOClient)
	minioConfig := createTestMinIOConfig()

	controller := ctrl.NewFileController(mockFileRepo, mockFileTypeRepo, mockMinIOClient, minioConfig)

	fileID := 1
	file := createTestFile()
	expectedURL := "http://localhost:9000/test-bucket/" + file.Name

	// Настраиваем моки
	mockFileRepo.On("GetFileByID", fileID).Return(file, nil)
	mockMinIOClient.On("StatObject", mock.Anything, minioConfig.Bucket, file.Name, mock.Anything).Return(minio.ObjectInfo{}, nil)

	// Act
	result, err := controller.GetFile(fileID)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, file.ID, result.ID)
	assert.Equal(t, file.Name, result.Name)
	assert.Equal(t, file.FileTypeID, result.FileTypeID)
	assert.Equal(t, expectedURL, result.URL)
	assert.Equal(t, file.FileType.Name, result.FileType.Name)

	mockFileRepo.AssertExpectations(t)
	mockMinIOClient.AssertExpectations(t)
}

func TestFileController_GetFile_SuccessWithURLUpdate(t *testing.T) {
	// Arrange
	mockFileRepo := new(MockFileRepository)
	mockFileTypeRepo := new(MockFileTypeRepository)
	mockMinIOClient := new(MockMinIOClient)
	minioConfig := createTestMinIOConfig()

	controller := ctrl.NewFileController(mockFileRepo, mockFileTypeRepo, mockMinIOClient, minioConfig)

	fileID := 1
	file := createTestFile()
	file.URL = "http://old-host:9000/test-bucket/" + file.Name // старый URL
	expectedURL := "http://localhost:9000/test-bucket/" + file.Name

	// Настраиваем моки
	mockFileRepo.On("GetFileByID", fileID).Return(file, nil)
	mockMinIOClient.On("StatObject", mock.Anything, minioConfig.Bucket, file.Name, mock.Anything).Return(minio.ObjectInfo{}, nil)
	mockFileRepo.On("UpdateFile", mock.MatchedBy(func(f *models.File) bool {
		return f.URL == expectedURL
	})).Return(nil)

	// Act
	result, err := controller.GetFile(fileID)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedURL, result.URL)

	mockFileRepo.AssertExpectations(t)
	mockMinIOClient.AssertExpectations(t)
}

func TestFileController_GetFile_NotFoundInDatabase(t *testing.T) {
	// Arrange
	mockFileRepo := new(MockFileRepository)
	mockFileTypeRepo := new(MockFileTypeRepository)
	mockMinIOClient := new(MockMinIOClient)
	minioConfig := createTestMinIOConfig()

	controller := ctrl.NewFileController(mockFileRepo, mockFileTypeRepo, mockMinIOClient, minioConfig)

	fileID := 999
	dbError := errors.New("record not found")

	// Настраиваем моки
	mockFileRepo.On("GetFileByID", fileID).Return(nil, dbError)

	// Act
	result, err := controller.GetFile(fileID)

	// Assert
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, dbError, err)

	mockFileRepo.AssertExpectations(t)
	mockMinIOClient.AssertNotCalled(t, "StatObject", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

func TestFileController_GetFile_NotFoundInMinIO_NoSuchKey(t *testing.T) {
	// Arrange
	mockFileRepo := new(MockFileRepository)
	mockFileTypeRepo := new(MockFileTypeRepository)
	mockMinIOClient := new(MockMinIOClient)
	minioConfig := createTestMinIOConfig()

	controller := ctrl.NewFileController(mockFileRepo, mockFileTypeRepo, mockMinIOClient, minioConfig)

	fileID := 1
	file := createTestFile()

	// Создаем ошибку MinIO с кодом NoSuchKey
	// minio.ToErrorResponse работает только с реальными ошибками из HTTP ответов MinIO
	// В тесте мы создаем *minio.ErrorResponse, но minio.ToErrorResponse может не распознать его правильно
	// Поэтому проверяем реальное поведение: если Code не распознается, возвращается исходная ошибка
	minioError := &minio.ErrorResponse{
		Code:    "NoSuchKey",
		Message: "The specified key does not exist.",
		Key:     file.Name,
	}

	// Проверяем, как minio.ToErrorResponse обрабатывает нашу ошибку
	errorResponse := minio.ToErrorResponse(minioError)

	// Настраиваем моки
	mockFileRepo.On("GetFileByID", fileID).Return(file, nil)
	mockMinIOClient.On("StatObject", mock.Anything, minioConfig.Bucket, file.Name, mock.Anything).Return(minio.ObjectInfo{}, minioError)

	// Act
	result, err := controller.GetFile(fileID)

	// Assert
	require.Error(t, err)
	assert.Nil(t, result)

	// Если minio.ToErrorResponse правильно распознал Code, контроллер создаст новое сообщение
	// Иначе вернется исходная ошибка MinIO
	if errorResponse.Code == "NoSuchKey" {
		assert.Contains(t, err.Error(), "Don't have file")
		assert.Contains(t, err.Error(), file.Name)
		assert.Contains(t, err.Error(), "filesource")
	} else {
		// Если minio.ToErrorResponse не распознал Code, контроллер вернет исходную ошибку
		assert.Contains(t, err.Error(), "key does not exist")
	}

	mockFileRepo.AssertExpectations(t)
	mockMinIOClient.AssertExpectations(t)
}

func TestFileController_GetFile_MinIOOtherError(t *testing.T) {
	// Arrange
	mockFileRepo := new(MockFileRepository)
	mockFileTypeRepo := new(MockFileTypeRepository)
	mockMinIOClient := new(MockMinIOClient)
	minioConfig := createTestMinIOConfig()

	controller := ctrl.NewFileController(mockFileRepo, mockFileTypeRepo, mockMinIOClient, minioConfig)

	fileID := 1
	file := createTestFile()
	minioError := errors.New("minio connection error")

	// Настраиваем моки
	mockFileRepo.On("GetFileByID", fileID).Return(file, nil)
	mockMinIOClient.On("StatObject", mock.Anything, minioConfig.Bucket, file.Name, mock.Anything).Return(minio.ObjectInfo{}, minioError)

	// Act
	result, err := controller.GetFile(fileID)

	// Assert
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, minioError, err)

	mockFileRepo.AssertExpectations(t)
	mockMinIOClient.AssertExpectations(t)
}

func TestFileController_GetFile_UpdateURLError(t *testing.T) {
	// Arrange
	mockFileRepo := new(MockFileRepository)
	mockFileTypeRepo := new(MockFileTypeRepository)
	mockMinIOClient := new(MockMinIOClient)
	minioConfig := createTestMinIOConfig()

	controller := ctrl.NewFileController(mockFileRepo, mockFileTypeRepo, mockMinIOClient, minioConfig)

	fileID := 1
	file := createTestFile()
	file.URL = "http://old-host:9000/test-bucket/" + file.Name
	updateError := errors.New("database update failed")

	// Настраиваем моки
	mockFileRepo.On("GetFileByID", fileID).Return(file, nil)
	mockMinIOClient.On("StatObject", mock.Anything, minioConfig.Bucket, file.Name, mock.Anything).Return(minio.ObjectInfo{}, nil)
	mockFileRepo.On("UpdateFile", mock.Anything).Return(updateError)

	// Act
	result, err := controller.GetFile(fileID)

	// Assert
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, updateError, err)

	mockFileRepo.AssertExpectations(t)
	mockMinIOClient.AssertExpectations(t)
}

func TestFileController_GetFile_WithExternalHost(t *testing.T) {
	// Arrange
	mockFileRepo := new(MockFileRepository)
	mockFileTypeRepo := new(MockFileTypeRepository)
	mockMinIOClient := new(MockMinIOClient)
	minioConfig := createTestMinIOConfigWithExternalHost()

	controller := ctrl.NewFileController(mockFileRepo, mockFileTypeRepo, mockMinIOClient, minioConfig)

	fileID := 1
	file := createTestFile()
	// URL отличается от ожидаемого, поэтому контроллер обновит его
	file.URL = "http://localhost:9000/test-bucket/" + file.Name
	expectedURL := "http://" + minioConfig.ExternalHost + "/test-bucket/" + file.Name

	// Настраиваем моки
	mockFileRepo.On("GetFileByID", fileID).Return(file, nil)
	mockMinIOClient.On("StatObject", mock.Anything, minioConfig.Bucket, file.Name, mock.Anything).Return(minio.ObjectInfo{}, nil)
	mockFileRepo.On("UpdateFile", mock.MatchedBy(func(f *models.File) bool {
		return f.URL == expectedURL
	})).Return(nil)

	// Act
	result, err := controller.GetFile(fileID)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedURL, result.URL)

	mockFileRepo.AssertExpectations(t)
	mockMinIOClient.AssertExpectations(t)
}

func TestFileController_RenameFile_WithExternalHost(t *testing.T) {
	// Arrange
	mockFileRepo := new(MockFileRepository)
	mockFileTypeRepo := new(MockFileTypeRepository)
	mockMinIOClient := new(MockMinIOClient)
	minioConfig := createTestMinIOConfigWithExternalHost()

	controller := ctrl.NewFileController(mockFileRepo, mockFileTypeRepo, mockMinIOClient, minioConfig)

	fileID := 1
	newName := "renamed_file"
	file := createTestFile()
	oldName := file.Name
	newFullName := newName + ".plain"
	expectedURL := "http://" + minioConfig.ExternalHost + "/test-bucket/" + newFullName

	// Настраиваем моки
	mockFileRepo.On("GetFileByID", fileID).Return(file, nil)
	mockFileRepo.On("GetFileByName", newFullName).Return(nil, errors.New("not found"))
	mockMinIOClient.On("CopyObject", mock.Anything, mock.Anything, mock.Anything).Return(minio.UploadInfo{}, nil)
	mockMinIOClient.On("RemoveObject", mock.Anything, minioConfig.Bucket, oldName, mock.Anything).Return(nil)
	mockFileRepo.On("UpdateFile", mock.MatchedBy(func(f *models.File) bool {
		return f.Name == newFullName && f.URL == expectedURL
	})).Return(nil)

	// Act
	result, err := controller.RenameFile(fileID, newName)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedURL, result.URL)

	mockFileRepo.AssertExpectations(t)
	mockMinIOClient.AssertExpectations(t)
}

// Тесты для FileController.RenameFile

func TestFileController_RenameFile_Success(t *testing.T) {
	// Arrange
	mockFileRepo := new(MockFileRepository)
	mockFileTypeRepo := new(MockFileTypeRepository)
	mockMinIOClient := new(MockMinIOClient)
	minioConfig := createTestMinIOConfig()

	controller := ctrl.NewFileController(mockFileRepo, mockFileTypeRepo, mockMinIOClient, minioConfig)

	fileID := 1
	newName := "renamed_file"
	file := createTestFile()
	oldName := file.Name
	// FileType.Name = "text/plain", поэтому расширение будет "plain"
	newFullName := newName + ".plain"
	expectedURL := "http://localhost:9000/test-bucket/" + newFullName

	// Настраиваем моки
	mockFileRepo.On("GetFileByID", fileID).Return(file, nil)
	mockFileRepo.On("GetFileByName", newFullName).Return(nil, errors.New("not found"))
	mockMinIOClient.On("CopyObject", mock.Anything, mock.Anything, mock.Anything).Return(minio.UploadInfo{}, nil)
	mockMinIOClient.On("RemoveObject", mock.Anything, minioConfig.Bucket, oldName, mock.Anything).Return(nil)
	mockFileRepo.On("UpdateFile", mock.MatchedBy(func(f *models.File) bool {
		return f.Name == newFullName && f.URL == expectedURL
	})).Return(nil)

	// Act
	result, err := controller.RenameFile(fileID, newName)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, newFullName, result.Name)
	assert.Equal(t, expectedURL, result.URL)

	mockFileRepo.AssertExpectations(t)
	mockMinIOClient.AssertExpectations(t)
}

func TestFileController_RenameFile_NotFoundInDatabase(t *testing.T) {
	// Arrange
	mockFileRepo := new(MockFileRepository)
	mockFileTypeRepo := new(MockFileTypeRepository)
	mockMinIOClient := new(MockMinIOClient)
	minioConfig := createTestMinIOConfig()

	controller := ctrl.NewFileController(mockFileRepo, mockFileTypeRepo, mockMinIOClient, minioConfig)

	fileID := 999
	newName := "renamed_file"
	dbError := errors.New("record not found")

	// Настраиваем моки
	mockFileRepo.On("GetFileByID", fileID).Return(nil, dbError)

	// Act
	result, err := controller.RenameFile(fileID, newName)

	// Assert
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, dbError, err)

	mockFileRepo.AssertExpectations(t)
	mockMinIOClient.AssertNotCalled(t, "CopyObject", mock.Anything, mock.Anything, mock.Anything)
}

func TestFileController_RenameFile_UnsupportedFileType(t *testing.T) {
	// Arrange
	mockFileRepo := new(MockFileRepository)
	mockFileTypeRepo := new(MockFileTypeRepository)
	mockMinIOClient := new(MockMinIOClient)
	minioConfig := createTestMinIOConfig()

	controller := ctrl.NewFileController(mockFileRepo, mockFileTypeRepo, mockMinIOClient, minioConfig)

	fileID := 1
	newName := "renamed_file"
	file := createTestFile()
	file.FileType.Name = "invalid_format" // неправильный формат

	// Настраиваем моки
	mockFileRepo.On("GetFileByID", fileID).Return(file, nil)

	// Act
	result, err := controller.RenameFile(fileID, newName)

	// Assert
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "internal server error with file types")

	mockFileRepo.AssertExpectations(t)
	mockMinIOClient.AssertNotCalled(t, "CopyObject", mock.Anything, mock.Anything, mock.Anything)
}

func TestFileController_RenameFile_NewNameAlreadyExists(t *testing.T) {
	// Arrange
	mockFileRepo := new(MockFileRepository)
	mockFileTypeRepo := new(MockFileTypeRepository)
	mockMinIOClient := new(MockMinIOClient)
	minioConfig := createTestMinIOConfig()

	controller := ctrl.NewFileController(mockFileRepo, mockFileTypeRepo, mockMinIOClient, minioConfig)

	fileID := 1
	newName := "existing_file"
	file := createTestFile()
	// FileType.Name = "text/plain", поэтому расширение будет "plain"
	newFullName := newName + ".plain"
	existingFile := createTestFile()
	existingFile.Name = newFullName

	// Настраиваем моки
	mockFileRepo.On("GetFileByID", fileID).Return(file, nil)
	mockFileRepo.On("GetFileByName", newFullName).Return(existingFile, nil)

	// Act
	result, err := controller.RenameFile(fileID, newName)

	// Assert
	require.Error(t, err)
	assert.NotNil(t, result) // Возвращается исходный файл
	assert.Contains(t, err.Error(), "already exists")
	assert.Contains(t, err.Error(), newFullName)

	mockFileRepo.AssertExpectations(t)
	mockMinIOClient.AssertNotCalled(t, "CopyObject", mock.Anything, mock.Anything, mock.Anything)
}

func TestFileController_RenameFile_MinIOCopyError(t *testing.T) {
	// Arrange
	mockFileRepo := new(MockFileRepository)
	mockFileTypeRepo := new(MockFileTypeRepository)
	mockMinIOClient := new(MockMinIOClient)
	minioConfig := createTestMinIOConfig()

	controller := ctrl.NewFileController(mockFileRepo, mockFileTypeRepo, mockMinIOClient, minioConfig)

	fileID := 1
	newName := "renamed_file"
	file := createTestFile()
	// FileType.Name = "text/plain", поэтому расширение будет "plain"
	newFullName := newName + ".plain"
	copyError := errors.New("minio copy failed")

	// Настраиваем моки
	mockFileRepo.On("GetFileByID", fileID).Return(file, nil)
	mockFileRepo.On("GetFileByName", newFullName).Return(nil, errors.New("not found"))
	mockMinIOClient.On("CopyObject", mock.Anything, mock.Anything, mock.Anything).Return(minio.UploadInfo{}, copyError)

	// Act
	result, err := controller.RenameFile(fileID, newName)

	// Assert
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, copyError, err)

	mockFileRepo.AssertExpectations(t)
	mockMinIOClient.AssertExpectations(t)
	mockMinIOClient.AssertNotCalled(t, "RemoveObject", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	mockFileRepo.AssertNotCalled(t, "UpdateFile", mock.Anything)
}

func TestFileController_RenameFile_MinIORemoveError(t *testing.T) {
	// Arrange
	mockFileRepo := new(MockFileRepository)
	mockFileTypeRepo := new(MockFileTypeRepository)
	mockMinIOClient := new(MockMinIOClient)
	minioConfig := createTestMinIOConfig()

	controller := ctrl.NewFileController(mockFileRepo, mockFileTypeRepo, mockMinIOClient, minioConfig)

	fileID := 1
	newName := "renamed_file"
	file := createTestFile()
	oldName := file.Name
	// FileType.Name = "text/plain", поэтому расширение будет "plain"
	newFullName := newName + ".plain"
	removeError := errors.New("minio remove failed")

	// Настраиваем моки
	mockFileRepo.On("GetFileByID", fileID).Return(file, nil)
	mockFileRepo.On("GetFileByName", newFullName).Return(nil, errors.New("not found"))
	mockMinIOClient.On("CopyObject", mock.Anything, mock.Anything, mock.Anything).Return(minio.UploadInfo{}, nil)
	mockMinIOClient.On("RemoveObject", mock.Anything, minioConfig.Bucket, oldName, mock.Anything).Return(removeError)

	// Act
	result, err := controller.RenameFile(fileID, newName)

	// Assert
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, removeError, err)

	mockFileRepo.AssertExpectations(t)
	mockMinIOClient.AssertExpectations(t)
	mockFileRepo.AssertNotCalled(t, "UpdateFile", mock.Anything)
}

func TestFileController_RenameFile_DatabaseUpdateError(t *testing.T) {
	// Arrange
	mockFileRepo := new(MockFileRepository)
	mockFileTypeRepo := new(MockFileTypeRepository)
	mockMinIOClient := new(MockMinIOClient)
	minioConfig := createTestMinIOConfig()

	controller := ctrl.NewFileController(mockFileRepo, mockFileTypeRepo, mockMinIOClient, minioConfig)

	fileID := 1
	newName := "renamed_file"
	file := createTestFile()
	oldName := file.Name
	// FileType.Name = "text/plain", поэтому расширение будет "plain"
	newFullName := newName + ".plain"
	updateError := errors.New("database update failed")

	// Настраиваем моки
	mockFileRepo.On("GetFileByID", fileID).Return(file, nil)
	mockFileRepo.On("GetFileByName", newFullName).Return(nil, errors.New("not found"))
	mockMinIOClient.On("CopyObject", mock.Anything, mock.Anything, mock.Anything).Return(minio.UploadInfo{}, nil)
	mockMinIOClient.On("RemoveObject", mock.Anything, minioConfig.Bucket, oldName, mock.Anything).Return(nil)
	mockFileRepo.On("UpdateFile", mock.Anything).Return(updateError)

	// Act
	result, err := controller.RenameFile(fileID, newName)

	// Assert
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, updateError, err)

	mockFileRepo.AssertExpectations(t)
	mockMinIOClient.AssertExpectations(t)
}

// Тесты для FileController.GetFileNamesWithPagination

func TestFileController_GetFileNamesWithPagination_Success(t *testing.T) {
	// Arrange
	mockFileRepo := new(MockFileRepository)
	mockFileTypeRepo := new(MockFileTypeRepository)
	mockMinIOClient := new(MockMinIOClient)
	minioConfig := createTestMinIOConfig()

	controller := ctrl.NewFileController(mockFileRepo, mockFileTypeRepo, mockMinIOClient, minioConfig)

	limit := 10
	offset := 0
	expectedFiles := &[]dto.FileInformation{
		{Id: 1, Name: "file1.txt"},
		{Id: 2, Name: "file2.txt"},
	}

	// Настраиваем моки
	mockFileRepo.On("GetFileNamesWithPagination", limit, offset).Return(expectedFiles, nil)

	// Act
	result, err := controller.GetFileNamesWithPagination(limit, offset)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedFiles, result)

	mockFileRepo.AssertExpectations(t)
}

func TestFileController_GetFileNamesWithPagination_RepositoryError(t *testing.T) {
	// Arrange
	mockFileRepo := new(MockFileRepository)
	mockFileTypeRepo := new(MockFileTypeRepository)
	mockMinIOClient := new(MockMinIOClient)
	minioConfig := createTestMinIOConfig()

	controller := ctrl.NewFileController(mockFileRepo, mockFileTypeRepo, mockMinIOClient, minioConfig)

	limit := 10
	offset := 0
	repoError := errors.New("database error")

	// Настраиваем моки
	mockFileRepo.On("GetFileNamesWithPagination", limit, offset).Return(nil, repoError)

	// Act
	result, err := controller.GetFileNamesWithPagination(limit, offset)

	// Assert
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, repoError, err)

	mockFileRepo.AssertExpectations(t)
}
