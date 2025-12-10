package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"chatService/internal/controllers"
	"chatService/internal/handlers"
	"chatService/internal/handlers/dto"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRolePermissionHandler_GetRoleByID_InvalidID_NonNumeric(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockRoleRepo := new(MockChatRoleRepository)
	mockPermissionRepo := new(MockChatPermissionRepository)

	controller := controllers.NewRolePermissionController(mockRoleRepo, mockPermissionRepo)
	handler := handlers.NewRolePermissionHandler(controller)

	router := gin.New()
	router.GET("/chat-roles/:role_id", handler.GetRoleByID)

	req, _ := http.NewRequest("GET", "/chat-roles/abc", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockRoleRepo.AssertNotCalled(t, "GetRoleByID", mock.Anything)
}

func TestRolePermissionHandler_DeleteRole_InvalidID_NonNumeric(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockRoleRepo := new(MockChatRoleRepository)
	mockPermissionRepo := new(MockChatPermissionRepository)

	controller := controllers.NewRolePermissionController(mockRoleRepo, mockPermissionRepo)
	handler := handlers.NewRolePermissionHandler(controller)

	router := gin.New()
	router.DELETE("/chat-roles/:role_id", handler.DeleteRole)

	req, _ := http.NewRequest("DELETE", "/chat-roles/abc", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockRoleRepo.AssertNotCalled(t, "DeleteRole", mock.Anything)
}

func TestRolePermissionHandler_UpdateRolePermissions_InvalidID_NonNumeric(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockRoleRepo := new(MockChatRoleRepository)
	mockPermissionRepo := new(MockChatPermissionRepository)

	controller := controllers.NewRolePermissionController(mockRoleRepo, mockPermissionRepo)
	handler := handlers.NewRolePermissionHandler(controller)

	router := gin.New()
	router.PATCH("/chat-roles/:role_id/permissions", handler.UpdateRolePermissions)

	updateReq := dto.UpdateRolePermissionsRequest{PermissionIDs: []int{1, 2}}
	body, _ := json.Marshal(updateReq)
	req, _ := http.NewRequest("PATCH", "/chat-roles/abc/permissions", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockRoleRepo.AssertNotCalled(t, "UpdateRolePermissions", mock.Anything, mock.Anything)
}

func TestRolePermissionHandler_DeletePermission_InvalidID_NonNumeric(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockRoleRepo := new(MockChatRoleRepository)
	mockPermissionRepo := new(MockChatPermissionRepository)

	controller := controllers.NewRolePermissionController(mockRoleRepo, mockPermissionRepo)
	handler := handlers.NewRolePermissionHandler(controller)

	router := gin.New()
	router.DELETE("/chat-permissions/:permission_id", handler.DeletePermission)

	req, _ := http.NewRequest("DELETE", "/chat-permissions/abc", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockPermissionRepo.AssertNotCalled(t, "DeletePermission", mock.Anything)
}
