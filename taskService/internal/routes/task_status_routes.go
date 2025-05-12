package routes

import (
	"github.com/gin-gonic/gin"
	"taskService/internal/handlers"
)

func RegisterTaskStatusRoutes(r *gin.Engine, handler *handlers.TaskStatusHandler) {
	v1 := r.Group("/api/v1")

	v1.POST("/tasks/statuses", handler.Create)
	v1.GET("/tasks/statuses/:id", handler.GetByID)
	v1.DELETE("/tasks/statuses/:id", handler.DeleteByID)
	v1.GET("/tasks/statuses", handler.GetAll)
}
