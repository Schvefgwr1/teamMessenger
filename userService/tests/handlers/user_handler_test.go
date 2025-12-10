package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	au "common/contracts/api-user"
	fc "common/contracts/file-contracts"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"userService/internal/custom_errors"
	"userService/internal/handlers"
	"userService/internal/handlers/dto"
	"userService/internal/models"
)

// MockUserController - мок для UserController
type MockUserController struct {
	mock.Mock
}

func (m *MockUserController) GetUserProfile(id uuid.UUID) (*models.User, *fc.File, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, nil, args.Error(2)
	}
	var file *fc.File
	if args.Get(1) != nil {
		file = args.Get(1).(*fc.File)
	}
	return args.Get(0).(*models.User), file, args.Error(2)
}

func (m *MockUserController) UpdateUserProfile(req *au.UpdateUserRequest, userId *uuid.UUID) error {
	args := m.Called(req, userId)
	return args.Error(0)
}

func (m *MockUserController) GetUserBrief(userID uuid.UUID, chatID string, requesterID string) (*dto.UserBriefResponse, error) {
	args := m.Called(userID, chatID, requesterID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.UserBriefResponse), args.Error(1)
}

func (m *MockUserController) SearchUsers(query string, limit int) (*dto.UserSearchResponse, error) {
	args := m.Called(query, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.UserSearchResponse), args.Error(1)
}

func (m *MockUserController) UpdateUserRole(userID uuid.UUID, roleID int) error {
	args := m.Called(userID, roleID)
	return args.Error(0)
}

// Вспомогательные функции для создания тестовых данных
func createTestUserModel() *models.User {
	userID := uuid.New()
	return &models.User{
		ID:       userID,
		Username: "testuser",
		Email:    "test@example.com",
		RoleID:   1,
	}
}

func createTestFileContract() *fc.File {
	return &fc.File{
		ID:  1,
		URL: "http://example.com/avatar.jpg",
	}
}

// Тесты для UserHandler.GetProfile

func TestUserHandler_GetProfile_Success(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockUserController)
	handler := handlers.NewUserHandler(mockController)

	userID := uuid.New()
	expectedUser := createTestUserModel()
	expectedUser.ID = userID
	expectedFile := createTestFileContract()

	mockController.On("GetUserProfile", userID).Return(expectedUser, expectedFile, nil)

	router := gin.New()
	router.GET("/users/:user_id", handler.GetProfile)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/users/"+userID.String(), nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.NotNil(t, response["user"])
	assert.NotNil(t, response["file"])

	mockController.AssertExpectations(t)
}

func TestUserHandler_GetProfile_InvalidUUID(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockUserController)
	handler := handlers.NewUserHandler(mockController)

	router := gin.New()
	router.GET("/users/:user_id", handler.GetProfile)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/users/invalid-uuid", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "Invalid user ID", response["error"])

	mockController.AssertNotCalled(t, "GetUserProfile", mock.Anything)
}

func TestUserHandler_GetProfile_NotFound(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockUserController)
	handler := handlers.NewUserHandler(mockController)

	userID := uuid.New()
	mockController.On("GetUserProfile", userID).Return(nil, nil, gorm.ErrRecordNotFound)

	router := gin.New()
	router.GET("/users/:user_id", handler.GetProfile)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/users/"+userID.String(), nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "User not found", response["error"])

	mockController.AssertExpectations(t)
}

func TestUserHandler_GetProfile_AvatarLoadError(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockUserController)
	handler := handlers.NewUserHandler(mockController)

	userID := uuid.New()
	expectedUser := createTestUserModel()
	expectedUser.ID = userID
	fileError := errors.New("file load error")

	mockController.On("GetUserProfile", userID).Return(expectedUser, nil, fileError)

	router := gin.New()
	router.GET("/users/:user_id", handler.GetProfile)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/users/"+userID.String(), nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.NotNil(t, response["user"])
	assert.Nil(t, response["file"])
	assert.Contains(t, response["error"], "failed to load avatar")

	mockController.AssertExpectations(t)
}

// Тесты для UserHandler.UpdateProfile

func TestUserHandler_UpdateProfile_Success(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockUserController)
	handler := handlers.NewUserHandler(mockController)

	userID := uuid.New()
	req := &au.UpdateUserRequest{
		Username: stringPtr("newusername"),
	}
	reqJSON, _ := json.Marshal(req)

	mockController.On("UpdateUserProfile", mock.AnythingOfType("*api_user.UpdateUserRequest"), &userID).Return(nil)

	router := gin.New()
	router.PUT("/users/:user_id", handler.UpdateProfile)

	// Act
	w := httptest.NewRecorder()
	reqHTTP := httptest.NewRequest("PUT", "/users/"+userID.String(), bytes.NewBuffer(reqJSON))
	reqHTTP.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, reqHTTP)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "Profile updated", response["message"])

	mockController.AssertExpectations(t)
}

func TestUserHandler_UpdateProfile_InvalidUUID(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockUserController)
	handler := handlers.NewUserHandler(mockController)

	req := &au.UpdateUserRequest{
		Username: stringPtr("newusername"),
	}
	reqJSON, _ := json.Marshal(req)

	router := gin.New()
	router.PUT("/users/:user_id", handler.UpdateProfile)

	// Act
	w := httptest.NewRecorder()
	reqHTTP := httptest.NewRequest("PUT", "/users/invalid-uuid", bytes.NewBuffer(reqJSON))
	reqHTTP.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, reqHTTP)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "Invalid user ID", response["error"])

	mockController.AssertNotCalled(t, "UpdateUserProfile", mock.Anything, mock.Anything)
}

func TestUserHandler_UpdateProfile_InvalidJSON(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockUserController)
	handler := handlers.NewUserHandler(mockController)

	userID := uuid.New()
	invalidJSON := []byte("{invalid json}")

	router := gin.New()
	router.PUT("/users/:user_id", handler.UpdateProfile)

	// Act
	w := httptest.NewRecorder()
	reqHTTP := httptest.NewRequest("PUT", "/users/"+userID.String(), bytes.NewBuffer(invalidJSON))
	reqHTTP.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, reqHTTP)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Contains(t, response["error"], "Invalid request payload")

	mockController.AssertNotCalled(t, "UpdateUserProfile", mock.Anything, mock.Anything)
}

func TestUserHandler_UpdateProfile_UsernameConflict(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockUserController)
	handler := handlers.NewUserHandler(mockController)

	userID := uuid.New()
	req := &au.UpdateUserRequest{
		Username: stringPtr("existinguser"),
	}
	reqJSON, _ := json.Marshal(req)
	conflictError := custom_errors.NewUserUsernameConflictError("existinguser")

	mockController.On("UpdateUserProfile", mock.AnythingOfType("*api_user.UpdateUserRequest"), &userID).Return(conflictError)

	router := gin.New()
	router.PUT("/users/:user_id", handler.UpdateProfile)

	// Act
	w := httptest.NewRecorder()
	reqHTTP := httptest.NewRequest("PUT", "/users/"+userID.String(), bytes.NewBuffer(reqJSON))
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

func TestUserHandler_UpdateProfile_InvalidCredentials(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockUserController)
	handler := handlers.NewUserHandler(mockController)

	userID := uuid.New()
	req := &au.UpdateUserRequest{
		Username: stringPtr("newusername"),
	}
	reqJSON, _ := json.Marshal(req)

	mockController.On("UpdateUserProfile", mock.AnythingOfType("*api_user.UpdateUserRequest"), &userID).Return(custom_errors.ErrInvalidCredentials)

	router := gin.New()
	router.PUT("/users/:user_id", handler.UpdateProfile)

	// Act
	w := httptest.NewRecorder()
	reqHTTP := httptest.NewRequest("PUT", "/users/"+userID.String(), bytes.NewBuffer(reqJSON))
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

func TestUserHandler_UpdateProfile_RoleNotFound(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockUserController)
	handler := handlers.NewUserHandler(mockController)

	userID := uuid.New()
	roleID := 999
	req := &au.UpdateUserRequest{
		RoleID: intPtr(roleID),
	}
	reqJSON, _ := json.Marshal(req)
	roleNotFoundError := custom_errors.NewRoleNotFoundError(roleID)

	mockController.On("UpdateUserProfile", mock.AnythingOfType("*api_user.UpdateUserRequest"), &userID).Return(roleNotFoundError)

	router := gin.New()
	router.PUT("/users/:user_id", handler.UpdateProfile)

	// Act
	w := httptest.NewRecorder()
	reqHTTP := httptest.NewRequest("PUT", "/users/"+userID.String(), bytes.NewBuffer(reqJSON))
	reqHTTP.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, reqHTTP)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Contains(t, response["error"], "does not exist")

	mockController.AssertExpectations(t)
}

func TestUserHandler_UpdateProfile_FileHTTPError(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockUserController)
	handler := handlers.NewUserHandler(mockController)

	userID := uuid.New()
	avatarID := 999
	req := &au.UpdateUserRequest{
		AvatarFileID: intPtr(avatarID),
	}
	reqJSON, _ := json.Marshal(req)
	fileHTTPError := custom_errors.NewGetFileHTTPError(avatarID, "connection error")

	mockController.On("UpdateUserProfile", mock.AnythingOfType("*api_user.UpdateUserRequest"), &userID).Return(fileHTTPError)

	router := gin.New()
	router.PUT("/users/:user_id", handler.UpdateProfile)

	// Act
	w := httptest.NewRecorder()
	reqHTTP := httptest.NewRequest("PUT", "/users/"+userID.String(), bytes.NewBuffer(reqJSON))
	reqHTTP.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, reqHTTP)

	// Assert
	assert.Equal(t, http.StatusBadGateway, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Contains(t, response["error"], "can't get file")

	mockController.AssertExpectations(t)
}

// Тесты для UserHandler.GetUserBrief

func TestUserHandler_GetUserBrief_Success(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockUserController)
	handler := handlers.NewUserHandler(mockController)

	userID := uuid.New()
	chatID := "chat-123"
	requesterID := "requester-123"
	expectedBrief := &dto.UserBriefResponse{
		Username:     "testuser",
		Email:        "test@example.com",
		ChatRoleName: "admin",
	}

	mockController.On("GetUserBrief", userID, chatID, requesterID).Return(expectedBrief, nil)

	router := gin.New()
	router.GET("/users/:user_id/brief", handler.GetUserBrief)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/users/"+userID.String()+"/brief?chatId="+chatID, nil)
	req.Header.Set("X-User-ID", requesterID)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response dto.UserBriefResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, expectedBrief.Username, response.Username)
	assert.Equal(t, expectedBrief.ChatRoleName, response.ChatRoleName)

	mockController.AssertExpectations(t)
}

func TestUserHandler_GetUserBrief_InvalidUUID(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockUserController)
	handler := handlers.NewUserHandler(mockController)

	router := gin.New()
	router.GET("/users/:user_id/brief", handler.GetUserBrief)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/users/invalid-uuid/brief?chatId=chat-123", nil)
	req.Header.Set("X-User-ID", "requester-123")
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "Invalid user ID", response["error"])

	mockController.AssertNotCalled(t, "GetUserBrief", mock.Anything, mock.Anything, mock.Anything)
}

func TestUserHandler_GetUserBrief_MissingChatID(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockUserController)
	handler := handlers.NewUserHandler(mockController)

	userID := uuid.New()

	router := gin.New()
	router.GET("/users/:user_id/brief", handler.GetUserBrief)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/users/"+userID.String()+"/brief", nil)
	req.Header.Set("X-User-ID", "requester-123")
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Contains(t, response["error"], "chatId query parameter is required")

	mockController.AssertNotCalled(t, "GetUserBrief", mock.Anything, mock.Anything, mock.Anything)
}

func TestUserHandler_GetUserBrief_MissingRequesterID(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockUserController)
	handler := handlers.NewUserHandler(mockController)

	userID := uuid.New()
	chatID := "chat-123"

	router := gin.New()
	router.GET("/users/:user_id/brief", handler.GetUserBrief)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/users/"+userID.String()+"/brief?chatId="+chatID, nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Contains(t, response["error"], "X-User-ID header is required")

	mockController.AssertNotCalled(t, "GetUserBrief", mock.Anything, mock.Anything, mock.Anything)
}

func TestUserHandler_GetUserBrief_UserNotFound(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockUserController)
	handler := handlers.NewUserHandler(mockController)

	userID := uuid.New()
	chatID := "chat-123"
	requesterID := "requester-123"

	mockController.On("GetUserBrief", userID, chatID, requesterID).Return(nil, custom_errors.ErrInvalidCredentials)

	router := gin.New()
	router.GET("/users/:user_id/brief", handler.GetUserBrief)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/users/"+userID.String()+"/brief?chatId="+chatID, nil)
	req.Header.Set("X-User-ID", requesterID)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "User not found", response["error"])

	mockController.AssertExpectations(t)
}

// Тесты для UserHandler.SearchUsers

func TestUserHandler_SearchUsers_Success(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockUserController)
	handler := handlers.NewUserHandler(mockController)

	query := "test"
	limit := 10
	expectedResponse := &dto.UserSearchResponse{
		Users: []dto.UserSearchResult{
			{ID: "1", Username: "testuser1", Email: "test1@example.com"},
			{ID: "2", Username: "testuser2", Email: "test2@example.com"},
		},
	}

	mockController.On("SearchUsers", query, limit).Return(expectedResponse, nil)

	router := gin.New()
	router.GET("/users/search", handler.SearchUsers)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/users/search?q="+query+"&limit="+strconv.Itoa(limit), nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response dto.UserSearchResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Len(t, response.Users, 2)

	mockController.AssertExpectations(t)
}

func TestUserHandler_SearchUsers_MissingQuery(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockUserController)
	handler := handlers.NewUserHandler(mockController)

	router := gin.New()
	router.GET("/users/search", handler.SearchUsers)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/users/search", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Contains(t, response["error"], "Query parameter 'q' is required")

	mockController.AssertNotCalled(t, "SearchUsers", mock.Anything, mock.Anything)
}

func TestUserHandler_SearchUsers_QueryTooShort(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockUserController)
	handler := handlers.NewUserHandler(mockController)

	router := gin.New()
	router.GET("/users/search", handler.SearchUsers)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/users/search?q=a", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Contains(t, response["error"], "Query must be at least 2 characters")

	mockController.AssertNotCalled(t, "SearchUsers", mock.Anything, mock.Anything)
}

func TestUserHandler_SearchUsers_ControllerError(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockUserController)
	handler := handlers.NewUserHandler(mockController)

	query := "test"
	limit := 10
	controllerError := errors.New("database error")

	mockController.On("SearchUsers", query, limit).Return(nil, controllerError)

	router := gin.New()
	router.GET("/users/search", handler.SearchUsers)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/users/search?q="+query+"&limit="+strconv.Itoa(limit), nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "database error", response["error"])

	mockController.AssertExpectations(t)
}

// Тесты для UserHandler.UpdateUserRole

func TestUserHandler_UpdateUserRole_Success(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockUserController)
	handler := handlers.NewUserHandler(mockController)

	userID := uuid.New()
	roleID := 2
	req := dto.UpdateUserRoleRequest{
		RoleID: roleID,
	}
	reqJSON, _ := json.Marshal(req)

	mockController.On("UpdateUserRole", userID, roleID).Return(nil)

	router := gin.New()
	router.PATCH("/users/:user_id/role", handler.UpdateUserRole)

	// Act
	w := httptest.NewRecorder()
	reqHTTP := httptest.NewRequest("PATCH", "/users/"+userID.String()+"/role", bytes.NewBuffer(reqJSON))
	reqHTTP.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, reqHTTP)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "User role updated successfully", response["message"])

	mockController.AssertExpectations(t)
}

func TestUserHandler_UpdateUserRole_InvalidUUID(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockUserController)
	handler := handlers.NewUserHandler(mockController)

	req := dto.UpdateUserRoleRequest{
		RoleID: 2,
	}
	reqJSON, _ := json.Marshal(req)

	router := gin.New()
	router.PATCH("/users/:user_id/role", handler.UpdateUserRole)

	// Act
	w := httptest.NewRecorder()
	reqHTTP := httptest.NewRequest("PATCH", "/users/invalid-uuid/role", bytes.NewBuffer(reqJSON))
	reqHTTP.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, reqHTTP)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "Invalid user ID", response["error"])

	mockController.AssertNotCalled(t, "UpdateUserRole", mock.Anything, mock.Anything)
}

func TestUserHandler_UpdateUserRole_InvalidJSON(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockUserController)
	handler := handlers.NewUserHandler(mockController)

	userID := uuid.New()
	invalidJSON := []byte("{invalid json}")

	router := gin.New()
	router.PATCH("/users/:user_id/role", handler.UpdateUserRole)

	// Act
	w := httptest.NewRecorder()
	reqHTTP := httptest.NewRequest("PATCH", "/users/"+userID.String()+"/role", bytes.NewBuffer(invalidJSON))
	reqHTTP.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, reqHTTP)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Contains(t, response["error"], "Invalid request payload")

	mockController.AssertNotCalled(t, "UpdateUserRole", mock.Anything, mock.Anything)
}

func TestUserHandler_UpdateUserRole_UserNotFound(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockUserController)
	handler := handlers.NewUserHandler(mockController)

	userID := uuid.New()
	roleID := 2
	req := dto.UpdateUserRoleRequest{
		RoleID: roleID,
	}
	reqJSON, _ := json.Marshal(req)

	mockController.On("UpdateUserRole", userID, roleID).Return(custom_errors.ErrInvalidCredentials)

	router := gin.New()
	router.PATCH("/users/:user_id/role", handler.UpdateUserRole)

	// Act
	w := httptest.NewRecorder()
	reqHTTP := httptest.NewRequest("PATCH", "/users/"+userID.String()+"/role", bytes.NewBuffer(reqJSON))
	reqHTTP.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, reqHTTP)

	// Assert
	assert.Equal(t, http.StatusNotFound, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "User not found", response["error"])

	mockController.AssertExpectations(t)
}

func TestUserHandler_UpdateUserRole_RoleNotFound(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockUserController)
	handler := handlers.NewUserHandler(mockController)

	userID := uuid.New()
	roleID := 999
	req := dto.UpdateUserRoleRequest{
		RoleID: roleID,
	}
	reqJSON, _ := json.Marshal(req)
	roleNotFoundError := custom_errors.NewRoleNotFoundError(roleID)

	mockController.On("UpdateUserRole", userID, roleID).Return(roleNotFoundError)

	router := gin.New()
	router.PATCH("/users/:user_id/role", handler.UpdateUserRole)

	// Act
	w := httptest.NewRecorder()
	reqHTTP := httptest.NewRequest("PATCH", "/users/"+userID.String()+"/role", bytes.NewBuffer(reqJSON))
	reqHTTP.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, reqHTTP)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Contains(t, response["error"], "does not exist")

	mockController.AssertExpectations(t)
}

// Вспомогательные функции
func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}
