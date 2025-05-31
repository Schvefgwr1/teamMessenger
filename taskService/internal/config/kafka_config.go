package config

// KafkaConfig конфигурация Kafka для taskService
type KafkaConfig struct {
	Brokers []string `yaml:"brokers"`
	Topic   string   `yaml:"topic"`
}

// TaskServiceConfig расширенная конфигурация для taskService
type TaskServiceConfig struct {
	Kafka KafkaConfig `yaml:"kafka"`
}
