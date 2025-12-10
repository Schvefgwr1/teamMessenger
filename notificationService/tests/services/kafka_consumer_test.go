package services

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	"common/models"
	"github.com/IBM/sarama"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"notificationService/internal/services"
)

// MockEmailService - мок для EmailService
type MockEmailService struct {
	mock.Mock
}

func (m *MockEmailService) SendNotification(notification interface{}) error {
	args := m.Called(notification)
	return args.Error(0)
}

// Тесты для KafkaConsumer.processMessage

func TestKafkaConsumer_ProcessMessage_NewTaskNotification_Success(t *testing.T) {
	t.Parallel()
	// Arrange
	mockEmailService := new(MockEmailService)
	consumer := &services.KafkaConsumer{
		EmailService: mockEmailService,
	}

	taskNotification := &models.NewTaskNotification{
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

	kafkaMsg := models.KafkaMessage{
		Type:    models.NotificationNewTask,
		Payload: taskNotification,
	}

	messageBytes, err := json.Marshal(kafkaMsg)
	require.NoError(t, err)

	message := &sarama.ConsumerMessage{
		Topic:     "notifications",
		Partition: 0,
		Offset:    1,
		Value:     messageBytes,
	}

	mockEmailService.On("SendNotification", mock.MatchedBy(func(n interface{}) bool {
		notif, ok := n.(*models.NewTaskNotification)
		if !ok {
			return false
		}
		return notif.TaskID == taskNotification.TaskID &&
			notif.TaskTitle == taskNotification.TaskTitle &&
			notif.Email == taskNotification.Email
	})).Return(nil)

	// Act
	err = consumer.ProcessMessage(message)

	// Assert
	require.NoError(t, err)
	mockEmailService.AssertExpectations(t)
}

func TestKafkaConsumer_ProcessMessage_NewChatNotification_Success(t *testing.T) {
	t.Parallel()
	// Arrange
	mockEmailService := new(MockEmailService)
	consumer := &services.KafkaConsumer{
		EmailService: mockEmailService,
	}

	chatNotification := &models.NewChatNotification{
		BaseNotification: models.BaseNotification{
			ID:        uuid.New(),
			Type:      models.NotificationNewChat,
			Email:     "user@example.com",
			CreatedAt: time.Now(),
		},
		ChatID:      uuid.New(),
		ChatName:    "Test Chat",
		CreatorName: "Test Creator",
		IsGroup:     true,
		Description: "Test Description",
	}

	kafkaMsg := models.KafkaMessage{
		Type:    models.NotificationNewChat,
		Payload: chatNotification,
	}

	messageBytes, err := json.Marshal(kafkaMsg)
	require.NoError(t, err)

	message := &sarama.ConsumerMessage{
		Topic:     "notifications",
		Partition: 0,
		Offset:    1,
		Value:     messageBytes,
	}

	mockEmailService.On("SendNotification", mock.MatchedBy(func(n interface{}) bool {
		notif, ok := n.(*models.NewChatNotification)
		if !ok {
			return false
		}
		return notif.ChatID == chatNotification.ChatID &&
			notif.ChatName == chatNotification.ChatName &&
			notif.Email == chatNotification.Email
	})).Return(nil)

	// Act
	err = consumer.ProcessMessage(message)

	// Assert
	require.NoError(t, err)
	mockEmailService.AssertExpectations(t)
}

func TestKafkaConsumer_ProcessMessage_LoginNotification_Success(t *testing.T) {
	t.Parallel()
	// Arrange
	mockEmailService := new(MockEmailService)
	consumer := &services.KafkaConsumer{
		EmailService: mockEmailService,
	}

	loginNotification := &models.LoginNotification{
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

	kafkaMsg := models.KafkaMessage{
		Type:    models.NotificationLogin,
		Payload: loginNotification,
	}

	messageBytes, err := json.Marshal(kafkaMsg)
	require.NoError(t, err)

	message := &sarama.ConsumerMessage{
		Topic:     "notifications",
		Partition: 0,
		Offset:    1,
		Value:     messageBytes,
	}

	mockEmailService.On("SendNotification", mock.MatchedBy(func(n interface{}) bool {
		notif, ok := n.(*models.LoginNotification)
		if !ok {
			return false
		}
		return notif.UserID == loginNotification.UserID &&
			notif.Username == loginNotification.Username &&
			notif.Email == loginNotification.Email
	})).Return(nil)

	// Act
	err = consumer.ProcessMessage(message)

	// Assert
	require.NoError(t, err)
	mockEmailService.AssertExpectations(t)
}

func TestKafkaConsumer_ProcessMessage_InvalidJSON(t *testing.T) {
	t.Parallel()
	// Arrange
	mockEmailService := new(MockEmailService)
	consumer := &services.KafkaConsumer{
		EmailService: mockEmailService,
	}

	message := &sarama.ConsumerMessage{
		Topic:     "notifications",
		Partition: 0,
		Offset:    1,
		Value:     []byte("invalid json"),
	}

	// Act
	err := consumer.ProcessMessage(message)

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to unmarshal kafka message")
	mockEmailService.AssertNotCalled(t, "SendNotification", mock.Anything)
}

func TestKafkaConsumer_ProcessMessage_UnknownNotificationType(t *testing.T) {
	t.Parallel()
	// Arrange
	mockEmailService := new(MockEmailService)
	consumer := &services.KafkaConsumer{
		EmailService: mockEmailService,
	}

	kafkaMsg := models.KafkaMessage{
		Type:    "unknown_type",
		Payload: map[string]interface{}{"test": "data"},
	}

	messageBytes, err := json.Marshal(kafkaMsg)
	require.NoError(t, err)

	message := &sarama.ConsumerMessage{
		Topic:     "notifications",
		Partition: 0,
		Offset:    1,
		Value:     messageBytes,
	}

	// Act
	err = consumer.ProcessMessage(message)

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unknown notification type")
	mockEmailService.AssertNotCalled(t, "SendNotification", mock.Anything)
}

func TestKafkaConsumer_ProcessMessage_EmailServiceError(t *testing.T) {
	t.Parallel()
	// Arrange
	mockEmailService := new(MockEmailService)
	consumer := &services.KafkaConsumer{
		EmailService: mockEmailService,
	}

	taskNotification := &models.NewTaskNotification{
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

	kafkaMsg := models.KafkaMessage{
		Type:    models.NotificationNewTask,
		Payload: taskNotification,
	}

	messageBytes, err := json.Marshal(kafkaMsg)
	require.NoError(t, err)

	message := &sarama.ConsumerMessage{
		Topic:     "notifications",
		Partition: 0,
		Offset:    1,
		Value:     messageBytes,
	}

	emailError := errors.New("failed to send email")
	mockEmailService.On("SendNotification", mock.Anything).Return(emailError)

	// Act
	err = consumer.ProcessMessage(message)

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to send email notification")
	mockEmailService.AssertExpectations(t)
}

// Тесты для KafkaConsumer.parseNotification

func TestKafkaConsumer_ParseNotification_NewTaskNotification_Success(t *testing.T) {
	t.Parallel()
	// Arrange
	mockEmailService := new(MockEmailService)
	consumer := &services.KafkaConsumer{
		EmailService: mockEmailService,
	}

	taskNotification := &models.NewTaskNotification{
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

	kafkaMsg := models.KafkaMessage{
		Type:    models.NotificationNewTask,
		Payload: taskNotification,
	}

	// Act
	result, err := consumer.ParseNotification(kafkaMsg)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	notif, ok := result.(*models.NewTaskNotification)
	require.True(t, ok)
	assert.Equal(t, taskNotification.TaskID, notif.TaskID)
	assert.Equal(t, taskNotification.TaskTitle, notif.TaskTitle)
	assert.Equal(t, taskNotification.Email, notif.Email)
}

func TestKafkaConsumer_ParseNotification_NewChatNotification_Success(t *testing.T) {
	t.Parallel()
	// Arrange
	mockEmailService := new(MockEmailService)
	consumer := &services.KafkaConsumer{
		EmailService: mockEmailService,
	}

	chatNotification := &models.NewChatNotification{
		BaseNotification: models.BaseNotification{
			ID:        uuid.New(),
			Type:      models.NotificationNewChat,
			Email:     "user@example.com",
			CreatedAt: time.Now(),
		},
		ChatID:      uuid.New(),
		ChatName:    "Test Chat",
		CreatorName: "Test Creator",
		IsGroup:     true,
		Description: "Test Description",
	}

	kafkaMsg := models.KafkaMessage{
		Type:    models.NotificationNewChat,
		Payload: chatNotification,
	}

	// Act
	result, err := consumer.ParseNotification(kafkaMsg)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	notif, ok := result.(*models.NewChatNotification)
	require.True(t, ok)
	assert.Equal(t, chatNotification.ChatID, notif.ChatID)
	assert.Equal(t, chatNotification.ChatName, notif.ChatName)
	assert.Equal(t, chatNotification.Email, notif.Email)
}

func TestKafkaConsumer_ParseNotification_LoginNotification_Success(t *testing.T) {
	t.Parallel()
	// Arrange
	mockEmailService := new(MockEmailService)
	consumer := &services.KafkaConsumer{
		EmailService: mockEmailService,
	}

	loginNotification := &models.LoginNotification{
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

	kafkaMsg := models.KafkaMessage{
		Type:    models.NotificationLogin,
		Payload: loginNotification,
	}

	// Act
	result, err := consumer.ParseNotification(kafkaMsg)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	notif, ok := result.(*models.LoginNotification)
	require.True(t, ok)
	assert.Equal(t, loginNotification.UserID, notif.UserID)
	assert.Equal(t, loginNotification.Username, notif.Username)
	assert.Equal(t, loginNotification.Email, notif.Email)
}

func TestKafkaConsumer_ParseNotification_UnknownType(t *testing.T) {
	t.Parallel()
	// Arrange
	mockEmailService := new(MockEmailService)
	consumer := &services.KafkaConsumer{
		EmailService: mockEmailService,
	}

	kafkaMsg := models.KafkaMessage{
		Type:    "unknown_type",
		Payload: map[string]interface{}{"test": "data"},
	}

	// Act
	result, err := consumer.ParseNotification(kafkaMsg)

	// Assert
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "unknown notification type")
}
