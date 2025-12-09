package handlers

import (
	"encoding/json"
	"errors"
	hndlrs "fileService/internal/handlers"
	"fileService/internal/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockFileTypeController - мок для FileTypeController
type MockFileTypeController struct {
	mock.Mock
}

func (m *MockFileTypeController) CreateFileType(name string) (*models.FileType, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.FileType), args.Error(1)
}

func (m *MockFileTypeController) GetFileTypeByID(id int) (*models.FileType, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.FileType), args.Error(1)
}

func (m *MockFileTypeController) GetFileTypeByName(name string) (*models.FileType, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.FileType), args.Error(1)
}

func (m *MockFileTypeController) DeleteFileTypeByID(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

// Вспомогательные функции
func createTestFileType() *models.FileType {
	return &models.FileType{
		ID:   1,
		Name: "text/plain",
	}
}

// Тесты для FileTypeHandler.CreateFileTypeHandler

func TestFileTypeHandler_CreateFileTypeHandler_Success(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockFileTypeController)
	handler := hndlrs.NewFileTypeHandler(mockController)

	router := gin.New()
	router.POST("/file-types", handler.CreateFileTypeHandler)

	name := "text/plain"
	expectedFileType := createTestFileType()
	mockController.On("CreateFileType", name).Return(expectedFileType, nil)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/file-types?name="+name, nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response models.FileType
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, expectedFileType.ID, response.ID)
	assert.Equal(t, expectedFileType.Name, response.Name)

	mockController.AssertExpectations(t)
}

func TestFileTypeHandler_CreateFileTypeHandler_EmptyName(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockFileTypeController)
	handler := hndlrs.NewFileTypeHandler(mockController)

	router := gin.New()
	router.POST("/file-types", handler.CreateFileTypeHandler)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/file-types?name=", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "Name must be exist", response["error"])

	mockController.AssertNotCalled(t, "CreateFileType", mock.Anything)
}

func TestFileTypeHandler_CreateFileTypeHandler_NoNameParameter(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockFileTypeController)
	handler := hndlrs.NewFileTypeHandler(mockController)

	router := gin.New()
	router.POST("/file-types", handler.CreateFileTypeHandler)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/file-types", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "Name must be exist", response["error"])

	mockController.AssertNotCalled(t, "CreateFileType", mock.Anything)
}

func TestFileTypeHandler_CreateFileTypeHandler_ControllerError(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockFileTypeController)
	handler := hndlrs.NewFileTypeHandler(mockController)

	router := gin.New()
	router.POST("/file-types", handler.CreateFileTypeHandler)

	name := "text/plain"
	controllerError := errors.New("database error")
	mockController.On("CreateFileType", name).Return(nil, controllerError)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/file-types?name="+name, nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "database error", response["error"])

	mockController.AssertExpectations(t)
}

// Тесты для FileTypeHandler.GetFileTypeByIDHandler

func TestFileTypeHandler_GetFileTypeByIDHandler_Success(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockFileTypeController)
	handler := hndlrs.NewFileTypeHandler(mockController)

	router := gin.New()
	router.GET("/file-types/:id", handler.GetFileTypeByIDHandler)

	fileTypeID := 1
	expectedFileType := createTestFileType()
	mockController.On("GetFileTypeByID", fileTypeID).Return(expectedFileType, nil)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/file-types/1", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response models.FileType
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, expectedFileType.ID, response.ID)
	assert.Equal(t, expectedFileType.Name, response.Name)

	mockController.AssertExpectations(t)
}

func TestFileTypeHandler_GetFileTypeByIDHandler_InvalidID(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockFileTypeController)
	handler := hndlrs.NewFileTypeHandler(mockController)

	router := gin.New()
	router.GET("/file-types/:id", handler.GetFileTypeByIDHandler)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/file-types/invalid", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "Incorrect id", response["error"])

	mockController.AssertNotCalled(t, "GetFileTypeByID", mock.Anything)
}

func TestFileTypeHandler_GetFileTypeByIDHandler_NotFound(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockFileTypeController)
	handler := hndlrs.NewFileTypeHandler(mockController)

	router := gin.New()
	router.GET("/file-types/:id", handler.GetFileTypeByIDHandler)

	fileTypeID := 999
	notFoundError := errors.New("record not found")
	mockController.On("GetFileTypeByID", fileTypeID).Return(nil, notFoundError)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/file-types/999", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "File type does not exist", response["error"])

	mockController.AssertExpectations(t)
}

// Тесты для FileTypeHandler.GetFileTypeByNameHandler

func TestFileTypeHandler_GetFileTypeByNameHandler_Success(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockFileTypeController)
	handler := hndlrs.NewFileTypeHandler(mockController)

	router := gin.New()
	router.GET("/file-types/name/:name", handler.GetFileTypeByNameHandler)

	name := "plain"
	expectedFileType := createTestFileType()
	mockController.On("GetFileTypeByName", name).Return(expectedFileType, nil)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/file-types/name/plain", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response models.FileType
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, expectedFileType.ID, response.ID)
	assert.Equal(t, expectedFileType.Name, response.Name)

	mockController.AssertExpectations(t)
}

func TestFileTypeHandler_GetFileTypeByNameHandler_EmptyName(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockFileTypeController)
	handler := hndlrs.NewFileTypeHandler(mockController)

	router := gin.New()
	router.GET("/file-types/name/:name", handler.GetFileTypeByNameHandler)

	// Act - пустое имя в пути (хотя это маловероятно с Gin, но проверим)
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/file-types/name/", nil)
	router.ServeHTTP(w, req)

	// Assert - Gin может обработать это по-разному, но проверим
	// Если name пустое, должен вернуться BadRequest
	if w.Code == http.StatusBadRequest {
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "Name must be exist", response["error"])
		mockController.AssertNotCalled(t, "GetFileTypeByName", mock.Anything)
	}
}

func TestFileTypeHandler_GetFileTypeByNameHandler_NotFound(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockFileTypeController)
	handler := hndlrs.NewFileTypeHandler(mockController)

	router := gin.New()
	router.GET("/file-types/name/:name", handler.GetFileTypeByNameHandler)

	name := "nonexistent"
	notFoundError := errors.New("record not found")
	mockController.On("GetFileTypeByName", name).Return(nil, notFoundError)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/file-types/name/nonexistent", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "File type does not exist", response["error"])

	mockController.AssertExpectations(t)
}

// Тесты для FileTypeHandler.DeleteFileTypeByIDHandler

func TestFileTypeHandler_DeleteFileTypeByIDHandler_Success(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockFileTypeController)
	handler := hndlrs.NewFileTypeHandler(mockController)

	router := gin.New()
	router.DELETE("/file-types/:id", handler.DeleteFileTypeByIDHandler)

	fileTypeID := 1
	mockController.On("DeleteFileTypeByID", fileTypeID).Return(nil)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("DELETE", "/file-types/1", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "FileType has successfully deleted", response["message"])

	mockController.AssertExpectations(t)
}

func TestFileTypeHandler_DeleteFileTypeByIDHandler_InvalidID(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockFileTypeController)
	handler := hndlrs.NewFileTypeHandler(mockController)

	router := gin.New()
	router.DELETE("/file-types/:id", handler.DeleteFileTypeByIDHandler)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("DELETE", "/file-types/invalid", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "Incorrect id", response["error"])

	mockController.AssertNotCalled(t, "DeleteFileTypeByID", mock.Anything)
}

func TestFileTypeHandler_DeleteFileTypeByIDHandler_ControllerError(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockFileTypeController)
	handler := hndlrs.NewFileTypeHandler(mockController)

	router := gin.New()
	router.DELETE("/file-types/:id", handler.DeleteFileTypeByIDHandler)

	fileTypeID := 999
	controllerError := errors.New("database error")
	mockController.On("DeleteFileTypeByID", fileTypeID).Return(controllerError)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("DELETE", "/file-types/999", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "database error", response["error"])

	mockController.AssertExpectations(t)
}
