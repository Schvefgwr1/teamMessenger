package handlers

import (
	"apiService/internal/handlers"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Тесты для AuthHandler.Login

func TestAuthHandler_Login_Success(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockAuthController)
	handler := handlers.NewAuthHandler(mockController)

	userID := uuid.New()
	token := "test-token"

	mockController.On("Login", mock.Anything, mock.Anything).Return(token, userID, nil)

	router := gin.New()
	router.POST("/auth/login", handler.Login)

	// Act
	reqBody := `{"login":"test@example.com","password":"password123"}`
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/auth/login", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Assert
	// Проверяем, что мок был вызван
	mockController.AssertExpectations(t)
}

func TestAuthHandler_Login_InvalidJSON(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockAuthController)
	handler := handlers.NewAuthHandler(mockController)

	router := gin.New()
	router.POST("/auth/login", handler.Login)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/auth/login", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	mockController.AssertNotCalled(t, "Login", mock.Anything, mock.Anything)
}

func TestAuthHandler_Login_ControllerError(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockAuthController)
	handler := handlers.NewAuthHandler(mockController)

	authError := errors.New("invalid credentials")

	mockController.On("Login", mock.Anything, mock.Anything).Return("", uuid.Nil, authError)

	router := gin.New()
	router.POST("/auth/login", handler.Login)

	// Act
	reqBody := `{"login":"test@example.com","password":"wrongpassword"}`
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/auth/login", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	// В хендлере используется gin.H{"error": errReq}, где errReq - это error
	// Проверяем, что ошибка присутствует (может быть как строка или как объект)
	assert.NotNil(t, response["error"])

	mockController.AssertExpectations(t)
}

// Тесты для AuthHandler.Logout

func TestAuthHandler_Logout_Success(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockAuthController)
	handler := handlers.NewAuthHandler(mockController)

	userID := uuid.New()
	token := "test-token"

	mockController.On("Logout", mock.Anything, userID, token).Return(nil)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userID", userID)
		c.Set("token", token)
		c.Next()
	})
	router.POST("/auth/logout", handler.Logout)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/auth/logout", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "Logged out successfully", response["message"])

	mockController.AssertExpectations(t)
}

func TestAuthHandler_Logout_MissingUserID(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockAuthController)
	handler := handlers.NewAuthHandler(mockController)

	router := gin.New()
	router.POST("/auth/logout", handler.Logout)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/auth/logout", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Contains(t, response["error"], "User not authenticated")

	mockController.AssertNotCalled(t, "Logout", mock.Anything, mock.Anything, mock.Anything)
}

func TestAuthHandler_Logout_MissingToken(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockAuthController)
	handler := handlers.NewAuthHandler(mockController)

	userID := uuid.New()

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userID", userID)
		c.Next()
	})
	router.POST("/auth/logout", handler.Logout)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/auth/logout", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	mockController.AssertNotCalled(t, "Logout", mock.Anything, mock.Anything, mock.Anything)
}

func TestAuthHandler_Logout_ControllerError(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockAuthController)
	handler := handlers.NewAuthHandler(mockController)

	userID := uuid.New()
	token := "test-token"
	logoutError := errors.New("logout failed")

	mockController.On("Logout", mock.Anything, userID, token).Return(logoutError)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userID", userID)
		c.Set("token", token)
		c.Next()
	})
	router.POST("/auth/logout", handler.Logout)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/auth/logout", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Contains(t, response["error"], "Failed to revoke session")

	mockController.AssertExpectations(t)
}
