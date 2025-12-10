package services

import (
	"errors"
	"testing"

	"common/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"userService/internal/services"
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

// Тесты для NotificationService.SendLoginNotification

func TestNotificationService_SendLoginNotification_Success(t *testing.T) {
	// Arrange
	mockProducer := new(MockNotificationProducer)
	service := services.NewNotificationServiceWithProducer(mockProducer)

	userID := uuid.New()
	username := "testuser"
	email := "test@example.com"
	ipAddress := "192.168.1.1"
	userAgent := "Mozilla/5.0"

	mockProducer.On("SendNotification", mock.MatchedBy(func(n interface{}) bool {
		notification, ok := n.(*models.LoginNotification)
		if !ok {
			return false
		}
		return notification.UserID == userID &&
			notification.Username == username &&
			notification.Email == email &&
			notification.IPAddress == ipAddress &&
			notification.UserAgent == userAgent
	})).Return(nil)

	// Act
	err := service.SendLoginNotification(userID, username, email, ipAddress, userAgent)

	// Assert
	require.NoError(t, err)
	mockProducer.AssertExpectations(t)
}

func TestNotificationService_SendLoginNotification_EmptyEmail(t *testing.T) {
	// Arrange
	mockProducer := new(MockNotificationProducer)
	service := services.NewNotificationServiceWithProducer(mockProducer)

	userID := uuid.New()
	username := "testuser"
	email := ""
	ipAddress := "192.168.1.1"
	userAgent := "Mozilla/5.0"

	// Act
	err := service.SendLoginNotification(userID, username, email, ipAddress, userAgent)

	// Assert
	require.NoError(t, err)
	// Producer не должен быть вызван при пустом email
	mockProducer.AssertNotCalled(t, "SendNotification", mock.Anything)
}

func TestNotificationService_SendLoginNotification_KafkaError(t *testing.T) {
	// Arrange
	mockProducer := new(MockNotificationProducer)
	service := services.NewNotificationServiceWithProducer(mockProducer)

	userID := uuid.New()
	username := "testuser"
	email := "test@example.com"
	ipAddress := "192.168.1.1"
	userAgent := "Mozilla/5.0"
	kafkaError := errors.New("kafka connection failed")

	mockProducer.On("SendNotification", mock.Anything).Return(kafkaError)

	// Act
	err := service.SendLoginNotification(userID, username, email, ipAddress, userAgent)

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to send login notification")
	mockProducer.AssertExpectations(t)
}

// Тесты для NotificationService.Close

func TestNotificationService_Close_Success(t *testing.T) {
	// Arrange
	mockProducer := new(MockNotificationProducer)
	service := services.NewNotificationServiceWithProducer(mockProducer)

	mockProducer.On("Close").Return(nil)

	// Act
	err := service.Close()

	// Assert
	require.NoError(t, err)
	mockProducer.AssertExpectations(t)
}

func TestNotificationService_Close_Error(t *testing.T) {
	// Arrange
	mockProducer := new(MockNotificationProducer)
	service := services.NewNotificationServiceWithProducer(mockProducer)

	closeError := errors.New("close error")
	mockProducer.On("Close").Return(closeError)

	// Act
	err := service.Close()

	// Assert
	require.Error(t, err)
	assert.Equal(t, closeError, err)
	mockProducer.AssertExpectations(t)
}
