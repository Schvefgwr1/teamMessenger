package main

import (
	"common/config"
	"common/kafka"
	"context"
	"github.com/joho/godotenv"
	"log"
	"notificationService/internal/services"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Загружаем переменные окружения из .env файла (если существует)
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found or error loading .env file: %v", err)
	}

	// Загружаем конфигурацию
	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Apply environment variable overrides
	config.ApplyAppEnvOverrides(cfg)
	config.ApplyKafkaEnvOverrides(cfg)
	config.ApplyEmailEnvOverrides(cfg)

	// Логируем email конфигурацию (без пароля)
	log.Printf("Email config: Host=%s, Port=%d, Username=%s, FromEmail=%s",
		cfg.Email.SMTPHost, cfg.Email.SMTPPort, cfg.Email.Username, cfg.Email.FromEmail)

	// Инициализируем Email Service
	emailService, err := services.NewEmailService(&cfg.Email)
	if err != nil {
		log.Fatalf("Failed to initialize email service: %v", err)
	}

	// Инициализируем Kafka Consumer
	notificationsConsumerConfig := &services.KeyUpdateConsumerConfig{
		Brokers: kafka.GetKafkaBrokers(),
		Topic:   kafka.GetNotificationsTopic(),
		GroupID: cfg.Kafka.GroupID,
	}

	kafkaConsumer, err := services.NewKafkaConsumer(notificationsConsumerConfig, emailService)
	if err != nil {
		log.Fatalf("Failed to initialize kafka consumer: %v", err)
	}

	// Контекст для graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Запускаем Kafka consumer в горутине
	go func() {
		log.Println("Starting Kafka consumer...")
		if err := kafkaConsumer.Start(ctx); err != nil {
			log.Printf("Kafka consumer error: %v", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	log.Println("Notification Service is running. Press Ctrl+C to stop.")
	<-sigChan

	log.Println("Shutting down Notification Service...")

	// Graceful shutdown
	cancel()

	// Даем время на завершение обработки текущих сообщений
	time.Sleep(5 * time.Second)

	// Закрываем Kafka consumer
	if err := kafkaConsumer.Close(); err != nil {
		log.Printf("Error closing kafka consumer: %v", err)
	}

	log.Println("Notification Service stopped.")
}
