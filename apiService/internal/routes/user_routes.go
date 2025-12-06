package routes

import (
	"apiService/internal/handlers"
	"apiService/internal/middlewares"
	"apiService/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func RegisterUserRoutes(
	router *gin.Engine,
	userHandler *handlers.UserHandler,
	publicKeyManager *services.PublicKeyManager,
	sessionService *services.SessionService,
	redisClient *redis.Client,
	rateLimitConfig middlewares.RateLimitConfig,
) {
	v1 := router.Group("/api/v1")

	searches := v1.Group("/searches")
	searches.Use(
		middlewares.JWTMiddlewareWithKeyManager(publicKeyManager, sessionService),
		middlewares.RateLimitMiddleware(redisClient, rateLimitConfig),
	)

	{
		searches.GET("/users", userHandler.SearchUsers)
	}

	// -------------------------
	// /api/v1/users
	// -------------------------
	users := v1.Group("/users")
	users.Use(
		middlewares.JWTMiddlewareWithKeyManager(publicKeyManager, sessionService),
		middlewares.RateLimitMiddleware(redisClient, rateLimitConfig),
	)

	// -------------------------
	// /api/v1/users/me
	// -------------------------
	me := users.Group("/me")
	me.Use(middlewares.RequirePermission("process_your_acc"))

	{
		me.GET("", userHandler.GetUser)    // Без trailing slash
		me.PUT("", userHandler.UpdateUser) // Без trailing slash
	}

	// -------------------------
	// /api/v1/users/:user_id
	// -------------------------
	otherUsers := users.Group("")
	otherUsers.Use(middlewares.RequirePermission("watch_users"))

	{
		otherUsers.GET("/:user_id", userHandler.GetUserProfileByID)
		otherUsers.GET("/:user_id/brief", userHandler.GetUserBrief)
	}

	// -------------------------
	// /api/v1/users/:user_id/role
	// -------------------------
	userRoleUpdate := users.Group("")
	userRoleUpdate.Use(middlewares.RequirePermission("process_users_roles"))

	{
		userRoleUpdate.PATCH("/:user_id/role", userHandler.UpdateUserRole)
	}

	// -------------------------
	// /api/v1/permissions
	// -------------------------
	permissions := v1.Group("/permissions")
	permissions.Use(
		middlewares.JWTMiddlewareWithKeyManager(publicKeyManager, sessionService),
		middlewares.RateLimitMiddleware(redisClient, rateLimitConfig),
		middlewares.RequirePermission("get_permissions"),
	)
	{
		permissions.GET("", userHandler.GetAllPermissions)
	}

	// -------------------------
	// /api/v1/roles
	// -------------------------
	roles := v1.Group("/roles")
	roles.Use(
		middlewares.JWTMiddlewareWithKeyManager(publicKeyManager, sessionService),
		middlewares.RateLimitMiddleware(redisClient, rateLimitConfig),
		middlewares.RequirePermission("process_roles"),
	)
	{
		roles.GET("", userHandler.GetAllRoles)
		roles.POST("", userHandler.CreateRole)
		roles.PATCH("/:role_id/permissions", userHandler.UpdateRolePermissions)
		roles.DELETE("/:role_id", userHandler.DeleteRole)
	}
}
