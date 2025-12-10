package services

import (
	"errors"
	"testing"

	"chatService/internal/services"
	"common/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockNotificationProducer - мок для NotificationProducerInterface
type MockNotificationProducer struct {
	mock.Mock
}

func (m *MockNotificationProducer) SendNotification(notification interface{}) error {
	args := m.Called(notification)
	return args.Error(0)
}

func (m *MockNotificationProducer) Close() error {
	args := m.Called()
	return args.Error(0)
}

func TestNotificationService_SendChatCreatedNotification_Success(t *testing.T) {
	mockProducer := new(MockNotificationProducer)
	service := services.NewNotificationServiceWithProducer(mockProducer)

	chatID := uuid.New()
	chatName := "Test"
	creator := "creator"
	isGroup := true
	desc := "desc"
	email := "user@example.com"

	mockProducer.On("SendNotification", mock.MatchedBy(func(n interface{}) bool {
		msg, ok := n.(*models.NewChatNotification)
		if !ok {
			return false
		}
		return msg.ChatID == chatID && msg.ChatName == chatName && msg.CreatorName == creator && msg.IsGroup == isGroup && msg.Description == desc && msg.Email == email
	})).Return(nil)

	err := service.SendChatCreatedNotification(chatID, chatName, creator, isGroup, desc, email)

	require.NoError(t, err)
	mockProducer.AssertExpectations(t)
}

func TestNotificationService_SendChatCreatedNotification_NoEmail(t *testing.T) {
	mockProducer := new(MockNotificationProducer)
	service := services.NewNotificationServiceWithProducer(mockProducer)

	err := service.SendChatCreatedNotification(uuid.New(), "Test", "creator", false, "", "")

	require.NoError(t, err)
	mockProducer.AssertNotCalled(t, "SendNotification", mock.Anything)
}

func TestNotificationService_SendChatCreatedNotification_ProducerError(t *testing.T) {
	mockProducer := new(MockNotificationProducer)
	service := services.NewNotificationServiceWithProducer(mockProducer)

	mockProducer.On("SendNotification", mock.Anything).Return(errors.New("kafka"))

	err := service.SendChatCreatedNotification(uuid.New(), "Test", "creator", false, "", "user@example.com")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to send chat created notification")
	mockProducer.AssertExpectations(t)
}

func TestNotificationService_Close_Success(t *testing.T) {
	mockProducer := new(MockNotificationProducer)
	service := services.NewNotificationServiceWithProducer(mockProducer)

	mockProducer.On("Close").Return(nil)

	err := service.Close()

	require.NoError(t, err)
	mockProducer.AssertExpectations(t)
}

func TestNotificationService_Close_Error(t *testing.T) {
	mockProducer := new(MockNotificationProducer)
	service := services.NewNotificationServiceWithProducer(mockProducer)

	closeError := errors.New("close error")
	mockProducer.On("Close").Return(closeError)

	err := service.Close()

	require.Error(t, err)
	assert.Equal(t, closeError, err)
	mockProducer.AssertExpectations(t)
}
