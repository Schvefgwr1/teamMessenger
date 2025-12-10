package services

import (
	"common/kafka"
	"common/models"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

type NotificationService struct {
	producer NotificationProducerInterface
}

func NewNotificationService(kafkaConfig *kafka.ProducerConfig) (*NotificationService, error) {
	producer, err := kafka.NewNotificationProducer(kafkaConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create notification producer: %w", err)
	}

	return &NotificationService{
		producer: producer,
	}, nil
}

// NewNotificationServiceWithProducer создаёт сервис с заранее переданным producer (для тестов)
func NewNotificationServiceWithProducer(producer NotificationProducerInterface) *NotificationService {
	return &NotificationService{producer: producer}
}

func (ns *NotificationService) SendChatCreatedNotification(
	chatID uuid.UUID,
	chatName string,
	creatorName string,
	isGroup bool,
	description string,
	userEmail string,
) error {
	if userEmail == "" {
		log.Printf("No email provided for chat %s notification, skipping", chatName)
		return nil
	}

	notification := &models.NewChatNotification{
		BaseNotification: models.BaseNotification{
			ID:        uuid.New(),
			Type:      models.NotificationNewChat,
			Email:     userEmail,
			CreatedAt: time.Now(),
		},
		ChatID:      chatID,
		ChatName:    chatName,
		CreatorName: creatorName,
		IsGroup:     isGroup,
		Description: description,
	}

	if err := ns.producer.SendNotification(notification); err != nil {
		return fmt.Errorf("failed to send chat created notification: %w", err)
	}

	log.Printf("Chat created notification sent for chat %s to %s", chatName, userEmail)
	return nil
}

func (ns *NotificationService) Close() error {
	return ns.producer.Close()
}
