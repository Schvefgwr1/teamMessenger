//go:build integration
// +build integration

package integration

import (
	"apiService/internal/controllers"
	"apiService/internal/dto"
	"apiService/internal/services"
	au "common/contracts/api-user"
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestUserController_GetUser_Integration тестирует получение пользователя с реальными интеграциями и кешированием
func TestUserController_GetUser_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange
	redisClient := setupTestRedis(t)
	_, _, _, _, userClient, _, _, fileClient, _ := setupTestHTTPClients(t)

	cacheService := services.NewCacheService(redisClient)
	sessionService := services.NewSessionService(redisClient)
	userController := controllers.NewUserController(fileClient, userClient, cacheService, sessionService)

	userID := uuid.New()
	ctx := context.Background()

	// Act - первый запрос (должен идти в сервис)
	user1, err := userController.GetUser(userID)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, user1)

	// Act - второй запрос (должен быть из кеша)
	user2, err := userController.GetUser(userID)

	// Assert - проверяем, что данные из кеша
	require.NoError(t, err)
	assert.NotNil(t, user2)
	if user1.User != nil && user2.User != nil {
		assert.Equal(t, user1.User.ID, user2.User.ID)
	}

	// Проверяем, что данные реально в кеше
	var cachedUser au.GetUserResponse
	err = cacheService.GetUserCache(ctx, userID.String(), &cachedUser)
	require.NoError(t, err)
	if user1.User != nil && cachedUser.User != nil {
		assert.Equal(t, user1.User.ID, cachedUser.User.ID)
	}
}

// TestUserController_UpdateUser_Integration тестирует обновление пользователя с реальными интеграциями
func TestUserController_UpdateUser_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange
	redisClient := setupTestRedis(t)
	_, _, _, _, userClient, _, _, fileClient, _ := setupTestHTTPClients(t)

	cacheService := services.NewCacheService(redisClient)
	sessionService := services.NewSessionService(redisClient)
	userController := controllers.NewUserController(fileClient, userClient, cacheService, sessionService)

	userID := uuid.New()
	username := "updated_user"
	age := 30
	updateReq := &dto.UpdateUserRequestGateway{
		Username: &username,
		Age:      &age,
	}

	// Act
	response, err := userController.UpdateUser(userID, updateReq, nil)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, response)

	// Проверяем, что кеш пользователя инвалидирован с улучшенными проверками
	cacheKey := cacheService.UserCacheKey(userID.String())
	assertCacheNotExists(t, cacheService, cacheKey)
}

// TestUserController_UpdateUser_Integration_WithFile тестирует обновление пользователя с загрузкой файла
func TestUserController_UpdateUser_Integration_WithFile(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange
	redisClient := setupTestRedis(t)
	_, _, _, _, userClient, _, _, fileClient, _ := setupTestHTTPClients(t)

	cacheService := services.NewCacheService(redisClient)
	sessionService := services.NewSessionService(redisClient)
	userController := controllers.NewUserController(fileClient, userClient, cacheService, sessionService)

	userID := uuid.New()
	updateReq := NewTestUpdateUserRequest()

	file := createTestFileHeader(t, "avatar.jpg", "image/jpeg", []byte("test image"))

	// Act
	response, err := userController.UpdateUser(userID, updateReq, file)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, response)
}

// TestUserController_GetAllPermissions_Integration тестирует получение всех разрешений с кешированием
func TestUserController_GetAllPermissions_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange
	redisClient := setupTestRedis(t)
	_, _, _, _, userClient, _, _, fileClient, _ := setupTestHTTPClients(t)

	cacheService := services.NewCacheService(redisClient)
	sessionService := services.NewSessionService(redisClient)
	userController := controllers.NewUserController(fileClient, userClient, cacheService, sessionService)

	// Act - первый запрос
	permissions1, err := userController.GetAllPermissions()

	// Assert
	require.NoError(t, err)
	assert.NotEmpty(t, permissions1)

	// Act - второй запрос (должен быть из кеша)
	permissions2, err := userController.GetAllPermissions()

	// Assert
	require.NoError(t, err)
	assert.Equal(t, len(permissions1), len(permissions2))
}

// TestUserController_GetAllRoles_Integration тестирует получение всех ролей с кешированием
func TestUserController_GetAllRoles_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange
	redisClient := setupTestRedis(t)
	_, _, _, _, userClient, _, _, fileClient, _ := setupTestHTTPClients(t)

	cacheService := services.NewCacheService(redisClient)
	sessionService := services.NewSessionService(redisClient)
	userController := controllers.NewUserController(fileClient, userClient, cacheService, sessionService)

	// Act - первый запрос
	roles1, err := userController.GetAllRoles()

	// Assert
	require.NoError(t, err)
	assert.NotEmpty(t, roles1)

	// Act - второй запрос (должен быть из кеша)
	roles2, err := userController.GetAllRoles()

	// Assert
	require.NoError(t, err)
	assert.Equal(t, len(roles1), len(roles2))
}

// TestUserController_CreateRole_Integration тестирует создание роли с инвалидацией кеша
func TestUserController_CreateRole_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange
	redisClient := setupTestRedis(t)
	_, _, _, _, userClient, _, _, fileClient, _ := setupTestHTTPClients(t)

	cacheService := services.NewCacheService(redisClient)
	sessionService := services.NewSessionService(redisClient)
	userController := controllers.NewUserController(fileClient, userClient, cacheService, sessionService)

	// Сначала получаем роли, чтобы заполнить кеш
	_, err := userController.GetAllRoles()
	require.NoError(t, err)

	createReq := &au.CreateRoleRequest{
		Name: "test_role_" + uuid.New().String()[:8],
	}

	// Act
	role, err := userController.CreateRole(createReq)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, role)

	// Проверяем, что кеш ролей инвалидирован
	ctx := context.Background()
	exists, err := cacheService.Exists(ctx, "roles:all")
	require.NoError(t, err)
	assert.False(t, exists, "Roles cache should be invalidated after create")
}

// TestUserController_UpdateUserRole_Integration тестирует изменение роли пользователя с отзывом сессий
func TestUserController_UpdateUserRole_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange
	redisClient := setupTestRedis(t)
	_, _, _, _, userClient, _, _, fileClient, _ := setupTestHTTPClients(t)

	cacheService := services.NewCacheService(redisClient)
	sessionService := services.NewSessionService(redisClient)
	userController := controllers.NewUserController(fileClient, userClient, cacheService, sessionService)

	userID := uuid.New()
	ctx := context.Background()

	// Создаем сессию пользователя
	token := "test_token_" + uuid.New().String()
	expiresAt := time.Now().Add(24 * time.Hour)
	err := sessionService.CreateSession(ctx, userID, token, expiresAt)
	require.NoError(t, err)

	// Act
	err = userController.UpdateUserRole(userID, 2)

	// Assert
	require.NoError(t, err)

	// Проверяем, что сессия отозвана
	session, err := sessionService.GetSession(ctx, userID, token)
	require.NoError(t, err)
	assert.Equal(t, services.SessionRevoked, session.Status)
}

// TestUserController_UpdateRolePermissions_Integration тестирует обновление разрешений роли
func TestUserController_UpdateRolePermissions_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange
	redisClient := setupTestRedis(t)
	_, _, _, _, userClient, _, _, fileClient, _ := setupTestHTTPClients(t)

	cacheService := services.NewCacheService(redisClient)
	sessionService := services.NewSessionService(redisClient)
	userController := controllers.NewUserController(fileClient, userClient, cacheService, sessionService)

	roleID := 1
	permissionIDs := []int{1, 2}

	// Act
	err := userController.UpdateRolePermissions(roleID, permissionIDs)

	// Assert
	require.NoError(t, err)

	// Проверяем, что кеш ролей инвалидирован
	ctx := context.Background()
	exists, err := cacheService.Exists(ctx, "roles:all")
	require.NoError(t, err)
	assert.False(t, exists, "Roles cache should be invalidated after update")
}

// TestUserController_DeleteRole_Integration тестирует удаление роли с инвалидацией кеша
func TestUserController_DeleteRole_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange
	redisClient := setupTestRedis(t)
	_, _, _, _, userClient, _, _, fileClient, _ := setupTestHTTPClients(t)

	cacheService := services.NewCacheService(redisClient)
	sessionService := services.NewSessionService(redisClient)
	userController := controllers.NewUserController(fileClient, userClient, cacheService, sessionService)

	roleID := 1

	// Act
	err := userController.DeleteRole(roleID)

	// Assert
	require.NoError(t, err)

	// Проверяем, что кеш ролей инвалидирован
	ctx := context.Background()
	exists, err := cacheService.Exists(ctx, "roles:all")
	require.NoError(t, err)
	assert.False(t, exists, "Roles cache should be invalidated after delete")
}

// TestUserController_GetUserProfileByID_Integration тестирует получение профиля пользователя
func TestUserController_GetUserProfileByID_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange
	redisClient := setupTestRedis(t)
	_, _, _, _, userClient, _, _, fileClient, _ := setupTestHTTPClients(t)

	cacheService := services.NewCacheService(redisClient)
	sessionService := services.NewSessionService(redisClient)
	userController := controllers.NewUserController(fileClient, userClient, cacheService, sessionService)

	userID := uuid.New()

	// Act
	profile, err := userController.GetUserProfileByID(userID)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, profile)
	if profile.User != nil {
		assert.Equal(t, userID.String(), profile.User.ID.String())
	}
}

// TestUserController_GetUserBrief_Integration тестирует получение краткой информации о пользователе
func TestUserController_GetUserBrief_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange
	redisClient := setupTestRedis(t)
	_, _, _, _, userClient, _, _, fileClient, _ := setupTestHTTPClients(t)

	cacheService := services.NewCacheService(redisClient)
	sessionService := services.NewSessionService(redisClient)
	userController := controllers.NewUserController(fileClient, userClient, cacheService, sessionService)

	userID := uuid.New()
	chatID := uuid.New().String()
	requesterID := uuid.New()

	// Act
	brief, err := userController.GetUserBrief(userID, chatID, requesterID)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, brief)
}

// TestUserController_SearchUsers_Integration тестирует поиск пользователей
func TestUserController_SearchUsers_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange
	redisClient := setupTestRedis(t)
	_, _, _, _, userClient, _, _, fileClient, _ := setupTestHTTPClients(t)

	cacheService := services.NewCacheService(redisClient)
	sessionService := services.NewSessionService(redisClient)
	userController := controllers.NewUserController(fileClient, userClient, cacheService, sessionService)

	query := "test"
	limit := 10

	// Act
	result, err := userController.SearchUsers(query, limit)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
}
