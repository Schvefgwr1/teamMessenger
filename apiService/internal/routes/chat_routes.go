package routes

import (
	"apiService/internal/handlers"
	"apiService/internal/middlewares"
	"apiService/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func RegisterChatRoutes(router *gin.Engine, chatHandler *handlers.ChatHandler, publicKeyManager *services.PublicKeyManager, sessionService *services.SessionService, redisClient *redis.Client, rateLimitConfig middlewares.RateLimitConfig) {
	chats := router.Group("api/v1/chats").Use(
		middlewares.JWTMiddlewareWithKeyManager(publicKeyManager, sessionService),
		middlewares.RateLimitMiddleware(redisClient, rateLimitConfig),
		middlewares.RequirePermission("process_chats"),
	)
	{
		chats.GET("/:user_id", chatHandler.GetUserChats)
		chats.POST("", chatHandler.CreateChat)
		chats.PUT("/:chat_id", chatHandler.UpdateChat)
		chats.DELETE("/:chat_id", chatHandler.DeleteChat)
		chats.PATCH("/:chat_id/ban/:user_id", chatHandler.BanUser)
		chats.PATCH("/:chat_id/roles/change", chatHandler.ChangeUserRole)
		chats.GET("/me/role/:chat_id", chatHandler.GetMyRoleInChat)
		chats.GET("/members/:chat_id", chatHandler.GetChatMembers)
		chats.POST("/messages/:chat_id", chatHandler.SendMessage)
		chats.GET("/messages/:chat_id", chatHandler.GetChatMessages)
		chats.GET("/search/:chat_id", chatHandler.SearchMessages)
	}
}

func RegisterRolePermissionRoutes(router *gin.Engine, rolePermissionHandler *handlers.ChatRolePermissionHandler, publicKeyManager *services.PublicKeyManager, sessionService *services.SessionService, redisClient *redis.Client, rateLimitConfig middlewares.RateLimitConfig) {
	// Роли чатов (глобальные)
	roles := router.Group("api/v1/chat-roles").Use(
		middlewares.JWTMiddlewareWithKeyManager(publicKeyManager, sessionService),
		middlewares.RateLimitMiddleware(redisClient, rateLimitConfig),
		middlewares.RequirePermission("process_chats_roles"),
	)
	{
		roles.GET("", rolePermissionHandler.GetAllRoles)
		roles.GET("/:role_id", rolePermissionHandler.GetRoleByID)
		roles.POST("", rolePermissionHandler.CreateRole)
		roles.DELETE("/:role_id", rolePermissionHandler.DeleteRole)
		roles.PATCH("/:role_id/permissions", rolePermissionHandler.UpdateRolePermissions)
	}

	// Permissions чатов (глобальные)
	permissions := router.Group("api/v1/chat-permissions").Use(
		middlewares.JWTMiddlewareWithKeyManager(publicKeyManager, sessionService),
		middlewares.RateLimitMiddleware(redisClient, rateLimitConfig),
		middlewares.RequirePermission("process_chats_permissions"),
	)
	{
		permissions.GET("", rolePermissionHandler.GetAllPermissions)
		permissions.POST("", rolePermissionHandler.CreatePermission)
		permissions.DELETE("/:permission_id", rolePermissionHandler.DeletePermission)
	}
}
