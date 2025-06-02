package config

import (
	"os"
	"strconv"
	"strings"
)

// ApplyDatabaseEnvOverrides применяет переменные окружения для базы данных
func ApplyDatabaseEnvOverrides(cfg *Config) {
	if host := os.Getenv("DB_HOST"); host != "" {
		cfg.Database.Host = host
	}
	if user := os.Getenv("DB_USER"); user != "" {
		cfg.Database.User = user
	}
	if password := os.Getenv("DB_PASSWORD"); password != "" {
		cfg.Database.Password = password
	}
	if name := os.Getenv("DB_NAME"); name != "" {
		cfg.Database.Name = name
	}
	if port := os.Getenv("DB_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			cfg.Database.Port = p
		}
	}
}

// ApplyRedisEnvOverrides применяет переменные окружения для Redis
func ApplyRedisEnvOverrides(cfg *Config) {
	if host := os.Getenv("REDIS_HOST"); host != "" {
		cfg.Redis.Host = host
	}
	if port := os.Getenv("REDIS_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			cfg.Redis.Port = p
		}
	}
	if password := os.Getenv("REDIS_PASSWORD"); password != "" {
		cfg.Redis.Password = password
	}
	if db := os.Getenv("REDIS_DB"); db != "" {
		if d, err := strconv.Atoi(db); err == nil {
			cfg.Redis.DB = d
		}
	}
}

// ApplyMinIOEnvOverrides применяет переменные окружения для MinIO
func ApplyMinIOEnvOverrides(cfg *Config) {
	if host := os.Getenv("MINIO_HOST"); host != "" {
		cfg.MinIO.Host = host
	}
	if externalHost := os.Getenv("MINIO_EXTERNAL_HOST"); externalHost != "" {
		cfg.MinIO.ExternalHost = externalHost
	}
	if bucket := os.Getenv("MINIO_BUCKET"); bucket != "" {
		cfg.MinIO.Bucket = bucket
	}
	if accessKey := os.Getenv("MINIO_ACCESS_KEY"); accessKey != "" {
		cfg.MinIO.AccessKey = accessKey
	}
	if secretKey := os.Getenv("MINIO_SECRET_KEY"); secretKey != "" {
		cfg.MinIO.SecretKey = secretKey
	}
}

// ApplyKafkaEnvOverrides применяет переменные окружения для Kafka
func ApplyKafkaEnvOverrides(cfg *Config) {
	if brokers := os.Getenv("KAFKA_BROKERS"); brokers != "" {
		cfg.Kafka.Brokers = strings.Split(brokers, ",")
	}
	if groupID := os.Getenv("KAFKA_GROUP_ID"); groupID != "" {
		cfg.Kafka.GroupID = groupID
	}
	if topic := os.Getenv("NOTIFICATIONS_TOPIC"); topic != "" {
		cfg.Kafka.Topics.Notifications = topic
	}
	if keyTopic := os.Getenv("KEY_UPDATES_TOPIC"); keyTopic != "" {
		cfg.Kafka.Topics.Keys = keyTopic
	}
}

// ApplyAppEnvOverrides применяет переменные окружения для приложения
func ApplyAppEnvOverrides(cfg *Config) {
	if port := os.Getenv("APP_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			cfg.App.Port = p
		}
	}
	if name := os.Getenv("APP_NAME"); name != "" {
		cfg.App.Name = name
	}
}

// ApplyKeysEnvOverrides применяет переменные окружения для ключей
func ApplyKeysEnvOverrides(cfg *Config) {
	if interval := os.Getenv("KEY_ROTATION_INTERVAL"); interval != "" {
		cfg.Keys.RotationInterval = interval
	}
}

// ApplyEmailEnvOverrides применяет переменные окружения для email
func ApplyEmailEnvOverrides(cfg *Config) {
	if host := os.Getenv("SMTP_HOST"); host != "" {
		cfg.Email.SMTPHost = host
	}
	if port := os.Getenv("SMTP_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			cfg.Email.SMTPPort = p
		}
	}
	if username := os.Getenv("SMTP_USERNAME"); username != "" {
		cfg.Email.Username = username
	}
	if password := os.Getenv("SMTP_PASSWORD"); password != "" {
		cfg.Email.Password = password
	}
	if fromEmail := os.Getenv("FROM_EMAIL"); fromEmail != "" {
		cfg.Email.FromEmail = fromEmail
	}
	if fromName := os.Getenv("FROM_NAME"); fromName != "" {
		cfg.Email.FromName = fromName
	}
	if templatePath := os.Getenv("TEMPLATE_PATH"); templatePath != "" {
		cfg.Email.TemplatePath = templatePath
	}
}
