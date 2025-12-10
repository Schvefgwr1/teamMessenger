package services

import (
	"common/kafka"
	"common/models"
	"fmt"
	"log"
	"time"
	"userService/internal/utils"

	"github.com/google/uuid"
)

type KeyManagementService struct {
	keyProducer KeyUpdateProducerInterface
	keyVersion  int
}

func NewKeyManagementService(kafkaConfig *kafka.ProducerConfig) (*KeyManagementService, error) {
	producer, err := kafka.NewKeyUpdateProducer(kafkaConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create key update producer: %w", err)
	}

	return &KeyManagementService{
		keyProducer: producer,
		keyVersion:  1, // Начинаем с версии 1
	}, nil
}

// NewKeyManagementServiceWithProducer создает сервис с указанным producer (для тестирования)
func NewKeyManagementServiceWithProducer(producer KeyUpdateProducerInterface, initialVersion int) *KeyManagementService {
	return &KeyManagementService{
		keyProducer: producer,
		keyVersion:  initialVersion,
	}
}

func (kms *KeyManagementService) RegenerateKeys() error {
	log.Println("Starting key regeneration process...")

	// Генерируем новые ключи и сохраняем их в файлы
	publicKey, err := utils.GenerateAndSaveNewKeys()
	if err != nil {
		return fmt.Errorf("failed to generate and save new keys: %w", err)
	}

	log.Println("New keys generated and saved successfully")

	// Конвертируем публичный ключ в PEM строку
	publicKeyPEM, err := utils.PublicKeyToPEM(publicKey)
	if err != nil {
		return fmt.Errorf("failed to convert public key to PEM: %w", err)
	}

	// Создаем сообщение об обновлении ключа
	keyUpdate := &models.PublicKeyUpdate{
		ID:           uuid.New(),
		PublicKeyPEM: publicKeyPEM,
		UpdatedAt:    time.Now(),
		ServiceName:  "userService",
		KeyVersion:   kms.keyVersion,
	}

	// Отправляем в Kafka
	if err := kms.keyProducer.SendKeyUpdate(keyUpdate); err != nil {
		return fmt.Errorf("failed to send key update to kafka: %w", err)
	}

	log.Printf("Public key update sent to Kafka (version %d)", kms.keyVersion)

	// Увеличиваем версию ключа
	kms.keyVersion++

	return nil
}

func (kms *KeyManagementService) GetCurrentKeyVersion() int {
	return kms.keyVersion
}

func (kms *KeyManagementService) Close() error {
	return kms.keyProducer.Close()
}
