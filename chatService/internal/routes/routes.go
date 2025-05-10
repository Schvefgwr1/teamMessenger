package routes

import (
	"chatService/internal/handlers"
	"chatService/internal/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterChatRoutes(router *gin.Engine, chatHandler *handlers.ChatHandler, messageHandler *handlers.MessageHandler, permissionMiddleware *middlewares.ChatPermissionMiddleware) {
	chats := router.Group("api/v1/chats")
	{
		chats.GET("/:user_id", chatHandler.GetUserChats)
		chats.POST("", chatHandler.CreateChat)

		chats.POST("/messages/:chat_id", permissionMiddleware.RequireChatPermission("send_message"), messageHandler.SendMessage)
		chats.GET("/messages/:chat_id", permissionMiddleware.RequireChatPermission("view_messages"), messageHandler.GetChatMessages)

		chats.GET("/search/:chat_id", permissionMiddleware.RequireChatPermission("view_messages"), messageHandler.SearchMessages) // поиск по всем сообщениям пользователя

		chats.PATCH("/:chat_id/roles/change", permissionMiddleware.RequireChatPermission("change_role"), chatHandler.ChangeUserRole)
		chats.PATCH("/:chat_id/ban/:user_id", permissionMiddleware.RequireChatPermission("ban_user"), chatHandler.BanUser)

		chats.PUT("/:chat_id", permissionMiddleware.RequireChatPermission("edit_chat"), chatHandler.UpdateChat)
		chats.DELETE("/:chat_id", permissionMiddleware.RequireChatPermission("delete_chat"), chatHandler.DeleteChat)
	}
}
