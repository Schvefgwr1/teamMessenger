package services

import (
	"errors"
	"testing"
	"time"

	"common/config"
	"common/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"notificationService/internal/services"
)

// Тесты для EmailService.SendNotification

func TestEmailService_SendNotification_NewTaskNotification_Success(t *testing.T) {
	t.Parallel()
	// Arrange
	mockSender := new(MockEmailSender)
	cfg := createTestEmailConfig()

	emailService, err := services.NewEmailServiceWithSender(cfg, mockSender)
	require.NoError(t, err)

	notification := &models.NewTaskNotification{
		BaseNotification: models.BaseNotification{
			ID:        uuid.New(),
			Type:      models.NotificationNewTask,
			Email:     "executor@example.com",
			CreatedAt: time.Now(),
		},
		TaskID:      1,
		TaskTitle:   "Test Task",
		CreatorName: "Test Creator",
		ExecutorID:  uuid.New(),
	}

	mockSender.On("DialAndSend", mock.Anything).Return(nil)

	// Act
	err = emailService.SendNotification(notification)

	// Assert
	require.NoError(t, err)
	mockSender.AssertExpectations(t)
}

func TestEmailService_SendNotification_NewChatNotification_Group_Success(t *testing.T) {
	t.Parallel()
	// Arrange
	mockSender := new(MockEmailSender)
	cfg := createTestEmailConfig()

	emailService, err := services.NewEmailServiceWithSender(cfg, mockSender)
	require.NoError(t, err)

	notification := &models.NewChatNotification{
		BaseNotification: models.BaseNotification{
			ID:        uuid.New(),
			Type:      models.NotificationNewChat,
			Email:     "user@example.com",
			CreatedAt: time.Now(),
		},
		ChatID:      uuid.New(),
		ChatName:    "Test Group",
		CreatorName: "Test Creator",
		IsGroup:     true,
		Description: "Test Description",
	}

	mockSender.On("DialAndSend", mock.Anything).Return(nil)

	// Act
	err = emailService.SendNotification(notification)

	// Assert
	require.NoError(t, err)
	mockSender.AssertExpectations(t)
}

func TestEmailService_SendNotification_NewChatNotification_Private_Success(t *testing.T) {
	t.Parallel()
	// Arrange
	mockSender := new(MockEmailSender)
	cfg := createTestEmailConfig()

	emailService, err := services.NewEmailServiceWithSender(cfg, mockSender)
	require.NoError(t, err)

	notification := &models.NewChatNotification{
		BaseNotification: models.BaseNotification{
			ID:        uuid.New(),
			Type:      models.NotificationNewChat,
			Email:     "user@example.com",
			CreatedAt: time.Now(),
		},
		ChatID:      uuid.New(),
		ChatName:    "Test Chat",
		CreatorName: "Test Creator",
		IsGroup:     false,
	}

	mockSender.On("DialAndSend", mock.Anything).Return(nil)

	// Act
	err = emailService.SendNotification(notification)

	// Assert
	require.NoError(t, err)
	mockSender.AssertExpectations(t)
}

func TestEmailService_SendNotification_LoginNotification_Success(t *testing.T) {
	t.Parallel()
	// Arrange
	mockSender := new(MockEmailSender)
	cfg := createTestEmailConfig()

	emailService, err := services.NewEmailServiceWithSender(cfg, mockSender)
	require.NoError(t, err)

	notification := &models.LoginNotification{
		BaseNotification: models.BaseNotification{
			ID:        uuid.New(),
			Type:      models.NotificationLogin,
			Email:     "user@example.com",
			CreatedAt: time.Now(),
		},
		UserID:    uuid.New(),
		Username:  "testuser",
		IPAddress: "192.168.1.1",
		LoginTime: time.Now(),
		UserAgent: "Mozilla/5.0",
	}

	mockSender.On("DialAndSend", mock.Anything).Return(nil)

	// Act
	err = emailService.SendNotification(notification)

	// Assert
	require.NoError(t, err)
	mockSender.AssertExpectations(t)
}

func TestEmailService_SendNotification_UnknownType(t *testing.T) {
	t.Parallel()
	// Arrange
	mockSender := new(MockEmailSender)
	cfg := createTestEmailConfig()

	emailService, err := services.NewEmailServiceWithSender(cfg, mockSender)
	require.NoError(t, err)

	unknownNotification := "unknown type"

	// Act
	err = emailService.SendNotification(unknownNotification)

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unknown notification type")
	mockSender.AssertNotCalled(t, "DialAndSend", mock.Anything)
}

func TestEmailService_SendNotification_SenderError(t *testing.T) {
	t.Parallel()
	// Arrange
	mockSender := new(MockEmailSender)
	cfg := createTestEmailConfig()

	emailService, err := services.NewEmailServiceWithSender(cfg, mockSender)
	require.NoError(t, err)

	notification := &models.NewTaskNotification{
		BaseNotification: models.BaseNotification{
			ID:        uuid.New(),
			Type:      models.NotificationNewTask,
			Email:     "executor@example.com",
			CreatedAt: time.Now(),
		},
		TaskID:      1,
		TaskTitle:   "Test Task",
		CreatorName: "Test Creator",
		ExecutorID:  uuid.New(),
	}

	senderError := errors.New("SMTP connection failed")
	mockSender.On("DialAndSend", mock.Anything).Return(senderError)

	// Act
	err = emailService.SendNotification(notification)

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to send email")
	mockSender.AssertExpectations(t)
}

// Тесты для EmailService.NewEmailService

func TestEmailService_NewEmailService_InvalidConfig_EmptyHost(t *testing.T) {
	t.Parallel()
	// Arrange
	cfg := &config.EmailConfig{
		SMTPHost:     "",
		SMTPPort:     587,
		Username:     "test@example.com",
		Password:     "password",
		FromName:     "Test Sender",
		FromEmail:    "sender@example.com",
		TemplatePath: "../../templates",
	}

	// Act
	service, err := services.NewEmailService(cfg)

	// Assert
	require.Error(t, err)
	assert.Nil(t, service)
	assert.Contains(t, err.Error(), "SMTP host is required")
}

func TestEmailService_NewEmailService_InvalidConfig_EmptyUsername(t *testing.T) {
	t.Parallel()
	// Arrange
	cfg := &config.EmailConfig{
		SMTPHost:     "smtp.example.com",
		SMTPPort:     587,
		Username:     "",
		Password:     "password",
		FromName:     "Test Sender",
		FromEmail:    "sender@example.com",
		TemplatePath: "../../templates",
	}

	// Act
	service, err := services.NewEmailService(cfg)

	// Assert
	require.Error(t, err)
	assert.Nil(t, service)
	assert.Contains(t, err.Error(), "SMTP username is required")
}

func TestEmailService_NewEmailService_InvalidConfig_EmptyPassword(t *testing.T) {
	t.Parallel()
	// Arrange
	cfg := &config.EmailConfig{
		SMTPHost:     "smtp.example.com",
		SMTPPort:     587,
		Username:     "test@example.com",
		Password:     "",
		FromName:     "Test Sender",
		FromEmail:    "sender@example.com",
		TemplatePath: "../../templates",
	}

	// Act
	service, err := services.NewEmailService(cfg)

	// Assert
	require.Error(t, err)
	assert.Nil(t, service)
	assert.Contains(t, err.Error(), "SMTP password is required")
}
