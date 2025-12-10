package services

import (
	"github.com/google/uuid"
)

// NotificationServiceInterface - интерфейс для NotificationService для возможности мокирования
type NotificationServiceInterface interface {
	SendTaskCreatedNotification(
		taskID int,
		taskTitle string,
		creatorName string,
		executorID uuid.UUID,
		executorEmail string,
	) error
	Close() error
}
