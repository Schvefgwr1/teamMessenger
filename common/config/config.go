package config

import (
	"gopkg.in/yaml.v3"
	"os"
	"time"
)

type AppConfig struct {
	Port int    `yaml:"port"`
	Name string `yaml:"name"`
}

type MinIO struct {
	Host      string `yaml:"host"`
	AccessKey string `yaml:"access_key"`
	SecretKey string `yaml:"secret_key"`
	Bucket    string `yaml:"bucket"`
}

type Redis struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"DB"`
}

type KeysConfig struct {
	RotationInterval string `yaml:"rotation_interval"`
}

type KafkaConsumerConfig struct {
	GroupID string `yaml:"group_id"`
}

type KafkaConfig struct {
	Consumer KafkaConsumerConfig `yaml:"consumer"`
}

type Config struct {
	Database struct {
		Host     string `yaml:"host"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Name     string `yaml:"name"`
		Port     int    `yaml:"port"`
	} `yaml:"db"`
	MinIO MinIO       `yaml:"minio"`
	Redis Redis       `yaml:"redis"`
	App   AppConfig   `yaml:"app"`
	Keys  KeysConfig  `yaml:"keys"`
	Kafka KafkaConfig `yaml:"kafka"`
}

// GetKeyRotationInterval возвращает интервал обновления ключей как time.Duration
func (c *Config) GetKeyRotationInterval() (time.Duration, error) {
	if c.Keys.RotationInterval == "" {
		return 24 * time.Hour, nil // Значение по умолчанию - 24 часа
	}
	return time.ParseDuration(c.Keys.RotationInterval)
}

func LoadConfig(filename string) (*Config, error) {
	// Чтение данных из YAML-файла
	data, err := os.ReadFile(filename)

	if err != nil {
		return nil, err
	}
	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
