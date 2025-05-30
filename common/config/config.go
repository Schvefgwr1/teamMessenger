package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

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

type Config struct {
	Database struct {
		Host     string `yaml:"host"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Name     string `yaml:"name"`
		Port     int    `yaml:"port"`
	} `yaml:"db"`
	MinIO MinIO `yaml:"minio"`
	Redis Redis `yaml:"redis"`
	App   struct {
		Port int `yaml:"port"`
	} `yaml:"app"`
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
