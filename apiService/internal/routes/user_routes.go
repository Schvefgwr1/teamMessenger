package routes

import (
	"apiService/internal/handlers"
	"apiService/internal/middlewares"
	"apiService/internal/services"
	"crypto/rsa"
	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(router *gin.Engine, userHandler *handlers.UserHandler, publicKey *rsa.PublicKey, sessionService *services.SessionService) {
	v1 := router.Group("/api/v1")
	{
		user := v1.Group("/users").Use(middlewares.JWTMiddlewareWithRedis(publicKey, sessionService))
		{
			user.GET("/me", userHandler.GetUser)
			user.PUT("/me", userHandler.UpdateUser)
		}
	}
}
