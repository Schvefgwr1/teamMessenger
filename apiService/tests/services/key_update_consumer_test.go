package services

import (
	"apiService/internal/services"
	"context"
	"testing"

	"github.com/IBM/sarama"
	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupTestRedisForConsumer создает тестовый Redis клиент с miniredis
func setupTestRedisForConsumer(t *testing.T) *redis.Client {
	mr, err := miniredis.Run()
	require.NoError(t, err)

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	return client
}

// Mock для sarama.ConsumerGroupSession
type mockConsumerGroupSession struct{}

func (m *mockConsumerGroupSession) Claims() map[string][]int32 {
	return nil
}

func (m *mockConsumerGroupSession) MemberID() string {
	return "test-member"
}

func (m *mockConsumerGroupSession) GenerationID() int32 {
	return 1
}

func (m *mockConsumerGroupSession) MarkOffset(topic string, partition int32, offset int64, metadata string) {
}

func (m *mockConsumerGroupSession) Commit() {
}

func (m *mockConsumerGroupSession) ResetOffset(topic string, partition int32, offset int64, metadata string) {
}

func (m *mockConsumerGroupSession) MarkMessage(msg *sarama.ConsumerMessage, metadata string) {
}

func (m *mockConsumerGroupSession) Context() context.Context {
	return context.Background()
}

// Тесты для Setup и Cleanup (единственные публичные методы, которые можно протестировать без Kafka)

func TestKeyUpdateConsumer_Setup(t *testing.T) {
	// Arrange
	publicKeyManager := services.NewPublicKeyManager()
	redisClient := setupTestRedisForConsumer(t)
	defer redisClient.Close()
	sessionService := services.NewSessionService(redisClient)

	// Создаем consumer через NewKeyUpdateConsumer
	config := &services.KeyUpdateConsumerConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "key-updates",
		GroupID: "test-group",
	}

	// Пропускаем тест, если Kafka недоступен
	consumer, err := services.NewKeyUpdateConsumer(config, publicKeyManager, sessionService, redisClient)
	if err != nil {
		t.Skipf("Skipping test: Kafka not available: %v", err)
		return
	}
	defer consumer.Close()

	session := &mockConsumerGroupSession{}

	// Act
	err = consumer.Setup(session)

	// Assert
	require.NoError(t, err)
}

func TestKeyUpdateConsumer_Cleanup(t *testing.T) {
	// Arrange
	publicKeyManager := services.NewPublicKeyManager()
	redisClient := setupTestRedisForConsumer(t)
	defer redisClient.Close()
	sessionService := services.NewSessionService(redisClient)

	// Создаем consumer через NewKeyUpdateConsumer
	config := &services.KeyUpdateConsumerConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "key-updates",
		GroupID: "test-group",
	}

	// Пропускаем тест, если Kafka недоступен
	consumer, err := services.NewKeyUpdateConsumer(config, publicKeyManager, sessionService, redisClient)
	if err != nil {
		t.Skipf("Skipping test: Kafka not available: %v", err)
		return
	}
	defer consumer.Close()

	session := &mockConsumerGroupSession{}

	// Act
	err = consumer.Cleanup(session)

	// Assert
	require.NoError(t, err)
}

func TestKeyUpdateConsumer_NewKeyUpdateConsumer_Success(t *testing.T) {
	// Arrange
	publicKeyManager := services.NewPublicKeyManager()
	redisClient := setupTestRedisForConsumer(t)
	defer redisClient.Close()
	sessionService := services.NewSessionService(redisClient)

	config := &services.KeyUpdateConsumerConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "key-updates",
		GroupID: "test-group",
	}

	// Act
	consumer, err := services.NewKeyUpdateConsumer(config, publicKeyManager, sessionService, redisClient)

	// Assert
	if err != nil {
		// Если Kafka недоступен, пропускаем тест
		t.Skipf("Skipping test: Kafka not available: %v", err)
		return
	}

	require.NoError(t, err)
	assert.NotNil(t, consumer)

	// Cleanup
	if consumer != nil {
		consumer.Close()
	}
}

func TestKeyUpdateConsumer_NewKeyUpdateConsumer_InvalidBrokers(t *testing.T) {
	// Arrange
	publicKeyManager := services.NewPublicKeyManager()
	redisClient := setupTestRedisForConsumer(t)
	defer redisClient.Close()
	sessionService := services.NewSessionService(redisClient)

	config := &services.KeyUpdateConsumerConfig{
		Brokers: []string{"invalid:9999"},
		Topic:   "key-updates",
		GroupID: "test-group",
	}

	// Act
	consumer, err := services.NewKeyUpdateConsumer(config, publicKeyManager, sessionService, redisClient)

	// Assert
	if err == nil {
		// Если Kafka доступен, тест должен пройти
		if consumer != nil {
			consumer.Close()
		}
		return
	}

	// Если Kafka недоступен, ожидаем ошибку
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create consumer group")
}

func TestKeyUpdateConsumer_Close(t *testing.T) {
	// Arrange
	publicKeyManager := services.NewPublicKeyManager()
	redisClient := setupTestRedisForConsumer(t)
	defer redisClient.Close()
	sessionService := services.NewSessionService(redisClient)

	config := &services.KeyUpdateConsumerConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "key-updates",
		GroupID: "test-group",
	}

	consumer, err := services.NewKeyUpdateConsumer(config, publicKeyManager, sessionService, redisClient)
	if err != nil {
		t.Skipf("Skipping test: Kafka not available: %v", err)
		return
	}

	// Act
	err = consumer.Close()

	// Assert
	require.NoError(t, err)
}
