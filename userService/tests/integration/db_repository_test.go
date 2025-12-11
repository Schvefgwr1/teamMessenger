//go:build integration
// +build integration

package integration

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"userService/internal/models"
	"userService/internal/repositories"
	"userService/internal/utils"
)

// TestUserRepository_Integration_CreateAndGetUser тестирует создание и получение пользователя из реальной БД
func TestUserRepository_Integration_CreateAndGetUser(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	t.Log("=== START TestUserRepository_Integration_CreateAndGetUser")

	// Arrange - настройка реальной БД
	db := setupTestDB(t)
	userRepo := repositories.NewUserRepository(db)

	// Act - создание пользователя
	testUserID := uuid.New()
	testUsername := "test_user_" + testUserID.String()[:8]
	testEmail := "test_" + testUserID.String()[:8] + "@example.com"
	hashedPassword, err := utils.HashPassword("testpassword123")
	require.NoError(t, err)

	// Получаем существующую роль (предполагаем, что роль с ID=1 существует)
	roleRepo := repositories.NewRoleRepository(db)
	role, err := roleRepo.GetRoleByID(1)
	require.NoError(t, err, "Role with ID=1 should exist in test database")

	user := &models.User{
		ID:           testUserID,
		Username:     testUsername,
		Email:        testEmail,
		PasswordHash: hashedPassword,
		RoleID:       *role.ID,
	}

	err = userRepo.CreateUser(user)
	require.NoError(t, err)

	// Assert - проверяем, что пользователь создан
	retrievedUser, err := userRepo.GetUserByID(testUserID)
	require.NoError(t, err)
	assert.Equal(t, testUsername, retrievedUser.Username)
	assert.Equal(t, testEmail, retrievedUser.Email)
	assert.Equal(t, hashedPassword, retrievedUser.PasswordHash)
	assert.Equal(t, *role.ID, retrievedUser.RoleID)
}

// TestUserRepository_Integration_GetUserByEmail тестирует получение пользователя по email из реальной БД
func TestUserRepository_Integration_GetUserByEmail(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	t.Log("=== START TestUserRepository_Integration_GetUserByEmail")

	// Arrange
	db := setupTestDB(t)
	userRepo := repositories.NewUserRepository(db)
	roleRepo := repositories.NewRoleRepository(db)

	testUserID := uuid.New()
	testUsername := "test_user_" + testUserID.String()[:8]
	testEmail := "test_" + testUserID.String()[:8] + "@example.com"
	hashedPassword, err := utils.HashPassword("testpassword123")
	require.NoError(t, err)

	role, err := roleRepo.GetRoleByID(1)
	require.NoError(t, err)

	user := &models.User{
		ID:           testUserID,
		Username:     testUsername,
		Email:        testEmail,
		PasswordHash: hashedPassword,
		RoleID:       *role.ID,
	}

	err = userRepo.CreateUser(user)
	require.NoError(t, err)

	// Act - получение пользователя по email
	retrievedUser, err := userRepo.GetUserByEmail(testEmail)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, testEmail, retrievedUser.Email)
	assert.Equal(t, testUsername, retrievedUser.Username)
}

// TestUserRepository_Integration_GetUserByUsername тестирует получение пользователя по username из реальной БД
func TestUserRepository_Integration_GetUserByUsername(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	t.Log("=== START TestUserRepository_Integration_GetUserByUsername")

	// Arrange
	db := setupTestDB(t)
	userRepo := repositories.NewUserRepository(db)
	roleRepo := repositories.NewRoleRepository(db)

	testUserID := uuid.New()
	testUsername := "test_user_" + testUserID.String()[:8]
	testEmail := "test_" + testUserID.String()[:8] + "@example.com"
	hashedPassword, err := utils.HashPassword("testpassword123")
	require.NoError(t, err)

	role, err := roleRepo.GetRoleByID(1)
	require.NoError(t, err)

	user := &models.User{
		ID:           testUserID,
		Username:     testUsername,
		Email:        testEmail,
		PasswordHash: hashedPassword,
		RoleID:       *role.ID,
	}

	err = userRepo.CreateUser(user)
	require.NoError(t, err)

	// Act - получение пользователя по username
	retrievedUser, err := userRepo.GetUserByUsername(testUsername)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, testUsername, retrievedUser.Username)
	assert.Equal(t, testEmail, retrievedUser.Email)
}

// TestUserRepository_Integration_UpdateUser тестирует обновление пользователя в реальной БД
func TestUserRepository_Integration_UpdateUser(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	t.Log("=== START TestUserRepository_Integration_UpdateUser")

	// Arrange
	db := setupTestDB(t)
	userRepo := repositories.NewUserRepository(db)
	roleRepo := repositories.NewRoleRepository(db)

	testUserID := uuid.New()
	testUsername := "test_user_" + testUserID.String()[:8]
	testEmail := "test_" + testUserID.String()[:8] + "@example.com"
	hashedPassword, err := utils.HashPassword("testpassword123")
	require.NoError(t, err)

	role, err := roleRepo.GetRoleByID(1)
	require.NoError(t, err)

	user := &models.User{
		ID:           testUserID,
		Username:     testUsername,
		Email:        testEmail,
		PasswordHash: hashedPassword,
		RoleID:       *role.ID,
	}

	err = userRepo.CreateUser(user)
	require.NoError(t, err)

	// Act - обновление пользователя
	newDescription := "Updated description"
	user.Description = &newDescription
	err = userRepo.UpdateUser(user)
	require.NoError(t, err)

	// Assert - проверяем обновление
	retrievedUser, err := userRepo.GetUserByID(testUserID)
	require.NoError(t, err)
	assert.Equal(t, newDescription, *retrievedUser.Description)
}

// TestUserRepository_Integration_SearchUsers тестирует поиск пользователей в реальной БД
func TestUserRepository_Integration_SearchUsers(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	t.Log("=== START TestUserRepository_Integration_SearchUsers")

	// Arrange
	db := setupTestDB(t)
	userRepo := repositories.NewUserRepository(db)
	roleRepo := repositories.NewRoleRepository(db)

	role, err := roleRepo.GetRoleByID(1)
	require.NoError(t, err)

	// Создаем несколько тестовых пользователей
	searchPrefix := "search_test_" + uuid.New().String()[:8]
	testUsers := []*models.User{}

	for i := 0; i < 3; i++ {
		testUserID := uuid.New()
		testUsername := searchPrefix + "_user_" + testUserID.String()[:8]
		testEmail := searchPrefix + "_" + testUserID.String()[:8] + "@example.com"
		hashedPassword, err := utils.HashPassword("testpassword123")
		require.NoError(t, err)

		user := &models.User{
			ID:           testUserID,
			Username:     testUsername,
			Email:        testEmail,
			PasswordHash: hashedPassword,
			RoleID:       *role.ID,
		}

		err = userRepo.CreateUser(user)
		require.NoError(t, err)
		testUsers = append(testUsers, user)
	}

	// Act - поиск пользователей
	results, err := userRepo.SearchUsers(searchPrefix, 10)

	// Assert
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(results), 3, "Should find at least 3 users")

	// Проверяем, что все найденные пользователи содержат префикс
	for _, result := range results {
		assert.True(t,
			contains(result.Username, searchPrefix) || contains(result.Email, searchPrefix),
			"User should match search prefix")
	}
}

// TestUserRepository_Integration_GetUserWithRoleAndPermissions тестирует загрузку роли и прав доступа
func TestUserRepository_Integration_GetUserWithRoleAndPermissions(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	t.Log("=== START TestUserRepository_Integration_GetUserWithRoleAndPermissions")

	// Arrange
	db := setupTestDB(t)
	userRepo := repositories.NewUserRepository(db)
	roleRepo := repositories.NewRoleRepository(db)

	testUserID := uuid.New()
	testUsername := "test_user_" + testUserID.String()[:8]
	testEmail := "test_" + testUserID.String()[:8] + "@example.com"
	hashedPassword, err := utils.HashPassword("testpassword123")
	require.NoError(t, err)

	role, err := roleRepo.GetRoleByID(1)
	require.NoError(t, err)

	user := &models.User{
		ID:           testUserID,
		Username:     testUsername,
		Email:        testEmail,
		PasswordHash: hashedPassword,
		RoleID:       *role.ID,
	}

	err = userRepo.CreateUser(user)
	require.NoError(t, err)

	// Act - получение пользователя с предзагрузкой роли и прав
	retrievedUser, err := userRepo.GetUserByID(testUserID)

	// Assert - проверяем, что роль и права загружены
	require.NoError(t, err)
	assert.NotNil(t, retrievedUser.Role)
	assert.Equal(t, *role.ID, retrievedUser.RoleID)
	// Проверяем, что права доступа загружены (если они есть у роли)
	assert.NotNil(t, retrievedUser.Role.Permissions)
}

// TestRoleRepository_Integration_GetRoleByID тестирует получение роли из реальной БД
func TestRoleRepository_Integration_GetRoleByID(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	t.Log("=== START TestRoleRepository_Integration_GetRoleByID")

	// Arrange
	db := setupTestDB(t)
	roleRepo := repositories.NewRoleRepository(db)

	// Act - получение роли (предполагаем, что роль с ID=1 существует)
	role, err := roleRepo.GetRoleByID(1)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, role)
	assert.Equal(t, 1, *role.ID)
}

// TestRoleRepository_Integration_GetAllRoles тестирует получение всех ролей из реальной БД
func TestRoleRepository_Integration_GetAllRoles(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	t.Log("=== START TestRoleRepository_Integration_GetAllRoles")

	// Arrange
	db := setupTestDB(t)
	roleRepo := repositories.NewRoleRepository(db)

	// Act - получение всех ролей
	roles, err := roleRepo.GetAllRoles()

	// Assert
	require.NoError(t, err)
	assert.NotEmpty(t, roles, "Should have at least one role in database")

	// Проверяем структуру роли
	for _, role := range roles {
		assert.NotNil(t, role.ID)
		assert.NotEmpty(t, role.Name)
	}
}

// Вспомогательная функция для проверки подстроки
func contains(s, substr string) bool {
	if len(s) < len(substr) {
		return false
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
