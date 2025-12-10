package services

// NotificationProducerInterface - интерфейс для Kafka NotificationProducer для возможности мокирования
type NotificationProducerInterface interface {
	SendNotification(notification interface{}) error
	Close() error
}
