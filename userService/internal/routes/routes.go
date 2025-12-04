package routes

import (
	"github.com/gin-gonic/gin"
	"userService/internal/handlers"
)

func RegisterRoutes(router *gin.Engine, authHandler *handlers.AuthHandler, userHandler *handlers.UserHandler, roleHandler *handlers.RoleHandler, permHandler *handlers.PermissionHandler, keyHandler handlers.KeyHandler) {
	v1 := router.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		users := v1.Group("/users")
		{
			// Поиск должен быть перед /:user_id чтобы избежать конфликта
			users.GET("/search", userHandler.SearchUsers)
			users.GET("/:user_id", userHandler.GetProfile)
			users.PUT("/:user_id", userHandler.UpdateProfile)
			users.GET("/:user_id/brief", userHandler.GetUserBrief)
		}

		roles := v1.Group("/roles")
		{
			roles.GET("/", roleHandler.GetRoles)
			roles.POST("/", roleHandler.CreateRole)
		}

		permissions := v1.Group("/permissions")
		{
			permissions.GET("/", permHandler.GetPermissions)
		}

		keys := v1.Group("/keys")
		{
			keys.GET("/public", keyHandler.GetPublicKey)
			keys.POST("/regenerate", keyHandler.RegenerateKeys)
		}
	}
}
