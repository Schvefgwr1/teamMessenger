package routes

import (
	"apiService/internal/handlers"
	"apiService/internal/middlewares"
	"apiService/internal/services"
	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(router *gin.Engine, authHandler *handlers.AuthHandler, publicKeyManager *services.PublicKeyManager, sessionService *services.SessionService) {
	v1 := router.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		// Защищенные routes для управления сессиями
		authProtected := v1.Group("/auth").Use(middlewares.JWTMiddlewareWithKeyManager(publicKeyManager, sessionService))
		{
			authProtected.POST("/logout", authHandler.Logout)
		}
	}
}
