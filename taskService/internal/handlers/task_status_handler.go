package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"taskService/internal/controllers"
)

type TaskStatusHandler struct {
	Controller *controllers.TaskStatusController
}

func NewTaskStatusHandler(controller *controllers.TaskStatusController) *TaskStatusHandler {
	return &TaskStatusHandler{Controller: controller}
}

type CreateStatusDTO struct {
	Name string `json:"name" binding:"required"`
}

// Create POST /api/v1/tasks/statuses
func (h *TaskStatusHandler) Create(c *gin.Context) {
	var dto CreateStatusDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	status, err := h.Controller.Create(dto.Name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, status)
}

// GetByID GET /api/v1/tasks/statuses/:id
func (h *TaskStatusHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	status, err := h.Controller.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, status)
}

// DeleteByID DELETE /api/v1/tasks/statuses/:id
func (h *TaskStatusHandler) DeleteByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	err = h.Controller.DeleteByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not delete task status"})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetAll GET /api/v1/tasks/statuses
func (h *TaskStatusHandler) GetAll(c *gin.Context) {
	statuses, err := h.Controller.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not get statuses"})
		return
	}

	c.JSON(http.StatusOK, statuses)
}
