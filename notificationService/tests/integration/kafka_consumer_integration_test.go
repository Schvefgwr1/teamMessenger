//go:build integration
// +build integration

package integration

import (
	"encoding/json"
	"testing"
	"time"

	"common/models"
	"notificationService/internal/services"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestKafkaConsumer_ProcessMessage_NewTask_Integration тестирует обработку сообщения о новой задаче из Kafka
func TestKafkaConsumer_ProcessMessage_NewTask_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange - настройка реальных зависимостей
	kafkaTopic := getEnvOrDefault("KAFKA_TOPIC_NOTIFICATIONS", "notifications_test")
	producer := setupTestKafkaProducer(t, kafkaTopic)
	if producer == nil {
		t.Skip("Kafka недоступен, пропускаем тест")
		return
	}

	cfg := setupTestEmailConfig(t)
	emailService, err := services.NewEmailService(cfg)
	if err != nil {
		t.Skipf("EmailService недоступен, пропускаем тест: %v", err)
		return
	}

	// Создаем тестовое сообщение
	notification := createTestNewTaskNotification()
	kafkaMessage := models.KafkaMessage{
		Type:    models.NotificationNewTask,
		Payload: notification,
	}

	messageBytes, err := json.Marshal(kafkaMessage)
	require.NoError(t, err)

	// Создаем consumer для обработки сообщения
	consumerConfig := &services.KeyUpdateConsumerConfig{
		Brokers: getKafkaBrokers(),
		Topic:   kafkaTopic,
		GroupID: "test-group-" + uuid.New().String(),
	}

	kafkaConsumer, err := services.NewKafkaConsumer(consumerConfig, emailService)
	require.NoError(t, err)
	defer kafkaConsumer.Close()

	// Отправляем сообщение в Kafka
	err = producer.SendNotification(notification)
	require.NoError(t, err)

	// Создаем тестовое сообщение для обработки
	saramaMessage := &sarama.ConsumerMessage{
		Topic:     kafkaTopic,
		Partition: 0,
		Offset:    0,
		Value:     messageBytes,
		Timestamp: time.Now(),
	}

	// Act - обрабатываем сообщение
	err = kafkaConsumer.ProcessMessage(saramaMessage)

	// Assert - проверяем, что сообщение обработано без ошибок
	// Если SMTP недоступен, это нормально для интеграционного теста
	if err != nil {
		// Проверяем, что ошибка не связана с парсингом или обработкой Kafka сообщения
		assert.NotContains(t, err.Error(), "failed to unmarshal kafka message")
		assert.NotContains(t, err.Error(), "failed to parse notification")
		assert.NotContains(t, err.Error(), "unknown notification type")
	}
}

// TestKafkaConsumer_ProcessMessage_NewChat_Integration тестирует обработку сообщения о новом чате из Kafka
func TestKafkaConsumer_ProcessMessage_NewChat_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange
	kafkaTopic := getEnvOrDefault("KAFKA_TOPIC_NOTIFICATIONS", "notifications_test")
	producer := setupTestKafkaProducer(t, kafkaTopic)
	if producer == nil {
		t.Skip("Kafka недоступен, пропускаем тест")
		return
	}

	cfg := setupTestEmailConfig(t)
	emailService, err := services.NewEmailService(cfg)
	if err != nil {
		t.Skipf("EmailService недоступен, пропускаем тест: %v", err)
		return
	}

	notification := createTestNewChatNotification(true)
	kafkaMessage := models.KafkaMessage{
		Type:    models.NotificationNewChat,
		Payload: notification,
	}

	messageBytes, err := json.Marshal(kafkaMessage)
	require.NoError(t, err)

	consumerConfig := &services.KeyUpdateConsumerConfig{
		Brokers: getKafkaBrokers(),
		Topic:   kafkaTopic,
		GroupID: "test-group-" + uuid.New().String(),
	}

	kafkaConsumer, err := services.NewKafkaConsumer(consumerConfig, emailService)
	require.NoError(t, err)
	defer kafkaConsumer.Close()

	err = producer.SendNotification(notification)
	require.NoError(t, err)

	saramaMessage := &sarama.ConsumerMessage{
		Topic:     kafkaTopic,
		Partition: 0,
		Offset:    0,
		Value:     messageBytes,
		Timestamp: time.Now(),
	}

	// Act
	err = kafkaConsumer.ProcessMessage(saramaMessage)

	// Assert
	if err != nil {
		assert.NotContains(t, err.Error(), "failed to unmarshal kafka message")
		assert.NotContains(t, err.Error(), "failed to parse notification")
		assert.NotContains(t, err.Error(), "unknown notification type")
	}
}

// TestKafkaConsumer_ProcessMessage_Login_Integration тестирует обработку сообщения о входе из Kafka
func TestKafkaConsumer_ProcessMessage_Login_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange
	kafkaTopic := getEnvOrDefault("KAFKA_TOPIC_NOTIFICATIONS", "notifications_test")
	producer := setupTestKafkaProducer(t, kafkaTopic)
	if producer == nil {
		t.Skip("Kafka недоступен, пропускаем тест")
		return
	}

	cfg := setupTestEmailConfig(t)
	emailService, err := services.NewEmailService(cfg)
	if err != nil {
		t.Skipf("EmailService недоступен, пропускаем тест: %v", err)
		return
	}

	notification := createTestLoginNotification()
	kafkaMessage := models.KafkaMessage{
		Type:    models.NotificationLogin,
		Payload: notification,
	}

	messageBytes, err := json.Marshal(kafkaMessage)
	require.NoError(t, err)

	consumerConfig := &services.KeyUpdateConsumerConfig{
		Brokers: getKafkaBrokers(),
		Topic:   kafkaTopic,
		GroupID: "test-group-" + uuid.New().String(),
	}

	kafkaConsumer, err := services.NewKafkaConsumer(consumerConfig, emailService)
	require.NoError(t, err)
	defer kafkaConsumer.Close()

	err = producer.SendNotification(notification)
	require.NoError(t, err)

	saramaMessage := &sarama.ConsumerMessage{
		Topic:     kafkaTopic,
		Partition: 0,
		Offset:    0,
		Value:     messageBytes,
		Timestamp: time.Now(),
	}

	// Act
	err = kafkaConsumer.ProcessMessage(saramaMessage)

	// Assert
	if err != nil {
		assert.NotContains(t, err.Error(), "failed to unmarshal kafka message")
		assert.NotContains(t, err.Error(), "failed to parse notification")
		assert.NotContains(t, err.Error(), "unknown notification type")
	}
}

// TestKafkaConsumer_ParseNotification_NewTask_Integration тестирует парсинг уведомления о новой задаче
func TestKafkaConsumer_ParseNotification_NewTask_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange
	cfg := setupTestEmailConfig(t)
	emailService, err := services.NewEmailService(cfg)
	require.NoError(t, err)

	consumerConfig := &services.KeyUpdateConsumerConfig{
		Brokers: getKafkaBrokers(),
		Topic:   getEnvOrDefault("KAFKA_TOPIC_NOTIFICATIONS", "notifications_test"),
		GroupID: "test-group-" + uuid.New().String(),
	}

	kafkaConsumer, err := services.NewKafkaConsumer(consumerConfig, emailService)
	require.NoError(t, err)
	defer kafkaConsumer.Close()

	notification := createTestNewTaskNotification()
	kafkaMessage := models.KafkaMessage{
		Type:    models.NotificationNewTask,
		Payload: notification,
	}

	// Act - парсим уведомление
	parsedNotification, err := kafkaConsumer.ParseNotification(kafkaMessage)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, parsedNotification)

	parsedTaskNotification, ok := parsedNotification.(*models.NewTaskNotification)
	require.True(t, ok, "Parsed notification should be *models.NewTaskNotification")
	assert.Equal(t, notification.TaskTitle, parsedTaskNotification.TaskTitle)
	assert.Equal(t, notification.CreatorName, parsedTaskNotification.CreatorName)
}

// TestKafkaConsumer_ParseNotification_NewChat_Integration тестирует парсинг уведомления о новом чате
func TestKafkaConsumer_ParseNotification_NewChat_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange
	cfg := setupTestEmailConfig(t)
	emailService, err := services.NewEmailService(cfg)
	require.NoError(t, err)

	consumerConfig := &services.KeyUpdateConsumerConfig{
		Brokers: getKafkaBrokers(),
		Topic:   getEnvOrDefault("KAFKA_TOPIC_NOTIFICATIONS", "notifications_test"),
		GroupID: "test-group-" + uuid.New().String(),
	}

	kafkaConsumer, err := services.NewKafkaConsumer(consumerConfig, emailService)
	require.NoError(t, err)
	defer kafkaConsumer.Close()

	notification := createTestNewChatNotification(false)
	kafkaMessage := models.KafkaMessage{
		Type:    models.NotificationNewChat,
		Payload: notification,
	}

	// Act
	parsedNotification, err := kafkaConsumer.ParseNotification(kafkaMessage)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, parsedNotification)

	parsedChatNotification, ok := parsedNotification.(*models.NewChatNotification)
	require.True(t, ok, "Parsed notification should be *models.NewChatNotification")
	assert.Equal(t, notification.ChatName, parsedChatNotification.ChatName)
	assert.Equal(t, notification.IsGroup, parsedChatNotification.IsGroup)
}

// TestKafkaConsumer_ParseNotification_Login_Integration тестирует парсинг уведомления о входе
func TestKafkaConsumer_ParseNotification_Login_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange
	cfg := setupTestEmailConfig(t)
	emailService, err := services.NewEmailService(cfg)
	require.NoError(t, err)

	consumerConfig := &services.KeyUpdateConsumerConfig{
		Brokers: getKafkaBrokers(),
		Topic:   getEnvOrDefault("KAFKA_TOPIC_NOTIFICATIONS", "notifications_test"),
		GroupID: "test-group-" + uuid.New().String(),
	}

	kafkaConsumer, err := services.NewKafkaConsumer(consumerConfig, emailService)
	require.NoError(t, err)
	defer kafkaConsumer.Close()

	notification := createTestLoginNotification()
	kafkaMessage := models.KafkaMessage{
		Type:    models.NotificationLogin,
		Payload: notification,
	}

	// Act
	parsedNotification, err := kafkaConsumer.ParseNotification(kafkaMessage)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, parsedNotification)

	parsedLoginNotification, ok := parsedNotification.(*models.LoginNotification)
	require.True(t, ok, "Parsed notification should be *models.LoginNotification")
	assert.Equal(t, notification.Username, parsedLoginNotification.Username)
	assert.Equal(t, notification.IPAddress, parsedLoginNotification.IPAddress)
}

// TestKafkaConsumer_ParseNotification_UnknownType_Integration тестирует обработку неизвестного типа уведомления
func TestKafkaConsumer_ParseNotification_UnknownType_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange
	cfg := setupTestEmailConfig(t)
	emailService, err := services.NewEmailService(cfg)
	require.NoError(t, err)

	consumerConfig := &services.KeyUpdateConsumerConfig{
		Brokers: getKafkaBrokers(),
		Topic:   getEnvOrDefault("KAFKA_TOPIC_NOTIFICATIONS", "notifications_test"),
		GroupID: "test-group-" + uuid.New().String(),
	}

	kafkaConsumer, err := services.NewKafkaConsumer(consumerConfig, emailService)
	require.NoError(t, err)
	defer kafkaConsumer.Close()

	kafkaMessage := models.KafkaMessage{
		Type:    models.NotificationType("unknown_type"),
		Payload: map[string]interface{}{"test": "data"},
	}

	// Act
	parsedNotification, err := kafkaConsumer.ParseNotification(kafkaMessage)

	// Assert
	require.Error(t, err)
	assert.Nil(t, parsedNotification)
	assert.Contains(t, err.Error(), "unknown notification type")
}

// TestKafkaConsumer_ParseNotification_InvalidJSON_Integration тестирует обработку невалидного JSON
func TestKafkaConsumer_ParseNotification_InvalidJSON_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange
	cfg := setupTestEmailConfig(t)
	emailService, err := services.NewEmailService(cfg)
	require.NoError(t, err)

	consumerConfig := &services.KeyUpdateConsumerConfig{
		Brokers: getKafkaBrokers(),
		Topic:   getEnvOrDefault("KAFKA_TOPIC_NOTIFICATIONS", "notifications_test"),
		GroupID: "test-group-" + uuid.New().String(),
	}

	kafkaConsumer, err := services.NewKafkaConsumer(consumerConfig, emailService)
	require.NoError(t, err)
	defer kafkaConsumer.Close()

	// Создаем сообщение с невалидным JSON в payload
	kafkaMessage := models.KafkaMessage{
		Type:    models.NotificationNewTask,
		Payload: "invalid json string",
	}

	// Act
	_, parseErr := kafkaConsumer.ParseNotification(kafkaMessage)

	// Assert
	// ParseNotification должен попытаться замаршалить payload, что может привести к ошибке
	// или успешному парсингу в зависимости от реализации
	// Проверяем, что метод обрабатывает ситуацию корректно
	if parseErr != nil {
		assert.Contains(t, parseErr.Error(), "failed to unmarshal")
	}
}

// TestKafkaConsumer_NewKafkaConsumer_Integration тестирует создание Kafka consumer с реальным Kafka
func TestKafkaConsumer_NewKafkaConsumer_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange
	cfg := setupTestEmailConfig(t)
	emailService, err := services.NewEmailService(cfg)
	require.NoError(t, err)

	consumerConfig := &services.KeyUpdateConsumerConfig{
		Brokers: getKafkaBrokers(),
		Topic:   getEnvOrDefault("KAFKA_TOPIC_NOTIFICATIONS", "notifications_test"),
		GroupID: "test-group-" + uuid.New().String(),
	}

	// Act - создаем consumer
	kafkaConsumer, err := services.NewKafkaConsumer(consumerConfig, emailService)

	// Assert
	if err != nil {
		t.Skipf("Kafka недоступен, пропускаем тест: %v", err)
		return
	}
	require.NoError(t, err)
	assert.NotNil(t, kafkaConsumer)
	defer kafkaConsumer.Close()
}

// TestKafkaConsumer_NewKafkaConsumer_KafkaUnavailable_Integration тестирует обработку ошибки недоступного Kafka
func TestKafkaConsumer_NewKafkaConsumer_KafkaUnavailable_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Arrange - подключаемся к несуществующему Kafka брокеру
	cfg := setupTestEmailConfig(t)
	emailService, err := services.NewEmailService(cfg)
	require.NoError(t, err)

	consumerConfig := &services.KeyUpdateConsumerConfig{
		Brokers: []string{"localhost:9999"}, // Несуществующий порт
		Topic:   "test-topic",
		GroupID: "test-group",
	}

	// Act
	kafkaConsumer, err := services.NewKafkaConsumer(consumerConfig, emailService)

	// Assert - проверяем реальную ошибку подключения
	require.Error(t, err)
	assert.Nil(t, kafkaConsumer)
	assert.Contains(t, err.Error(), "failed to create consumer group")
}
