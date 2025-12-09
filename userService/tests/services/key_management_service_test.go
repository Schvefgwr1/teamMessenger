package services

import (
	"errors"
	"testing"

	"common/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"userService/internal/services"
)

// MockKeyUpdateProducer - мок для KeyUpdateProducerInterface
type MockKeyUpdateProducer struct {
	mock.Mock
}

func (m *MockKeyUpdateProducer) SendKeyUpdate(keyUpdate *models.PublicKeyUpdate) error {
	args := m.Called(keyUpdate)
	return args.Error(0)
}

func (m *MockKeyUpdateProducer) Close() error {
	args := m.Called()
	return args.Error(0)
}

// MockUtils - мок для утилит (для тестирования без реальной генерации ключей)
// В реальных тестах мы можем использовать реальные утилиты или создать моки

// Тесты для KeyManagementService.RegenerateKeys

func TestKeyManagementService_RegenerateKeys_Success(t *testing.T) {
	// Arrange
	mockProducer := new(MockKeyUpdateProducer)
	initialVersion := 1
	service := services.NewKeyManagementServiceWithProducer(mockProducer, initialVersion)

	mockProducer.On("SendKeyUpdate", mock.MatchedBy(func(update *models.PublicKeyUpdate) bool {
		return update.ServiceName == "userService" &&
			update.KeyVersion == initialVersion &&
			update.PublicKeyPEM != ""
	})).Return(nil)

	// Act
	err := service.RegenerateKeys()

	// Assert
	require.NoError(t, err)
	assert.Equal(t, initialVersion+1, service.GetCurrentKeyVersion())
	mockProducer.AssertExpectations(t)
}

func TestKeyManagementService_RegenerateKeys_KafkaError(t *testing.T) {
	// Arrange
	mockProducer := new(MockKeyUpdateProducer)
	initialVersion := 1
	service := services.NewKeyManagementServiceWithProducer(mockProducer, initialVersion)

	kafkaError := errors.New("kafka connection failed")
	mockProducer.On("SendKeyUpdate", mock.Anything).Return(kafkaError)

	// Act
	err := service.RegenerateKeys()

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to send key update to kafka")
	assert.Equal(t, initialVersion, service.GetCurrentKeyVersion()) // Версия не должна измениться при ошибке
	mockProducer.AssertExpectations(t)
}

// Тесты для KeyManagementService.GetCurrentKeyVersion

func TestKeyManagementService_GetCurrentKeyVersion(t *testing.T) {
	// Arrange
	mockProducer := new(MockKeyUpdateProducer)
	initialVersion := 5
	service := services.NewKeyManagementServiceWithProducer(mockProducer, initialVersion)

	// Act
	version := service.GetCurrentKeyVersion()

	// Assert
	assert.Equal(t, initialVersion, version)
}

func TestKeyManagementService_GetCurrentKeyVersion_AfterRegenerate(t *testing.T) {
	// Arrange
	mockProducer := new(MockKeyUpdateProducer)
	initialVersion := 1
	service := services.NewKeyManagementServiceWithProducer(mockProducer, initialVersion)

	mockProducer.On("SendKeyUpdate", mock.Anything).Return(nil)

	// Act
	err := service.RegenerateKeys()
	require.NoError(t, err)

	version := service.GetCurrentKeyVersion()

	// Assert
	assert.Equal(t, initialVersion+1, version)
	mockProducer.AssertExpectations(t)
}

// Тесты для KeyManagementService.Close

func TestKeyManagementService_Close_Success(t *testing.T) {
	// Arrange
	mockProducer := new(MockKeyUpdateProducer)
	service := services.NewKeyManagementServiceWithProducer(mockProducer, 1)

	mockProducer.On("Close").Return(nil)

	// Act
	err := service.Close()

	// Assert
	require.NoError(t, err)
	mockProducer.AssertExpectations(t)
}

func TestKeyManagementService_Close_Error(t *testing.T) {
	// Arrange
	mockProducer := new(MockKeyUpdateProducer)
	service := services.NewKeyManagementServiceWithProducer(mockProducer, 1)

	closeError := errors.New("close error")
	mockProducer.On("Close").Return(closeError)

	// Act
	err := service.Close()

	// Assert
	require.Error(t, err)
	assert.Equal(t, closeError, err)
	mockProducer.AssertExpectations(t)
}

// Тесты для проверки версионирования ключей

func TestKeyManagementService_MultipleRegenerations(t *testing.T) {
	// Arrange
	mockProducer := new(MockKeyUpdateProducer)
	initialVersion := 1
	service := services.NewKeyManagementServiceWithProducer(mockProducer, initialVersion)

	mockProducer.On("SendKeyUpdate", mock.MatchedBy(func(update *models.PublicKeyUpdate) bool {
		return update.KeyVersion == 1
	})).Return(nil).Once()

	mockProducer.On("SendKeyUpdate", mock.MatchedBy(func(update *models.PublicKeyUpdate) bool {
		return update.KeyVersion == 2
	})).Return(nil).Once()

	// Act
	err1 := service.RegenerateKeys()
	require.NoError(t, err1)
	assert.Equal(t, 2, service.GetCurrentKeyVersion())

	err2 := service.RegenerateKeys()
	require.NoError(t, err2)
	assert.Equal(t, 3, service.GetCurrentKeyVersion())

	// Assert
	mockProducer.AssertExpectations(t)
}
