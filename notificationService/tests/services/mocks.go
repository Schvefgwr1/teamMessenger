package services

import (
	"common/config"
	"github.com/stretchr/testify/mock"
	"gopkg.in/gomail.v2"
)

// MockEmailSender - мок для EmailSender интерфейса
type MockEmailSender struct {
	mock.Mock
}

func (m *MockEmailSender) DialAndSend(messages ...*gomail.Message) error {
	args := m.Called(messages)
	return args.Error(0)
}

// Вспомогательные функции для создания тестовых данных

func createTestEmailConfig() *config.EmailConfig {
	return &config.EmailConfig{
		SMTPHost:     "smtp.example.com",
		SMTPPort:     587,
		Username:     "test@example.com",
		Password:     "password",
		FromName:     "Test Sender",
		FromEmail:    "sender@example.com",
		TemplatePath: "../../templates",
	}
}
