package handlers

import (
	"apiService/internal/dto"
	"apiService/internal/handlers"
	"bytes"
	au "common/contracts/api-user"
	uc "common/contracts/user-contracts"
	"encoding/json"
	"errors"
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

// Тесты для UserHandler.GetUser

func TestUserHandler_GetUser_Success(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockUserController)
	handler := handlers.NewUserHandler(mockController)

	userID := uuid.New()
	expectedUser := &au.GetUserResponse{
		User: &uc.User{
			ID:       userID,
			Username: "testuser",
			Email:    "test@example.com",
		},
	}

	mockController.On("GetUser", userID).Return(expectedUser, nil)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userID", userID)
		c.Next()
	})
	router.GET("/users/me", handler.GetUser)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/users/me", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response au.GetUserResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, expectedUser.User.Username, response.User.Username)

	mockController.AssertExpectations(t)
}

func TestUserHandler_GetUser_MissingUserID(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockUserController)
	handler := handlers.NewUserHandler(mockController)

	router := gin.New()
	router.GET("/users/me", handler.GetUser)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/users/me", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Contains(t, response["error"], "User ID not found")

	mockController.AssertNotCalled(t, "GetUser", mock.Anything)
}

func TestUserHandler_GetUser_InvalidUserIDType(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockUserController)
	handler := handlers.NewUserHandler(mockController)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userID", "invalid-type")
		c.Next()
	})
	router.GET("/users/me", handler.GetUser)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/users/me", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	mockController.AssertNotCalled(t, "GetUser", mock.Anything)
}

func TestUserHandler_GetUser_ControllerError(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockUserController)
	handler := handlers.NewUserHandler(mockController)

	userID := uuid.New()
	serviceError := errors.New("service error")

	mockController.On("GetUser", userID).Return(nil, serviceError)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userID", userID)
		c.Next()
	})
	router.GET("/users/me", handler.GetUser)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/users/me", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, serviceError.Error(), response["error"])

	mockController.AssertExpectations(t)
}

// Тесты для UserHandler.GetUserProfileByID

func TestUserHandler_GetUserProfileByID_Success(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockUserController)
	handler := handlers.NewUserHandler(mockController)

	userID := uuid.New()
	expectedUser := &au.GetUserResponse{
		User: &uc.User{
			ID:       userID,
			Username: "testuser",
		},
	}

	mockController.On("GetUserProfileByID", userID).Return(expectedUser, nil)

	router := gin.New()
	router.GET("/users/:user_id", handler.GetUserProfileByID)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/users/"+userID.String(), nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	mockController.AssertExpectations(t)
}

func TestUserHandler_GetUserProfileByID_InvalidUUID(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockUserController)
	handler := handlers.NewUserHandler(mockController)

	router := gin.New()
	router.GET("/users/:user_id", handler.GetUserProfileByID)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/users/invalid-uuid", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	mockController.AssertNotCalled(t, "GetUserProfileByID", mock.Anything)
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
			{Username: "testuser", Email: "test@example.com"},
		},
	}

	mockController.On("SearchUsers", query, limit).Return(expectedResponse, nil)

	router := gin.New()
	router.GET("/users/search", handler.SearchUsers)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/users/search?q="+query+"&limit=10", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response dto.UserSearchResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Len(t, response.Users, 1)

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

	mockController.AssertNotCalled(t, "SearchUsers", mock.Anything, mock.Anything)
}

// Тесты для UserHandler.GetAllPermissions

func TestUserHandler_GetAllPermissions_Success(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockUserController)
	handler := handlers.NewUserHandler(mockController)

	expectedPermissions := []*uc.Permission{
		{ID: 1, Name: "read"},
		{ID: 2, Name: "write"},
	}

	mockController.On("GetAllPermissions").Return(expectedPermissions, nil)

	router := gin.New()
	router.GET("/permissions", handler.GetAllPermissions)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/permissions", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response []*uc.Permission
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Len(t, response, 2)

	mockController.AssertExpectations(t)
}

// Тесты для UserHandler.GetAllRoles

func TestUserHandler_GetAllRoles_Success(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockUserController)
	handler := handlers.NewUserHandler(mockController)

	expectedRoles := []*uc.Role{
		{ID: 1, Name: "admin"},
	}

	mockController.On("GetAllRoles").Return(expectedRoles, nil)

	router := gin.New()
	router.GET("/roles", handler.GetAllRoles)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/roles", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	mockController.AssertExpectations(t)
}

// Тесты для UserHandler.CreateRole

func TestUserHandler_CreateRole_Success(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockUserController)
	handler := handlers.NewUserHandler(mockController)

	expectedRole := &uc.Role{
		ID:   1,
		Name: "newrole",
	}

	mockController.On("CreateRole", mock.Anything).Return(expectedRole, nil)

	router := gin.New()
	router.POST("/roles", handler.CreateRole)

	// Act
	reqBody := `{"name":"newrole","permissionIds":[1,2]}`
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/roles", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Assert
	// Проверяем, что мок был вызван
	mockController.AssertExpectations(t)
}

// Тесты для UserHandler.UpdateUser

func TestUserHandler_UpdateUser_Success(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockUserController)
	handler := handlers.NewUserHandler(mockController)

	userID := uuid.New()
	expectedResponse := &au.UpdateUserResponse{
		Message: stringPtr("User updated successfully"),
	}

	mockController.On("UpdateUser", userID, mock.Anything, mock.Anything).Return(expectedResponse, nil)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userID", userID)
		c.Next()
	})
	router.PATCH("/users/me", handler.UpdateUser)

	// Act - создаем multipart форму
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	jsonData := `{"username":"updateduser","description":"Updated description"}`
	writer.WriteField("data", jsonData)

	// Добавляем файл
	part, err := writer.CreateFormFile("file", "avatar.jpg")
	require.NoError(t, err)
	_, err = part.Write([]byte("fake image content"))
	require.NoError(t, err)

	writer.Close()

	w := httptest.NewRecorder()
	req := httptest.NewRequest("PATCH", "/users/me", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusCreated, w.Code)

	mockController.AssertExpectations(t)
}

func TestUserHandler_UpdateUser_WithoutFile(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockUserController)
	handler := handlers.NewUserHandler(mockController)

	userID := uuid.New()
	expectedResponse := &au.UpdateUserResponse{
		Message: stringPtr("User updated successfully"),
	}

	mockController.On("UpdateUser", userID, mock.Anything, mock.Anything).Return(expectedResponse, nil)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userID", userID)
		c.Next()
	})
	router.PATCH("/users/me", handler.UpdateUser)

	// Act - создаем multipart форму без файла
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	jsonData := `{"username":"updateduser"}`
	writer.WriteField("data", jsonData)
	writer.Close()

	w := httptest.NewRecorder()
	req := httptest.NewRequest("PATCH", "/users/me", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusCreated, w.Code)

	mockController.AssertExpectations(t)
}

func TestUserHandler_UpdateUser_InvalidJSON(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockUserController)
	handler := handlers.NewUserHandler(mockController)

	userID := uuid.New()

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userID", userID)
		c.Next()
	})
	router.PATCH("/users/me", handler.UpdateUser)

	// Act
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.WriteField("data", "invalid json")
	writer.Close()

	w := httptest.NewRecorder()
	req := httptest.NewRequest("PATCH", "/users/me", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	mockController.AssertNotCalled(t, "UpdateUser", mock.Anything, mock.Anything, mock.Anything)
}

func TestUserHandler_UpdateUser_MissingUserID(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockUserController)
	handler := handlers.NewUserHandler(mockController)

	router := gin.New()
	router.PATCH("/users/me", handler.UpdateUser)

	// Act
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	jsonData := `{"username":"updateduser"}`
	writer.WriteField("data", jsonData)
	writer.Close()

	w := httptest.NewRecorder()
	req := httptest.NewRequest("PATCH", "/users/me", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	mockController.AssertNotCalled(t, "UpdateUser", mock.Anything, mock.Anything, mock.Anything)
}

// Вспомогательная функция для создания указателя на строку
func stringPtr(s string) *string {
	return &s
}
