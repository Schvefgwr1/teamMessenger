package kafka

import (
	"common/models"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"log"
)

type NotificationProducer struct {
	producer sarama.SyncProducer
	topic    string
}

type ProducerConfig struct {
	Brokers []string
	Topic   string
}

func NewNotificationProducer(config *ProducerConfig) (*NotificationProducer, error) {
	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Producer.Return.Successes = true
	kafkaConfig.Producer.Return.Errors = true
	kafkaConfig.Producer.RequiredAcks = sarama.WaitForAll

	producer, err := sarama.NewSyncProducer(config.Brokers, kafkaConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create kafka producer: %w", err)
	}

	return &NotificationProducer{
		producer: producer,
		topic:    config.Topic,
	}, nil
}

func (p *NotificationProducer) SendNotification(notification interface{}) error {
	var notificationType models.NotificationType

	// Определяем тип уведомления
	switch notification.(type) {
	case *models.NewTaskNotification:
		notificationType = models.NotificationNewTask
	case *models.NewChatNotification:
		notificationType = models.NotificationNewChat
	case *models.LoginNotification:
		notificationType = models.NotificationLogin
	default:
		return fmt.Errorf("unknown notification type: %T", notification)
	}

	// Создаем сообщение для Kafka
	kafkaMessage := models.KafkaMessage{
		Type:    notificationType,
		Payload: notification,
	}

	messageBytes, err := json.Marshal(kafkaMessage)
	if err != nil {
		return fmt.Errorf("failed to marshal kafka message: %w", err)
	}

	// Отправляем сообщение
	partition, offset, err := p.producer.SendMessage(&sarama.ProducerMessage{
		Topic: p.topic,
		Value: sarama.StringEncoder(messageBytes),
	})

	if err != nil {
		return fmt.Errorf("failed to send message to kafka: %w", err)
	}

	log.Printf("Message sent to Kafka: topic=%s partition=%d offset=%d type=%s",
		p.topic, partition, offset, notificationType)

	return nil
}

func (p *NotificationProducer) Close() error {
	return p.producer.Close()
}
