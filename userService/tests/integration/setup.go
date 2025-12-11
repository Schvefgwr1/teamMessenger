//go:build integration
// +build integration

package integration

import (
	"common/config"
	"common/db"
	"common/kafka"
	"context"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"common/models"
	"userService/internal/services"
	"userService/internal/utils"

	"gorm.io/gorm"
)

// logStep - единый helper для читаемых логов в тестах
func logStep(t *testing.T, format string, args ...interface{}) {
	t.Helper()
	t.Logf("==> "+format, args...)
}

// setupTestDB создает подключение к тестовой PostgreSQL базе данных
func setupTestDB(t *testing.T) *gorm.DB {
	// Используем переменные окружения или значения по умолчанию
	dbHost := getEnvOrDefault("DB_HOST", "localhost")
	dbUser := getEnvOrDefault("DB_USER", "postgres")
	dbPassword := getEnvOrDefault("DB_PASSWORD", "postgres")
	dbName := getEnvOrDefault("DB_NAME", "team_messenger_test")
	dbPort := getEnvOrDefaultInt("DB_PORT", 5432)

	logStep(t, "Подключение к Postgres host=%s port=%d db=%s user=%s", dbHost, dbPort, dbName, dbUser)

	ensureTestKeys(t)

	cfg := &config.Config{}
	cfg.Database.Host = dbHost
	cfg.Database.User = dbUser
	cfg.Database.Password = dbPassword
	cfg.Database.Name = dbName
	cfg.Database.Port = dbPort

	database, err := db.InitDB(cfg)
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Проверяем подключение
	sqlDB, err := database.DB()
	if err != nil {
		t.Fatalf("Failed to get underlying sql.DB: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		t.Fatalf("Failed to ping database: %v", err)
	}

	// Очищаем тестовые данные перед тестом
	logStep(t, "Очистка тестовых данных в Postgres")
	cleanupTestData(t, database)

	t.Cleanup(func() {
		logStep(t, "Очистка после теста и закрытие соединения Postgres")
		cleanupTestData(t, database)
		sqlDB.Close()
	})

	return database
}

// ensureTestKeys гарантирует наличие RSA ключей для JWT
func ensureTestKeys(t *testing.T) {
	t.Helper()
	if _, err := utils.ExtractPublicKeyFromFile(); err == nil {
		return
	}

	logStep(t, "Генерируем тестовые RSA ключи")
	if _, err := utils.GenerateAndSaveNewKeys(); err != nil {
		t.Fatalf("Failed to generate test keys: %v", err)
	}
}

// cleanupTestData очищает тестовые данные из базы
func cleanupTestData(t *testing.T, db *gorm.DB) {
	// Удаляем тестовых пользователей (сначала удаляем связи)
	db.Exec("DELETE FROM user_service.users WHERE username LIKE 'test_%' OR email LIKE 'test_%'")
	// Очищаем другие тестовые данные при необходимости
}

// setupTestKafkaProducer создает Kafka producer для тестов
func setupTestKafkaProducer(t *testing.T, topic string) services.NotificationProducerInterface {
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

	return safeProducer
}

// setupTestKafkaKeyUpdateProducer создает Kafka producer для обновлений ключей
func setupTestKafkaKeyUpdateProducer(t *testing.T, topic string) (services.KeyUpdateProducerInterface, error) {
	brokers := getKafkaBrokers()

	logStep(t, "Создаем Kafka KeyUpdateProducer brokers=%v topic=%s", brokers, topic)

	config := &kafka.ProducerConfig{
		Brokers: brokers,
		Topic:   topic,
	}

	producer, err := kafka.NewKeyUpdateProducer(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create key update producer: %w", err)
	}

	safeProducer := &safeKeyUpdateProducer{inner: producer}

	t.Cleanup(func() {
		if safeProducer != nil {
			_ = safeProducer.Close()
		}
	})

	return safeProducer, nil
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

// safeKeyUpdateProducer делает Close идемпотентным
type safeKeyUpdateProducer struct {
	inner *kafka.KeyUpdateProducer
	once  sync.Once
}

func (s *safeKeyUpdateProducer) SendKeyUpdate(keyUpdate *models.PublicKeyUpdate) error {
	return s.inner.SendKeyUpdate(keyUpdate)
}

func (s *safeKeyUpdateProducer) Close() error {
	var err error
	s.once.Do(func() {
		err = s.inner.Close()
	})
	return err
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
