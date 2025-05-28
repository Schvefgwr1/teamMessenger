package routes

import (
	"apiService/internal/handlers"
	"apiService/internal/middlewares"
	"crypto/rsa"
	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(router *gin.Engine, userHandler *handlers.UserHandler, publicKey *rsa.PublicKey) {
	v1 := router.Group("/api/v1")
	{
		user := v1.Group("/users").Use(middlewares.JWTMiddleware(publicKey))
		{
			user.GET("/me", userHandler.GetUser)
			user.PUT("/me", userHandler.UpdateUser)
		}
	}
}
