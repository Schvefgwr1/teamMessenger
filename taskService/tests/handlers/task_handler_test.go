package handlers

import (
	"bytes"
	fc "common/contracts/file-contracts"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"taskService/internal/custom_errors"
	"taskService/internal/handlers"
	"taskService/internal/handlers/dto"
	"taskService/internal/models"
)

// MockTaskController - мок для TaskController
type MockTaskController struct {
	mock.Mock
}

func (m *MockTaskController) Create(taskDTO *dto.CreateTaskDTO) (*models.Task, error) {
	args := m.Called(taskDTO)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Task), args.Error(1)
}

func (m *MockTaskController) UpdateStatus(taskID, statusID int) error {
	args := m.Called(taskID, statusID)
	return args.Error(0)
}

func (m *MockTaskController) GetByID(taskID int) (*dto.TaskResponse, error) {
	args := m.Called(taskID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.TaskResponse), args.Error(1)
}

func (m *MockTaskController) GetUserTasks(userID string, limit, offset int) (*[]dto.TaskToList, error) {
	args := m.Called(userID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*[]dto.TaskToList), args.Error(1)
}

// Вспомогательные функции для создания тестовых данных
func createTestTaskModel() *models.Task {
	return &models.Task{
		ID:          1,
		Title:       "Test Task",
		Description: "Test Description",
		StatusID:    1,
		CreatorID:   uuid.New(),
		ExecutorID:  uuid.New(),
		Status: &models.TaskStatus{
			ID:   1,
			Name: "created",
		},
	}
}

func createTestTaskResponse() *dto.TaskResponse {
	task := createTestTaskModel()
	return &dto.TaskResponse{
		Task:  task,
		Files: &[]fc.File{},
	}
}

func createTestTaskToList() *[]dto.TaskToList {
	return &[]dto.TaskToList{
		{
			ID:     1,
			Title:  "Test Task 1",
			Status: "created",
		},
		{
			ID:     2,
			Title:  "Test Task 2",
			Status: "in_progress",
		},
	}
}

// Тесты для TaskHandler.CreateTask

func TestTaskHandler_CreateTask_Success(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockTaskController)
	handler := handlers.NewTaskHandler(mockController)

	taskDTO := dto.CreateTaskDTO{
		Title:       "Test Task",
		Description: stringPtr("Test Description"),
		CreatorID:   uuid.New(),
		ExecutorID:  uuid.New(),
		FileIDs:     []int{1, 2},
	}
	reqJSON, _ := json.Marshal(taskDTO)
	expectedTask := createTestTaskModel()

	mockController.On("Create", mock.AnythingOfType("*dto.CreateTaskDTO")).Return(expectedTask, nil)

	router := gin.New()
	router.POST("/tasks", handler.CreateTask)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/tasks", bytes.NewBuffer(reqJSON))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response models.Task
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, expectedTask.ID, response.ID)
	assert.Equal(t, expectedTask.Title, response.Title)

	mockController.AssertExpectations(t)
}

func TestTaskHandler_CreateTask_InvalidJSON(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockTaskController)
	handler := handlers.NewTaskHandler(mockController)

	invalidJSON := []byte("{invalid json}")

	router := gin.New()
	router.POST("/tasks", handler.CreateTask)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/tasks", bytes.NewBuffer(invalidJSON))
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

func TestTaskHandler_CreateTask_StatusNotFound(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockTaskController)
	handler := handlers.NewTaskHandler(mockController)

	taskDTO := dto.CreateTaskDTO{
		Title:      "Test Task",
		CreatorID:  uuid.New(),
		ExecutorID: uuid.New(),
	}
	reqJSON, _ := json.Marshal(taskDTO)
	statusError := custom_errors.NewTaskStatusNotFoundError("created")

	mockController.On("Create", mock.AnythingOfType("*dto.CreateTaskDTO")).Return(nil, statusError)

	router := gin.New()
	router.POST("/tasks", handler.CreateTask)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/tasks", bytes.NewBuffer(reqJSON))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Contains(t, response["error"], "not found")

	mockController.AssertExpectations(t)
}

func TestTaskHandler_CreateTask_GetUserHTTPError(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockTaskController)
	handler := handlers.NewTaskHandler(mockController)

	creatorID := uuid.New()
	taskDTO := dto.CreateTaskDTO{
		Title:      "Test Task",
		CreatorID:  creatorID,
		ExecutorID: uuid.New(),
	}
	reqJSON, _ := json.Marshal(taskDTO)
	userError := custom_errors.NewGetUserHTTPError(creatorID.String(), "user not found")

	mockController.On("Create", mock.AnythingOfType("*dto.CreateTaskDTO")).Return(nil, userError)

	router := gin.New()
	router.POST("/tasks", handler.CreateTask)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/tasks", bytes.NewBuffer(reqJSON))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadGateway, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Contains(t, response["error"], "can't get user")

	mockController.AssertExpectations(t)
}

func TestTaskHandler_CreateTask_GetChatHTTPError(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockTaskController)
	handler := handlers.NewTaskHandler(mockController)

	chatID := uuid.New()
	taskDTO := dto.CreateTaskDTO{
		Title:      "Test Task",
		CreatorID:  uuid.New(),
		ExecutorID: uuid.New(),
		ChatID:     chatID,
	}
	reqJSON, _ := json.Marshal(taskDTO)
	chatError := custom_errors.NewGetChatHTTPError(chatID.String(), "chat not found")

	mockController.On("Create", mock.AnythingOfType("*dto.CreateTaskDTO")).Return(nil, chatError)

	router := gin.New()
	router.POST("/tasks", handler.CreateTask)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/tasks", bytes.NewBuffer(reqJSON))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadGateway, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Contains(t, response["error"], "can't get chat")

	mockController.AssertExpectations(t)
}

func TestTaskHandler_CreateTask_GetFileHTTPError(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockTaskController)
	handler := handlers.NewTaskHandler(mockController)

	taskDTO := dto.CreateTaskDTO{
		Title:      "Test Task",
		CreatorID:  uuid.New(),
		ExecutorID: uuid.New(),
		FileIDs:    []int{1},
	}
	reqJSON, _ := json.Marshal(taskDTO)
	fileError := custom_errors.NewGetFileHTTPError(1, "file not found")

	mockController.On("Create", mock.AnythingOfType("*dto.CreateTaskDTO")).Return(nil, fileError)

	router := gin.New()
	router.POST("/tasks", handler.CreateTask)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/tasks", bytes.NewBuffer(reqJSON))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadGateway, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Contains(t, response["error"], "can't get file")

	mockController.AssertExpectations(t)
}

func TestTaskHandler_CreateTask_InternalServerError(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockTaskController)
	handler := handlers.NewTaskHandler(mockController)

	taskDTO := dto.CreateTaskDTO{
		Title:      "Test Task",
		CreatorID:  uuid.New(),
		ExecutorID: uuid.New(),
	}
	reqJSON, _ := json.Marshal(taskDTO)
	internalError := errors.New("database error")

	mockController.On("Create", mock.AnythingOfType("*dto.CreateTaskDTO")).Return(nil, internalError)

	router := gin.New()
	router.POST("/tasks", handler.CreateTask)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/tasks", bytes.NewBuffer(reqJSON))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "internal server error", response["error"])

	mockController.AssertExpectations(t)
}

// Тесты для TaskHandler.UpdateTaskStatus

func TestTaskHandler_UpdateTaskStatus_Success(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockTaskController)
	handler := handlers.NewTaskHandler(mockController)

	taskID := 1
	statusID := 2

	mockController.On("UpdateStatus", taskID, statusID).Return(nil)

	router := gin.New()
	router.PATCH("/tasks/:task_id/status/:status_id", handler.UpdateTaskStatus)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("PATCH", "/tasks/"+strconv.Itoa(taskID)+"/status/"+strconv.Itoa(statusID), nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	mockController.AssertExpectations(t)
}

func TestTaskHandler_UpdateTaskStatus_InvalidTaskID(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockTaskController)
	handler := handlers.NewTaskHandler(mockController)

	router := gin.New()
	router.PATCH("/tasks/:task_id/status/:status_id", handler.UpdateTaskStatus)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("PATCH", "/tasks/invalid/status/2", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "invalid task ID", response["error"])

	mockController.AssertNotCalled(t, "UpdateStatus", mock.Anything, mock.Anything)
}

func TestTaskHandler_UpdateTaskStatus_InvalidStatusID(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockTaskController)
	handler := handlers.NewTaskHandler(mockController)

	taskID := 1

	router := gin.New()
	router.PATCH("/tasks/:task_id/status/:status_id", handler.UpdateTaskStatus)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("PATCH", "/tasks/"+strconv.Itoa(taskID)+"/status/invalid", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "invalid status ID", response["error"])

	mockController.AssertNotCalled(t, "UpdateStatus", mock.Anything, mock.Anything)
}

func TestTaskHandler_UpdateTaskStatus_StatusNotFound(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockTaskController)
	handler := handlers.NewTaskHandler(mockController)

	taskID := 1
	statusID := 999
	statusError := custom_errors.NewTaskStatusNotFoundError(strconv.Itoa(statusID))

	mockController.On("UpdateStatus", taskID, statusID).Return(statusError)

	router := gin.New()
	router.PATCH("/tasks/:task_id/status/:status_id", handler.UpdateTaskStatus)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("PATCH", "/tasks/"+strconv.Itoa(taskID)+"/status/"+strconv.Itoa(statusID), nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Contains(t, response["error"], "not found")

	mockController.AssertExpectations(t)
}

func TestTaskHandler_UpdateTaskStatus_TaskNotFound(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockTaskController)
	handler := handlers.NewTaskHandler(mockController)

	taskID := 999
	statusID := 2
	taskError := custom_errors.NewTaskNotFoundError(taskID)

	mockController.On("UpdateStatus", taskID, statusID).Return(taskError)

	router := gin.New()
	router.PATCH("/tasks/:task_id/status/:status_id", handler.UpdateTaskStatus)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("PATCH", "/tasks/"+strconv.Itoa(taskID)+"/status/"+strconv.Itoa(statusID), nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Contains(t, response["error"], "not found")

	mockController.AssertExpectations(t)
}

func TestTaskHandler_UpdateTaskStatus_InternalServerError(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockTaskController)
	handler := handlers.NewTaskHandler(mockController)

	taskID := 1
	statusID := 2
	internalError := errors.New("database error")

	mockController.On("UpdateStatus", taskID, statusID).Return(internalError)

	router := gin.New()
	router.PATCH("/tasks/:task_id/status/:status_id", handler.UpdateTaskStatus)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("PATCH", "/tasks/"+strconv.Itoa(taskID)+"/status/"+strconv.Itoa(statusID), nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "internal server error", response["error"])

	mockController.AssertExpectations(t)
}

// Тесты для TaskHandler.GetTaskByID

func TestTaskHandler_GetTaskByID_Success(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockTaskController)
	handler := handlers.NewTaskHandler(mockController)

	taskID := 1
	expectedResponse := createTestTaskResponse()

	mockController.On("GetByID", taskID).Return(expectedResponse, nil)

	router := gin.New()
	router.GET("/tasks/:task_id", handler.GetTaskByID)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/tasks/"+strconv.Itoa(taskID), nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response dto.TaskResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, expectedResponse.Task.ID, response.Task.ID)

	mockController.AssertExpectations(t)
}

func TestTaskHandler_GetTaskByID_InvalidTaskID(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockTaskController)
	handler := handlers.NewTaskHandler(mockController)

	router := gin.New()
	router.GET("/tasks/:task_id", handler.GetTaskByID)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/tasks/invalid", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "invalid task ID", response["error"])

	mockController.AssertNotCalled(t, "GetByID", mock.Anything)
}

func TestTaskHandler_GetTaskByID_GetFileHTTPError(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockTaskController)
	handler := handlers.NewTaskHandler(mockController)

	taskID := 1
	fileError := custom_errors.NewGetFileHTTPError(1, "file not found")

	mockController.On("GetByID", taskID).Return(nil, fileError)

	router := gin.New()
	router.GET("/tasks/:task_id", handler.GetTaskByID)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/tasks/"+strconv.Itoa(taskID), nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadGateway, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Contains(t, response["error"], "can't get file")

	mockController.AssertExpectations(t)
}

func TestTaskHandler_GetTaskByID_InternalServerError(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockTaskController)
	handler := handlers.NewTaskHandler(mockController)

	taskID := 1
	internalError := errors.New("database error")

	mockController.On("GetByID", taskID).Return(nil, internalError)

	router := gin.New()
	router.GET("/tasks/:task_id", handler.GetTaskByID)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/tasks/"+strconv.Itoa(taskID), nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "internal server error", response["error"])

	mockController.AssertExpectations(t)
}

// Тесты для TaskHandler.GetUserTasks

func TestTaskHandler_GetUserTasks_Success(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockTaskController)
	handler := handlers.NewTaskHandler(mockController)

	userID := uuid.New().String()
	limit := 10
	offset := 0
	expectedTasks := createTestTaskToList()

	mockController.On("GetUserTasks", userID, limit, offset).Return(expectedTasks, nil)

	router := gin.New()
	router.GET("/users/:user_id/tasks", handler.GetUserTasks)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/users/"+userID+"/tasks?limit="+strconv.Itoa(limit)+"&offset="+strconv.Itoa(offset), nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response []dto.TaskToList
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Len(t, response, 2)

	mockController.AssertExpectations(t)
}

func TestTaskHandler_GetUserTasks_DefaultValues(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockTaskController)
	handler := handlers.NewTaskHandler(mockController)

	userID := uuid.New().String()
	limit := 20
	offset := 0
	expectedTasks := createTestTaskToList()

	mockController.On("GetUserTasks", userID, limit, offset).Return(expectedTasks, nil)

	router := gin.New()
	router.GET("/users/:user_id/tasks", handler.GetUserTasks)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/users/"+userID+"/tasks", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response []dto.TaskToList
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Len(t, response, 2)

	mockController.AssertExpectations(t)
}

func TestTaskHandler_GetUserTasks_EmptyUserID(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockTaskController)
	handler := handlers.NewTaskHandler(mockController)

	router := gin.New()
	router.GET("/users/:user_id/tasks", handler.GetUserTasks)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/users//tasks", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "user ID is required", response["error"])

	mockController.AssertNotCalled(t, "GetUserTasks", mock.Anything, mock.Anything, mock.Anything)
}

func TestTaskHandler_GetUserTasks_InvalidLimit(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockTaskController)
	handler := handlers.NewTaskHandler(mockController)

	userID := uuid.New().String()

	router := gin.New()
	router.GET("/users/:user_id/tasks", handler.GetUserTasks)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/users/"+userID+"/tasks?limit=invalid", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "invalid limit", response["error"])

	mockController.AssertNotCalled(t, "GetUserTasks", mock.Anything, mock.Anything, mock.Anything)
}

func TestTaskHandler_GetUserTasks_InvalidLimitZero(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockTaskController)
	handler := handlers.NewTaskHandler(mockController)

	userID := uuid.New().String()

	router := gin.New()
	router.GET("/users/:user_id/tasks", handler.GetUserTasks)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/users/"+userID+"/tasks?limit=0", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "invalid limit", response["error"])

	mockController.AssertNotCalled(t, "GetUserTasks", mock.Anything, mock.Anything, mock.Anything)
}

func TestTaskHandler_GetUserTasks_InvalidOffset(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockTaskController)
	handler := handlers.NewTaskHandler(mockController)

	userID := uuid.New().String()

	router := gin.New()
	router.GET("/users/:user_id/tasks", handler.GetUserTasks)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/users/"+userID+"/tasks?limit=10&offset=invalid", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "invalid offset", response["error"])

	mockController.AssertNotCalled(t, "GetUserTasks", mock.Anything, mock.Anything, mock.Anything)
}

func TestTaskHandler_GetUserTasks_InvalidOffsetNegative(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockTaskController)
	handler := handlers.NewTaskHandler(mockController)

	userID := uuid.New().String()

	router := gin.New()
	router.GET("/users/:user_id/tasks", handler.GetUserTasks)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/users/"+userID+"/tasks?limit=10&offset=-1", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "invalid offset", response["error"])

	mockController.AssertNotCalled(t, "GetUserTasks", mock.Anything, mock.Anything, mock.Anything)
}

func TestTaskHandler_GetUserTasks_ControllerError(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockTaskController)
	handler := handlers.NewTaskHandler(mockController)

	userID := uuid.New().String()
	limit := 10
	offset := 0
	controllerError := errors.New("database error")

	mockController.On("GetUserTasks", userID, limit, offset).Return(nil, controllerError)

	router := gin.New()
	router.GET("/users/:user_id/tasks", handler.GetUserTasks)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/users/"+userID+"/tasks?limit="+strconv.Itoa(limit)+"&offset="+strconv.Itoa(offset), nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "internal server error", response["error"])

	mockController.AssertExpectations(t)
}

// Вспомогательные функции
func stringPtr(s string) *string {
	return &s
}
