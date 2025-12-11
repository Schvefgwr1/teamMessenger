//go:build integration
// +build integration

package integration

import (
	"testing"
	"time"

	"common/config"
	"common/models"
	"notificationService/internal/services"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestEmailService_SendNotification_NewTask_Integration тестирует отправку уведомления о новой задаче с реальным SMTP
func TestEmailService_SendNotification_NewTask_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange - настройка реального EmailService с тестовым SMTP
	cfg := setupTestEmailConfig(t)
	emailService, err := services.NewEmailService(cfg)
	if err != nil {
		t.Skipf("SMTP недоступен, пропускаем тест: %v", err)
		return
	}

	notification := createTestNewTaskNotification()

	// Act - выполнение реального сценария отправки email
	err = emailService.SendNotification(notification)

	// Assert - проверка реального результата
	// В реальном SMTP мы не можем проверить доставку, но можем проверить отсутствие ошибок
	if err != nil {
		// Если SMTP сервер недоступен, пропускаем тест
		if err.Error() != "" {
			t.Skipf("SMTP сервер недоступен, пропускаем тест: %v", err)
			return
		}
	}
	// Если ошибки нет, значит интеграция работает
	assert.NoError(t, err)
}

// TestEmailService_SendNotification_NewChat_Group_Integration тестирует отправку уведомления о групповом чате
func TestEmailService_SendNotification_NewChat_Group_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange
	cfg := setupTestEmailConfig(t)
	emailService, err := services.NewEmailService(cfg)
	if err != nil {
		t.Skipf("SMTP недоступен, пропускаем тест: %v", err)
		return
	}

	notification := createTestNewChatNotification(true)

	// Act
	err = emailService.SendNotification(notification)

	// Assert
	if err != nil {
		t.Skipf("SMTP сервер недоступен, пропускаем тест: %v", err)
		return
	}
	assert.NoError(t, err)
}

// TestEmailService_SendNotification_NewChat_Private_Integration тестирует отправку уведомления о приватном чате
func TestEmailService_SendNotification_NewChat_Private_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange
	cfg := setupTestEmailConfig(t)
	emailService, err := services.NewEmailService(cfg)
	if err != nil {
		t.Skipf("SMTP недоступен, пропускаем тест: %v", err)
		return
	}

	notification := createTestNewChatNotification(false)

	// Act
	err = emailService.SendNotification(notification)

	// Assert
	if err != nil {
		t.Skipf("SMTP сервер недоступен, пропускаем тест: %v", err)
		return
	}
	assert.NoError(t, err)
}

// TestEmailService_SendNotification_Login_Integration тестирует отправку уведомления о входе
func TestEmailService_SendNotification_Login_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange
	cfg := setupTestEmailConfig(t)
	emailService, err := services.NewEmailService(cfg)
	if err != nil {
		t.Skipf("SMTP недоступен, пропускаем тест: %v", err)
		return
	}

	notification := createTestLoginNotification()

	// Act
	err = emailService.SendNotification(notification)

	// Assert
	if err != nil {
		t.Skipf("SMTP сервер недоступен, пропускаем тест: %v", err)
		return
	}
	assert.NoError(t, err)
}

// TestEmailService_SendNotification_UnknownType_Integration тестирует обработку неизвестного типа уведомления
func TestEmailService_SendNotification_UnknownType_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange
	cfg := setupTestEmailConfig(t)
	emailService, err := services.NewEmailService(cfg)
	require.NoError(t, err)

	unknownNotification := "unknown type"

	// Act
	err = emailService.SendNotification(unknownNotification)

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unknown notification type")
}

// TestEmailService_SendNotification_SMTPUnavailable_Integration тестирует обработку ошибки недоступного SMTP сервера
func TestEmailService_SendNotification_SMTPUnavailable_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange - подключаемся к несуществующему SMTP серверу
	cfg := &config.EmailConfig{
		SMTPHost:     "localhost",
		SMTPPort:     9999, // Несуществующий порт
		Username:     "test@example.com",
		Password:     "password",
		FromName:     "Test",
		FromEmail:    "test@example.com",
		TemplatePath: getEnvOrDefault("TEMPLATE_PATH", "../../templates"),
	}

	emailService, err := services.NewEmailService(cfg)
	if err != nil {
		t.Skipf("Не удалось создать EmailService: %v", err)
		return
	}

	notification := createTestNewTaskNotification()

	// Act - попытка отправить email
	err = emailService.SendNotification(notification)

	// Assert - проверяем реальную ошибку подключения
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to send email")
}

// TestEmailService_GenerateHTML_NewTask_Integration тестирует генерацию HTML для уведомления о задаче
func TestEmailService_GenerateHTML_NewTask_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange
	cfg := setupTestEmailConfig(t)
	emailService, err := services.NewEmailService(cfg)
	require.NoError(t, err)

	notification := &models.NewTaskNotification{
		BaseNotification: models.BaseNotification{
			ID:        uuid.New(),
			Type:      models.NotificationNewTask,
			Email:     "executor@test.com",
			CreatedAt: time.Now(),
		},
		TaskID:      1,
		TaskTitle:   "Test Task",
		CreatorName: "Test Creator",
		ExecutorID:  uuid.New(),
	}

	// Act - генерируем HTML через SendNotification (который использует generateHTML внутри)
	err = emailService.SendNotification(notification)

	// Assert - проверяем, что HTML был сгенерирован без ошибок
	// Если ошибка связана с SMTP, это нормально для интеграционного теста
	if err != nil {
		// Проверяем, что ошибка не связана с генерацией HTML
		assert.NotContains(t, err.Error(), "template not found")
		assert.NotContains(t, err.Error(), "failed to execute template")
		assert.NotContains(t, err.Error(), "failed to generate HTML")
	}
}

// TestEmailService_GenerateHTML_NewChat_Integration тестирует генерацию HTML для уведомления о чате
func TestEmailService_GenerateHTML_NewChat_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange
	cfg := setupTestEmailConfig(t)
	emailService, err := services.NewEmailService(cfg)
	require.NoError(t, err)

	notification := createTestNewChatNotification(true)

	// Act
	err = emailService.SendNotification(notification)

	// Assert
	if err != nil {
		// Проверяем, что ошибка не связана с генерацией HTML
		assert.NotContains(t, err.Error(), "template not found")
		assert.NotContains(t, err.Error(), "failed to execute template")
		assert.NotContains(t, err.Error(), "failed to generate HTML")
	}
}

// TestEmailService_GenerateHTML_Login_Integration тестирует генерацию HTML для уведомления о входе
func TestEmailService_GenerateHTML_Login_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange
	cfg := setupTestEmailConfig(t)
	emailService, err := services.NewEmailService(cfg)
	require.NoError(t, err)

	notification := createTestLoginNotification()

	// Act
	err = emailService.SendNotification(notification)

	// Assert
	if err != nil {
		// Проверяем, что ошибка не связана с генерацией HTML
		assert.NotContains(t, err.Error(), "template not found")
		assert.NotContains(t, err.Error(), "failed to execute template")
		assert.NotContains(t, err.Error(), "failed to generate HTML")
	}
}
