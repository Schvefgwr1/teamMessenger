package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"chatService/internal/controllers"
	"chatService/internal/handlers"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRolePermissionHandler_CreateRole_InvalidJSON_Validation(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockRoleRepo := new(MockChatRoleRepository)
	mockPermissionRepo := new(MockChatPermissionRepository)

	controller := controllers.NewRolePermissionController(mockRoleRepo, mockPermissionRepo)
	handler := handlers.NewRolePermissionHandler(controller)

	router := gin.New()
	router.POST("/chat-roles", handler.CreateRole)

	req, _ := http.NewRequest("POST", "/chat-roles", bytes.NewBuffer([]byte(`{invalid json`)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockRoleRepo.AssertNotCalled(t, "CreateRole", mock.Anything, mock.Anything)
}

func TestRolePermissionHandler_CreateRole_MissingName_Validation(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockRoleRepo := new(MockChatRoleRepository)
	mockPermissionRepo := new(MockChatPermissionRepository)

	controller := controllers.NewRolePermissionController(mockRoleRepo, mockPermissionRepo)
	handler := handlers.NewRolePermissionHandler(controller)

	router := gin.New()
	router.POST("/chat-roles", handler.CreateRole)

	reqBody := `{"permissionIds": [1, 2]}`
	req, _ := http.NewRequest("POST", "/chat-roles", bytes.NewBuffer([]byte(reqBody)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockRoleRepo.AssertNotCalled(t, "CreateRole", mock.Anything, mock.Anything)
}

func TestRolePermissionHandler_UpdateRolePermissions_InvalidJSON_Validation(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockRoleRepo := new(MockChatRoleRepository)
	mockPermissionRepo := new(MockChatPermissionRepository)

	controller := controllers.NewRolePermissionController(mockRoleRepo, mockPermissionRepo)
	handler := handlers.NewRolePermissionHandler(controller)

	router := gin.New()
	router.PATCH("/chat-roles/:role_id/permissions", handler.UpdateRolePermissions)

	req, _ := http.NewRequest("PATCH", "/chat-roles/1/permissions", bytes.NewBuffer([]byte(`{invalid json`)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockRoleRepo.AssertNotCalled(t, "UpdateRolePermissions", mock.Anything, mock.Anything)
}

func TestRolePermissionHandler_UpdateRolePermissions_MissingPermissionIDs_Validation(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockRoleRepo := new(MockChatRoleRepository)
	mockPermissionRepo := new(MockChatPermissionRepository)

	controller := controllers.NewRolePermissionController(mockRoleRepo, mockPermissionRepo)
	handler := handlers.NewRolePermissionHandler(controller)

	router := gin.New()
	router.PATCH("/chat-roles/:role_id/permissions", handler.UpdateRolePermissions)

	reqBody := `{}`
	req, _ := http.NewRequest("PATCH", "/chat-roles/1/permissions", bytes.NewBuffer([]byte(reqBody)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockRoleRepo.AssertNotCalled(t, "UpdateRolePermissions", mock.Anything, mock.Anything)
}

func TestRolePermissionHandler_CreatePermission_InvalidJSON_Validation(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockRoleRepo := new(MockChatRoleRepository)
	mockPermissionRepo := new(MockChatPermissionRepository)

	controller := controllers.NewRolePermissionController(mockRoleRepo, mockPermissionRepo)
	handler := handlers.NewRolePermissionHandler(controller)

	router := gin.New()
	router.POST("/chat-permissions", handler.CreatePermission)

	req, _ := http.NewRequest("POST", "/chat-permissions", bytes.NewBuffer([]byte(`{invalid json`)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockPermissionRepo.AssertNotCalled(t, "CreatePermission", mock.Anything)
}

func TestRolePermissionHandler_CreatePermission_MissingName_Validation(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockRoleRepo := new(MockChatRoleRepository)
	mockPermissionRepo := new(MockChatPermissionRepository)

	controller := controllers.NewRolePermissionController(mockRoleRepo, mockPermissionRepo)
	handler := handlers.NewRolePermissionHandler(controller)

	router := gin.New()
	router.POST("/chat-permissions", handler.CreatePermission)

	reqBody := `{}`
	req, _ := http.NewRequest("POST", "/chat-permissions", bytes.NewBuffer([]byte(reqBody)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockPermissionRepo.AssertNotCalled(t, "CreatePermission", mock.Anything)
}
