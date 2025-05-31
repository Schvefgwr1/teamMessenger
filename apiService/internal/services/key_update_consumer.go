package services

import (
	"common/models"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/IBM/sarama"
	"github.com/redis/go-redis/v9"
)

type KeyUpdateConsumer struct {
	consumer         sarama.ConsumerGroup
	publicKeyManager *PublicKeyManager
	sessionService   *SessionService
	redisClient      *redis.Client
	brokers          []string
	topic            string
	groupID          string
}

type KeyUpdateConsumerConfig struct {
	Brokers []string
	Topic   string
	GroupID string
}

func NewKeyUpdateConsumer(
	config *KeyUpdateConsumerConfig,
	publicKeyManager *PublicKeyManager,
	sessionService *SessionService,
	redisClient *redis.Client,
) (*KeyUpdateConsumer, error) {
	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	kafkaConfig.Consumer.Offsets.Initial = sarama.OffsetNewest
	kafkaConfig.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumerGroup(config.Brokers, config.GroupID, kafkaConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer group: %w", err)
	}

	return &KeyUpdateConsumer{
		consumer:         consumer,
		publicKeyManager: publicKeyManager,
		sessionService:   sessionService,
		redisClient:      redisClient,
		brokers:          config.Brokers,
		topic:            config.Topic,
		groupID:          config.GroupID,
	}, nil
}

func (kuc *KeyUpdateConsumer) Start(ctx context.Context) error {
	topics := []string{kuc.topic}

	// Горутина для обработки ошибок
	go func() {
		for err := range kuc.consumer.Errors() {
			log.Printf("Kafka key update consumer error: %v", err)
		}
	}()

	// Основной цикл обработки сообщений
	for {
		select {
		case <-ctx.Done():
			log.Println("Key update consumer context cancelled")
			return ctx.Err()
		default:
			if err := kuc.consumer.Consume(ctx, topics, kuc); err != nil {
				log.Printf("Error from key update consumer: %v", err)
				time.Sleep(5 * time.Second) // Пауза перед переподключением
			}
		}
	}
}

func (kuc *KeyUpdateConsumer) Close() error {
	return kuc.consumer.Close()
}

// Реализация интерфейса sarama.ConsumerGroupHandler

func (kuc *KeyUpdateConsumer) Setup(sarama.ConsumerGroupSession) error {
	log.Println("Key update consumer setup")
	return nil
}

func (kuc *KeyUpdateConsumer) Cleanup(sarama.ConsumerGroupSession) error {
	log.Println("Key update consumer cleanup")
	return nil
}

func (kuc *KeyUpdateConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message := <-claim.Messages():
			if message == nil {
				return nil
			}

			log.Printf("Received key update message: topic=%s partition=%d offset=%d",
				message.Topic, message.Partition, message.Offset)

			if err := kuc.processKeyUpdate(message); err != nil {
				log.Printf("Error processing key update message: %v", err)
				// Продолжаем обработку, не прерывая процесс
			}

			// Помечаем сообщение как обработанное
			session.MarkMessage(message, "")

		case <-session.Context().Done():
			return nil
		}
	}
}

func (kuc *KeyUpdateConsumer) processKeyUpdate(message *sarama.ConsumerMessage) error {
	var keyUpdate models.PublicKeyUpdate
	if err := json.Unmarshal(message.Value, &keyUpdate); err != nil {
		return fmt.Errorf("failed to unmarshal key update message: %w", err)
	}

	log.Printf("Processing key update from %s, version %d", keyUpdate.ServiceName, keyUpdate.KeyVersion)

	// Обновляем публичный ключ
	if err := kuc.publicKeyManager.UpdateKey(keyUpdate.PublicKeyPEM, keyUpdate.KeyVersion); err != nil {
		return fmt.Errorf("failed to update public key: %w", err)
	}

	// Инвалидируем все активные сессии
	if err := kuc.invalidateAllSessions(); err != nil {
		log.Printf("Warning: Failed to invalidate all sessions: %v", err)
		// Не возвращаем ошибку, так как ключ уже обновлен
	}

	log.Printf("Key update processed successfully: service=%s version=%d",
		keyUpdate.ServiceName, keyUpdate.KeyVersion)

	return nil
}

func (kuc *KeyUpdateConsumer) invalidateAllSessions() error {
	ctx := context.Background()

	// Получаем все ключи сессий
	pattern := "session:*"
	iter := kuc.redisClient.Scan(ctx, 0, pattern, 0).Iterator()

	var sessionsInvalidated int

	for iter.Next(ctx) {
		key := iter.Val()

		// Удаляем сессию
		if err := kuc.redisClient.Del(ctx, key).Err(); err != nil {
			log.Printf("Failed to delete session %s: %v", key, err)
			continue
		}

		sessionsInvalidated++
	}

	if err := iter.Err(); err != nil {
		return fmt.Errorf("error scanning session keys: %w", err)
	}

	log.Printf("Invalidated %d sessions due to key update", sessionsInvalidated)
	return nil
}
