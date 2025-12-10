package handlers

import (
	"bytes"
	fc "common/contracts/file-contracts"
	"encoding/json"
	"errors"
	"fileService/internal/dto"
	hndlrs "fileService/internal/handlers"
	"fileService/internal/models"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockFileController - мок для FileController
type MockFileController struct {
	mock.Mock
}

func (m *MockFileController) UploadFile(fileHeader *multipart.FileHeader) (*models.File, error) {
	args := m.Called(fileHeader)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.File), args.Error(1)
}

func (m *MockFileController) GetFile(id int) (*fc.File, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*fc.File), args.Error(1)
}

func (m *MockFileController) RenameFile(id int, newName string) (*models.File, error) {
	args := m.Called(id, newName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.File), args.Error(1)
}

func (m *MockFileController) GetFileNamesWithPagination(limit, offset int) (*[]dto.FileInformation, error) {
	args := m.Called(limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*[]dto.FileInformation), args.Error(1)
}

// Вспомогательные функции для создания тестовых данных
func createTestFileModel() *models.File {
	return &models.File{
		ID:         1,
		Name:       "test.txt",
		FileTypeID: 1,
		URL:        "http://localhost:9000/test-bucket/test.txt",
		CreatedAt:  time.Now(),
	}
}

func createTestFileContract() *fc.File {
	return &fc.File{
		ID:         1,
		Name:       "test.txt",
		FileTypeID: 1,
		URL:        "http://localhost:9000/test-bucket/test.txt",
		CreatedAt:  time.Now(),
		FileType: fc.FileType{
			ID:   1,
			Name: "text/plain",
		},
	}
}

func createTestFileInformationList() *[]dto.FileInformation {
	return &[]dto.FileInformation{
		{Id: 1, Name: "file1.txt"},
		{Id: 2, Name: "file2.txt"},
	}
}

// Тесты для FileHandler.UploadFileHandler

func TestFileHandler_UploadFileHandler_Success(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockFileController)
	handler := hndlrs.NewFileHandler(mockController)

	router := gin.New()
	router.POST("/files", handler.UploadFileHandler)

	// Создаем multipart форму с файлом
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "test.txt")
	require.NoError(t, err)
	_, err = part.Write([]byte("test content"))
	require.NoError(t, err)
	writer.Close()

	expectedFile := createTestFileModel()
	mockController.On("UploadFile", mock.AnythingOfType("*multipart.FileHeader")).Return(expectedFile, nil)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/files", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusCreated, w.Code)

	var response models.File
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, expectedFile.ID, response.ID)
	assert.Equal(t, expectedFile.Name, response.Name)

	mockController.AssertExpectations(t)
}

func TestFileHandler_UploadFileHandler_NoFileInForm(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockFileController)
	handler := hndlrs.NewFileHandler(mockController)

	router := gin.New()
	router.POST("/files", handler.UploadFileHandler)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/files", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "File in form does not exist", response["error"])

	mockController.AssertNotCalled(t, "UploadFile", mock.Anything)
}

func TestFileHandler_UploadFileHandler_UnsupportedFileType(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockFileController)
	handler := hndlrs.NewFileHandler(mockController)

	router := gin.New()
	router.POST("/files", handler.UploadFileHandler)

	// Создаем multipart форму
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "test.xyz")
	require.NoError(t, err)
	_, err = part.Write([]byte("test content"))
	require.NoError(t, err)
	writer.Close()

	uploadError := errors.New("Unsupported file type: application/xyz")
	mockController.On("UploadFile", mock.AnythingOfType("*multipart.FileHeader")).Return(nil, uploadError)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/files", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusUnsupportedMediaType, w.Code)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Contains(t, response["error"], "Unsupported file type")

	mockController.AssertExpectations(t)
}

func TestFileHandler_UploadFileHandler_FileExistsInDatabase(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockFileController)
	handler := hndlrs.NewFileHandler(mockController)

	router := gin.New()
	router.POST("/files", handler.UploadFileHandler)

	// Создаем multipart форму
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "test.txt")
	require.NoError(t, err)
	_, err = part.Write([]byte("test content"))
	require.NoError(t, err)
	writer.Close()

	uploadError := errors.New("file test.txt already exists in database")
	mockController.On("UploadFile", mock.AnythingOfType("*multipart.FileHeader")).Return(nil, uploadError)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/files", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusConflict, w.Code)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Contains(t, response["error"], "already exists in database")

	mockController.AssertExpectations(t)
}

func TestFileHandler_UploadFileHandler_FileExistsInMinIO(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockFileController)
	handler := hndlrs.NewFileHandler(mockController)

	router := gin.New()
	router.POST("/files", handler.UploadFileHandler)

	// Создаем multipart форму
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "test.txt")
	require.NoError(t, err)
	_, err = part.Write([]byte("test content"))
	require.NoError(t, err)
	writer.Close()

	uploadError := errors.New("file test.txt already exists in MinIO")
	mockController.On("UploadFile", mock.AnythingOfType("*multipart.FileHeader")).Return(nil, uploadError)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/files", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusConflict, w.Code)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Contains(t, response["error"], "already exists in MinIO")

	mockController.AssertExpectations(t)
}

func TestFileHandler_UploadFileHandler_InternalServerError(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockFileController)
	handler := hndlrs.NewFileHandler(mockController)

	router := gin.New()
	router.POST("/files", handler.UploadFileHandler)

	// Создаем multipart форму
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "test.txt")
	require.NoError(t, err)
	_, err = part.Write([]byte("test content"))
	require.NoError(t, err)
	writer.Close()

	uploadError := errors.New("internal server error")
	mockController.On("UploadFile", mock.AnythingOfType("*multipart.FileHeader")).Return(nil, uploadError)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/files", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "internal server error", response["error"])

	mockController.AssertExpectations(t)
}

// Тесты для FileHandler.GetFileHandler

func TestFileHandler_GetFileHandler_Success(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockFileController)
	handler := hndlrs.NewFileHandler(mockController)

	router := gin.New()
	router.GET("/files/:file_id", handler.GetFileHandler)

	fileID := 1
	expectedFile := createTestFileContract()
	mockController.On("GetFile", fileID).Return(expectedFile, nil)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/files/1", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response fc.File
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, expectedFile.ID, response.ID)
	assert.Equal(t, expectedFile.Name, response.Name)

	mockController.AssertExpectations(t)
}

func TestFileHandler_GetFileHandler_InvalidFileID(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockFileController)
	handler := hndlrs.NewFileHandler(mockController)

	router := gin.New()
	router.GET("/files/:file_id", handler.GetFileHandler)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/files/invalid", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "Incorrect file_id", response["error"])

	mockController.AssertNotCalled(t, "GetFile", mock.Anything)
}

func TestFileHandler_GetFileHandler_NotFound(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockFileController)
	handler := hndlrs.NewFileHandler(mockController)

	router := gin.New()
	router.GET("/files/:file_id", handler.GetFileHandler)

	fileID := 999
	notFoundError := errors.New("file not found")
	mockController.On("GetFile", fileID).Return(nil, notFoundError)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/files/999", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "file not found", response["error"])

	mockController.AssertExpectations(t)
}

// Тесты для FileHandler.RenameFileHandler

func TestFileHandler_RenameFileHandler_Success(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockFileController)
	handler := hndlrs.NewFileHandler(mockController)

	router := gin.New()
	router.PUT("/files/:file_id/rename", handler.RenameFileHandler)

	fileID := 1
	newName := "renamed_file"
	expectedFile := createTestFileModel()
	expectedFile.Name = "renamed_file.plain"
	mockController.On("RenameFile", fileID, newName).Return(expectedFile, nil)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("PUT", "/files/1/rename?name=renamed_file", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response models.File
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, expectedFile.Name, response.Name)

	mockController.AssertExpectations(t)
}

func TestFileHandler_RenameFileHandler_InvalidFileID(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockFileController)
	handler := hndlrs.NewFileHandler(mockController)

	router := gin.New()
	router.PUT("/files/:file_id/rename", handler.RenameFileHandler)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("PUT", "/files/invalid/rename?name=test", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "Incorrect file_id", response["error"])

	mockController.AssertNotCalled(t, "RenameFile", mock.Anything, mock.Anything)
}

func TestFileHandler_RenameFileHandler_EmptyName(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockFileController)
	handler := hndlrs.NewFileHandler(mockController)

	router := gin.New()
	router.PUT("/files/:file_id/rename", handler.RenameFileHandler)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("PUT", "/files/1/rename?name=", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "New name cannot be empty", response["error"])

	mockController.AssertNotCalled(t, "RenameFile", mock.Anything, mock.Anything)
}

func TestFileHandler_RenameFileHandler_NoNameParameter(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockFileController)
	handler := hndlrs.NewFileHandler(mockController)

	router := gin.New()
	router.PUT("/files/:file_id/rename", handler.RenameFileHandler)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("PUT", "/files/1/rename", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "New name cannot be empty", response["error"])

	mockController.AssertNotCalled(t, "RenameFile", mock.Anything, mock.Anything)
}

func TestFileHandler_RenameFileHandler_ControllerError(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockFileController)
	handler := hndlrs.NewFileHandler(mockController)

	router := gin.New()
	router.PUT("/files/:file_id/rename", handler.RenameFileHandler)

	fileID := 1
	newName := "renamed_file"
	controllerError := errors.New("rename failed")
	mockController.On("RenameFile", fileID, newName).Return(nil, controllerError)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("PUT", "/files/1/rename?name=renamed_file", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "rename failed", response["error"])

	mockController.AssertExpectations(t)
}

// Тесты для FileHandler.GetFileNamesHandler

func TestFileHandler_GetFileNamesHandler_Success(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockFileController)
	handler := hndlrs.NewFileHandler(mockController)

	router := gin.New()
	router.GET("/files/names", handler.GetFileNamesHandler)

	limit := 10
	offset := 0
	expectedFiles := createTestFileInformationList()
	mockController.On("GetFileNamesWithPagination", limit, offset).Return(expectedFiles, nil)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/files/names?limit=10&offset=0", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response []dto.FileInformation
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Len(t, response, 2)
	assert.Equal(t, (*expectedFiles)[0].Id, response[0].Id)

	mockController.AssertExpectations(t)
}

func TestFileHandler_GetFileNamesHandler_DefaultValues(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockFileController)
	handler := hndlrs.NewFileHandler(mockController)

	router := gin.New()
	router.GET("/files/names", handler.GetFileNamesHandler)

	limit := 10
	offset := 0
	expectedFiles := createTestFileInformationList()
	mockController.On("GetFileNamesWithPagination", limit, offset).Return(expectedFiles, nil)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/files/names", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response []dto.FileInformation
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Len(t, response, 2)

	mockController.AssertExpectations(t)
}

func TestFileHandler_GetFileNamesHandler_InvalidLimit(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockFileController)
	handler := hndlrs.NewFileHandler(mockController)

	router := gin.New()
	router.GET("/files/names", handler.GetFileNamesHandler)

	// Act - limit не является числом
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/files/names?limit=invalid", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "Limit must be a natural number", response["error"])

	mockController.AssertNotCalled(t, "GetFileNamesWithPagination", mock.Anything, mock.Anything)
}

func TestFileHandler_GetFileNamesHandler_InvalidLimitZero(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockFileController)
	handler := hndlrs.NewFileHandler(mockController)

	router := gin.New()
	router.GET("/files/names", handler.GetFileNamesHandler)

	// Act - limit = 0
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/files/names?limit=0", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "Limit must be a natural number", response["error"])

	mockController.AssertNotCalled(t, "GetFileNamesWithPagination", mock.Anything, mock.Anything)
}

func TestFileHandler_GetFileNamesHandler_InvalidLimitNegative(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockFileController)
	handler := hndlrs.NewFileHandler(mockController)

	router := gin.New()
	router.GET("/files/names", handler.GetFileNamesHandler)

	// Act - limit отрицательный
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/files/names?limit=-1", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "Limit must be a natural number", response["error"])

	mockController.AssertNotCalled(t, "GetFileNamesWithPagination", mock.Anything, mock.Anything)
}

func TestFileHandler_GetFileNamesHandler_InvalidOffset(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockFileController)
	handler := hndlrs.NewFileHandler(mockController)

	router := gin.New()
	router.GET("/files/names", handler.GetFileNamesHandler)

	// Act - offset не является числом
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/files/names?limit=10&offset=invalid", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "Offset must be a natural number", response["error"])

	mockController.AssertNotCalled(t, "GetFileNamesWithPagination", mock.Anything, mock.Anything)
}

func TestFileHandler_GetFileNamesHandler_InvalidOffsetNegative(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockFileController)
	handler := hndlrs.NewFileHandler(mockController)

	router := gin.New()
	router.GET("/files/names", handler.GetFileNamesHandler)

	// Act - offset отрицательный
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/files/names?limit=10&offset=-1", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "Offset must be a natural number", response["error"])

	mockController.AssertNotCalled(t, "GetFileNamesWithPagination", mock.Anything, mock.Anything)
}

func TestFileHandler_GetFileNamesHandler_ControllerError(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockFileController)
	handler := hndlrs.NewFileHandler(mockController)

	router := gin.New()
	router.GET("/files/names", handler.GetFileNamesHandler)

	limit := 10
	offset := 0
	controllerError := errors.New("database error")
	mockController.On("GetFileNamesWithPagination", limit, offset).Return(nil, controllerError)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/files/names?limit=10&offset=0", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "database error", response["error"])

	mockController.AssertExpectations(t)
}
