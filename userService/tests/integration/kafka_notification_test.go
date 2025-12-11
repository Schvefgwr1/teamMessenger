//go:build integration
// +build integration

package integration

import (
	"testing"
	"time"

	"common/kafka"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"userService/internal/services"
)

// TestNotificationService_Integration_SendLoginNotification тестирует отправку уведомления о входе в реальный Kafka
func TestNotificationService_Integration_SendLoginNotification(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	t.Log("=== START TestNotificationService_Integration_SendLoginNotification")

	// Arrange - настройка реального Kafka producer
	topic := getEnvOrDefault("KAFKA_TOPIC_NOTIFICATIONS", "notifications_test")
	producer := setupTestKafkaProducer(t, topic)
	if producer == nil {
		return // Kafka недоступен, тест пропущен
	}

	notificationService := services.NewNotificationServiceWithProducer(producer)

	// Act - отправка уведомления о входе
	userID := uuid.New()
	username := "testuser"
	email := "testuser@example.com"
	ipAddress := "192.168.1.1"
	userAgent := "Mozilla/5.0"

	err := notificationService.SendLoginNotification(
		userID,
		username,
		email,
		ipAddress,
		userAgent,
	)

	// Assert - проверяем, что уведомление отправлено без ошибок
	require.NoError(t, err, "Should send login notification successfully")
}

// TestNotificationService_Integration_SendLoginNotification_EmptyEmail тестирует отправку уведомления с пустым email
func TestNotificationService_Integration_SendLoginNotification_EmptyEmail(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	t.Log("=== START TestNotificationService_Integration_SendLoginNotification_EmptyEmail")

	// Arrange
	topic := getEnvOrDefault("KAFKA_TOPIC_NOTIFICATIONS", "notifications_test")
	producer := setupTestKafkaProducer(t, topic)
	if producer == nil {
		return
	}

	notificationService := services.NewNotificationServiceWithProducer(producer)

	// Act - отправка уведомления с пустым email (должно быть пропущено)
	userID := uuid.New()
	username := "testuser"
	email := ""
	ipAddress := "192.168.1.1"
	userAgent := "Mozilla/5.0"

	err := notificationService.SendLoginNotification(
		userID,
		username,
		email,
		ipAddress,
		userAgent,
	)

	// Assert - не должно быть ошибки, но уведомление не отправляется
	require.NoError(t, err, "Should handle empty email gracefully")
}

// TestNotificationService_Integration_NewNotificationService тестирует создание NotificationService с реальным Kafka
func TestNotificationService_Integration_NewNotificationService(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	t.Log("=== START TestNotificationService_Integration_NewNotificationService")

	// Arrange
	brokers := getKafkaBrokers()
	topic := getEnvOrDefault("KAFKA_TOPIC_NOTIFICATIONS", "notifications_test")

	kafkaConfig := &kafka.ProducerConfig{
		Brokers: brokers,
		Topic:   topic,
	}

	// Act - создание сервиса с реальным Kafka
	notificationService, err := services.NewNotificationService(kafkaConfig)

	// Assert
	if err != nil {
		t.Skipf("Kafka недоступен, пропускаем тест: %v", err)
		return
	}

	require.NoError(t, err)
	assert.NotNil(t, notificationService)

	// Проверяем, что сервис можно закрыть
	t.Cleanup(func() {
		if notificationService != nil {
			notificationService.Close()
		}
	})
}

// TestNotificationService_Integration_KafkaUnavailable тестирует обработку ошибки при недоступном Kafka
func TestNotificationService_Integration_KafkaUnavailable(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	t.Log("=== START TestNotificationService_Integration_KafkaUnavailable")

	// Arrange - подключаемся к несуществующему Kafka
	kafkaConfig := &kafka.ProducerConfig{
		Brokers: []string{"localhost:9999"}, // Несуществующий порт
		Topic:   "test-topic",
	}

	// Act - попытка создать сервис с недоступным Kafka
	notificationService, err := services.NewNotificationService(kafkaConfig)

	// Assert - должна быть ошибка
	require.Error(t, err)
	assert.Nil(t, notificationService)
	assert.Contains(t, err.Error(), "failed to create notification producer")
}

// TestNotificationService_Integration_Close тестирует закрытие соединения с Kafka
func TestNotificationService_Integration_Close(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	t.Log("=== START TestNotificationService_Integration_Close")

	// Arrange
	topic := getEnvOrDefault("KAFKA_TOPIC_NOTIFICATIONS", "notifications_test")
	producer := setupTestKafkaProducer(t, topic)
	if producer == nil {
		return
	}

	notificationService := services.NewNotificationServiceWithProducer(producer)

	// Act - закрытие сервиса
	err := notificationService.Close()

	// Assert
	require.NoError(t, err, "Should close notification service successfully")
}

// TestNotificationService_Integration_MultipleNotifications тестирует отправку нескольких уведомлений подряд
func TestNotificationService_Integration_MultipleNotifications(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	t.Log("=== START TestNotificationService_Integration_MultipleNotifications")

	// Arrange
	topic := getEnvOrDefault("KAFKA_TOPIC_NOTIFICATIONS", "notifications_test")
	producer := setupTestKafkaProducer(t, topic)
	if producer == nil {
		return
	}

	notificationService := services.NewNotificationServiceWithProducer(producer)

	// Act - отправка нескольких уведомлений
	for i := 0; i < 3; i++ {
		userID := uuid.New()
		username := "testuser" + string(rune('0'+i))
		email := "testuser" + string(rune('0'+i)) + "@example.com"
		ipAddress := "192.168.1." + string(rune('0'+i+1))
		userAgent := "Mozilla/5.0"

		err := notificationService.SendLoginNotification(
			userID,
			username,
			email,
			ipAddress,
			userAgent,
		)

		// Assert - каждое уведомление должно быть отправлено успешно
		require.NoError(t, err, "Should send notification %d successfully", i)
	}

	// Даем время Kafka обработать сообщения
	time.Sleep(1 * time.Second)
}
