//go:build integration
// +build integration

package integration

import (
	"apiService/internal/controllers"
	"apiService/internal/dto"
	"apiService/internal/services"
	ac "common/contracts/api-chat"
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestChatController_GetUserChats_Integration тестирует получение чатов пользователя с кешированием
func TestChatController_GetUserChats_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange
	redisClient := setupTestRedis(t)
	_, chatServer, _, _, _, chatClient, _, fileClient, _ := setupTestHTTPClients(t)
	defer chatServer.Close()

	cacheService := services.NewCacheService(redisClient)
	chatController := controllers.NewChatController(chatClient, fileClient, cacheService)

	userID := uuid.New()

	// Act - первый запрос (должен идти в сервис)
	chats1, err := chatController.GetUserChats(userID)

	// Assert
	require.NoError(t, err)
	assert.NotEmpty(t, chats1)

	// Act - второй запрос (должен быть из кеша)
	chats2, err := chatController.GetUserChats(userID)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, len(chats1), len(chats2))

	// Проверяем, что данные реально в кеше с улучшенными проверками
	ctx := context.Background()
	cacheKey := cacheService.UserChatListCacheKey(userID.String())
	assertCacheExists(t, cacheService, cacheKey)

	var cachedChats []*ac.ChatResponse
	err = cacheService.GetUserChatListCache(ctx, userID.String(), &cachedChats)
	require.NoError(t, err)
	assert.Equal(t, len(chats1), len(cachedChats))

	// Проверяем TTL кеша (должен быть установлен на 15 минут)
	ttl := getCacheTTL(t, cacheService, cacheKey)
	assert.Greater(t, ttl, time.Duration(0), "Cache TTL should be greater than 0")
	assert.LessOrEqual(t, ttl, 15*time.Minute+time.Second, "Cache TTL should be <= 15 minutes")
}

// TestChatController_CreateChat_Integration тестирует создание чата с инвалидацией кеша
func TestChatController_CreateChat_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange
	redisClient := setupTestRedis(t)
	_, chatServer, _, _, _, chatClient, _, fileClient, _ := setupTestHTTPClients(t)
	defer chatServer.Close()

	cacheService := services.NewCacheService(redisClient)
	chatController := controllers.NewChatController(chatClient, fileClient, cacheService)

	ownerID := uuid.New()
	userID1 := uuid.New()
	userID2 := uuid.New()

	// Создаем кеш для пользователей перед созданием чата
	ctx := context.Background()
	cacheService.SetUserChatListCache(ctx, ownerID.String(), []*ac.ChatResponse{})
	cacheService.SetUserChatListCache(ctx, userID1.String(), []*ac.ChatResponse{})

	createReq := &dto.CreateChatRequestGateway{
		Name:    "test_chat",
		OwnerID: ownerID.String(),
		UserIDs: []string{userID1.String(), userID2.String()},
	}

	// Act
	response, err := chatController.CreateChat(createReq, ownerID, []uuid.UUID{userID1, userID2})

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, response)

	// Проверяем, что кеш списков чатов инвалидирован для всех участников
	exists1, _ := cacheService.Exists(ctx, cacheService.UserChatListCacheKey(ownerID.String()))
	exists2, _ := cacheService.Exists(ctx, cacheService.UserChatListCacheKey(userID1.String()))
	assert.False(t, exists1, "Owner chat list cache should be invalidated")
	assert.False(t, exists2, "User chat list cache should be invalidated")
}

// TestChatController_SendMessage_Integration тестирует отправку сообщения с инвалидацией кеша
func TestChatController_SendMessage_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange
	redisClient := setupTestRedis(t)
	_, chatServer, _, _, _, chatClient, _, fileClient, _ := setupTestHTTPClients(t)
	defer chatServer.Close()

	cacheService := services.NewCacheService(redisClient)
	chatController := controllers.NewChatController(chatClient, fileClient, cacheService)

	chatID := uuid.New()
	senderID := uuid.New()

	// Создаем кеш сообщений перед отправкой
	ctx := context.Background()
	cacheService.SetChatMessagesCache(ctx, chatID.String(), []*ac.GetChatMessage{})

	sendReq := &dto.SendMessageRequestGateway{
		Content: "test message",
	}

	// Act
	message, err := chatController.SendMessage(chatID, senderID, sendReq)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, message)

	// Проверяем, что кеш сообщений инвалидирован
	exists, _ := cacheService.Exists(ctx, cacheService.ChatMessagesCacheKey(chatID.String()))
	assert.False(t, exists, "Messages cache should be invalidated after send")
}

// TestChatController_GetChatMessages_Integration тестирует получение сообщений чата с кешированием
func TestChatController_GetChatMessages_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange
	redisClient := setupTestRedis(t)
	_, chatServer, _, _, _, chatClient, _, fileClient, _ := setupTestHTTPClients(t)
	defer chatServer.Close()

	cacheService := services.NewCacheService(redisClient)
	chatController := controllers.NewChatController(chatClient, fileClient, cacheService)

	chatID := uuid.New()
	userID := uuid.New()

	// Act - первый запрос (offset=0, limit=20, должен кешироваться)
	messages1, err := chatController.GetChatMessages(chatID, userID, 0, 20)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, messages1)

	// Act - второй запрос (должен быть из кеша)
	messages2, err := chatController.GetChatMessages(chatID, userID, 0, 20)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, len(messages1), len(messages2))

	// Проверяем, что данные реально в кеше
	ctx := context.Background()
	var cachedMessages []*ac.GetChatMessage
	err = cacheService.GetChatMessagesCache(ctx, chatID.String(), &cachedMessages)
	require.NoError(t, err)
}

// TestChatController_SearchMessages_Integration тестирует поиск сообщений с кешированием
func TestChatController_SearchMessages_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange
	redisClient := setupTestRedis(t)
	_, chatServer, _, _, _, chatClient, _, fileClient, _ := setupTestHTTPClients(t)
	defer chatServer.Close()

	cacheService := services.NewCacheService(redisClient)
	chatController := controllers.NewChatController(chatClient, fileClient, cacheService)

	userID := uuid.New()
	chatID := uuid.New()
	query := "test"

	// Act - первый запрос (offset=0, limit=20, должен кешироваться)
	result1, err := chatController.SearchMessages(userID, chatID, query, 0, 20)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result1)

	// Act - второй запрос (должен быть из кеша)
	result2, err := chatController.SearchMessages(userID, chatID, query, 0, 20)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result2)
}

// TestChatController_UpdateChat_Integration тестирует обновление чата с инвалидацией кеша
func TestChatController_UpdateChat_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange
	redisClient := setupTestRedis(t)
	_, chatServer, _, _, _, chatClient, _, fileClient, _ := setupTestHTTPClients(t)
	defer chatServer.Close()

	cacheService := services.NewCacheService(redisClient)
	chatController := controllers.NewChatController(chatClient, fileClient, cacheService)

	chatID := uuid.New()
	userID := uuid.New()

	// Создаем кеш перед обновлением
	ctx := context.Background()
	cacheService.SetChatInfoCache(ctx, chatID.String(), &ac.ChatResponse{ID: chatID})

	updateReq := &dto.UpdateChatRequestGateway{
		Name: stringPtr("updated_chat"),
	}
	updateServiceReq := &ac.UpdateChatRequest{
		Name: stringPtr("updated_chat"),
	}

	// Act
	result, err := chatController.UpdateChat(chatID, updateReq, updateServiceReq, userID)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)

	// Проверяем, что кеш информации о чате инвалидирован
	exists, _ := cacheService.Exists(ctx, cacheService.ChatInfoCacheKey(chatID.String()))
	assert.False(t, exists, "Chat info cache should be invalidated after update")
}

// TestChatController_DeleteChat_Integration тестирует удаление чата с инвалидацией кеша
func TestChatController_DeleteChat_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange
	redisClient := setupTestRedis(t)
	_, chatServer, _, _, _, chatClient, _, fileClient, _ := setupTestHTTPClients(t)
	defer chatServer.Close()

	cacheService := services.NewCacheService(redisClient)
	chatController := controllers.NewChatController(chatClient, fileClient, cacheService)

	chatID := uuid.New()
	userID := uuid.New()

	// Создаем кеш перед удалением
	ctx := context.Background()
	cacheService.SetChatInfoCache(ctx, chatID.String(), &ac.ChatResponse{ID: chatID})
	cacheService.SetChatMessagesCache(ctx, chatID.String(), []*ac.GetChatMessage{})

	// Act
	err := chatController.DeleteChat(chatID, userID)

	// Assert
	require.NoError(t, err)

	// Проверяем, что кеш инвалидирован
	exists1, _ := cacheService.Exists(ctx, cacheService.ChatInfoCacheKey(chatID.String()))
	exists2, _ := cacheService.Exists(ctx, cacheService.ChatMessagesCacheKey(chatID.String()))
	assert.False(t, exists1, "Chat info cache should be invalidated after delete")
	assert.False(t, exists2, "Messages cache should be invalidated after delete")
}

// TestChatController_BanUser_Integration тестирует бан пользователя с инвалидацией кеша
func TestChatController_BanUser_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange
	redisClient := setupTestRedis(t)
	_, chatServer, _, _, _, chatClient, _, fileClient, _ := setupTestHTTPClients(t)
	defer chatServer.Close()

	cacheService := services.NewCacheService(redisClient)
	chatController := controllers.NewChatController(chatClient, fileClient, cacheService)

	chatID := uuid.New()
	userID := uuid.New()
	ownerID := uuid.New()

	// Создаем кеш перед баном
	ctx := context.Background()
	cacheService.Set(ctx, cacheService.ChatMembersCacheKey(chatID.String()), []*ac.ChatMember{}, 10*time.Minute)
	cacheService.SetChatUserRoleCache(ctx, chatID.String(), userID.String(), &ac.MyRoleResponse{RoleName: "main"})

	// Act
	err := chatController.BanUser(chatID, userID, ownerID)

	// Assert
	require.NoError(t, err)

	// Проверяем, что кеш инвалидирован
	exists1, _ := cacheService.Exists(ctx, cacheService.ChatMembersCacheKey(chatID.String()))
	exists2, _ := cacheService.Exists(ctx, cacheService.ChatUserRoleCacheKey(chatID.String(), userID.String()))
	assert.False(t, exists1, "Chat members cache should be invalidated after ban")
	assert.False(t, exists2, "User role cache should be invalidated after ban")
}

// TestChatController_ChangeUserRole_Integration тестирует изменение роли пользователя с инвалидацией кеша
func TestChatController_ChangeUserRole_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange
	redisClient := setupTestRedis(t)
	_, chatServer, _, _, _, chatClient, _, fileClient, _ := setupTestHTTPClients(t)
	defer chatServer.Close()

	cacheService := services.NewCacheService(redisClient)
	chatController := controllers.NewChatController(chatClient, fileClient, cacheService)

	chatID := uuid.New()
	ownerID := uuid.New()
	userID := uuid.New()

	// Создаем кеш перед изменением роли
	ctx := context.Background()
	cacheService.SetChatUserRoleCache(ctx, chatID.String(), userID.String(), &ac.MyRoleResponse{RoleName: "main"})

	changeRoleReq := &ac.ChangeRoleRequest{
		UserID: userID,
		RoleID: 2,
	}

	// Act
	err := chatController.ChangeUserRole(chatID, ownerID, changeRoleReq)

	// Assert
	require.NoError(t, err)

	// Проверяем, что кеш роли инвалидирован
	exists, _ := cacheService.Exists(ctx, cacheService.ChatUserRoleCacheKey(chatID.String(), userID.String()))
	assert.False(t, exists, "User role cache should be invalidated after change")
}

// TestChatController_GetMyRoleInChat_Integration тестирует получение роли пользователя с кешированием
func TestChatController_GetMyRoleInChat_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange
	redisClient := setupTestRedis(t)
	_, chatServer, _, _, _, chatClient, _, fileClient, _ := setupTestHTTPClients(t)
	defer chatServer.Close()

	cacheService := services.NewCacheService(redisClient)
	chatController := controllers.NewChatController(chatClient, fileClient, cacheService)

	chatID := uuid.New()
	userID := uuid.New()

	// Act - первый запрос
	role1, err := chatController.GetMyRoleInChat(chatID, userID)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, role1)

	// Act - второй запрос (должен быть из кеша)
	role2, err := chatController.GetMyRoleInChat(chatID, userID)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, role1.RoleName, role2.RoleName)
}

// TestChatController_GetChatMembers_Integration тестирует получение участников чата
func TestChatController_GetChatMembers_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange
	redisClient := setupTestRedis(t)
	_, chatServer, _, _, _, chatClient, _, fileClient, _ := setupTestHTTPClients(t)
	defer chatServer.Close()

	cacheService := services.NewCacheService(redisClient)
	chatController := controllers.NewChatController(chatClient, fileClient, cacheService)

	chatID := uuid.New()

	// Act
	members, err := chatController.GetChatMembers(chatID)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, members)
}
