package services

// EmailServiceInterface интерфейс для EmailService (создан для мокирования в тестах)
type EmailServiceInterface interface {
	SendNotification(notification interface{}) error
}
