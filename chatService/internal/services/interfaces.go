package services

import (
	"github.com/google/uuid"
)

// NotificationServiceInterface - интерфейс для NotificationService для возможности мокирования
type NotificationServiceInterface interface {
	SendChatCreatedNotification(
		chatID uuid.UUID,
		chatName string,
		creatorName string,
		isGroup bool,
		description string,
		userEmail string,
	) error
	Close() error
}

// NotificationProducerInterface - интерфейс для Kafka producer (для мокирования)
type NotificationProducerInterface interface {
	SendNotification(notification interface{}) error
	Close() error
}
