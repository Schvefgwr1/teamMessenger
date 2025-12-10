package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"chatService/internal/controllers"
	"chatService/internal/handlers"
	"chatService/internal/handlers/dto"
	"chatService/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Моки для репозиториев
type MockChatRoleRepository struct {
	mock.Mock
}

func (m *MockChatRoleRepository) GetRoleByID(roleID int) (*models.ChatRole, error) {
	args := m.Called(roleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ChatRole), args.Error(1)
}

func (m *MockChatRoleRepository) GetRoleByName(roleName string) (*models.ChatRole, error) {
	args := m.Called(roleName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ChatRole), args.Error(1)
}

func (m *MockChatRoleRepository) GetAllRoles() ([]models.ChatRole, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.ChatRole), args.Error(1)
}

func (m *MockChatRoleRepository) CreateRole(role *models.ChatRole, permissionIDs []int) error {
	args := m.Called(role, permissionIDs)
	return args.Error(0)
}

func (m *MockChatRoleRepository) DeleteRole(roleID int) error {
	args := m.Called(roleID)
	return args.Error(0)
}

func (m *MockChatRoleRepository) UpdateRolePermissions(roleID int, permissionIDs []int) error {
	args := m.Called(roleID, permissionIDs)
	return args.Error(0)
}

type MockChatPermissionRepository struct {
	mock.Mock
}

func (m *MockChatPermissionRepository) GetAllPermissions() ([]models.ChatPermission, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.ChatPermission), args.Error(1)
}

func (m *MockChatPermissionRepository) GetPermissionByID(id int) (*models.ChatPermission, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ChatPermission), args.Error(1)
}

func (m *MockChatPermissionRepository) GetPermissionByName(name string) (*models.ChatPermission, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ChatPermission), args.Error(1)
}

func (m *MockChatPermissionRepository) CreatePermission(permission *models.ChatPermission) error {
	args := m.Called(permission)
	return args.Error(0)
}

func (m *MockChatPermissionRepository) DeletePermission(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

// Вспомогательные функции
func createTestChatPermissionWithID(id int, name string) *models.ChatPermission {
	return &models.ChatPermission{
		ID:   id,
		Name: name,
	}
}

func createTestChatRoleWithPermissionsAndID(roleID int, roleName string, permissions []models.ChatPermission) *models.ChatRole {
	return &models.ChatRole{
		ID:          roleID,
		Name:        roleName,
		Permissions: permissions,
	}
}

// ==================== Тесты для Roles ====================

func TestRolePermissionHandler_GetAllRoles_Success(t *testing.T) {
	// Arrange
	mockRoleRepo := new(MockChatRoleRepository)
	mockPermissionRepo := new(MockChatPermissionRepository)

	roles := []models.ChatRole{
		*createTestChatRoleWithPermissionsAndID(1, "owner", []models.ChatPermission{
			*createTestChatPermissionWithID(1, "send_message"),
		}),
		*createTestChatRoleWithPermissionsAndID(2, "main", []models.ChatPermission{}),
	}

	mockRoleRepo.On("GetAllRoles").Return(roles, nil)

	controller := controllers.NewRolePermissionController(mockRoleRepo, mockPermissionRepo)
	handler := handlers.NewRolePermissionHandler(controller)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/chat-roles", handler.GetAllRoles)

	req, _ := http.NewRequest("GET", "/chat-roles", nil)
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	var result []dto.RoleResponse
	err := json.Unmarshal(w.Body.Bytes(), &result)
	require.NoError(t, err)
	assert.Len(t, result, 2)

	mockRoleRepo.AssertExpectations(t)
}

func TestRolePermissionHandler_GetAllRoles_ControllerError(t *testing.T) {
	// Arrange
	mockRoleRepo := new(MockChatRoleRepository)
	mockPermissionRepo := new(MockChatPermissionRepository)

	mockRoleRepo.On("GetAllRoles").Return(nil, errors.New("database error"))

	controller := controllers.NewRolePermissionController(mockRoleRepo, mockPermissionRepo)
	handler := handlers.NewRolePermissionHandler(controller)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/chat-roles", handler.GetAllRoles)

	req, _ := http.NewRequest("GET", "/chat-roles", nil)
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	mockRoleRepo.AssertExpectations(t)
}

func TestRolePermissionHandler_GetAllRoles_EmptyList(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockRoleRepo := new(MockChatRoleRepository)
	mockPermissionRepo := new(MockChatPermissionRepository)

	emptyRoles := []models.ChatRole{}
	mockRoleRepo.On("GetAllRoles").Return(emptyRoles, nil)

	controller := controllers.NewRolePermissionController(mockRoleRepo, mockPermissionRepo)
	handler := handlers.NewRolePermissionHandler(controller)

	router := gin.New()
	router.GET("/chat-roles", handler.GetAllRoles)

	req, _ := http.NewRequest("GET", "/chat-roles", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var result []dto.RoleResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &result))
	assert.Len(t, result, 0)
	mockRoleRepo.AssertExpectations(t)
}

func TestRolePermissionHandler_GetAllPermissions_EmptyList(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockRoleRepo := new(MockChatRoleRepository)
	mockPermissionRepo := new(MockChatPermissionRepository)

	emptyPermissions := []models.ChatPermission{}
	mockPermissionRepo.On("GetAllPermissions").Return(emptyPermissions, nil)

	controller := controllers.NewRolePermissionController(mockRoleRepo, mockPermissionRepo)
	handler := handlers.NewRolePermissionHandler(controller)

	router := gin.New()
	router.GET("/chat-permissions", handler.GetAllPermissions)

	req, _ := http.NewRequest("GET", "/chat-permissions", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var result []dto.PermissionResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &result))
	assert.Len(t, result, 0)
	mockPermissionRepo.AssertExpectations(t)
}

func TestRolePermissionHandler_GetRoleByID_Success(t *testing.T) {
	// Arrange
	mockRoleRepo := new(MockChatRoleRepository)
	mockPermissionRepo := new(MockChatPermissionRepository)

	roleID := 1
	role := createTestChatRoleWithPermissionsAndID(roleID, "owner", []models.ChatPermission{
		*createTestChatPermissionWithID(1, "send_message"),
	})

	mockRoleRepo.On("GetRoleByID", roleID).Return(role, nil)

	controller := controllers.NewRolePermissionController(mockRoleRepo, mockPermissionRepo)
	handler := handlers.NewRolePermissionHandler(controller)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/chat-roles/:role_id", handler.GetRoleByID)

	req, _ := http.NewRequest("GET", "/chat-roles/1", nil)
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	var result dto.RoleResponse
	err := json.Unmarshal(w.Body.Bytes(), &result)
	require.NoError(t, err)
	assert.Equal(t, roleID, result.ID)
	assert.Equal(t, "owner", result.Name)

	mockRoleRepo.AssertExpectations(t)
}

func TestRolePermissionHandler_GetRoleByID_InvalidID(t *testing.T) {
	// Arrange
	mockRoleRepo := new(MockChatRoleRepository)
	mockPermissionRepo := new(MockChatPermissionRepository)

	controller := controllers.NewRolePermissionController(mockRoleRepo, mockPermissionRepo)
	handler := handlers.NewRolePermissionHandler(controller)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/chat-roles/:role_id", handler.GetRoleByID)

	req, _ := http.NewRequest("GET", "/chat-roles/invalid", nil)
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRolePermissionHandler_GetRoleByID_NotFound(t *testing.T) {
	// Arrange
	mockRoleRepo := new(MockChatRoleRepository)
	mockPermissionRepo := new(MockChatPermissionRepository)

	roleID := 999
	mockRoleRepo.On("GetRoleByID", roleID).Return(nil, errors.New("role not found"))

	controller := controllers.NewRolePermissionController(mockRoleRepo, mockPermissionRepo)
	handler := handlers.NewRolePermissionHandler(controller)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/chat-roles/:role_id", handler.GetRoleByID)

	req, _ := http.NewRequest("GET", "/chat-roles/999", nil)
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, w.Code)

	mockRoleRepo.AssertExpectations(t)
}

func TestRolePermissionHandler_CreateRole_Success(t *testing.T) {
	// Arrange
	mockRoleRepo := new(MockChatRoleRepository)
	mockPermissionRepo := new(MockChatPermissionRepository)

	createReq := dto.CreateRoleRequest{
		Name:          "admin",
		PermissionIDs: []int{1, 2},
	}

	mockRoleRepo.On("CreateRole", mock.Anything, createReq.PermissionIDs).Return(nil).Run(func(args mock.Arguments) {
		role := args.Get(0).(*models.ChatRole)
		role.ID = 3
		role.Permissions = []models.ChatPermission{
			*createTestChatPermissionWithID(1, "send_message"),
			*createTestChatPermissionWithID(2, "delete_message"),
		}
	})

	controller := controllers.NewRolePermissionController(mockRoleRepo, mockPermissionRepo)
	handler := handlers.NewRolePermissionHandler(controller)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/chat-roles", handler.CreateRole)

	body, _ := json.Marshal(createReq)
	req, _ := http.NewRequest("POST", "/chat-roles", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusCreated, w.Code)
	var result dto.RoleResponse
	err := json.Unmarshal(w.Body.Bytes(), &result)
	require.NoError(t, err)
	assert.Equal(t, "admin", result.Name)

	mockRoleRepo.AssertExpectations(t)
}

func TestRolePermissionHandler_CreateRole_InvalidJSON(t *testing.T) {
	// Arrange
	mockRoleRepo := new(MockChatRoleRepository)
	mockPermissionRepo := new(MockChatPermissionRepository)

	controller := controllers.NewRolePermissionController(mockRoleRepo, mockPermissionRepo)
	handler := handlers.NewRolePermissionHandler(controller)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/chat-roles", handler.CreateRole)

	req, _ := http.NewRequest("POST", "/chat-roles", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRolePermissionHandler_CreateRole_ControllerError(t *testing.T) {
	// Arrange
	mockRoleRepo := new(MockChatRoleRepository)
	mockPermissionRepo := new(MockChatPermissionRepository)

	createReq := dto.CreateRoleRequest{
		Name:          "admin",
		PermissionIDs: []int{1, 2},
	}

	mockRoleRepo.On("CreateRole", mock.Anything, createReq.PermissionIDs).Return(errors.New("database error"))

	controller := controllers.NewRolePermissionController(mockRoleRepo, mockPermissionRepo)
	handler := handlers.NewRolePermissionHandler(controller)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/chat-roles", handler.CreateRole)

	body, _ := json.Marshal(createReq)
	req, _ := http.NewRequest("POST", "/chat-roles", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	mockRoleRepo.AssertExpectations(t)
}

func TestRolePermissionHandler_DeleteRole_Success(t *testing.T) {
	// Arrange
	mockRoleRepo := new(MockChatRoleRepository)
	mockPermissionRepo := new(MockChatPermissionRepository)

	roleID := 1
	mockRoleRepo.On("DeleteRole", roleID).Return(nil)

	controller := controllers.NewRolePermissionController(mockRoleRepo, mockPermissionRepo)
	handler := handlers.NewRolePermissionHandler(controller)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.DELETE("/chat-roles/:role_id", handler.DeleteRole)

	req, _ := http.NewRequest("DELETE", "/chat-roles/1", nil)
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusNoContent, w.Code)

	mockRoleRepo.AssertExpectations(t)
}

func TestRolePermissionHandler_DeleteRole_InvalidID(t *testing.T) {
	// Arrange
	mockRoleRepo := new(MockChatRoleRepository)
	mockPermissionRepo := new(MockChatPermissionRepository)

	controller := controllers.NewRolePermissionController(mockRoleRepo, mockPermissionRepo)
	handler := handlers.NewRolePermissionHandler(controller)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.DELETE("/chat-roles/:role_id", handler.DeleteRole)

	req, _ := http.NewRequest("DELETE", "/chat-roles/invalid", nil)
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRolePermissionHandler_DeleteRole_ControllerError(t *testing.T) {
	// Arrange
	mockRoleRepo := new(MockChatRoleRepository)
	mockPermissionRepo := new(MockChatPermissionRepository)

	roleID := 999
	mockRoleRepo.On("DeleteRole", roleID).Return(errors.New("role not found"))

	controller := controllers.NewRolePermissionController(mockRoleRepo, mockPermissionRepo)
	handler := handlers.NewRolePermissionHandler(controller)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.DELETE("/chat-roles/:role_id", handler.DeleteRole)

	req, _ := http.NewRequest("DELETE", "/chat-roles/999", nil)
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	mockRoleRepo.AssertExpectations(t)
}

func TestRolePermissionHandler_UpdateRolePermissions_Success(t *testing.T) {
	// Arrange
	mockRoleRepo := new(MockChatRoleRepository)
	mockPermissionRepo := new(MockChatPermissionRepository)

	roleID := 1
	updateReq := dto.UpdateRolePermissionsRequest{
		PermissionIDs: []int{1, 3},
	}

	updatedRole := createTestChatRoleWithPermissionsAndID(roleID, "admin", []models.ChatPermission{
		*createTestChatPermissionWithID(1, "send_message"),
		*createTestChatPermissionWithID(3, "manage_users"),
	})

	mockRoleRepo.On("UpdateRolePermissions", roleID, updateReq.PermissionIDs).Return(nil)
	mockRoleRepo.On("GetRoleByID", roleID).Return(updatedRole, nil)

	controller := controllers.NewRolePermissionController(mockRoleRepo, mockPermissionRepo)
	handler := handlers.NewRolePermissionHandler(controller)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.PATCH("/chat-roles/:role_id/permissions", handler.UpdateRolePermissions)

	body, _ := json.Marshal(updateReq)
	req, _ := http.NewRequest("PATCH", "/chat-roles/1/permissions", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	var result dto.RoleResponse
	err := json.Unmarshal(w.Body.Bytes(), &result)
	require.NoError(t, err)
	assert.Equal(t, roleID, result.ID)
	assert.Len(t, result.Permissions, 2)

	mockRoleRepo.AssertExpectations(t)
}

func TestRolePermissionHandler_UpdateRolePermissions_InvalidID(t *testing.T) {
	// Arrange
	mockRoleRepo := new(MockChatRoleRepository)
	mockPermissionRepo := new(MockChatPermissionRepository)

	controller := controllers.NewRolePermissionController(mockRoleRepo, mockPermissionRepo)
	handler := handlers.NewRolePermissionHandler(controller)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.PATCH("/chat-roles/:role_id/permissions", handler.UpdateRolePermissions)

	updateReq := dto.UpdateRolePermissionsRequest{PermissionIDs: []int{1, 2}}
	body, _ := json.Marshal(updateReq)
	req, _ := http.NewRequest("PATCH", "/chat-roles/invalid/permissions", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRolePermissionHandler_UpdateRolePermissions_InvalidJSON(t *testing.T) {
	// Arrange
	mockRoleRepo := new(MockChatRoleRepository)
	mockPermissionRepo := new(MockChatPermissionRepository)

	controller := controllers.NewRolePermissionController(mockRoleRepo, mockPermissionRepo)
	handler := handlers.NewRolePermissionHandler(controller)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.PATCH("/chat-roles/:role_id/permissions", handler.UpdateRolePermissions)

	req, _ := http.NewRequest("PATCH", "/chat-roles/1/permissions", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRolePermissionHandler_UpdateRolePermissions_ControllerError(t *testing.T) {
	// Arrange
	mockRoleRepo := new(MockChatRoleRepository)
	mockPermissionRepo := new(MockChatPermissionRepository)

	roleID := 999
	updateReq := dto.UpdateRolePermissionsRequest{
		PermissionIDs: []int{1, 2},
	}

	mockRoleRepo.On("UpdateRolePermissions", roleID, updateReq.PermissionIDs).Return(errors.New("role not found"))

	controller := controllers.NewRolePermissionController(mockRoleRepo, mockPermissionRepo)
	handler := handlers.NewRolePermissionHandler(controller)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.PATCH("/chat-roles/:role_id/permissions", handler.UpdateRolePermissions)

	body, _ := json.Marshal(updateReq)
	req, _ := http.NewRequest("PATCH", "/chat-roles/999/permissions", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	mockRoleRepo.AssertExpectations(t)
}

// ==================== Тесты для Permissions ====================

func TestRolePermissionHandler_GetAllPermissions_Success(t *testing.T) {
	// Arrange
	mockRoleRepo := new(MockChatRoleRepository)
	mockPermissionRepo := new(MockChatPermissionRepository)

	permissions := []models.ChatPermission{
		*createTestChatPermissionWithID(1, "send_message"),
		*createTestChatPermissionWithID(2, "delete_message"),
	}

	mockPermissionRepo.On("GetAllPermissions").Return(permissions, nil)

	controller := controllers.NewRolePermissionController(mockRoleRepo, mockPermissionRepo)
	handler := handlers.NewRolePermissionHandler(controller)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/chat-permissions", handler.GetAllPermissions)

	req, _ := http.NewRequest("GET", "/chat-permissions", nil)
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	var result []dto.PermissionResponse
	err := json.Unmarshal(w.Body.Bytes(), &result)
	require.NoError(t, err)
	assert.Len(t, result, 2)

	mockPermissionRepo.AssertExpectations(t)
}

func TestRolePermissionHandler_GetAllPermissions_ControllerError(t *testing.T) {
	// Arrange
	mockRoleRepo := new(MockChatRoleRepository)
	mockPermissionRepo := new(MockChatPermissionRepository)

	mockPermissionRepo.On("GetAllPermissions").Return(nil, errors.New("database error"))

	controller := controllers.NewRolePermissionController(mockRoleRepo, mockPermissionRepo)
	handler := handlers.NewRolePermissionHandler(controller)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/chat-permissions", handler.GetAllPermissions)

	req, _ := http.NewRequest("GET", "/chat-permissions", nil)
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	mockPermissionRepo.AssertExpectations(t)
}

func TestRolePermissionHandler_CreatePermission_Success(t *testing.T) {
	// Arrange
	mockRoleRepo := new(MockChatRoleRepository)
	mockPermissionRepo := new(MockChatPermissionRepository)

	createReq := dto.CreatePermissionRequest{
		Name: "edit_message",
	}

	mockPermissionRepo.On("CreatePermission", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		permission := args.Get(0).(*models.ChatPermission)
		permission.ID = 4
	})

	controller := controllers.NewRolePermissionController(mockRoleRepo, mockPermissionRepo)
	handler := handlers.NewRolePermissionHandler(controller)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/chat-permissions", handler.CreatePermission)

	body, _ := json.Marshal(createReq)
	req, _ := http.NewRequest("POST", "/chat-permissions", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusCreated, w.Code)
	var result dto.PermissionResponse
	err := json.Unmarshal(w.Body.Bytes(), &result)
	require.NoError(t, err)
	assert.Equal(t, "edit_message", result.Name)

	mockPermissionRepo.AssertExpectations(t)
}

func TestRolePermissionHandler_CreatePermission_InvalidJSON(t *testing.T) {
	// Arrange
	mockRoleRepo := new(MockChatRoleRepository)
	mockPermissionRepo := new(MockChatPermissionRepository)

	controller := controllers.NewRolePermissionController(mockRoleRepo, mockPermissionRepo)
	handler := handlers.NewRolePermissionHandler(controller)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/chat-permissions", handler.CreatePermission)

	req, _ := http.NewRequest("POST", "/chat-permissions", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRolePermissionHandler_CreatePermission_ControllerError(t *testing.T) {
	// Arrange
	mockRoleRepo := new(MockChatRoleRepository)
	mockPermissionRepo := new(MockChatPermissionRepository)

	createReq := dto.CreatePermissionRequest{
		Name: "edit_message",
	}

	mockPermissionRepo.On("CreatePermission", mock.Anything).Return(errors.New("permission already exists"))

	controller := controllers.NewRolePermissionController(mockRoleRepo, mockPermissionRepo)
	handler := handlers.NewRolePermissionHandler(controller)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/chat-permissions", handler.CreatePermission)

	body, _ := json.Marshal(createReq)
	req, _ := http.NewRequest("POST", "/chat-permissions", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	mockPermissionRepo.AssertExpectations(t)
}

func TestRolePermissionHandler_DeletePermission_Success(t *testing.T) {
	// Arrange
	mockRoleRepo := new(MockChatRoleRepository)
	mockPermissionRepo := new(MockChatPermissionRepository)

	permissionID := 1
	mockPermissionRepo.On("DeletePermission", permissionID).Return(nil)

	controller := controllers.NewRolePermissionController(mockRoleRepo, mockPermissionRepo)
	handler := handlers.NewRolePermissionHandler(controller)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.DELETE("/chat-permissions/:permission_id", handler.DeletePermission)

	req, _ := http.NewRequest("DELETE", "/chat-permissions/1", nil)
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusNoContent, w.Code)

	mockPermissionRepo.AssertExpectations(t)
}

func TestRolePermissionHandler_DeletePermission_InvalidID(t *testing.T) {
	// Arrange
	mockRoleRepo := new(MockChatRoleRepository)
	mockPermissionRepo := new(MockChatPermissionRepository)

	controller := controllers.NewRolePermissionController(mockRoleRepo, mockPermissionRepo)
	handler := handlers.NewRolePermissionHandler(controller)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.DELETE("/chat-permissions/:permission_id", handler.DeletePermission)

	req, _ := http.NewRequest("DELETE", "/chat-permissions/invalid", nil)
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRolePermissionHandler_DeletePermission_ControllerError(t *testing.T) {
	// Arrange
	mockRoleRepo := new(MockChatRoleRepository)
	mockPermissionRepo := new(MockChatPermissionRepository)

	permissionID := 999
	mockPermissionRepo.On("DeletePermission", permissionID).Return(errors.New("permission not found"))

	controller := controllers.NewRolePermissionController(mockRoleRepo, mockPermissionRepo)
	handler := handlers.NewRolePermissionHandler(controller)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.DELETE("/chat-permissions/:permission_id", handler.DeletePermission)

	req, _ := http.NewRequest("DELETE", "/chat-permissions/999", nil)
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	mockPermissionRepo.AssertExpectations(t)
}
