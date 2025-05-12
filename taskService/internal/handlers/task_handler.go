package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"taskService/internal/controllers"
	"taskService/internal/custom_errors"
	"taskService/internal/handlers/dto"
)

type TaskHandler struct {
	TaskController *controllers.TaskController
}

func NewTaskHandler(controller *controllers.TaskController) *TaskHandler {
	return &TaskHandler{TaskController: controller}
}

// CreateTask POST /api/v1/tasks
func (h *TaskHandler) CreateTask(c *gin.Context) {
	var taskDTO dto.CreateTaskDTO
	if err := c.ShouldBindJSON(&taskDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	task, err := h.TaskController.Create(&taskDTO)
	if err != nil {
		var userErr *custom_errors.GetUserHTTPError
		var chatErr *custom_errors.GetChatHTTPError
		var fileErr *custom_errors.GetFileHTTPError
		var statusErr *custom_errors.TaskStatusNotFoundError

		switch {
		case errors.As(err, &userErr),
			errors.As(err, &chatErr),
			errors.As(err, &fileErr):
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		case errors.As(err, &statusErr):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, task)
}

// UpdateTaskStatus PATCH /api/v1/tasks/:task_id/status/:status_id
func (h *TaskHandler) UpdateTaskStatus(c *gin.Context) {
	taskID, err := strconv.Atoi(c.Param("task_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task ID"})
		return
	}

	statusID, err := strconv.Atoi(c.Param("status_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status ID"})
		return
	}

	err = h.TaskController.UpdateStatus(taskID, statusID)
	if err != nil {
		var statusErr *custom_errors.TaskStatusNotFoundError
		var taskErr *custom_errors.TaskNotFoundError
		if errors.As(err, &statusErr) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			if errors.As(err, &taskErr) {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			}
		}
		return
	}

	c.Status(http.StatusOK)
}

// GetTaskByID GET /api/v1/tasks/:task_id
func (h *TaskHandler) GetTaskByID(c *gin.Context) {
	taskID, err := strconv.Atoi(c.Param("task_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task ID"})
		return
	}

	task, err := h.TaskController.GetByID(taskID)
	if err != nil {
		var httpClientError *custom_errors.GetFileHTTPError
		if errors.As(err, &httpClientError) {
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, task)
}

// GetUserTasks GET /api/v1/users/:user_id/tasks
func (h *TaskHandler) GetUserTasks(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user ID is required"})
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if err != nil || limit <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit"})
		return
	}

	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil || offset < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid offset"})
		return
	}

	tasks, err := h.TaskController.GetUserTasks(userID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, tasks)
}
