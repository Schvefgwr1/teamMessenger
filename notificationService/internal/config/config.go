package config

import "common/config"

type NotificationConfig struct {
	App   config.AppConfig `yaml:"app"`
	Kafka KafkaConfig      `yaml:"kafka"`
	Email EmailConfig      `yaml:"email"`
}

type KafkaConfig struct {
	Brokers []string `yaml:"brokers"`
	GroupID string   `yaml:"group_id"`
	Topics  Topics   `yaml:"topics"`
}

type Topics struct {
	Notifications string `yaml:"notifications"`
}

type EmailConfig struct {
	SMTPHost     string `yaml:"smtp_host"`
	SMTPPort     int    `yaml:"smtp_port"`
	Username     string `yaml:"username"`
	Password     string `yaml:"password"`
	FromEmail    string `yaml:"from_email"`
	FromName     string `yaml:"from_name"`
	TemplatePath string `yaml:"template_path"`
}
