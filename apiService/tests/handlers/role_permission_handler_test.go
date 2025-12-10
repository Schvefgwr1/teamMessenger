package handlers

import (
	"apiService/internal/dto"
	"apiService/internal/handlers"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Тесты для ChatRolePermissionHandler.GetAllRoles

func TestChatRolePermissionHandler_GetAllRoles_Success(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockChatRolePermissionController)
	handler := handlers.NewRolePermissionHandler(mockController)

	expectedRoles := []dto.ChatRoleResponseGateway{
		{ID: 1, Name: "Admin", Permissions: []dto.ChatPermissionResponseGateway{{ID: 1, Name: "read"}}},
	}

	mockController.On("GetAllRoles").Return(expectedRoles, nil)

	router := gin.New()
	router.GET("/chat-roles", handler.GetAllRoles)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/chat-roles", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response []dto.ChatRoleResponseGateway
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Len(t, response, 1)

	mockController.AssertExpectations(t)
}

func TestChatRolePermissionHandler_GetAllRoles_ControllerError(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockChatRolePermissionController)
	handler := handlers.NewRolePermissionHandler(mockController)

	serviceError := errors.New("service error")

	mockController.On("GetAllRoles").Return(nil, serviceError)

	router := gin.New()
	router.GET("/chat-roles", handler.GetAllRoles)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/chat-roles", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, serviceError.Error(), response["error"])

	mockController.AssertExpectations(t)
}

// Тесты для ChatRolePermissionHandler.GetRoleByID

func TestChatRolePermissionHandler_GetRoleByID_Success(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockChatRolePermissionController)
	handler := handlers.NewRolePermissionHandler(mockController)

	roleID := 1
	expectedRole := &dto.ChatRoleResponseGateway{
		ID:   roleID,
		Name: "Admin",
	}

	mockController.On("GetRoleByID", roleID).Return(expectedRole, nil)

	router := gin.New()
	router.GET("/chat-roles/:role_id", handler.GetRoleByID)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/chat-roles/1", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response dto.ChatRoleResponseGateway
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, expectedRole.Name, response.Name)

	mockController.AssertExpectations(t)
}

func TestChatRolePermissionHandler_GetRoleByID_InvalidID(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockChatRolePermissionController)
	handler := handlers.NewRolePermissionHandler(mockController)

	router := gin.New()
	router.GET("/chat-roles/:role_id", handler.GetRoleByID)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/chat-roles/invalid", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	mockController.AssertNotCalled(t, "GetRoleByID", mock.Anything)
}

func TestChatRolePermissionHandler_GetRoleByID_NotFound(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockChatRolePermissionController)
	handler := handlers.NewRolePermissionHandler(mockController)

	roleID := 999
	serviceError := errors.New("role not found")

	mockController.On("GetRoleByID", roleID).Return(nil, serviceError)

	router := gin.New()
	router.GET("/chat-roles/:role_id", handler.GetRoleByID)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/chat-roles/999", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, w.Code)

	mockController.AssertExpectations(t)
}

// Тесты для ChatRolePermissionHandler.GetAllPermissions

func TestChatRolePermissionHandler_GetAllPermissions_Success(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockChatRolePermissionController)
	handler := handlers.NewRolePermissionHandler(mockController)

	expectedPermissions := []dto.ChatPermissionResponseGateway{
		{ID: 1, Name: "read"},
		{ID: 2, Name: "write"},
	}

	mockController.On("GetAllPermissions").Return(expectedPermissions, nil)

	router := gin.New()
	router.GET("/chat-permissions", handler.GetAllPermissions)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/chat-permissions", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response []dto.ChatPermissionResponseGateway
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Len(t, response, 2)

	mockController.AssertExpectations(t)
}

// Тесты для ChatRolePermissionHandler.DeleteRole

func TestChatRolePermissionHandler_DeleteRole_Success(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockChatRolePermissionController)
	handler := handlers.NewRolePermissionHandler(mockController)

	roleID := 1

	mockController.On("DeleteRole", roleID).Return(nil)

	router := gin.New()
	router.DELETE("/chat-roles/:role_id", handler.DeleteRole)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("DELETE", "/chat-roles/1", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusNoContent, w.Code)

	mockController.AssertExpectations(t)
}

func TestChatRolePermissionHandler_DeleteRole_InvalidID(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockChatRolePermissionController)
	handler := handlers.NewRolePermissionHandler(mockController)

	router := gin.New()
	router.DELETE("/chat-roles/:role_id", handler.DeleteRole)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("DELETE", "/chat-roles/invalid", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	mockController.AssertNotCalled(t, "DeleteRole", mock.Anything)
}

// Тесты для ChatRolePermissionHandler.DeletePermission

func TestChatRolePermissionHandler_DeletePermission_Success(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockChatRolePermissionController)
	handler := handlers.NewRolePermissionHandler(mockController)

	permissionID := 1

	mockController.On("DeletePermission", permissionID).Return(nil)

	router := gin.New()
	router.DELETE("/chat-permissions/:permission_id", handler.DeletePermission)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("DELETE", "/chat-permissions/1", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusNoContent, w.Code)

	mockController.AssertExpectations(t)
}

// Тесты для ChatRolePermissionHandler.CreateRole

func TestChatRolePermissionHandler_CreateRole_Success(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockChatRolePermissionController)
	handler := handlers.NewRolePermissionHandler(mockController)

	expectedRole := &dto.ChatRoleResponseGateway{
		ID:   1,
		Name: "New Role",
	}

	mockController.On("CreateRole", mock.Anything).Return(expectedRole, nil)

	router := gin.New()
	router.POST("/chat-roles", handler.CreateRole)

	// Act
	reqBody := `{"name":"New Role","permissionIds":[1,2]}`
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/chat-roles", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusCreated, w.Code)

	var response dto.ChatRoleResponseGateway
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, expectedRole.Name, response.Name)

	mockController.AssertExpectations(t)
}

func TestChatRolePermissionHandler_CreateRole_InvalidJSON(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockChatRolePermissionController)
	handler := handlers.NewRolePermissionHandler(mockController)

	router := gin.New()
	router.POST("/chat-roles", handler.CreateRole)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/chat-roles", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	mockController.AssertNotCalled(t, "CreateRole", mock.Anything)
}

// Тесты для ChatRolePermissionHandler.UpdateRolePermissions

func TestChatRolePermissionHandler_UpdateRolePermissions_Success(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockChatRolePermissionController)
	handler := handlers.NewRolePermissionHandler(mockController)

	roleID := 1
	expectedRole := &dto.ChatRoleResponseGateway{
		ID:   roleID,
		Name: "Updated Role",
	}

	mockController.On("UpdateRolePermissions", roleID, mock.Anything).Return(expectedRole, nil)

	router := gin.New()
	router.PATCH("/chat-roles/:role_id/permissions", handler.UpdateRolePermissions)

	// Act
	reqBody := `{"permissionIds":[1,2,3]}`
	w := httptest.NewRecorder()
	req := httptest.NewRequest("PATCH", "/chat-roles/1/permissions", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response dto.ChatRoleResponseGateway
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, expectedRole.Name, response.Name)

	mockController.AssertExpectations(t)
}

func TestChatRolePermissionHandler_UpdateRolePermissions_InvalidRoleID(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockChatRolePermissionController)
	handler := handlers.NewRolePermissionHandler(mockController)

	router := gin.New()
	router.PATCH("/chat-roles/:role_id/permissions", handler.UpdateRolePermissions)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("PATCH", "/chat-roles/invalid/permissions", bytes.NewBufferString(`{"permissionIds":[1]}`))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	mockController.AssertNotCalled(t, "UpdateRolePermissions", mock.Anything, mock.Anything)
}

// Тесты для ChatRolePermissionHandler.CreatePermission

func TestChatRolePermissionHandler_CreatePermission_Success(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockChatRolePermissionController)
	handler := handlers.NewRolePermissionHandler(mockController)

	expectedPermission := &dto.ChatPermissionResponseGateway{
		ID:   1,
		Name: "New Permission",
	}

	mockController.On("CreatePermission", mock.Anything).Return(expectedPermission, nil)

	router := gin.New()
	router.POST("/chat-permissions", handler.CreatePermission)

	// Act
	reqBody := `{"name":"New Permission"}`
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/chat-permissions", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusCreated, w.Code)

	var response dto.ChatPermissionResponseGateway
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, expectedPermission.Name, response.Name)

	mockController.AssertExpectations(t)
}

func TestChatRolePermissionHandler_CreatePermission_InvalidJSON(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockChatRolePermissionController)
	handler := handlers.NewRolePermissionHandler(mockController)

	router := gin.New()
	router.POST("/chat-permissions", handler.CreatePermission)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/chat-permissions", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	mockController.AssertNotCalled(t, "CreatePermission", mock.Anything)
}
