package routes

import (
	"apiService/internal/handlers"
	"apiService/internal/middlewares"
	"apiService/internal/services"
	"github.com/gin-gonic/gin"
)

func RegisterChatRoutes(router *gin.Engine, chatHandler *handlers.ChatHandler, publicKeyManager *services.PublicKeyManager, sessionService *services.SessionService) {
	chats := router.Group("api/v1/chats").Use(middlewares.JWTMiddlewareWithKeyManager(publicKeyManager, sessionService))
	{
		chats.GET("/:user_id", chatHandler.GetUserChats)
		chats.POST("", chatHandler.CreateChat)
		chats.POST("/messages/:chat_id", chatHandler.SendMessage)
		chats.GET("/messages/:chat_id", chatHandler.GetChatMessages)
		chats.GET("/search/:chat_id", chatHandler.SearchMessages)
	}
}
