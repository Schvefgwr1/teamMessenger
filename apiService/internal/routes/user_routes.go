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
		}
	}
}
