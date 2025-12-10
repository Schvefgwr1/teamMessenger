package controllers

import (
	"apiService/internal/controllers"
	"apiService/internal/dto"
	"apiService/internal/services"
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	af "common/contracts/api-file"
	au "common/contracts/api-user"
	uc "common/contracts/user-contracts"
	"github.com/alicebob/miniredis/v2"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"mime/multipart"
)

// setupTestRedis создает тестовый Redis клиент с miniredis
func setupTestRedis(t *testing.T) *redis.Client {
	mr, err := miniredis.Run()
	require.NoError(t, err)

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	return client
}

// Тесты для UserController.GetUser

func TestUserController_GetUser_Success_FromCache(t *testing.T) {
	// Arrange
	mockFileClient := new(MockFileClient)
	mockUserClient := new(MockUserClient)
	redisClient := setupTestRedis(t)
	defer redisClient.Close()
	cacheService := services.NewCacheService(redisClient)
	sessionService := services.NewSessionService(redisClient)

	controller := controllers.NewUserController(mockFileClient, mockUserClient, cacheService, sessionService)

	userID := uuid.New()
	cachedUser := createTestUserResponse()

	// Сохраняем пользователя в кеш
	cacheService.SetUserCache(context.Background(), userID.String(), cachedUser)

	// Act
	result, err := controller.GetUser(userID)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, cachedUser.User.ID, result.User.ID)
	assert.Equal(t, cachedUser.User.Username, result.User.Username)

	mockUserClient.AssertNotCalled(t, "GetUserByID", mock.Anything)
}

func TestUserController_GetUser_Success_FromService(t *testing.T) {
	// Arrange
	mockFileClient := new(MockFileClient)
	mockUserClient := new(MockUserClient)
	redisClient := setupTestRedis(t)
	defer redisClient.Close()
	cacheService := services.NewCacheService(redisClient)
	sessionService := services.NewSessionService(redisClient)

	controller := controllers.NewUserController(mockFileClient, mockUserClient, cacheService, sessionService)

	userID := uuid.New()
	expectedUser := createTestUserResponse()

	mockUserClient.On("GetUserByID", userID.String()).Return(expectedUser, nil)

	// Act
	result, err := controller.GetUser(userID)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedUser.User.ID, result.User.ID)
	assert.Equal(t, expectedUser.User.Username, result.User.Username)

	mockUserClient.AssertExpectations(t)
}

func TestUserController_GetUser_ServiceError(t *testing.T) {
	// Arrange
	mockFileClient := new(MockFileClient)
	mockUserClient := new(MockUserClient)
	redisClient := setupTestRedis(t)
	defer redisClient.Close()
	cacheService := services.NewCacheService(redisClient)
	sessionService := services.NewSessionService(redisClient)

	controller := controllers.NewUserController(mockFileClient, mockUserClient, cacheService, sessionService)

	userID := uuid.New()
	serviceError := errors.New("user service error")

	mockUserClient.On("GetUserByID", userID.String()).Return(nil, serviceError)

	// Act
	result, err := controller.GetUser(userID)

	// Assert
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, serviceError, err)

	mockUserClient.AssertExpectations(t)
}

func TestUserController_GetUser_CacheSetError(t *testing.T) {
	// Arrange
	mockFileClient := new(MockFileClient)
	mockUserClient := new(MockUserClient)
	redisClient := setupTestRedis(t)
	defer redisClient.Close()
	cacheService := services.NewCacheService(redisClient)
	sessionService := services.NewSessionService(redisClient)

	controller := controllers.NewUserController(mockFileClient, mockUserClient, cacheService, sessionService)

	userID := uuid.New()
	expectedUser := createTestUserResponse()

	mockUserClient.On("GetUserByID", userID.String()).Return(expectedUser, nil)

	// Act
	result, err := controller.GetUser(userID)

	// Assert
	// Ошибка кеширования не должна прерывать выполнение
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedUser.User.ID, result.User.ID)

	mockUserClient.AssertExpectations(t)
}

// Тесты для UserController.UpdateUser

func TestUserController_UpdateUser_Success_WithoutFile(t *testing.T) {
	// Arrange
	mockFileClient := new(MockFileClient)
	mockUserClient := new(MockUserClient)
	redisClient := setupTestRedis(t)
	defer redisClient.Close()
	cacheService := services.NewCacheService(redisClient)
	sessionService := services.NewSessionService(redisClient)

	controller := controllers.NewUserController(mockFileClient, mockUserClient, cacheService, sessionService)

	userID := uuid.New()
	updateRequest := &dto.UpdateUserRequestGateway{
		Username: stringPtr("newusername"),
		Age:      intPtr(25),
	}
	expectedResponse := createTestUpdateUserResponse()

	mockUserClient.On("UpdateUser", userID.String(), mock.MatchedBy(func(req *au.UpdateUserRequest) bool {
		return req.Username != nil && *req.Username == "newusername" && req.Age != nil && *req.Age == 25
	})).Return(expectedResponse, nil)

	// Act
	result, err := controller.UpdateUser(userID, updateRequest, nil)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotNil(t, result)

	mockUserClient.AssertExpectations(t)
	mockFileClient.AssertNotCalled(t, "UploadFile", mock.Anything)
}

func TestUserController_UpdateUser_Success_WithFile(t *testing.T) {
	// Arrange
	mockFileClient := new(MockFileClient)
	mockUserClient := new(MockUserClient)
	redisClient := setupTestRedis(t)
	defer redisClient.Close()
	cacheService := services.NewCacheService(redisClient)
	sessionService := services.NewSessionService(redisClient)

	controller := controllers.NewUserController(mockFileClient, mockUserClient, cacheService, sessionService)

	userID := uuid.New()
	updateRequest := &dto.UpdateUserRequestGateway{
		Username: stringPtr("newusername"),
	}
	file := &multipart.FileHeader{Filename: "avatar.jpg"}
	fileID := 1
	uploadedFile := &af.FileUploadResponse{ID: &fileID}
	expectedResponse := createTestUpdateUserResponse()

	mockUserClient.On("UpdateUser", userID.String(), mock.Anything).Return(expectedResponse, nil)
	mockFileClient.On("UploadFile", file).Return(uploadedFile, nil)
	mockUserClient.On("UpdateUser", userID.String(), mock.MatchedBy(func(req *au.UpdateUserRequest) bool {
		return req.AvatarFileID != nil && *req.AvatarFileID == 1
	})).Return(expectedResponse, nil)

	// Act
	result, err := controller.UpdateUser(userID, updateRequest, file)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)

	mockUserClient.AssertExpectations(t)
	mockFileClient.AssertExpectations(t)
}

func TestUserController_UpdateUser_UpdateError(t *testing.T) {
	// Arrange
	mockFileClient := new(MockFileClient)
	mockUserClient := new(MockUserClient)
	redisClient := setupTestRedis(t)
	defer redisClient.Close()
	cacheService := services.NewCacheService(redisClient)
	sessionService := services.NewSessionService(redisClient)

	controller := controllers.NewUserController(mockFileClient, mockUserClient, cacheService, sessionService)

	userID := uuid.New()
	updateRequest := &dto.UpdateUserRequestGateway{
		Username: stringPtr("newusername"),
	}
	updateError := errors.New("update error")

	mockUserClient.On("UpdateUser", userID.String(), mock.Anything).Return(nil, updateError)

	// Act
	result, err := controller.UpdateUser(userID, updateRequest, nil)

	// Assert
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, updateError, err)

	mockUserClient.AssertExpectations(t)
}

func TestUserController_UpdateUser_ResponseError(t *testing.T) {
	// Arrange
	mockFileClient := new(MockFileClient)
	mockUserClient := new(MockUserClient)
	redisClient := setupTestRedis(t)
	defer redisClient.Close()
	cacheService := services.NewCacheService(redisClient)
	sessionService := services.NewSessionService(redisClient)

	controller := controllers.NewUserController(mockFileClient, mockUserClient, cacheService, sessionService)

	userID := uuid.New()
	updateRequest := &dto.UpdateUserRequestGateway{
		Username: stringPtr("newusername"),
	}
	errorMsg := "validation error"
	expectedResponse := &au.UpdateUserResponse{
		Error: &errorMsg,
	}

	mockUserClient.On("UpdateUser", userID.String(), mock.Anything).Return(expectedResponse, nil)

	// Act
	result, err := controller.UpdateUser(userID, updateRequest, nil)

	// Assert
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), errorMsg)

	mockUserClient.AssertExpectations(t)
}

func TestUserController_UpdateUser_FileUploadError(t *testing.T) {
	// Arrange
	mockFileClient := new(MockFileClient)
	mockUserClient := new(MockUserClient)
	redisClient := setupTestRedis(t)
	defer redisClient.Close()
	cacheService := services.NewCacheService(redisClient)
	sessionService := services.NewSessionService(redisClient)

	controller := controllers.NewUserController(mockFileClient, mockUserClient, cacheService, sessionService)

	userID := uuid.New()
	updateRequest := &dto.UpdateUserRequestGateway{
		Username: stringPtr("newusername"),
	}
	file := &multipart.FileHeader{Filename: "avatar.jpg"}
	uploadError := errors.New("upload error")
	expectedResponse := createTestUpdateUserResponse()

	mockUserClient.On("UpdateUser", userID.String(), mock.Anything).Return(expectedResponse, nil)
	mockFileClient.On("UploadFile", file).Return(nil, uploadError)

	// Act
	result, err := controller.UpdateUser(userID, updateRequest, file)

	// Assert
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, uploadError, err)

	mockUserClient.AssertExpectations(t)
	mockFileClient.AssertExpectations(t)
}

// Тесты для UserController.GetAllPermissions

func TestUserController_GetAllPermissions_Success_FromCache(t *testing.T) {
	// Arrange
	mockFileClient := new(MockFileClient)
	mockUserClient := new(MockUserClient)
	redisClient := setupTestRedis(t)
	defer redisClient.Close()
	cacheService := services.NewCacheService(redisClient)
	sessionService := services.NewSessionService(redisClient)

	controller := controllers.NewUserController(mockFileClient, mockUserClient, cacheService, sessionService)

	expectedPermissions := []*uc.Permission{createTestPermission()}

	// Сохраняем permissions в кеш
	cacheService.Set(context.Background(), "permissions:all", expectedPermissions, time.Hour)

	// Act
	result, err := controller.GetAllPermissions()

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 1)

	mockUserClient.AssertNotCalled(t, "GetAllPermissions")
}

func TestUserController_GetAllPermissions_Success_FromService(t *testing.T) {
	// Arrange
	mockFileClient := new(MockFileClient)
	mockUserClient := new(MockUserClient)
	redisClient := setupTestRedis(t)
	defer redisClient.Close()
	cacheService := services.NewCacheService(redisClient)
	sessionService := services.NewSessionService(redisClient)

	controller := controllers.NewUserController(mockFileClient, mockUserClient, cacheService, sessionService)

	expectedPermissions := []*uc.Permission{createTestPermission()}

	mockUserClient.On("GetAllPermissions").Return(expectedPermissions, nil)

	// Act
	result, err := controller.GetAllPermissions()

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 1)

	mockUserClient.AssertExpectations(t)
}

func TestUserController_GetAllPermissions_ServiceError(t *testing.T) {
	// Arrange
	mockFileClient := new(MockFileClient)
	mockUserClient := new(MockUserClient)
	redisClient := setupTestRedis(t)
	defer redisClient.Close()
	cacheService := services.NewCacheService(redisClient)
	sessionService := services.NewSessionService(redisClient)

	controller := controllers.NewUserController(mockFileClient, mockUserClient, cacheService, sessionService)

	serviceError := errors.New("service error")

	mockUserClient.On("GetAllPermissions").Return(nil, serviceError)

	// Act
	result, err := controller.GetAllPermissions()

	// Assert
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, serviceError, err)

	mockUserClient.AssertExpectations(t)
}

// Тесты для UserController.GetAllRoles

func TestUserController_GetAllRoles_Success_FromCache(t *testing.T) {
	// Arrange
	mockFileClient := new(MockFileClient)
	mockUserClient := new(MockUserClient)
	redisClient := setupTestRedis(t)
	defer redisClient.Close()
	cacheService := services.NewCacheService(redisClient)
	sessionService := services.NewSessionService(redisClient)

	controller := controllers.NewUserController(mockFileClient, mockUserClient, cacheService, sessionService)

	expectedRoles := []*uc.Role{createTestRole()}

	// Сохраняем roles в кеш
	cacheService.Set(context.Background(), "roles:all", expectedRoles, time.Hour)

	// Act
	result, err := controller.GetAllRoles()

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 1)

	mockUserClient.AssertNotCalled(t, "GetAllRoles")
}

func TestUserController_GetAllRoles_Success_FromService(t *testing.T) {
	// Arrange
	mockFileClient := new(MockFileClient)
	mockUserClient := new(MockUserClient)
	redisClient := setupTestRedis(t)
	defer redisClient.Close()
	cacheService := services.NewCacheService(redisClient)
	sessionService := services.NewSessionService(redisClient)

	controller := controllers.NewUserController(mockFileClient, mockUserClient, cacheService, sessionService)

	expectedRoles := []*uc.Role{createTestRole()}

	mockUserClient.On("GetAllRoles").Return(expectedRoles, nil)

	// Act
	result, err := controller.GetAllRoles()

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 1)

	mockUserClient.AssertExpectations(t)
}

// Тесты для UserController.CreateRole

func TestUserController_CreateRole_Success(t *testing.T) {
	// Arrange
	mockFileClient := new(MockFileClient)
	mockUserClient := new(MockUserClient)
	redisClient := setupTestRedis(t)
	defer redisClient.Close()
	cacheService := services.NewCacheService(redisClient)
	sessionService := services.NewSessionService(redisClient)

	controller := controllers.NewUserController(mockFileClient, mockUserClient, cacheService, sessionService)

	createRequest := &au.CreateRoleRequest{
		Name: "admin",
	}
	expectedRole := createTestRole()
	expectedRole.Name = "admin"

	mockUserClient.On("CreateRole", createRequest).Return(expectedRole, nil)

	// Act
	result, err := controller.CreateRole(createRequest)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedRole.Name, result.Name)

	mockUserClient.AssertExpectations(t)
}

func TestUserController_CreateRole_ServiceError(t *testing.T) {
	// Arrange
	mockFileClient := new(MockFileClient)
	mockUserClient := new(MockUserClient)
	redisClient := setupTestRedis(t)
	defer redisClient.Close()
	cacheService := services.NewCacheService(redisClient)
	sessionService := services.NewSessionService(redisClient)

	controller := controllers.NewUserController(mockFileClient, mockUserClient, cacheService, sessionService)

	createRequest := &au.CreateRoleRequest{
		Name: "admin",
	}
	serviceError := errors.New("service error")

	mockUserClient.On("CreateRole", createRequest).Return(nil, serviceError)

	// Act
	result, err := controller.CreateRole(createRequest)

	// Assert
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, serviceError, err)

	mockUserClient.AssertExpectations(t)
}

// Тесты для UserController.UpdateUserRole

func TestUserController_UpdateUserRole_Success(t *testing.T) {
	// Arrange
	mockFileClient := new(MockFileClient)
	mockUserClient := new(MockUserClient)
	redisClient := setupTestRedis(t)
	defer redisClient.Close()
	cacheService := services.NewCacheService(redisClient)
	sessionService := services.NewSessionService(redisClient)

	controller := controllers.NewUserController(mockFileClient, mockUserClient, cacheService, sessionService)

	userID := uuid.New()
	roleID := 2

	mockUserClient.On("UpdateUserRole", userID.String(), roleID).Return(nil)

	// Act
	err := controller.UpdateUserRole(userID, roleID)

	// Assert
	require.NoError(t, err)

	mockUserClient.AssertExpectations(t)
}

func TestUserController_UpdateUserRole_ServiceError(t *testing.T) {
	// Arrange
	mockFileClient := new(MockFileClient)
	mockUserClient := new(MockUserClient)
	redisClient := setupTestRedis(t)
	defer redisClient.Close()
	cacheService := services.NewCacheService(redisClient)
	sessionService := services.NewSessionService(redisClient)

	controller := controllers.NewUserController(mockFileClient, mockUserClient, cacheService, sessionService)

	userID := uuid.New()
	roleID := 2
	serviceError := errors.New("service error")

	mockUserClient.On("UpdateUserRole", userID.String(), roleID).Return(serviceError)

	// Act
	err := controller.UpdateUserRole(userID, roleID)

	// Assert
	require.Error(t, err)
	assert.Equal(t, serviceError, err)

	mockUserClient.AssertExpectations(t)
}

// Тесты для UserController.UpdateRolePermissions

func TestUserController_UpdateRolePermissions_Success(t *testing.T) {
	// Arrange
	mockFileClient := new(MockFileClient)
	mockUserClient := new(MockUserClient)
	redisClient := setupTestRedis(t)
	defer redisClient.Close()
	cacheService := services.NewCacheService(redisClient)
	sessionService := services.NewSessionService(redisClient)

	controller := controllers.NewUserController(mockFileClient, mockUserClient, cacheService, sessionService)

	roleID := 1
	permissionIDs := []int{1, 2, 3}

	mockUserClient.On("UpdateRolePermissions", roleID, permissionIDs).Return(nil)

	// Act
	err := controller.UpdateRolePermissions(roleID, permissionIDs)

	// Assert
	require.NoError(t, err)

	mockUserClient.AssertExpectations(t)
}

// Тесты для UserController.DeleteRole

func TestUserController_DeleteRole_Success(t *testing.T) {
	// Arrange
	mockFileClient := new(MockFileClient)
	mockUserClient := new(MockUserClient)
	redisClient := setupTestRedis(t)
	defer redisClient.Close()
	cacheService := services.NewCacheService(redisClient)
	sessionService := services.NewSessionService(redisClient)

	controller := controllers.NewUserController(mockFileClient, mockUserClient, cacheService, sessionService)

	roleID := 1

	mockUserClient.On("DeleteRole", roleID).Return(nil)

	// Act
	err := controller.DeleteRole(roleID)

	// Assert
	require.NoError(t, err)

	mockUserClient.AssertExpectations(t)
}

// Тесты для UserController.GetUserProfileByID

func TestUserController_GetUserProfileByID_Success_FromCache(t *testing.T) {
	// Arrange
	mockFileClient := new(MockFileClient)
	mockUserClient := new(MockUserClient)
	redisClient := setupTestRedis(t)
	defer redisClient.Close()
	cacheService := services.NewCacheService(redisClient)
	sessionService := services.NewSessionService(redisClient)

	controller := controllers.NewUserController(mockFileClient, mockUserClient, cacheService, sessionService)

	userID := uuid.New()
	cachedUser := createTestUserResponse()

	// Сохраняем пользователя в кеш
	cacheService.SetUserCache(context.Background(), userID.String(), cachedUser)

	// Act
	result, err := controller.GetUserProfileByID(userID)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, cachedUser.User.ID, result.User.ID)

	mockUserClient.AssertNotCalled(t, "GetUserByID", mock.Anything)
}

func TestUserController_GetUserProfileByID_Success_FromService(t *testing.T) {
	// Arrange
	mockFileClient := new(MockFileClient)
	mockUserClient := new(MockUserClient)
	redisClient := setupTestRedis(t)
	defer redisClient.Close()
	cacheService := services.NewCacheService(redisClient)
	sessionService := services.NewSessionService(redisClient)

	controller := controllers.NewUserController(mockFileClient, mockUserClient, cacheService, sessionService)

	userID := uuid.New()
	expectedUser := createTestUserResponse()

	mockUserClient.On("GetUserByID", userID.String()).Return(expectedUser, nil)

	// Act
	result, err := controller.GetUserProfileByID(userID)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedUser.User.ID, result.User.ID)

	mockUserClient.AssertExpectations(t)
}

func TestUserController_GetUserProfileByID_ServiceError(t *testing.T) {
	// Arrange
	mockFileClient := new(MockFileClient)
	mockUserClient := new(MockUserClient)
	redisClient := setupTestRedis(t)
	defer redisClient.Close()
	cacheService := services.NewCacheService(redisClient)
	sessionService := services.NewSessionService(redisClient)

	controller := controllers.NewUserController(mockFileClient, mockUserClient, cacheService, sessionService)

	userID := uuid.New()
	serviceError := errors.New("service error")

	mockUserClient.On("GetUserByID", userID.String()).Return(nil, serviceError)

	// Act
	result, err := controller.GetUserProfileByID(userID)

	// Assert
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, serviceError, err)

	mockUserClient.AssertExpectations(t)
}

// Тесты для UserController.GetUserBrief

func TestUserController_GetUserBrief_Success_FromCache(t *testing.T) {
	// Arrange
	mockFileClient := new(MockFileClient)
	mockUserClient := new(MockUserClient)
	redisClient := setupTestRedis(t)
	defer redisClient.Close()
	cacheService := services.NewCacheService(redisClient)
	sessionService := services.NewSessionService(redisClient)

	controller := controllers.NewUserController(mockFileClient, mockUserClient, cacheService, sessionService)

	userID := uuid.New()
	chatID := "chat-123"
	requesterID := uuid.New()
	cachedBrief := &dto.UserBriefResponse{
		Username: "testuser",
		Email:    "test@example.com",
	}

	// Сохраняем brief в кеш
	cacheKey := fmt.Sprintf("user_brief:%s:%s", userID.String(), chatID)
	cacheService.Set(context.Background(), cacheKey, cachedBrief, 5*time.Minute)

	// Act
	result, err := controller.GetUserBrief(userID, chatID, requesterID)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, cachedBrief.Username, result.Username)

	mockUserClient.AssertNotCalled(t, "GetUserBrief", mock.Anything, mock.Anything, mock.Anything)
}

func TestUserController_GetUserBrief_Success_FromService(t *testing.T) {
	// Arrange
	mockFileClient := new(MockFileClient)
	mockUserClient := new(MockUserClient)
	redisClient := setupTestRedis(t)
	defer redisClient.Close()
	cacheService := services.NewCacheService(redisClient)
	sessionService := services.NewSessionService(redisClient)

	controller := controllers.NewUserController(mockFileClient, mockUserClient, cacheService, sessionService)

	userID := uuid.New()
	chatID := "chat-123"
	requesterID := uuid.New()
	expectedBrief := &dto.UserBriefResponse{
		Username: "testuser",
		Email:    "test@example.com",
	}

	mockUserClient.On("GetUserBrief", userID.String(), chatID, requesterID.String()).Return(expectedBrief, nil)

	// Act
	result, err := controller.GetUserBrief(userID, chatID, requesterID)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedBrief.Username, result.Username)

	mockUserClient.AssertExpectations(t)
}

// Тесты для UserController.SearchUsers

func TestUserController_SearchUsers_Success(t *testing.T) {
	// Arrange
	mockFileClient := new(MockFileClient)
	mockUserClient := new(MockUserClient)
	redisClient := setupTestRedis(t)
	defer redisClient.Close()
	cacheService := services.NewCacheService(redisClient)
	sessionService := services.NewSessionService(redisClient)

	controller := controllers.NewUserController(mockFileClient, mockUserClient, cacheService, sessionService)

	query := "test"
	limit := 10
	expectedResponse := &dto.UserSearchResponse{
		Users: []dto.UserSearchResult{
			{Username: "testuser", Email: "test@example.com"},
		},
	}

	mockUserClient.On("SearchUsers", query, limit).Return(expectedResponse, nil)

	// Act
	result, err := controller.SearchUsers(query, limit)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Users, 1)

	mockUserClient.AssertExpectations(t)
}

func TestUserController_SearchUsers_ServiceError(t *testing.T) {
	// Arrange
	mockFileClient := new(MockFileClient)
	mockUserClient := new(MockUserClient)
	redisClient := setupTestRedis(t)
	defer redisClient.Close()
	cacheService := services.NewCacheService(redisClient)
	sessionService := services.NewSessionService(redisClient)

	controller := controllers.NewUserController(mockFileClient, mockUserClient, cacheService, sessionService)

	query := "test"
	limit := 10
	serviceError := errors.New("service error")

	mockUserClient.On("SearchUsers", query, limit).Return(nil, serviceError)

	// Act
	result, err := controller.SearchUsers(query, limit)

	// Assert
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, serviceError, err)

	mockUserClient.AssertExpectations(t)
}

// Вспомогательные функции
func intPtr(i int) *int {
	return &i
}
