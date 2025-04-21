package routes

import (
	"github.com/gin-gonic/gin"
	"userService/internal/handlers"
)

func RegisterRoutes(router *gin.Engine, authHandler *handlers.AuthHandler, userHandler *handlers.UserHandler, roleHandler *handlers.RoleHandler, permHandler *handlers.PermissionHandler) {
	v1 := router.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		users := v1.Group("/users")
		{
			users.GET("/:user_id", userHandler.GetProfile)
			users.PUT("/:user_id", userHandler.UpdateProfile)
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
	}
}
