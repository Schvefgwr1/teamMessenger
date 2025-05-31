package config

import "os"

// GetEnvOrDefault получает переменную окружения или возвращает значение по умолчанию
func GetEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
