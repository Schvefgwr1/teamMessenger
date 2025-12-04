package routes

import (
	"chatService/internal/handlers"
	"chatService/internal/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterChatRoutes(router *gin.Engine, chatHandler *handlers.ChatHandler, messageHandler *handlers.MessageHandler, permissionMiddleware *middlewares.ChatPermissionMiddleware) {
	chats := router.Group("api/v1/chats")
	{
		// Общие роуты
		chats.POST("", chatHandler.CreateChat)

		// Получение списка чатов пользователя - используем /user/:user_id для избежания конфликта
		chats.GET("/user/:user_id", chatHandler.GetUserChats)

		// Роуты с префиксом /messages и /search
		chats.POST("/messages/:chat_id", permissionMiddleware.RequireChatPermission("send_message"), messageHandler.SendMessage)
		chats.GET("/messages/:chat_id", permissionMiddleware.RequireChatPermission("view_messages"), messageHandler.GetChatMessages)
		chats.GET("/search/:chat_id", permissionMiddleware.RequireChatPermission("view_messages"), messageHandler.SearchMessages)

		// Роуты с /:chat_id
		chatID := chats.Group("/:chat_id")
		{
			// Получение роли пользователя в чате (без проверки permissions, но с проверкой членства в чате)
			chatID.GET("/user-roles/:user_id", chatHandler.GetUserRoleInChat)
			// Получение своей роли с permissions (для текущего пользователя)
			chatID.GET("/me/role", chatHandler.GetMyRoleInChat)
			// Получение списка участников чата
			chatID.GET("/members", chatHandler.GetChatMembers)
			chatID.PATCH("/roles/change", permissionMiddleware.RequireChatPermission("change_role"), chatHandler.ChangeUserRole)
			chatID.PATCH("/ban/:user_id", permissionMiddleware.RequireChatPermission("ban_user"), chatHandler.BanUser)
			chatID.PUT("", permissionMiddleware.RequireChatPermission("edit_chat"), chatHandler.UpdateChat)
			chatID.DELETE("", permissionMiddleware.RequireChatPermission("delete_chat"), chatHandler.DeleteChat)
		}
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
