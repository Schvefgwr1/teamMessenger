package routes

import (
	"github.com/gin-gonic/gin"
	"taskService/internal/handlers"
)

func RegisterTaskRoutes(r *gin.Engine, handler *handlers.TaskHandler) {
	v1 := r.Group("/api/v1")

	tasks := v1.Group("/tasks")
	{
		tasks.POST("", handler.CreateTask)
		tasks.PATCH("/:task_id/status/:status_id", handler.UpdateTaskStatus)
		tasks.GET("/:task_id", handler.GetTaskByID)
	}

	users := v1.Group("/users")
	{
		users.GET("/:user_id/tasks", handler.GetUserTasks)
	}
}
