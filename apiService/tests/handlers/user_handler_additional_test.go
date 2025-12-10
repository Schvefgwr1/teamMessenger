package handlers

import (
	"apiService/internal/dto"
	"apiService/internal/handlers"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Тесты для UserHandler.UpdateUserRole

func TestUserHandler_UpdateUserRole_Success(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockUserController)
	handler := handlers.NewUserHandler(mockController)

	userID := uuid.New()
	roleID := 2

	mockController.On("UpdateUserRole", userID, roleID).Return(nil)

	router := gin.New()
	router.PATCH("/users/:user_id/role", handler.UpdateUserRole)

	// Act
	reqBody := `{"role_id":2}`
	w := httptest.NewRecorder()
	req := httptest.NewRequest("PATCH", "/users/"+userID.String()+"/role", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "User role updated successfully", response["message"])

	mockController.AssertExpectations(t)
}

func TestUserHandler_UpdateUserRole_InvalidUserID(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockUserController)
	handler := handlers.NewUserHandler(mockController)

	router := gin.New()
	router.PATCH("/users/:user_id/role", handler.UpdateUserRole)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("PATCH", "/users/invalid/role", bytes.NewBufferString(`{"role_id":2}`))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	mockController.AssertNotCalled(t, "UpdateUserRole", mock.Anything, mock.Anything)
}

// Тесты для UserHandler.UpdateRolePermissions

func TestUserHandler_UpdateRolePermissions_Success(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockUserController)
	handler := handlers.NewUserHandler(mockController)

	roleID := 1
	permissionIDs := []int{1, 2, 3}

	mockController.On("UpdateRolePermissions", roleID, permissionIDs).Return(nil)

	router := gin.New()
	router.PATCH("/roles/:role_id/permissions", handler.UpdateRolePermissions)

	// Act
	reqBody := `{"permission_ids":[1,2,3]}`
	w := httptest.NewRecorder()
	req := httptest.NewRequest("PATCH", "/roles/1/permissions", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "Role permissions updated successfully", response["message"])

	mockController.AssertExpectations(t)
}

func TestUserHandler_UpdateRolePermissions_InvalidRoleID(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockUserController)
	handler := handlers.NewUserHandler(mockController)

	router := gin.New()
	router.PATCH("/roles/:role_id/permissions", handler.UpdateRolePermissions)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("PATCH", "/roles/invalid/permissions", bytes.NewBufferString(`{"permission_ids":[1]}`))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	mockController.AssertNotCalled(t, "UpdateRolePermissions", mock.Anything, mock.Anything)
}

// Тесты для UserHandler.DeleteRole

func TestUserHandler_DeleteRole_Success(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockUserController)
	handler := handlers.NewUserHandler(mockController)

	roleID := 1

	mockController.On("DeleteRole", roleID).Return(nil)

	router := gin.New()
	router.DELETE("/roles/:role_id", handler.DeleteRole)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("DELETE", "/roles/1", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "Role deleted successfully", response["message"])

	mockController.AssertExpectations(t)
}

func TestUserHandler_DeleteRole_InvalidRoleID(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockUserController)
	handler := handlers.NewUserHandler(mockController)

	router := gin.New()
	router.DELETE("/roles/:role_id", handler.DeleteRole)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("DELETE", "/roles/invalid", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	mockController.AssertNotCalled(t, "DeleteRole", mock.Anything)
}

// Тесты для UserHandler.GetUserBrief

func TestUserHandler_GetUserBrief_Success(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockUserController)
	handler := handlers.NewUserHandler(mockController)

	userID := uuid.New()
	chatID := uuid.New().String()
	requesterID := uuid.New()

	expectedBrief := &dto.UserBriefResponse{
		Username: "testuser",
		Email:    "test@example.com",
	}

	mockController.On("GetUserBrief", userID, chatID, requesterID).Return(expectedBrief, nil)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userID", requesterID)
		c.Next()
	})
	router.GET("/users/:user_id/brief", handler.GetUserBrief)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/users/"+userID.String()+"/brief?chatId="+chatID, nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response dto.UserBriefResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, expectedBrief.Username, response.Username)

	mockController.AssertExpectations(t)
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
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	mockController.AssertNotCalled(t, "GetUserBrief", mock.Anything, mock.Anything, mock.Anything)
}

func TestUserHandler_GetUserBrief_InvalidUserID(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockUserController)
	handler := handlers.NewUserHandler(mockController)

	router := gin.New()
	router.GET("/users/:user_id/brief", handler.GetUserBrief)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/users/invalid/brief?chatId="+uuid.New().String(), nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	mockController.AssertNotCalled(t, "GetUserBrief", mock.Anything, mock.Anything, mock.Anything)
}
