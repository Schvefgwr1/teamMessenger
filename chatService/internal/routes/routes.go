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

func RegisterRolePermissionRoutes(router *gin.Engine, rolePermissionHandler *handlers.RolePermissionHandler) {
	// Роли чатов (глобальные)
	roles := router.Group("api/v1/chat-roles")
	{
		roles.GET("", rolePermissionHandler.GetAllRoles)
		roles.GET("/:role_id", rolePermissionHandler.GetRoleByID)
		roles.POST("", rolePermissionHandler.CreateRole)
		roles.DELETE("/:role_id", rolePermissionHandler.DeleteRole)
		roles.PATCH("/:role_id/permissions", rolePermissionHandler.UpdateRolePermissions)
	}

	// Permissions чатов (глобальные)
	permissions := router.Group("api/v1/chat-permissions")
	{
		permissions.GET("", rolePermissionHandler.GetAllPermissions)
		permissions.POST("", rolePermissionHandler.CreatePermission)
		permissions.DELETE("/:permission_id", rolePermissionHandler.DeletePermission)
	}
}
