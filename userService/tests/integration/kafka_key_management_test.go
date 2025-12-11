//go:build integration
// +build integration

package integration

import (
	"testing"
	"time"

	"common/kafka"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"userService/internal/services"
)

// TestKeyManagementService_Integration_RegenerateKeys тестирует генерацию и отправку ключей в реальный Kafka
func TestKeyManagementService_Integration_RegenerateKeys(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	t.Log("=== START TestKeyManagementService_Integration_RegenerateKeys")

	// Arrange - настройка реального Kafka producer для обновлений ключей
	topic := getEnvOrDefault("KAFKA_TOPIC_KEY_UPDATES", "key_updates_test")

	producer, err := setupTestKafkaKeyUpdateProducer(t, topic)
	if err != nil {
		t.Skipf("Kafka недоступен, пропускаем тест: %v", err)
		return
	}

	keyManagementService := services.NewKeyManagementServiceWithProducer(producer, 1)

	// Act - генерация и отправка новых ключей
	err = keyManagementService.RegenerateKeys()

	// Assert - проверяем, что ключи сгенерированы и отправлены
	require.NoError(t, err, "Should regenerate keys successfully")

	// Проверяем, что версия ключа увеличилась
	version := keyManagementService.GetCurrentKeyVersion()
	assert.Equal(t, 2, version, "Key version should be incremented")
}

// TestKeyManagementService_Integration_NewKeyManagementService тестирует создание KeyManagementService с реальным Kafka
func TestKeyManagementService_Integration_NewKeyManagementService(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	t.Log("=== START TestKeyManagementService_Integration_NewKeyManagementService")

	// Arrange
	topic := getEnvOrDefault("KAFKA_TOPIC_KEY_UPDATES", "key_updates_test")

	kafkaConfig := &kafka.ProducerConfig{
		Brokers: getKafkaBrokers(),
		Topic:   topic,
	}

	// Act - создание сервиса с реальным Kafka
	keyManagementService, err := services.NewKeyManagementService(kafkaConfig)

	// Assert
	if err != nil {
		t.Skipf("Kafka недоступен, пропускаем тест: %v", err)
		return
	}

	require.NoError(t, err)
	assert.NotNil(t, keyManagementService)

	// Проверяем начальную версию ключа
	version := keyManagementService.GetCurrentKeyVersion()
	assert.Equal(t, 1, version, "Initial key version should be 1")

	// Проверяем, что сервис можно закрыть
	t.Cleanup(func() {
		if keyManagementService != nil {
			keyManagementService.Close()
		}
	})
}

// TestKeyManagementService_Integration_KafkaUnavailable тестирует обработку ошибки при недоступном Kafka
func TestKeyManagementService_Integration_KafkaUnavailable(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	t.Log("=== START TestKeyManagementService_Integration_KafkaUnavailable")

	// Arrange - подключаемся к несуществующему Kafka
	kafkaConfig := &kafka.ProducerConfig{
		Brokers: []string{"localhost:9999"}, // Несуществующий порт
		Topic:   "test-topic",
	}

	// Act - попытка создать сервис с недоступным Kafka
	keyManagementService, err := services.NewKeyManagementService(kafkaConfig)

	// Assert - должна быть ошибка
	require.Error(t, err)
	assert.Nil(t, keyManagementService)
	assert.Contains(t, err.Error(), "failed to create key update producer")
}

// TestKeyManagementService_Integration_MultipleKeyRegenerations тестирует несколько генераций ключей подряд
func TestKeyManagementService_Integration_MultipleKeyRegenerations(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	t.Log("=== START TestKeyManagementService_Integration_MultipleKeyRegenerations")

	// Arrange
	topic := getEnvOrDefault("KAFKA_TOPIC_KEY_UPDATES", "key_updates_test")

	producer, err := setupTestKafkaKeyUpdateProducer(t, topic)
	if err != nil {
		t.Skipf("Kafka недоступен, пропускаем тест: %v", err)
		return
	}

	keyManagementService := services.NewKeyManagementServiceWithProducer(producer, 1)

	// Act - несколько генераций ключей
	initialVersion := keyManagementService.GetCurrentKeyVersion()
	assert.Equal(t, 1, initialVersion)

	// Первая генерация
	err = keyManagementService.RegenerateKeys()
	require.NoError(t, err)
	assert.Equal(t, 2, keyManagementService.GetCurrentKeyVersion())

	// Вторая генерация
	err = keyManagementService.RegenerateKeys()
	require.NoError(t, err)
	assert.Equal(t, 3, keyManagementService.GetCurrentKeyVersion())

	// Третья генерация
	err = keyManagementService.RegenerateKeys()
	require.NoError(t, err)
	assert.Equal(t, 4, keyManagementService.GetCurrentKeyVersion())

	// Даем время Kafka обработать сообщения
	time.Sleep(1 * time.Second)
}

// TestKeyManagementService_Integration_Close тестирует закрытие соединения с Kafka
func TestKeyManagementService_Integration_Close(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	t.Log("=== START TestKeyManagementService_Integration_Close")

	// Arrange
	topic := getEnvOrDefault("KAFKA_TOPIC_KEY_UPDATES", "key_updates_test")

	producer, err := setupTestKafkaKeyUpdateProducer(t, topic)
	if err != nil {
		t.Skipf("Kafka недоступен, пропускаем тест: %v", err)
		return
	}

	keyManagementService := services.NewKeyManagementServiceWithProducer(producer, 1)

	// Act - закрытие сервиса
	err = keyManagementService.Close()

	// Assert
	require.NoError(t, err, "Should close key management service successfully")
}

// TestKeyManagementService_Integration_KeyVersionIncrement тестирует правильное увеличение версии ключа
func TestKeyManagementService_Integration_KeyVersionIncrement(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	t.Log("=== START TestKeyManagementService_Integration_KeyVersionIncrement")

	// Arrange
	topic := getEnvOrDefault("KAFKA_TOPIC_KEY_UPDATES", "key_updates_test")

	producer, err := setupTestKafkaKeyUpdateProducer(t, topic)
	if err != nil {
		t.Skipf("Kafka недоступен, пропускаем тест: %v", err)
		return
	}

	initialVersion := 5
	keyManagementService := services.NewKeyManagementServiceWithProducer(producer, initialVersion)

	// Act - генерация ключей
	err = keyManagementService.RegenerateKeys()
	require.NoError(t, err)

	// Assert - проверяем, что версия увеличилась
	newVersion := keyManagementService.GetCurrentKeyVersion()
	assert.Equal(t, initialVersion+1, newVersion, "Key version should be incremented by 1")
}
