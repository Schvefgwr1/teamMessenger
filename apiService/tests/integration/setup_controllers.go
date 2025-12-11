//go:build integration
// +build integration

package integration

import (
	"apiService/internal/controllers"
	"apiService/internal/services"
	"testing"

	"github.com/redis/go-redis/v9"
)

// setupTestControllers создает контроллеры с реальными зависимостями
func setupTestControllers(t *testing.T) (
	authController *controllers.AuthController,
	userController *controllers.UserController,
	chatController *controllers.ChatController,
	taskController *controllers.TaskController,
	rolePermissionController *controllers.ChatRolePermissionController,
) {
	redisClient := setupTestRedis(t)
	_, _, _, _, userClient, chatClient, taskClient, fileClient, rolePermissionClient := setupTestHTTPClients(t)

	sessionService := services.NewSessionService(redisClient)
	cacheService := services.NewCacheService(redisClient)

	authController = controllers.NewAuthController(fileClient, userClient, sessionService)
	userController = controllers.NewUserController(fileClient, userClient, cacheService, sessionService)
	chatController = controllers.NewChatController(chatClient, fileClient, cacheService)
	taskController = controllers.NewTaskController(taskClient, fileClient, cacheService)
	rolePermissionController = controllers.NewRolePermissionController(rolePermissionClient, cacheService)

	return authController, userController, chatController, taskController, rolePermissionController
}

// setupAuthController создает AuthController с реальными зависимостями
func setupAuthController(t *testing.T, redisClient *redis.Client) *controllers.AuthController {
	_, _, _, _, userClient, _, _, fileClient, _ := setupTestHTTPClients(t)
	sessionService := services.NewSessionService(redisClient)
	return controllers.NewAuthController(fileClient, userClient, sessionService)
}

// setupUserController создает UserController с реальными зависимостями
func setupUserController(t *testing.T, redisClient *redis.Client) *controllers.UserController {
	_, _, _, _, userClient, _, _, fileClient, _ := setupTestHTTPClients(t)
	cacheService := services.NewCacheService(redisClient)
	sessionService := services.NewSessionService(redisClient)
	return controllers.NewUserController(fileClient, userClient, cacheService, sessionService)
}

// setupChatController создает ChatController с реальными зависимостями
func setupChatController(t *testing.T, redisClient *redis.Client) *controllers.ChatController {
	_, chatServer, _, _, _, chatClient, _, fileClient, _ := setupTestHTTPClients(t)
	defer chatServer.Close()
	cacheService := services.NewCacheService(redisClient)
	return controllers.NewChatController(chatClient, fileClient, cacheService)
}

// setupTaskController создает TaskController с реальными зависимостями
func setupTaskController(t *testing.T, redisClient *redis.Client) *controllers.TaskController {
	_, _, taskServer, _, _, _, taskClient, fileClient, _ := setupTestHTTPClients(t)
	defer taskServer.Close()
	cacheService := services.NewCacheService(redisClient)
	return controllers.NewTaskController(taskClient, fileClient, cacheService)
}

// setupRolePermissionController создает ChatRolePermissionController с реальными зависимостями
func setupRolePermissionController(t *testing.T, redisClient *redis.Client) *controllers.ChatRolePermissionController {
	_, chatServer, _, _, _, _, _, _, rolePermissionClient := setupTestHTTPClients(t)
	defer chatServer.Close()
	cacheService := services.NewCacheService(redisClient)
	return controllers.NewRolePermissionController(rolePermissionClient, cacheService)
}
