package models

import (
	"github.com/google/uuid"
	"time"
)

// NotificationType определяет тип уведомления
type NotificationType string

const (
	NotificationNewTask NotificationType = "new_task"
	NotificationNewChat NotificationType = "new_chat"
	NotificationLogin   NotificationType = "user_login"
)

// BaseNotification базовая структура уведомления
type BaseNotification struct {
	ID        uuid.UUID        `json:"id"`
	Type      NotificationType `json:"type"`
	Email     string           `json:"email"`
	CreatedAt time.Time        `json:"created_at"`
}

// NewTaskNotification уведомление о новой задаче
type NewTaskNotification struct {
	BaseNotification
	TaskID      int       `json:"task_id"`
	TaskTitle   string    `json:"task_title"`
	CreatorName string    `json:"creator_name"`
	ExecutorID  uuid.UUID `json:"executor_id,omitempty"`
}

// NewChatNotification уведомление о новом чате
type NewChatNotification struct {
	BaseNotification
	ChatID      uuid.UUID `json:"chat_id"`
	ChatName    string    `json:"chat_name"`
	CreatorName string    `json:"creator_name"`
	IsGroup     bool      `json:"is_group"`
	Description string    `json:"description,omitempty"`
}

// LoginNotification уведомление о входе
type LoginNotification struct {
	BaseNotification
	UserID    uuid.UUID `json:"user_id"`
	Username  string    `json:"username"`
	IPAddress string    `json:"ip_address"`
	LoginTime time.Time `json:"login_time"`
	UserAgent string    `json:"user_agent"`
}

// KafkaMessage обертка для сообщений в Kafka
type KafkaMessage struct {
	Type    NotificationType `json:"type"`
	Payload interface{}      `json:"payload"`
}
