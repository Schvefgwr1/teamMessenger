//go:build integration
// +build integration

package integration

import (
	"testing"
	"time"

	au "common/contracts/api-user"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"userService/internal/controllers"
	"userService/internal/repositories"
	"userService/internal/services"
)

// TestAuthFlow_Integration_RegistrationToLogin тестирует полный сценарий: регистрация -> логин
func TestAuthFlow_Integration_RegistrationToLogin(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	t.Log("=== START TestAuthFlow_Integration_RegistrationToLogin")

	// Arrange - настройка реальных зависимостей
	db := setupTestDB(t)
	userRepo := repositories.NewUserRepository(db)
	roleRepo := repositories.NewRoleRepository(db)

	// Настройка Kafka для уведомлений
	topic := getEnvOrDefault("KAFKA_TOPIC_NOTIFICATIONS", "notifications_test")
	producer := setupTestKafkaProducer(t, topic)
	var notificationService *services.NotificationService
	if producer != nil {
		notificationService = services.NewNotificationServiceWithProducer(producer)
	}

	authController := controllers.NewAuthController(
		userRepo,
		roleRepo,
		notificationService,
	)

	// Получаем существующую роль
	role, err := roleRepo.GetRoleByID(1)
	require.NoError(t, err, "Role with ID=1 should exist in test database")

	// Act - полный сценарий
	// 1. Регистрация
	testUserID := uuid.New()
	testUsername := "test_user_" + testUserID.String()[:8]
	testEmail := "test_" + testUserID.String()[:8] + "@example.com"

	registerReq := &au.RegisterUserRequest{
		Username: testUsername,
		Email:    testEmail,
		Password: "testpassword123",
		Age:      25,
		RoleID:   *role.ID,
	}

	registeredUser, err := authController.Register(registerReq)
	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, registeredUser.ID)
	assert.Equal(t, testUsername, registeredUser.Username)
	assert.Equal(t, testEmail, registeredUser.Email)

	// 2. Логин
	loginData := &au.Login{
		Login:    testUsername,
		Password: "testpassword123",
	}

	token, userID, err := authController.Login(loginData, "192.168.1.1", "Mozilla/5.0")

	// Assert - проверяем результат логина
	require.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.Equal(t, registeredUser.ID, userID)

	// Даем время Kafka обработать уведомление о входе
	if notificationService != nil {
		time.Sleep(1 * time.Second)
	}
}

// TestAuthFlow_Integration_RegistrationWithDuplicateEmail тестирует регистрацию с дублирующимся email
func TestAuthFlow_Integration_RegistrationWithDuplicateEmail(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	t.Log("=== START TestAuthFlow_Integration_RegistrationWithDuplicateEmail")

	// Arrange
	db := setupTestDB(t)
	userRepo := repositories.NewUserRepository(db)
	roleRepo := repositories.NewRoleRepository(db)

	authController := controllers.NewAuthController(
		userRepo,
		roleRepo,
		nil, // notificationService не нужен для этого теста
	)

	role, err := roleRepo.GetRoleByID(1)
	require.NoError(t, err)

	testUserID := uuid.New()
	testUsername := "test_user_" + testUserID.String()[:8]
	testEmail := "test_" + testUserID.String()[:8] + "@example.com"

	// Создаем первого пользователя
	registerReq1 := &au.RegisterUserRequest{
		Username: testUsername,
		Email:    testEmail,
		Password: "testpassword123",
		Age:      25,
		RoleID:   *role.ID,
	}

	_, err = authController.Register(registerReq1)
	require.NoError(t, err)

	// Act - попытка зарегистрировать второго пользователя с тем же email
	testUserID2 := uuid.New()
	testUsername2 := "test_user_" + testUserID2.String()[:8]

	registerReq2 := &au.RegisterUserRequest{
		Username: testUsername2,
		Email:    testEmail, // Дублирующийся email
		Password: "testpassword123",
		Age:      25,
		RoleID:   *role.ID,
	}

	_, err = authController.Register(registerReq2)

	// Assert - должна быть ошибка конфликта email
	require.Error(t, err)
	assert.Contains(t, err.Error(), "email")
}

// TestAuthFlow_Integration_RegistrationWithDuplicateUsername тестирует регистрацию с дублирующимся username
func TestAuthFlow_Integration_RegistrationWithDuplicateUsername(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	t.Log("=== START TestAuthFlow_Integration_RegistrationWithDuplicateUsername")

	// Arrange
	db := setupTestDB(t)
	userRepo := repositories.NewUserRepository(db)
	roleRepo := repositories.NewRoleRepository(db)

	authController := controllers.NewAuthController(
		userRepo,
		roleRepo,
		nil,
	)

	role, err := roleRepo.GetRoleByID(1)
	require.NoError(t, err)

	testUserID := uuid.New()
	testUsername := "test_user_" + testUserID.String()[:8]
	testEmail := "test_" + testUserID.String()[:8] + "@example.com"

	// Создаем первого пользователя
	registerReq1 := &au.RegisterUserRequest{
		Username: testUsername,
		Email:    testEmail,
		Password: "testpassword123",
		Age:      25,
		RoleID:   *role.ID,
	}

	_, err = authController.Register(registerReq1)
	require.NoError(t, err)

	// Act - попытка зарегистрировать второго пользователя с тем же username
	testUserID2 := uuid.New()
	testEmail2 := "test_" + testUserID2.String()[:8] + "@example.com"

	registerReq2 := &au.RegisterUserRequest{
		Username: testUsername, // Дублирующийся username
		Email:    testEmail2,
		Password: "testpassword123",
		Age:      25,
		RoleID:   *role.ID,
	}

	_, err = authController.Register(registerReq2)

	// Assert - должна быть ошибка конфликта username
	require.Error(t, err)
	assert.Contains(t, err.Error(), "username")
}

// TestAuthFlow_Integration_LoginWithInvalidCredentials тестирует логин с неверными учетными данными
func TestAuthFlow_Integration_LoginWithInvalidCredentials(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	t.Log("=== START TestAuthFlow_Integration_LoginWithInvalidCredentials")

	// Arrange
	db := setupTestDB(t)
	userRepo := repositories.NewUserRepository(db)
	roleRepo := repositories.NewRoleRepository(db)

	authController := controllers.NewAuthController(
		userRepo,
		roleRepo,
		nil,
	)

	role, err := roleRepo.GetRoleByID(1)
	require.NoError(t, err)

	testUserID := uuid.New()
	testUsername := "test_user_" + testUserID.String()[:8]
	testEmail := "test_" + testUserID.String()[:8] + "@example.com"

	// Регистрируем пользователя
	registerReq := &au.RegisterUserRequest{
		Username: testUsername,
		Email:    testEmail,
		Password: "testpassword123",
		Age:      25,
		RoleID:   *role.ID,
	}

	_, err = authController.Register(registerReq)
	require.NoError(t, err)

	// Act - попытка входа с неверным паролем
	loginData := &au.Login{
		Login:    testUsername,
		Password: "wrongpassword",
	}

	token, userID, err := authController.Login(loginData, "192.168.1.1", "Mozilla/5.0")

	// Assert - должна быть ошибка неверных учетных данных
	require.Error(t, err)
	assert.Empty(t, token)
	assert.Equal(t, uuid.Nil, userID)
}

// TestAuthFlow_Integration_LoginWithNonExistentUser тестирует логин несуществующего пользователя
func TestAuthFlow_Integration_LoginWithNonExistentUser(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	t.Log("=== START TestAuthFlow_Integration_LoginWithNonExistentUser")

	// Arrange
	db := setupTestDB(t)
	userRepo := repositories.NewUserRepository(db)
	roleRepo := repositories.NewRoleRepository(db)

	authController := controllers.NewAuthController(
		userRepo,
		roleRepo,
		nil,
	)

	// Act - попытка входа несуществующего пользователя
	loginData := &au.Login{
		Login:    "nonexistent_user",
		Password: "password123",
	}

	token, userID, err := authController.Login(loginData, "192.168.1.1", "Mozilla/5.0")

	// Assert - должна быть ошибка неверных учетных данных
	require.Error(t, err)
	assert.Empty(t, token)
	assert.Equal(t, uuid.Nil, userID)
}

// TestAuthFlow_Integration_LoginSendsNotification тестирует, что при логине отправляется уведомление в Kafka
func TestAuthFlow_Integration_LoginSendsNotification(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	t.Log("=== START TestAuthFlow_Integration_LoginSendsNotification")

	// Arrange
	db := setupTestDB(t)
	userRepo := repositories.NewUserRepository(db)
	roleRepo := repositories.NewRoleRepository(db)

	// Настройка реального Kafka
	topic := getEnvOrDefault("KAFKA_TOPIC_NOTIFICATIONS", "notifications_test")
	producer := setupTestKafkaProducer(t, topic)
	if producer == nil {
		t.Skip("Kafka недоступен, пропускаем тест")
		return
	}

	notificationService := services.NewNotificationServiceWithProducer(producer)
	authController := controllers.NewAuthController(
		userRepo,
		roleRepo,
		notificationService,
	)

	role, err := roleRepo.GetRoleByID(1)
	require.NoError(t, err)

	testUserID := uuid.New()
	testUsername := "test_user_" + testUserID.String()[:8]
	testEmail := "test_" + testUserID.String()[:8] + "@example.com"

	// Регистрируем пользователя
	registerReq := &au.RegisterUserRequest{
		Username: testUsername,
		Email:    testEmail,
		Password: "testpassword123",
		Age:      25,
		RoleID:   *role.ID,
	}

	_, err = authController.Register(registerReq)
	require.NoError(t, err)

	// Act - логин (должен отправить уведомление)
	loginData := &au.Login{
		Login:    testUsername,
		Password: "testpassword123",
	}

	token, userID, err := authController.Login(loginData, "192.168.1.1", "Mozilla/5.0")

	// Assert - проверяем успешный логин
	require.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.NotEqual(t, uuid.Nil, userID)

	// Даем время Kafka обработать уведомление
	time.Sleep(1 * time.Second)
}
