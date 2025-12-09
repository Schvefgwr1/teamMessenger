package services

import (
	"common/models"
)

// NotificationProducerInterface - интерфейс для Kafka NotificationProducer для возможности мокирования
type NotificationProducerInterface interface {
	SendNotification(notification interface{}) error
	Close() error
}

// KeyUpdateProducerInterface - интерфейс для Kafka KeyUpdateProducer для возможности мокирования
type KeyUpdateProducerInterface interface {
	SendKeyUpdate(keyUpdate *models.PublicKeyUpdate) error
	Close() error
}
