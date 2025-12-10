package handlers

import (
	"apiService/internal/handlers"
	"bytes"
	at "common/contracts/api-task"
	"encoding/json"
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

// Тесты для TaskHandler.GetTaskByID

func TestTaskHandler_GetTaskByID_Success(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockTaskController)
	handler := handlers.NewTaskHandler(mockController)

	taskID := 1
	expectedTask := &at.TaskServiceResponse{
		Task: &at.TaskResponse{
			ID:    taskID,
			Title: "Test Task",
		},
	}

	mockController.On("GetTaskByID", taskID).Return(expectedTask, nil)

	router := gin.New()
	router.GET("/tasks/:task_id", handler.GetTaskByID)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/tasks/1", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response at.TaskServiceResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, expectedTask.Task.Title, response.Task.Title)

	mockController.AssertExpectations(t)
}

func TestTaskHandler_GetTaskByID_InvalidID(t *testing.T) {
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

	mockController.AssertNotCalled(t, "GetTaskByID", mock.Anything)
}

// Тесты для TaskHandler.GetUserTasks

func TestTaskHandler_GetUserTasks_Success(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockTaskController)
	handler := handlers.NewTaskHandler(mockController)

	userID := uuid.New().String()
	expectedTasks := []at.TaskToList{
		{ID: 1, Title: "Task 1"},
		{ID: 2, Title: "Task 2"},
	}

	mockController.On("GetUserTasks", userID, 20, 0).Return(&expectedTasks, nil)

	router := gin.New()
	router.GET("/users/:user_id/tasks", handler.GetUserTasks)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/users/"+userID+"/tasks?limit=20&offset=0", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response []at.TaskToList
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Len(t, response, 2)

	mockController.AssertExpectations(t)
}

func TestTaskHandler_GetUserTasks_MissingUserID(t *testing.T) {
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
	// Может быть 404 или 400 в зависимости от роутера
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
	req := httptest.NewRequest("GET", "/users/"+userID+"/tasks?limit=101", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	mockController.AssertNotCalled(t, "GetUserTasks", mock.Anything, mock.Anything, mock.Anything)
}

// Тесты для TaskHandler.GetAllStatuses

func TestTaskHandler_GetAllStatuses_Success(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockTaskController)
	handler := handlers.NewTaskHandler(mockController)

	expectedStatuses := []at.TaskStatus{
		{ID: 1, Name: "Open"},
		{ID: 2, Name: "Closed"},
	}

	mockController.On("GetAllStatuses").Return(expectedStatuses, nil)

	router := gin.New()
	router.GET("/tasks/statuses", handler.GetAllStatuses)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/tasks/statuses", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response []at.TaskStatus
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Len(t, response, 2)

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

	mockController.On("UpdateTaskStatus", taskID, statusID).Return(nil)

	router := gin.New()
	router.PATCH("/tasks/:task_id/status/:status_id", handler.UpdateTaskStatus)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("PATCH", "/tasks/1/status/2", nil)
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

	mockController.AssertNotCalled(t, "UpdateTaskStatus", mock.Anything, mock.Anything)
}

// Тесты для TaskHandler.CreateStatus

func TestTaskHandler_CreateStatus_Success(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockTaskController)
	handler := handlers.NewTaskHandler(mockController)

	expectedStatus := &at.TaskStatus{
		ID:   1,
		Name: "New Status",
	}

	mockController.On("CreateStatus", "New Status").Return(expectedStatus, nil)

	router := gin.New()
	router.POST("/tasks/statuses", handler.CreateStatus)

	// Act
	reqBody := `{"name":"New Status"}`
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/tasks/statuses", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusCreated, w.Code)

	var response at.TaskStatus
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, expectedStatus.Name, response.Name)

	mockController.AssertExpectations(t)
}

func TestTaskHandler_CreateStatus_InvalidJSON(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockTaskController)
	handler := handlers.NewTaskHandler(mockController)

	router := gin.New()
	router.POST("/tasks/statuses", handler.CreateStatus)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/tasks/statuses", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	mockController.AssertNotCalled(t, "CreateStatus", mock.Anything)
}

// Тесты для TaskHandler.GetStatusByID

func TestTaskHandler_GetStatusByID_Success(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockTaskController)
	handler := handlers.NewTaskHandler(mockController)

	statusID := 1
	expectedStatus := &at.TaskStatus{
		ID:   statusID,
		Name: "Open",
	}

	mockController.On("GetStatusByID", statusID).Return(expectedStatus, nil)

	router := gin.New()
	router.GET("/tasks/statuses/:status_id", handler.GetStatusByID)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/tasks/statuses/1", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response at.TaskStatus
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, expectedStatus.Name, response.Name)

	mockController.AssertExpectations(t)
}

func TestTaskHandler_GetStatusByID_InvalidID(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockTaskController)
	handler := handlers.NewTaskHandler(mockController)

	router := gin.New()
	router.GET("/tasks/statuses/:status_id", handler.GetStatusByID)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/tasks/statuses/invalid", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	mockController.AssertNotCalled(t, "GetStatusByID", mock.Anything)
}

// Тесты для TaskHandler.DeleteStatus

func TestTaskHandler_DeleteStatus_Success(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockTaskController)
	handler := handlers.NewTaskHandler(mockController)

	statusID := 1

	mockController.On("DeleteStatus", statusID).Return(nil)

	router := gin.New()
	router.DELETE("/tasks/statuses/:status_id", handler.DeleteStatus)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("DELETE", "/tasks/statuses/1", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "status deleted successfully", response["message"])

	mockController.AssertExpectations(t)
}

func TestTaskHandler_DeleteStatus_InvalidID(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockTaskController)
	handler := handlers.NewTaskHandler(mockController)

	router := gin.New()
	router.DELETE("/tasks/statuses/:status_id", handler.DeleteStatus)

	// Act
	w := httptest.NewRecorder()
	req := httptest.NewRequest("DELETE", "/tasks/statuses/invalid", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	mockController.AssertNotCalled(t, "DeleteStatus", mock.Anything)
}

// Тесты для TaskHandler.CreateTask

func TestTaskHandler_CreateTask_Success(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockTaskController)
	handler := handlers.NewTaskHandler(mockController)

	userID := uuid.New()
	taskID := 1
	expectedTask := &at.TaskResponse{
		ID:    taskID,
		Title: "New Task",
	}

	mockController.On("CreateTask", mock.Anything, userID).Return(expectedTask, nil)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userID", userID)
		c.Next()
	})
	router.POST("/tasks", handler.CreateTask)

	// Act - создаем multipart форму
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.WriteField("title", "New Task")
	writer.WriteField("description", "Task description")
	writer.WriteField("executor_id", uuid.New().String())
	writer.Close()

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/tasks", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusCreated, w.Code)

	var response at.TaskResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, expectedTask.Title, response.Title)

	mockController.AssertExpectations(t)
}

func TestTaskHandler_CreateTask_MissingUserID(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockTaskController)
	handler := handlers.NewTaskHandler(mockController)

	router := gin.New()
	router.POST("/tasks", handler.CreateTask)

	// Act
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.WriteField("title", "New Task")
	writer.Close()

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/tasks", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	mockController.AssertNotCalled(t, "CreateTask", mock.Anything, mock.Anything)
}

func TestTaskHandler_CreateTask_InvalidForm(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockController := new(MockTaskController)
	handler := handlers.NewTaskHandler(mockController)

	userID := uuid.New()

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userID", userID)
		c.Next()
	})
	router.POST("/tasks", handler.CreateTask)

	// Act - отправляем пустое тело
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/tasks", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	mockController.AssertNotCalled(t, "CreateTask", mock.Anything, mock.Anything)
}
