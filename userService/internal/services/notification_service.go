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

// NewNotificationServiceWithProducer создает сервис с указанным producer (для тестирования)
func NewNotificationServiceWithProducer(producer NotificationProducerInterface) *NotificationService {
	return &NotificationService{
		producer: producer,
	}
}

func (ns *NotificationService) SendLoginNotification(
	userID uuid.UUID,
	username string,
	email string,
	ipAddress string,
	userAgent string,
) error {
	if email == "" {
		log.Printf("No email provided for user %s, skipping login notification", username)
		return nil
	}

	notification := &models.LoginNotification{
		BaseNotification: models.BaseNotification{
			ID:        uuid.New(),
			Type:      models.NotificationLogin,
			Email:     email,
			CreatedAt: time.Now(),
		},
		UserID:    userID,
		Username:  username,
		IPAddress: ipAddress,
		UserAgent: userAgent,
		LoginTime: time.Now(),
	}

	if err := ns.producer.SendNotification(notification); err != nil {
		return fmt.Errorf("failed to send login notification: %w", err)
	}

	log.Printf("Login notification sent for user %s to %s", username, email)
	return nil
}

func (ns *NotificationService) Close() error {
	return ns.producer.Close()
}
