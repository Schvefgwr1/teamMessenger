package routes

import (
	"apiService/internal/handlers"
	"apiService/internal/middlewares"
	"apiService/internal/services"
	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(router *gin.Engine, userHandler *handlers.UserHandler, publicKeyManager *services.PublicKeyManager, sessionService *services.SessionService) {
	v1 := router.Group("/api/v1")
	{
		user := v1.Group("/users").Use(middlewares.JWTMiddlewareWithKeyManager(publicKeyManager, sessionService))
		{
			user.GET("/me", userHandler.GetUser)
			user.PUT("/me", userHandler.UpdateUser)
			user.GET("/:user_id", userHandler.GetUserProfileByID)
		}

		// Разрешения и роли
		permissions := v1.Group("/permissions").Use(middlewares.JWTMiddlewareWithKeyManager(publicKeyManager, sessionService))
		{
			permissions.GET("", userHandler.GetAllPermissions)
		}

		roles := v1.Group("/roles").Use(middlewares.JWTMiddlewareWithKeyManager(publicKeyManager, sessionService))
		{
			roles.GET("", userHandler.GetAllRoles)
			roles.POST("", userHandler.CreateRole)
		}
	}
}
