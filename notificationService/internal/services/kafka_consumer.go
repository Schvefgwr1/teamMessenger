package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"common/models"
	"github.com/IBM/sarama"
)

type KeyUpdateConsumerConfig struct {
	Brokers []string
	Topic   string
	GroupID string
}

type KafkaConsumer struct {
	consumer     sarama.ConsumerGroup
	emailService *EmailService
	config       *KeyUpdateConsumerConfig
}

func NewKafkaConsumer(cfg *KeyUpdateConsumerConfig, emailService *EmailService) (*KafkaConsumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	config.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumerGroup(cfg.Brokers, cfg.GroupID, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer group: %w", err)
	}

	return &KafkaConsumer{
		consumer:     consumer,
		emailService: emailService,
		config:       cfg,
	}, nil
}

func (kc *KafkaConsumer) Start(ctx context.Context) error {
	topics := []string{kc.config.Topic}

	// Горутина для обработки ошибок
	go func() {
		for err := range kc.consumer.Errors() {
			log.Printf("Kafka consumer error: %v", err)
		}
	}()

	// Основной цикл обработки сообщений
	for {
		select {
		case <-ctx.Done():
			log.Println("Kafka consumer context cancelled")
			return ctx.Err()
		default:
			if err := kc.consumer.Consume(ctx, topics, kc); err != nil {
				log.Printf("Error from consumer: %v", err)
				time.Sleep(5 * time.Second) // Пауза перед переподключением
			}
		}
	}
}

func (kc *KafkaConsumer) Close() error {
	return kc.consumer.Close()
}

// Реализация интерфейса sarama.ConsumerGroupHandler

func (kc *KafkaConsumer) Setup(sarama.ConsumerGroupSession) error {
	log.Println("Kafka consumer setup")
	return nil
}

func (kc *KafkaConsumer) Cleanup(sarama.ConsumerGroupSession) error {
	log.Println("Kafka consumer cleanup")
	return nil
}

func (kc *KafkaConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message := <-claim.Messages():
			if message == nil {
				return nil
			}

			log.Printf("Received message: topic=%s partition=%d offset=%d",
				message.Topic, message.Partition, message.Offset)

			if err := kc.processMessage(message); err != nil {
				log.Printf("Error processing message: %v", err)
				// Продолжаем обработку, не прерывая процесс
			}

			// Помечаем сообщение как обработанное
			session.MarkMessage(message, "")

		case <-session.Context().Done():
			return nil
		}
	}
}

func (kc *KafkaConsumer) processMessage(message *sarama.ConsumerMessage) error {
	var kafkaMsg models.KafkaMessage
	if err := json.Unmarshal(message.Value, &kafkaMsg); err != nil {
		return fmt.Errorf("failed to unmarshal kafka message: %w", err)
	}

	notification, err := kc.parseNotification(kafkaMsg)
	if err != nil {
		return fmt.Errorf("failed to parse notification: %w", err)
	}

	if err := kc.emailService.SendNotification(notification); err != nil {
		return fmt.Errorf("failed to send email notification: %w", err)
	}

	log.Printf("Successfully processed notification of type: %s", kafkaMsg.Type)
	return nil
}

func (kc *KafkaConsumer) parseNotification(kafkaMsg models.KafkaMessage) (interface{}, error) {
	// Преобразуем payload обратно в JSON для парсинга в конкретную структуру
	payloadBytes, err := json.Marshal(kafkaMsg.Payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	switch kafkaMsg.Type {
	case models.NotificationNewTask:
		var notification models.NewTaskNotification
		if err := json.Unmarshal(payloadBytes, &notification); err != nil {
			return nil, fmt.Errorf("failed to unmarshal new task notification: %w", err)
		}
		return &notification, nil

	case models.NotificationNewChat:
		var notification models.NewChatNotification
		if err := json.Unmarshal(payloadBytes, &notification); err != nil {
			return nil, fmt.Errorf("failed to unmarshal new chat notification: %w", err)
		}
		return &notification, nil

	case models.NotificationLogin:
		var notification models.LoginNotification
		if err := json.Unmarshal(payloadBytes, &notification); err != nil {
			return nil, fmt.Errorf("failed to unmarshal login notification: %w", err)
		}
		return &notification, nil

	default:
		return nil, fmt.Errorf("unknown notification type: %s", kafkaMsg.Type)
	}
}
