package routes

import (
	"apiService/internal/handlers"
	"apiService/internal/middlewares"
	"apiService/internal/services"
	"crypto/rsa"
	"github.com/gin-gonic/gin"
)

func RegisterTaskRoutes(router *gin.Engine, taskHandler *handlers.TaskHandler, publicKey *rsa.PublicKey, sessionService *services.SessionService) {
	tasks := router.Group("api/v1/tasks").Use(middlewares.JWTMiddlewareWithRedis(publicKey, sessionService))
	{
		tasks.POST("", taskHandler.CreateTask)
		tasks.PATCH("/:task_id/status/:status_id", taskHandler.UpdateTaskStatus)
		tasks.GET("/:task_id", taskHandler.GetTaskByID)
	}

	users := router.Group("api/v1/users").Use(middlewares.JWTMiddlewareWithRedis(publicKey, sessionService))
	{
		users.GET("/:user_id/tasks", taskHandler.GetUserTasks)
	}
}
