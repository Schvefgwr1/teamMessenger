package main

import (
	"context"
	"gopkg.in/yaml.v3"
	"log"
	"notificationService/internal/config"
	"notificationService/internal/services"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Загружаем конфигурацию
	cfg, err := loadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Инициализируем Email Service
	emailService, err := services.NewEmailService(&cfg.Email)
	if err != nil {
		log.Fatalf("Failed to initialize email service: %v", err)
	}

	// Инициализируем Kafka Consumer
	kafkaConsumer, err := services.NewKafkaConsumer(&cfg.Kafka, emailService)
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

func loadConfig() (*config.NotificationConfig, error) {
	return loadConfigFromFile("config/config.yaml")
}

func loadConfigFromFile(path string) (*config.NotificationConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var notificationConfig config.NotificationConfig
	if err := yaml.Unmarshal(data, &notificationConfig); err != nil {
		return nil, err
	}

	return &notificationConfig, nil
}
