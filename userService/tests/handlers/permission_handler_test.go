package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"userService/internal/handlers"
	"userService/internal/models"
)

// MockPermissionController - мок для PermissionController
type MockPermissionController struct {
	mock.Mock
}

func (m *MockPermissionController) GetPermissions() ([]models.Permission, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Permission), args.Error(1)
}

// Тесты для PermissionHandler.GetPermissions

func TestPermissionHandler_GetPermissions_Success(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockPermissionController)
	handler := handlers.NewPermissionHandler(mockController)

	expectedPermissions := []models.Permission{
		{ID: 1, Name: "read", Description: "Read permission"},
		{ID: 2, Name: "write", Description: "Write permission"},
	}

	mockController.On("GetPermissions").Return(expectedPermissions, nil)

	router := gin.New()
	router.GET("/permissions", handler.GetPermissions)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/permissions", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response []models.Permission
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Len(t, response, 2)
	assert.Equal(t, expectedPermissions[0].Name, response[0].Name)

	mockController.AssertExpectations(t)
}

func TestPermissionHandler_GetPermissions_ControllerError(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockPermissionController)
	handler := handlers.NewPermissionHandler(mockController)

	controllerError := errors.New("database error")
	mockController.On("GetPermissions").Return(nil, controllerError)

	router := gin.New()
	router.GET("/permissions", handler.GetPermissions)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/permissions", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "database error", response["error"])

	mockController.AssertExpectations(t)
}
