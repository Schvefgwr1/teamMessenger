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
	producer *kafka.NotificationProducer
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

func (ns *NotificationService) SendTaskCreatedNotification(
	taskID int,
	taskTitle string,
	creatorName string,
	executorID uuid.UUID,
	executorEmail string,
) error {
	if executorEmail == "" {
		log.Printf("No executor email provided for task %d, skipping notification", taskID)
		return nil
	}

	notification := &models.NewTaskNotification{
		BaseNotification: models.BaseNotification{
			ID:        uuid.New(),
			Type:      models.NotificationNewTask,
			Email:     executorEmail,
			CreatedAt: time.Now(),
		},
		TaskID:      taskID,
		TaskTitle:   taskTitle,
		CreatorName: creatorName,
		ExecutorID:  executorID,
	}

	if err := ns.producer.SendNotification(notification); err != nil {
		return fmt.Errorf("failed to send task created notification: %w", err)
	}

	log.Printf("Task created notification sent for task %d to %s", taskID, executorEmail)
	return nil
}

func (ns *NotificationService) Close() error {
	return ns.producer.Close()
}
