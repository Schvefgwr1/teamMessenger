package services

import (
	"apiService/internal/dto"
	"apiService/internal/services"
	au "common/contracts/api-user"
	uc "common/contracts/user-contracts"
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockUserClientForKeyService - мок для UserClient (только для тестирования LoadPublicKeyFromService)
type MockUserClientForKeyService struct {
	mock.Mock
}

func (m *MockUserClientForKeyService) GetPublicKey() (*rsa.PublicKey, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*rsa.PublicKey), args.Error(1)
}

// Заглушки для остальных методов интерфейса (не используются в тестах)
func (m *MockUserClientForKeyService) RegisterUser(data au.RegisterUserRequest) (*au.RegisterUserResponse, error) {
	return nil, nil
}

func (m *MockUserClientForKeyService) Login(body *au.Login) (string, uuid.UUID, error) {
	return "", uuid.Nil, nil
}

func (m *MockUserClientForKeyService) GetUserByID(s string) (*au.GetUserResponse, error) {
	return nil, nil
}

func (m *MockUserClientForKeyService) UpdateUser(userID string, req *au.UpdateUserRequest) (*au.UpdateUserResponse, error) {
	return nil, nil
}

func (m *MockUserClientForKeyService) GetAllPermissions() ([]*uc.Permission, error) {
	return nil, nil
}

func (m *MockUserClientForKeyService) GetAllRoles() ([]*uc.Role, error) {
	return nil, nil
}

func (m *MockUserClientForKeyService) CreateRole(req *au.CreateRoleRequest) (*uc.Role, error) {
	return nil, nil
}

func (m *MockUserClientForKeyService) UpdateUserRole(userID string, roleID int) error {
	return nil
}

func (m *MockUserClientForKeyService) UpdateRolePermissions(roleID int, permissionIDs []int) error {
	return nil
}

func (m *MockUserClientForKeyService) DeleteRole(roleID int) error {
	return nil
}

func (m *MockUserClientForKeyService) GetUserBrief(userID, chatID, requesterID string) (*dto.UserBriefResponse, error) {
	return nil, nil
}

func (m *MockUserClientForKeyService) SearchUsers(query string, limit int) (*dto.UserSearchResponse, error) {
	return nil, nil
}

// generateTestRSAKeyPair генерирует тестовую пару RSA ключей
func generateTestRSAKeyPairForService(t *testing.T) (*rsa.PrivateKey, *rsa.PublicKey) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)
	return privateKey, &privateKey.PublicKey
}

// Тесты для LoadPublicKeyFromService

func TestLoadPublicKeyFromService_SuccessFirstAttempt(t *testing.T) {
	// Arrange
	mockClient := new(MockUserClientForKeyService)
	publicKeyManager := services.NewPublicKeyManager()

	_, publicKey := generateTestRSAKeyPairForService(t)

	mockClient.On("GetPublicKey").Return(publicKey, nil).Once()

	// Act
	err := services.LoadPublicKeyFromService(mockClient, publicKeyManager)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, publicKey, publicKeyManager.GetCurrentKey())
	assert.True(t, publicKeyManager.HasKey())
	mockClient.AssertExpectations(t)
}

func TestLoadPublicKeyFromService_SuccessAfterRetries(t *testing.T) {
	// Arrange
	mockClient := new(MockUserClientForKeyService)
	publicKeyManager := services.NewPublicKeyManager()

	_, publicKey := generateTestRSAKeyPairForService(t)

	// Первые 2 попытки неудачны, третья успешна
	mockClient.On("GetPublicKey").Return(nil, errors.New("connection error")).Twice()
	mockClient.On("GetPublicKey").Return(publicKey, nil).Once()

	// Act
	err := services.LoadPublicKeyFromService(mockClient, publicKeyManager)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, publicKey, publicKeyManager.GetCurrentKey())
	assert.True(t, publicKeyManager.HasKey())
	mockClient.AssertExpectations(t)
}

func TestLoadPublicKeyFromService_AllRetriesFail(t *testing.T) {
	// Arrange
	mockClient := new(MockUserClientForKeyService)
	publicKeyManager := services.NewPublicKeyManager()

	expectedError := errors.New("service unavailable")

	// Все 10 попыток неудачны
	mockClient.On("GetPublicKey").Return(nil, expectedError).Times(10)

	// Act
	err := services.LoadPublicKeyFromService(mockClient, publicKeyManager)

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to load public key from userService after 10 attempts")
	assert.Nil(t, publicKeyManager.GetCurrentKey())
	assert.False(t, publicKeyManager.HasKey())
	mockClient.AssertExpectations(t)
}

func TestLoadPublicKeyFromService_MaxRetries(t *testing.T) {
	// Arrange
	mockClient := new(MockUserClientForKeyService)
	publicKeyManager := services.NewPublicKeyManager()

	expectedError := errors.New("timeout")

	// Все попытки неудачны
	mockClient.On("GetPublicKey").Return(nil, expectedError).Times(10)

	// Act
	start := time.Now()
	err := services.LoadPublicKeyFromService(mockClient, publicKeyManager)
	duration := time.Since(start)

	// Assert
	require.Error(t, err)
	// Проверяем, что было сделано 10 попыток с задержками (минимум 9 задержок по 2 секунды = 18 секунд)
	// Но учитываем что тесты могут быть быстрее, поэтому проверяем что прошло достаточно времени
	assert.GreaterOrEqual(t, duration, 9*time.Second, "Should have retried with delays")
	mockClient.AssertExpectations(t)
}
