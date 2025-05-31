package kafka

import (
	"common/models"
	"encoding/json"
	"fmt"
	"log"

	"github.com/IBM/sarama"
)

type KeyUpdateProducer struct {
	producer sarama.SyncProducer
	topic    string
}

func NewKeyUpdateProducer(config *ProducerConfig) (*KeyUpdateProducer, error) {
	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Producer.Return.Successes = true
	kafkaConfig.Producer.Return.Errors = true
	kafkaConfig.Producer.RequiredAcks = sarama.WaitForAll

	producer, err := sarama.NewSyncProducer(config.Brokers, kafkaConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create kafka key update producer: %w", err)
	}

	return &KeyUpdateProducer{
		producer: producer,
		topic:    config.Topic,
	}, nil
}

func (p *KeyUpdateProducer) SendKeyUpdate(keyUpdate *models.PublicKeyUpdate) error {
	messageBytes, err := json.Marshal(keyUpdate)
	if err != nil {
		return fmt.Errorf("failed to marshal key update message: %w", err)
	}

	// Отправляем сообщение
	partition, offset, err := p.producer.SendMessage(&sarama.ProducerMessage{
		Topic: p.topic,
		Value: sarama.StringEncoder(messageBytes),
	})

	if err != nil {
		return fmt.Errorf("failed to send key update to kafka: %w", err)
	}

	log.Printf("Key update sent to Kafka: topic=%s partition=%d offset=%d service=%s",
		p.topic, partition, offset, keyUpdate.ServiceName)

	return nil
}

func (p *KeyUpdateProducer) Close() error {
	return p.producer.Close()
}
