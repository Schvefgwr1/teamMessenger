package kafka

import (
	"os"
	"strings"
)

// GetKafkaBrokers получает список Kafka брокеров из переменных окружения
func GetKafkaBrokers() []string {
	brokers := os.Getenv("KAFKA_BROKERS")
	if brokers == "" {
		return []string{"localhost:9092"} // значение по умолчанию
	}

	// Поддержка нескольких брокеров через запятую
	brokerList := strings.Split(brokers, ",")
	for i, broker := range brokerList {
		brokerList[i] = strings.TrimSpace(broker)
	}
	return brokerList
}

// GetNotificationsTopic получает топик для уведомлений из переменных окружения
func GetNotificationsTopic() string {
	topic := os.Getenv("KAFKA_TOPIC_NOTIFICATIONS")
	if topic == "" {
		return "notifications" // значение по умолчанию
	}
	return topic
}

// GetKeyUpdatesTopic получает топик для обновлений ключей из переменных окружения
func GetKeyUpdatesTopic() string {
	topic := os.Getenv("KAFKA_TOPIC_KEY_UPDATES")
	if topic == "" {
		return "key_updates" // значение по умолчанию
	}
	return topic
}
