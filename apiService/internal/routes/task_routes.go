package routes

import (
	"apiService/internal/handlers"
	"apiService/internal/middlewares"
	"apiService/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func RegisterTaskRoutes(
	router *gin.Engine,
	taskHandler *handlers.TaskHandler,
	publicKeyManager *services.PublicKeyManager,
	sessionService *services.SessionService,
	redisClient *redis.Client,
	rateLimitConfig middlewares.RateLimitConfig,
) {

	// --- MAIN TASKS GROUP ---
	tasks := router.Group("api/v1/tasks")
	tasks.Use(
		middlewares.JWTMiddlewareWithKeyManager(publicKeyManager, sessionService),
		middlewares.RateLimitMiddleware(redisClient, rateLimitConfig),
		middlewares.RequirePermission("process_tasks"),
	)

	{
		// Основные маршруты задач
		tasks.POST("", taskHandler.CreateTask)
		tasks.PATCH("/:task_id/status/:status_id", taskHandler.UpdateTaskStatus)
		tasks.GET("/:task_id", taskHandler.GetTaskByID)

		// == /api/v1/tasks/statuses ==
		statuses := tasks.Group("/statuses")
		statuses.Use(middlewares.RequirePermission("process_tasks_statuses"))

		{
			statuses.GET("", taskHandler.GetAllStatuses)
			statuses.POST("", taskHandler.CreateStatus)
			statuses.GET("/:status_id", taskHandler.GetStatusByID)
			statuses.DELETE("/:status_id", taskHandler.DeleteStatus)
		}
	}

	// --- USER TASKS ---
	users := router.Group("api/v1/users")
	users.Use(
		middlewares.JWTMiddlewareWithKeyManager(publicKeyManager, sessionService),
		middlewares.RateLimitMiddleware(redisClient, rateLimitConfig),
		middlewares.RequirePermission("process_tasks"),
	)

	{
		// /api/v1/users/:user_id/tasks
		users.GET("/:user_id/tasks", taskHandler.GetUserTasks)
	}
}
