package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	au "common/contracts/api-user"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"userService/internal/custom_errors"
	"userService/internal/handlers"
	"userService/internal/models"
)

// MockAuthController - мок для AuthController
type MockAuthController struct {
	mock.Mock
}

func (m *MockAuthController) Register(req *au.RegisterUserRequest) (*models.User, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockAuthController) Login(req *au.Login, ipAddress, userAgent string) (string, uuid.UUID, error) {
	args := m.Called(req, ipAddress, userAgent)
	if args.Get(0) == nil {
		return "", uuid.Nil, args.Error(2)
	}
	return args.Get(0).(string), args.Get(1).(uuid.UUID), args.Error(2)
}

// Тесты для AuthHandler.Register

func TestAuthHandler_Register_Success(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockAuthController)
	// Используем интерфейс напрямую
	handler := handlers.NewAuthHandler(mockController)

	userID := uuid.New()
	expectedUser := &models.User{
		ID:       userID,
		Username: "newuser",
		Email:    "newuser@example.com",
		RoleID:   1,
	}

	req := &au.RegisterUserRequest{
		Username: "newuser",
		Email:    "newuser@example.com",
		Password: "password123",
		RoleID:   1,
		Gender:   "male",
		Age:      25,
	}
	reqJSON, _ := json.Marshal(req)

	mockController.On("Register", mock.AnythingOfType("*api_user.RegisterUserRequest")).Return(expectedUser, nil)

	router := gin.New()
	router.POST("/auth/register", handler.Register)

	// Act
	w := httptest.NewRecorder()
	reqHTTP := httptest.NewRequest("POST", "/auth/register", bytes.NewBuffer(reqJSON))
	reqHTTP.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, reqHTTP)

	// Assert
	assert.Equal(t, http.StatusCreated, w.Code)

	var response models.User
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, expectedUser.ID, response.ID)
	assert.Equal(t, expectedUser.Username, response.Username)

	mockController.AssertExpectations(t)
}

func TestAuthHandler_Register_InvalidJSON(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockAuthController)
	// Используем интерфейс напрямую
	handler := handlers.NewAuthHandler(mockController)

	invalidJSON := []byte("{invalid json}")

	router := gin.New()
	router.POST("/auth/register", handler.Register)

	// Act
	w := httptest.NewRecorder()
	reqHTTP := httptest.NewRequest("POST", "/auth/register", bytes.NewBuffer(invalidJSON))
	reqHTTP.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, reqHTTP)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Contains(t, response["error"], "Invalid request payload")

	mockController.AssertNotCalled(t, "Register", mock.Anything)
}

func TestAuthHandler_Register_EmailConflict(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockAuthController)
	// Используем интерфейс напрямую
	handler := handlers.NewAuthHandler(mockController)

	req := &au.RegisterUserRequest{
		Username: "newuser",
		Email:    "existing@example.com",
		Password: "password123",
		RoleID:   1,
		Gender:   "male",
		Age:      25,
	}
	reqJSON, _ := json.Marshal(req)
	conflictError := custom_errors.NewUserEmailConflictError("existing@example.com")

	mockController.On("Register", mock.AnythingOfType("*api_user.RegisterUserRequest")).Return(nil, conflictError)

	router := gin.New()
	router.POST("/auth/register", handler.Register)

	// Act
	w := httptest.NewRecorder()
	reqHTTP := httptest.NewRequest("POST", "/auth/register", bytes.NewBuffer(reqJSON))
	reqHTTP.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, reqHTTP)

	// Assert
	assert.Equal(t, http.StatusConflict, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Contains(t, response["error"], "already exists")

	mockController.AssertExpectations(t)
}

func TestAuthHandler_Register_UsernameConflict(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockAuthController)
	// Используем интерфейс напрямую
	handler := handlers.NewAuthHandler(mockController)

	req := &au.RegisterUserRequest{
		Username: "existinguser",
		Email:    "new@example.com",
		Password: "password123",
		RoleID:   1,
		Gender:   "male",
		Age:      25,
	}
	reqJSON, _ := json.Marshal(req)
	conflictError := custom_errors.NewUserUsernameConflictError("existinguser")

	mockController.On("Register", mock.AnythingOfType("*api_user.RegisterUserRequest")).Return(nil, conflictError)

	router := gin.New()
	router.POST("/auth/register", handler.Register)

	// Act
	w := httptest.NewRecorder()
	reqHTTP := httptest.NewRequest("POST", "/auth/register", bytes.NewBuffer(reqJSON))
	reqHTTP.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, reqHTTP)

	// Assert
	assert.Equal(t, http.StatusConflict, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Contains(t, response["error"], "already exists")

	mockController.AssertExpectations(t)
}

// Тесты для AuthHandler.Login

func TestAuthHandler_Login_Success(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockAuthController)
	// Используем интерфейс напрямую
	handler := handlers.NewAuthHandler(mockController)

	userID := uuid.New()
	token := "test.jwt.token"
	req := &au.Login{
		Login:    "testuser",
		Password: "password123",
	}
	reqJSON, _ := json.Marshal(req)

	mockController.On("Login", mock.AnythingOfType("*api_user.Login"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(token, userID, nil)

	router := gin.New()
	router.POST("/auth/login", handler.Login)

	// Act
	w := httptest.NewRecorder()
	reqHTTP := httptest.NewRequest("POST", "/auth/login", bytes.NewBuffer(reqJSON))
	reqHTTP.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, reqHTTP)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, token, response["token"])
	assert.Equal(t, userID.String(), response["userID"])

	mockController.AssertExpectations(t)
}

func TestAuthHandler_Login_InvalidJSON(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockAuthController)
	// Используем интерфейс напрямую
	handler := handlers.NewAuthHandler(mockController)

	invalidJSON := []byte("{invalid json}")

	router := gin.New()
	router.POST("/auth/login", handler.Login)

	// Act
	w := httptest.NewRecorder()
	reqHTTP := httptest.NewRequest("POST", "/auth/login", bytes.NewBuffer(invalidJSON))
	reqHTTP.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, reqHTTP)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Contains(t, response["error"], "Invalid request payload")

	mockController.AssertNotCalled(t, "Login", mock.Anything, mock.Anything, mock.Anything)
}

func TestAuthHandler_Login_InvalidCredentials(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockAuthController)
	// Используем интерфейс напрямую
	handler := handlers.NewAuthHandler(mockController)

	req := &au.Login{
		Login:    "testuser",
		Password: "wrongpassword",
	}
	reqJSON, _ := json.Marshal(req)

	mockController.On("Login", mock.AnythingOfType("*api_user.Login"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return("", uuid.Nil, custom_errors.ErrInvalidCredentials)

	router := gin.New()
	router.POST("/auth/login", handler.Login)

	// Act
	w := httptest.NewRecorder()
	reqHTTP := httptest.NewRequest("POST", "/auth/login", bytes.NewBuffer(reqJSON))
	reqHTTP.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, reqHTTP)

	// Assert
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Contains(t, response["error"], "invalid credentials")

	mockController.AssertExpectations(t)
}

func TestAuthHandler_Login_TokenGenerationError(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockAuthController)
	// Используем интерфейс напрямую
	handler := handlers.NewAuthHandler(mockController)

	req := &au.Login{
		Login:    "testuser",
		Password: "password123",
	}
	reqJSON, _ := json.Marshal(req)

	mockController.On("Login", mock.AnythingOfType("*api_user.Login"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return("", uuid.Nil, custom_errors.ErrTokenGeneration)

	router := gin.New()
	router.POST("/auth/login", handler.Login)

	// Act
	w := httptest.NewRecorder()
	reqHTTP := httptest.NewRequest("POST", "/auth/login", bytes.NewBuffer(reqJSON))
	reqHTTP.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, reqHTTP)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Contains(t, response["error"], "token generation failed")

	mockController.AssertExpectations(t)
}
