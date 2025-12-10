package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"userService/internal/handlers"
	"userService/internal/handlers/dto"
	"userService/internal/models"
)

// MockRoleController - мок для RoleController
type MockRoleController struct {
	mock.Mock
}

func (m *MockRoleController) GetRoles() ([]models.Role, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Role), args.Error(1)
}

func (m *MockRoleController) CreateRole(roleDTO *dto.CreateRole) error {
	args := m.Called(roleDTO)
	return args.Error(0)
}

func (m *MockRoleController) DeleteRole(roleID int) error {
	args := m.Called(roleID)
	return args.Error(0)
}

func (m *MockRoleController) UpdateRolePermissions(roleID int, permissionIDs []int) error {
	args := m.Called(roleID, permissionIDs)
	return args.Error(0)
}

// Тесты для RoleHandler.GetRoles

func TestRoleHandler_GetRoles_Success(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockRoleController)
	handler := handlers.NewRoleHandler(mockController)

	expectedRoles := []models.Role{
		{ID: intPtr(1), Name: "admin", Description: "Admin role"},
		{ID: intPtr(2), Name: "user", Description: "User role"},
	}

	mockController.On("GetRoles").Return(expectedRoles, nil)

	router := gin.New()
	router.GET("/roles", handler.GetRoles)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/roles", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response []models.Role
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Len(t, response, 2)
	assert.Equal(t, expectedRoles[0].Name, response[0].Name)

	mockController.AssertExpectations(t)
}

func TestRoleHandler_GetRoles_ControllerError(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockRoleController)
	handler := handlers.NewRoleHandler(mockController)

	controllerError := errors.New("database error")
	mockController.On("GetRoles").Return(nil, controllerError)

	router := gin.New()
	router.GET("/roles", handler.GetRoles)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/roles", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "database error", response["error"])

	mockController.AssertExpectations(t)
}

// Тесты для RoleHandler.CreateRole

func TestRoleHandler_CreateRole_Success(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockRoleController)
	handler := handlers.NewRoleHandler(mockController)

	roleDTO := dto.CreateRole{
		Name:          "moderator",
		Description:   "Moderator role",
		PermissionIds: []int{1, 2},
	}
	reqJSON, _ := json.Marshal(roleDTO)

	mockController.On("CreateRole", mock.AnythingOfType("*dto.CreateRole")).Return(nil)

	router := gin.New()
	router.POST("/roles", handler.CreateRole)

	// Act
	w := httptest.NewRecorder()
	reqHTTP := httptest.NewRequest("POST", "/roles", bytes.NewBuffer(reqJSON))
	reqHTTP.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, reqHTTP)

	// Assert
	assert.Equal(t, http.StatusCreated, w.Code)

	var response dto.CreateRole
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, roleDTO.Name, response.Name)

	mockController.AssertExpectations(t)
}

func TestRoleHandler_CreateRole_InvalidJSON(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockRoleController)
	handler := handlers.NewRoleHandler(mockController)

	invalidJSON := []byte("{invalid json}")

	router := gin.New()
	router.POST("/roles", handler.CreateRole)

	// Act
	w := httptest.NewRecorder()
	reqHTTP := httptest.NewRequest("POST", "/roles", bytes.NewBuffer(invalidJSON))
	reqHTTP.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, reqHTTP)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	mockController.AssertNotCalled(t, "CreateRole", mock.Anything)
}

func TestRoleHandler_CreateRole_ControllerError(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockRoleController)
	handler := handlers.NewRoleHandler(mockController)

	roleDTO := dto.CreateRole{
		Name:          "moderator",
		Description:   "Moderator role",
		PermissionIds: []int{1, 2},
	}
	reqJSON, _ := json.Marshal(roleDTO)
	controllerError := errors.New("database error")

	mockController.On("CreateRole", mock.AnythingOfType("*dto.CreateRole")).Return(controllerError)

	router := gin.New()
	router.POST("/roles", handler.CreateRole)

	// Act
	w := httptest.NewRecorder()
	reqHTTP := httptest.NewRequest("POST", "/roles", bytes.NewBuffer(reqJSON))
	reqHTTP.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, reqHTTP)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "database error", response["error"])

	mockController.AssertExpectations(t)
}

// Тесты для RoleHandler.DeleteRole

func TestRoleHandler_DeleteRole_Success(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockRoleController)
	handler := handlers.NewRoleHandler(mockController)

	roleID := 1

	mockController.On("DeleteRole", roleID).Return(nil)

	router := gin.New()
	router.DELETE("/roles/:role_id", handler.DeleteRole)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("DELETE", fmt.Sprintf("/roles/%d", roleID), nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "Role deleted successfully", response["message"])

	mockController.AssertExpectations(t)
}

func TestRoleHandler_DeleteRole_InvalidID(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockRoleController)
	handler := handlers.NewRoleHandler(mockController)

	router := gin.New()
	router.DELETE("/roles/:role_id", handler.DeleteRole)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("DELETE", "/roles/invalid", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "Invalid role ID", response["error"])

	mockController.AssertNotCalled(t, "DeleteRole", mock.Anything)
}

func TestRoleHandler_DeleteRole_NotFound(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockRoleController)
	handler := handlers.NewRoleHandler(mockController)

	roleID := 999

	mockController.On("DeleteRole", roleID).Return(gorm.ErrRecordNotFound)

	router := gin.New()
	router.DELETE("/roles/:role_id", handler.DeleteRole)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("DELETE", fmt.Sprintf("/roles/%d", roleID), nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "Role not found", response["error"])

	mockController.AssertExpectations(t)
}

// Тесты для RoleHandler.UpdateRolePermissions

func TestRoleHandler_UpdateRolePermissions_Success(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockRoleController)
	handler := handlers.NewRoleHandler(mockController)

	roleID := 1
	req := dto.UpdateRolePermissionsRequest{
		PermissionIds: []int{1, 2, 3},
	}
	reqJSON, _ := json.Marshal(req)

	mockController.On("UpdateRolePermissions", roleID, req.PermissionIds).Return(nil)

	router := gin.New()
	router.PATCH("/roles/:role_id/permissions", handler.UpdateRolePermissions)

	// Act
	w := httptest.NewRecorder()
	reqHTTP := httptest.NewRequest("PATCH", fmt.Sprintf("/roles/%d/permissions", roleID), bytes.NewBuffer(reqJSON))
	reqHTTP.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, reqHTTP)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "Role permissions updated successfully", response["message"])

	mockController.AssertExpectations(t)
}

func TestRoleHandler_UpdateRolePermissions_InvalidRoleID(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockRoleController)
	handler := handlers.NewRoleHandler(mockController)

	req := dto.UpdateRolePermissionsRequest{
		PermissionIds: []int{1, 2},
	}
	reqJSON, _ := json.Marshal(req)

	router := gin.New()
	router.PATCH("/roles/:role_id/permissions", handler.UpdateRolePermissions)

	// Act
	w := httptest.NewRecorder()
	reqHTTP := httptest.NewRequest("PATCH", "/roles/invalid/permissions", bytes.NewBuffer(reqJSON))
	reqHTTP.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, reqHTTP)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "Invalid role ID", response["error"])

	mockController.AssertNotCalled(t, "UpdateRolePermissions", mock.Anything, mock.Anything)
}

func TestRoleHandler_UpdateRolePermissions_InvalidJSON(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockRoleController)
	handler := handlers.NewRoleHandler(mockController)

	roleID := 1
	invalidJSON := []byte("{invalid json}")

	router := gin.New()
	router.PATCH("/roles/:role_id/permissions", handler.UpdateRolePermissions)

	// Act
	w := httptest.NewRecorder()
	reqHTTP := httptest.NewRequest("PATCH", fmt.Sprintf("/roles/%d/permissions", roleID), bytes.NewBuffer(invalidJSON))
	reqHTTP.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, reqHTTP)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	mockController.AssertNotCalled(t, "UpdateRolePermissions", mock.Anything, mock.Anything)
}

func TestRoleHandler_UpdateRolePermissions_RoleNotFound(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockRoleController)
	handler := handlers.NewRoleHandler(mockController)

	roleID := 999
	req := dto.UpdateRolePermissionsRequest{
		PermissionIds: []int{1, 2},
	}
	reqJSON, _ := json.Marshal(req)

	mockController.On("UpdateRolePermissions", roleID, req.PermissionIds).Return(gorm.ErrRecordNotFound)

	router := gin.New()
	router.PATCH("/roles/:role_id/permissions", handler.UpdateRolePermissions)

	// Act
	w := httptest.NewRecorder()
	reqHTTP := httptest.NewRequest("PATCH", fmt.Sprintf("/roles/%d/permissions", roleID), bytes.NewBuffer(reqJSON))
	reqHTTP.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, reqHTTP)

	// Assert
	assert.Equal(t, http.StatusNotFound, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "Role not found", response["error"])

	mockController.AssertExpectations(t)
}

// Вспомогательные функции уже определены в user_handler_test.go
