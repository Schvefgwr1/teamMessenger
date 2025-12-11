//go:build integration
// +build integration

package integration

import (
	"common/config"
	"common/kafka"
	"common/models"
	"context"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
)

// logStep - единый helper для читаемых логов в тестах
func logStep(t *testing.T, format string, args ...interface{}) {
	t.Helper()
	t.Logf("==> "+format, args...)
}

// setupTestKafkaProducer создает Kafka producer для тестов
func setupTestKafkaProducer(t *testing.T, topic string) *kafka.NotificationProducer {
	brokers := getKafkaBrokers()

	logStep(t, "Создаем Kafka NotificationProducer brokers=%v topic=%s", brokers, topic)

	config := &kafka.ProducerConfig{
		Brokers: brokers,
		Topic:   topic,
	}

	producer, err := kafka.NewNotificationProducer(config)
	if err != nil {
		t.Skipf("Kafka недоступен, пропускаем тест: %v", err)
		return nil
	}

	safeProducer := &safeNotificationProducer{inner: producer}

	t.Cleanup(func() {
		if safeProducer != nil {
			_ = safeProducer.Close()
		}
	})

	return safeProducer.inner
}

// setupTestKafkaConsumer создает Kafka consumer для тестов
func setupTestKafkaConsumer(t *testing.T, topic string, groupID string, emailService interface{}) (sarama.ConsumerGroup, error) {
	brokers := getKafkaBrokers()

	logStep(t, "Создаем Kafka Consumer brokers=%v topic=%s groupID=%s", brokers, topic, groupID)

	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	config.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer group: %w", err)
	}

	t.Cleanup(func() {
		if consumer != nil {
			_ = consumer.Close()
		}
	})

	return consumer, nil
}

// setupTestEmailConfig создает тестовую конфигурацию для EmailService
func setupTestEmailConfig(t *testing.T) *config.EmailConfig {
	templatePath := getEnvOrDefault("TEMPLATE_PATH", "../../templates")
	smtpHost := getEnvOrDefault("SMTP_HOST", "localhost")
	smtpPort := getEnvOrDefaultInt("SMTP_PORT", 1025) // MailHog по умолчанию
	smtpUser := getEnvOrDefault("SMTP_USERNAME", "test@example.com")
	smtpPass := getEnvOrDefault("SMTP_PASSWORD", "testpassword")

	return &config.EmailConfig{
		SMTPHost:     smtpHost,
		SMTPPort:     smtpPort,
		Username:     smtpUser,
		Password:     smtpPass,
		FromName:     "Test TeamMessenger",
		FromEmail:    "noreply@test.teammessenger.com",
		TemplatePath: templatePath,
	}
}

// getKafkaBrokers получает список Kafka брокеров из переменных окружения
func getKafkaBrokers() []string {
	brokers := os.Getenv("KAFKA_BROKERS")
	if brokers == "" {
		return []string{"localhost:9092"} // значение по умолчанию
	}
	return []string{brokers}
}

// getEnvOrDefault получает значение переменной окружения или возвращает значение по умолчанию
func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// getEnvOrDefaultInt получает значение переменной окружения как int или возвращает значение по умолчанию
func getEnvOrDefaultInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	var intValue int
	n, err := fmt.Sscanf(value, "%d", &intValue)
	if err != nil || n != 1 {
		return defaultValue
	}
	return intValue
}

// safeNotificationProducer делает Close идемпотентным
type safeNotificationProducer struct {
	inner *kafka.NotificationProducer
	once  sync.Once
}

func (s *safeNotificationProducer) SendNotification(notification interface{}) error {
	return s.inner.SendNotification(notification)
}

func (s *safeNotificationProducer) Close() error {
	var err error
	s.once.Do(func() {
		err = s.inner.Close()
	})
	return err
}

// createTestNewTaskNotification создает тестовое уведомление о новой задаче
func createTestNewTaskNotification() *models.NewTaskNotification {
	return &models.NewTaskNotification{
		BaseNotification: models.BaseNotification{
			ID:        uuid.New(),
			Type:      models.NotificationNewTask,
			Email:     "executor@test.com",
			CreatedAt: time.Now(),
		},
		TaskID:      1,
		TaskTitle:   "Test Task Title",
		CreatorName: "Test Creator",
		ExecutorID:  uuid.New(),
	}
}

// createTestNewChatNotification создает тестовое уведомление о новом чате
func createTestNewChatNotification(isGroup bool) *models.NewChatNotification {
	return &models.NewChatNotification{
		BaseNotification: models.BaseNotification{
			ID:        uuid.New(),
			Type:      models.NotificationNewChat,
			Email:     "user@test.com",
			CreatedAt: time.Now(),
		},
		ChatID:      uuid.New(),
		ChatName:    "Test Chat",
		CreatorName: "Test Creator",
		IsGroup:     isGroup,
		Description: "Test Description",
	}
}

// createTestLoginNotification создает тестовое уведомление о входе
func createTestLoginNotification() *models.LoginNotification {
	return &models.LoginNotification{
		BaseNotification: models.BaseNotification{
			ID:        uuid.New(),
			Type:      models.NotificationLogin,
			Email:     "user@test.com",
			CreatedAt: time.Now(),
		},
		UserID:    uuid.New(),
		Username:  "testuser",
		IPAddress: "192.168.1.1",
		LoginTime: time.Now(),
		UserAgent: "Mozilla/5.0",
	}
}

// waitForKafkaMessage ждет получения сообщения из Kafka (для тестов consumer)
func waitForKafkaMessage(ctx context.Context, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			time.Sleep(100 * time.Millisecond)
		}
	}
	return nil
}
