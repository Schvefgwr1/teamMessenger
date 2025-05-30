package routes

import (
	"apiService/internal/handlers"
	"apiService/internal/middlewares"
	"crypto/rsa"
	"github.com/gin-gonic/gin"
)

func RegisterTaskRoutes(router *gin.Engine, taskHandler *handlers.TaskHandler, publicKey *rsa.PublicKey) {
	tasks := router.Group("api/v1/tasks").Use(middlewares.JWTMiddleware(publicKey))
	{
		tasks.POST("", taskHandler.CreateTask)
		tasks.PATCH("/:task_id/status/:status_id", taskHandler.UpdateTaskStatus)
		tasks.GET("/:task_id", taskHandler.GetTaskByID)
	}

	users := router.Group("api/v1/users").Use(middlewares.JWTMiddleware(publicKey))
	{
		users.GET("/:user_id/tasks", taskHandler.GetUserTasks)
	}
}
