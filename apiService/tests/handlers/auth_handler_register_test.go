package handlers

import (
	"apiService/internal/dto"
	"apiService/internal/handlers"
	"bytes"
	au "common/contracts/api-user"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Тесты для AuthHandler.Register

func TestAuthHandler_Register_Success(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockAuthController)
	handler := handlers.NewAuthHandler(mockController)

	userID := uuid.New()
	expectedResponse := &dto.RegisterUserResponseGateway{
		User: &au.RegisterUserResponse{
			ID:       &userID,
			Username: stringPtr("testuser"),
			Email:    stringPtr("test@example.com"),
		},
	}

	mockController.On("Register", mock.Anything, mock.Anything).Return(expectedResponse)

	router := gin.New()
	router.POST("/auth/register", handler.Register)

	// Act - создаем multipart форму с данными и файлом
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Добавляем JSON данные
	jsonData := `{"username":"testuser","email":"test@example.com","password":"password123","age":25,"roleID":1}`
	writer.WriteField("data", jsonData)

	// Добавляем файл (опционально)
	part, err := writer.CreateFormFile("file", "avatar.jpg")
	require.NoError(t, err)
	_, err = part.Write([]byte("fake image content"))
	require.NoError(t, err)

	writer.Close()

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/auth/register", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusCreated, w.Code)

	mockController.AssertExpectations(t)
}

func TestAuthHandler_Register_WithoutFile(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockAuthController)
	handler := handlers.NewAuthHandler(mockController)

	userID := uuid.New()
	expectedResponse := &dto.RegisterUserResponseGateway{
		User: &au.RegisterUserResponse{
			ID:       &userID,
			Username: stringPtr("testuser"),
			Email:    stringPtr("test@example.com"),
		},
	}

	mockController.On("Register", mock.Anything, mock.Anything).Return(expectedResponse)

	router := gin.New()
	router.POST("/auth/register", handler.Register)

	// Act - создаем multipart форму только с данными (без файла)
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	jsonData := `{"username":"testuser","email":"test@example.com","password":"password123","age":25,"roleID":1}`
	writer.WriteField("data", jsonData)
	writer.Close()

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/auth/register", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusCreated, w.Code)

	mockController.AssertExpectations(t)
}

func TestAuthHandler_Register_InvalidJSON(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockAuthController)
	handler := handlers.NewAuthHandler(mockController)

	router := gin.New()
	router.POST("/auth/register", handler.Register)

	// Act
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.WriteField("data", "invalid json")
	writer.Close()

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/auth/register", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	mockController.AssertNotCalled(t, "Register", mock.Anything, mock.Anything)
}

func TestAuthHandler_Register_ControllerError(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockAuthController)
	handler := handlers.NewAuthHandler(mockController)

	errorMsg := "user already exists"
	expectedResponse := &dto.RegisterUserResponseGateway{
		Error: &errorMsg,
	}

	mockController.On("Register", mock.Anything, mock.Anything).Return(expectedResponse)

	router := gin.New()
	router.POST("/auth/register", handler.Register)

	// Act
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	jsonData := `{"username":"testuser","email":"test@example.com","password":"password123","age":25,"roleID":1}`
	writer.WriteField("data", jsonData)
	writer.Close()

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/auth/register", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	mockController.AssertExpectations(t)
}
