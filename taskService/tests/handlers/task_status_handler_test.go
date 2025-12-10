package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"taskService/internal/custom_errors"
	"taskService/internal/handlers"
	"taskService/internal/models"
)

// MockTaskStatusController - мок для TaskStatusController
type MockTaskStatusController struct {
	mock.Mock
}

func (m *MockTaskStatusController) Create(name string) (*models.TaskStatus, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.TaskStatus), args.Error(1)
}

func (m *MockTaskStatusController) GetByID(id int) (*models.TaskStatus, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.TaskStatus), args.Error(1)
}

func (m *MockTaskStatusController) DeleteByID(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockTaskStatusController) GetAll() ([]models.TaskStatus, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.TaskStatus), args.Error(1)
}

// Вспомогательные функции для создания тестовых данных
func createTestTaskStatusModel() *models.TaskStatus {
	return &models.TaskStatus{
		ID:   1,
		Name: "created",
	}
}

func createTestTaskStatusList() []models.TaskStatus {
	return []models.TaskStatus{
		{ID: 1, Name: "created"},
		{ID: 2, Name: "in_progress"},
		{ID: 3, Name: "completed"},
	}
}

// Тесты для TaskStatusHandler.Create

func TestTaskStatusHandler_Create_Success(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockTaskStatusController)
	handler := handlers.NewTaskStatusHandler(mockController)

	statusName := "new_status"
	expectedStatus := createTestTaskStatusModel()
	expectedStatus.Name = statusName

	reqDTO := handlers.CreateStatusDTO{
		Name: statusName,
	}
	reqJSON, _ := json.Marshal(reqDTO)

	mockController.On("Create", statusName).Return(expectedStatus, nil)

	router := gin.New()
	router.POST("/tasks/statuses", handler.Create)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/tasks/statuses", bytes.NewBuffer(reqJSON))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusCreated, w.Code)

	var response models.TaskStatus
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, expectedStatus.ID, response.ID)
	assert.Equal(t, expectedStatus.Name, response.Name)

	mockController.AssertExpectations(t)
}

func TestTaskStatusHandler_Create_InvalidJSON(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockTaskStatusController)
	handler := handlers.NewTaskStatusHandler(mockController)

	invalidJSON := []byte("{invalid json}")

	router := gin.New()
	router.POST("/tasks/statuses", handler.Create)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/tasks/statuses", bytes.NewBuffer(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "invalid request body", response["error"])

	mockController.AssertNotCalled(t, "Create", mock.Anything)
}

func TestTaskStatusHandler_Create_StatusAlreadyExists(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockTaskStatusController)
	handler := handlers.NewTaskStatusHandler(mockController)

	statusName := "existing_status"
	reqDTO := handlers.CreateStatusDTO{
		Name: statusName,
	}
	reqJSON, _ := json.Marshal(reqDTO)
	statusError := custom_errors.ErrStatusAlreadyExists

	mockController.On("Create", statusName).Return(nil, statusError)

	router := gin.New()
	router.POST("/tasks/statuses", handler.Create)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/tasks/statuses", bytes.NewBuffer(reqJSON))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Contains(t, response["error"], "already exists")

	mockController.AssertExpectations(t)
}

// Тесты для TaskStatusHandler.GetByID

func TestTaskStatusHandler_GetByID_Success(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockTaskStatusController)
	handler := handlers.NewTaskStatusHandler(mockController)

	statusID := 1
	expectedStatus := createTestTaskStatusModel()

	mockController.On("GetByID", statusID).Return(expectedStatus, nil)

	router := gin.New()
	router.GET("/tasks/statuses/:id", handler.GetByID)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/tasks/statuses/"+strconv.Itoa(statusID), nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response models.TaskStatus
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, expectedStatus.ID, response.ID)
	assert.Equal(t, expectedStatus.Name, response.Name)

	mockController.AssertExpectations(t)
}

func TestTaskStatusHandler_GetByID_InvalidID(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockTaskStatusController)
	handler := handlers.NewTaskStatusHandler(mockController)

	router := gin.New()
	router.GET("/tasks/statuses/:id", handler.GetByID)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/tasks/statuses/invalid", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "invalid ID", response["error"])

	mockController.AssertNotCalled(t, "GetByID", mock.Anything)
}

func TestTaskStatusHandler_GetByID_NotFound(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockTaskStatusController)
	handler := handlers.NewTaskStatusHandler(mockController)

	statusID := 999
	statusError := custom_errors.NewTaskStatusNotFoundError(strconv.Itoa(statusID))

	mockController.On("GetByID", statusID).Return(nil, statusError)

	router := gin.New()
	router.GET("/tasks/statuses/:id", handler.GetByID)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/tasks/statuses/"+strconv.Itoa(statusID), nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Contains(t, response["error"], "not found")

	mockController.AssertExpectations(t)
}

// Тесты для TaskStatusHandler.DeleteByID

func TestTaskStatusHandler_DeleteByID_Success(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockTaskStatusController)
	handler := handlers.NewTaskStatusHandler(mockController)

	statusID := 1

	mockController.On("DeleteByID", statusID).Return(nil)

	router := gin.New()
	router.DELETE("/tasks/statuses/:id", handler.DeleteByID)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("DELETE", "/tasks/statuses/"+strconv.Itoa(statusID), nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusNoContent, w.Code)

	mockController.AssertExpectations(t)
}

func TestTaskStatusHandler_DeleteByID_InvalidID(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockTaskStatusController)
	handler := handlers.NewTaskStatusHandler(mockController)

	router := gin.New()
	router.DELETE("/tasks/statuses/:id", handler.DeleteByID)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("DELETE", "/tasks/statuses/invalid", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "invalid ID", response["error"])

	mockController.AssertNotCalled(t, "DeleteByID", mock.Anything)
}

func TestTaskStatusHandler_DeleteByID_ControllerError(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockTaskStatusController)
	handler := handlers.NewTaskStatusHandler(mockController)

	statusID := 1
	controllerError := errors.New("database error")

	mockController.On("DeleteByID", statusID).Return(controllerError)

	router := gin.New()
	router.DELETE("/tasks/statuses/:id", handler.DeleteByID)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("DELETE", "/tasks/statuses/"+strconv.Itoa(statusID), nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "could not delete task status", response["error"])

	mockController.AssertExpectations(t)
}

// Тесты для TaskStatusHandler.GetAll

func TestTaskStatusHandler_GetAll_Success(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockTaskStatusController)
	handler := handlers.NewTaskStatusHandler(mockController)

	expectedStatuses := createTestTaskStatusList()

	mockController.On("GetAll").Return(expectedStatuses, nil)

	router := gin.New()
	router.GET("/tasks/statuses", handler.GetAll)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/tasks/statuses", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response []models.TaskStatus
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Len(t, response, 3)
	assert.Equal(t, expectedStatuses[0].ID, response[0].ID)

	mockController.AssertExpectations(t)
}

func TestTaskStatusHandler_GetAll_ControllerError(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockTaskStatusController)
	handler := handlers.NewTaskStatusHandler(mockController)

	controllerError := errors.New("database error")

	mockController.On("GetAll").Return(nil, controllerError)

	router := gin.New()
	router.GET("/tasks/statuses", handler.GetAll)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/tasks/statuses", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "could not get statuses", response["error"])

	mockController.AssertExpectations(t)
}
