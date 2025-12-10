package services

import (
	"apiService/internal/services"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// generateTestRSAKeyPair генерирует тестовую пару RSA ключей
func generateTestRSAKeyPair(t *testing.T) (*rsa.PrivateKey, *rsa.PublicKey) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)
	return privateKey, &privateKey.PublicKey
}

// publicKeyToPEM конвертирует RSA публичный ключ в PEM строку
func publicKeyToPEM(publicKey *rsa.PublicKey) (string, error) {
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return "", err
	}

	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	})

	return string(publicKeyPEM), nil
}

// Тесты для PublicKeyManager.NewPublicKeyManager

func TestPublicKeyManager_NewPublicKeyManager(t *testing.T) {
	// Act
	pkm := services.NewPublicKeyManager()

	// Assert
	require.NotNil(t, pkm)
	assert.Nil(t, pkm.GetCurrentKey())
	assert.Equal(t, 0, pkm.GetKeyVersion())
	assert.False(t, pkm.HasKey())
}

// Тесты для PublicKeyManager.SetInitialKey

func TestPublicKeyManager_SetInitialKey_Success(t *testing.T) {
	// Arrange
	pkm := services.NewPublicKeyManager()
	_, publicKey := generateTestRSAKeyPair(t)

	// Act
	pkm.SetInitialKey(publicKey)

	// Assert
	assert.Equal(t, publicKey, pkm.GetCurrentKey())
	assert.Equal(t, 0, pkm.GetKeyVersion())
	assert.True(t, pkm.HasKey())
}

func TestPublicKeyManager_SetInitialKey_Overwrite(t *testing.T) {
	// Arrange
	pkm := services.NewPublicKeyManager()
	_, publicKey1 := generateTestRSAKeyPair(t)
	_, publicKey2 := generateTestRSAKeyPair(t)

	// Act
	pkm.SetInitialKey(publicKey1)
	pkm.SetInitialKey(publicKey2)

	// Assert
	assert.Equal(t, publicKey2, pkm.GetCurrentKey())
	assert.Equal(t, 0, pkm.GetKeyVersion())
	assert.True(t, pkm.HasKey())
}

// Тесты для PublicKeyManager.GetCurrentKey

func TestPublicKeyManager_GetCurrentKey_NoKey(t *testing.T) {
	// Arrange
	pkm := services.NewPublicKeyManager()

	// Act
	key := pkm.GetCurrentKey()

	// Assert
	assert.Nil(t, key)
}

func TestPublicKeyManager_GetCurrentKey_WithKey(t *testing.T) {
	// Arrange
	pkm := services.NewPublicKeyManager()
	_, publicKey := generateTestRSAKeyPair(t)
	pkm.SetInitialKey(publicKey)

	// Act
	key := pkm.GetCurrentKey()

	// Assert
	assert.Equal(t, publicKey, key)
}

// Тесты для PublicKeyManager.GetKeyVersion

func TestPublicKeyManager_GetKeyVersion_Initial(t *testing.T) {
	// Arrange
	pkm := services.NewPublicKeyManager()

	// Act
	version := pkm.GetKeyVersion()

	// Assert
	assert.Equal(t, 0, version)
}

func TestPublicKeyManager_GetKeyVersion_AfterSetInitial(t *testing.T) {
	// Arrange
	pkm := services.NewPublicKeyManager()
	_, publicKey := generateTestRSAKeyPair(t)
	pkm.SetInitialKey(publicKey)

	// Act
	version := pkm.GetKeyVersion()

	// Assert
	assert.Equal(t, 0, version)
}

func TestPublicKeyManager_GetKeyVersion_AfterUpdate(t *testing.T) {
	// Arrange
	pkm := services.NewPublicKeyManager()
	_, publicKey := generateTestRSAKeyPair(t)
	publicKeyPEM, err := publicKeyToPEM(publicKey)
	require.NoError(t, err)

	err = pkm.UpdateKey(publicKeyPEM, 5)
	require.NoError(t, err)

	// Act
	version := pkm.GetKeyVersion()

	// Assert
	assert.Equal(t, 5, version)
}

// Тесты для PublicKeyManager.HasKey

func TestPublicKeyManager_HasKey_NoKey(t *testing.T) {
	// Arrange
	pkm := services.NewPublicKeyManager()

	// Act
	hasKey := pkm.HasKey()

	// Assert
	assert.False(t, hasKey)
}

func TestPublicKeyManager_HasKey_WithKey(t *testing.T) {
	// Arrange
	pkm := services.NewPublicKeyManager()
	_, publicKey := generateTestRSAKeyPair(t)
	pkm.SetInitialKey(publicKey)

	// Act
	hasKey := pkm.HasKey()

	// Assert
	assert.True(t, hasKey)
}

// Тесты для PublicKeyManager.UpdateKey

func TestPublicKeyManager_UpdateKey_Success(t *testing.T) {
	// Arrange
	pkm := services.NewPublicKeyManager()
	_, publicKey := generateTestRSAKeyPair(t)
	publicKeyPEM, err := publicKeyToPEM(publicKey)
	require.NoError(t, err)

	// Act
	err = pkm.UpdateKey(publicKeyPEM, 3)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, publicKey, pkm.GetCurrentKey())
	assert.Equal(t, 3, pkm.GetKeyVersion())
	assert.True(t, pkm.HasKey())
}

func TestPublicKeyManager_UpdateKey_InvalidPEM(t *testing.T) {
	// Arrange
	pkm := services.NewPublicKeyManager()
	invalidPEM := "invalid PEM data"

	// Act
	err := pkm.UpdateKey(invalidPEM, 1)

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid PEM block")
	assert.Nil(t, pkm.GetCurrentKey())
	assert.False(t, pkm.HasKey())
}

func TestPublicKeyManager_UpdateKey_EmptyPEM(t *testing.T) {
	// Arrange
	pkm := services.NewPublicKeyManager()

	// Act
	err := pkm.UpdateKey("", 1)

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid PEM block")
}

func TestPublicKeyManager_UpdateKey_WrongPEMType(t *testing.T) {
	// Arrange
	pkm := services.NewPublicKeyManager()
	// Создаем PEM блок с неправильным типом
	wrongPEM := `-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQC7VJTUt9Us8cKj
MzEfYyjiWA4R4/M2bN1Kp8O5Z8vC5Z5Z5Z5Z5Z5Z5Z5Z5Z5Z5Z5Z5Z5Z5Z5Z5Z5Z5Z
-----END PRIVATE KEY-----`

	// Act
	err := pkm.UpdateKey(wrongPEM, 1)

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid PEM block")
}

func TestPublicKeyManager_UpdateKey_NonRSAPublicKey(t *testing.T) {
	// Arrange
	pkm := services.NewPublicKeyManager()
	// Создаем валидный PEM блок, но не RSA ключ
	// Для этого используем простой текст в формате PEM
	nonRSAPEM := `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA1234567890
-----END PUBLIC KEY-----`

	// Act
	err := pkm.UpdateKey(nonRSAPEM, 1)

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "parse public key")
}

func TestPublicKeyManager_UpdateKey_MultipleUpdates(t *testing.T) {
	// Arrange
	pkm := services.NewPublicKeyManager()
	_, publicKey1 := generateTestRSAKeyPair(t)
	_, publicKey2 := generateTestRSAKeyPair(t)

	publicKeyPEM1, err := publicKeyToPEM(publicKey1)
	require.NoError(t, err)
	publicKeyPEM2, err := publicKeyToPEM(publicKey2)
	require.NoError(t, err)

	// Act
	err = pkm.UpdateKey(publicKeyPEM1, 1)
	require.NoError(t, err)
	assert.Equal(t, 1, pkm.GetKeyVersion())

	err = pkm.UpdateKey(publicKeyPEM2, 2)
	require.NoError(t, err)

	// Assert
	assert.Equal(t, publicKey2, pkm.GetCurrentKey())
	assert.Equal(t, 2, pkm.GetKeyVersion())
}

// Тесты для thread-safety (параллельный доступ)

func TestPublicKeyManager_ConcurrentAccess(t *testing.T) {
	// Arrange
	pkm := services.NewPublicKeyManager()
	_, publicKey := generateTestRSAKeyPair(t)
	publicKeyPEM, err := publicKeyToPEM(publicKey)
	require.NoError(t, err)

	pkm.SetInitialKey(publicKey)

	// Act - запускаем параллельные чтения и записи
	done := make(chan bool)
	errors := make(chan error, 10)

	// Параллельные чтения
	for i := 0; i < 5; i++ {
		go func() {
			for j := 0; j < 100; j++ {
				key := pkm.GetCurrentKey()
				if key == nil {
					errors <- assert.AnError
				}
				version := pkm.GetKeyVersion()
				if version < 0 {
					errors <- assert.AnError
				}
				hasKey := pkm.HasKey()
				if !hasKey {
					errors <- assert.AnError
				}
			}
			done <- true
		}()
	}

	// Параллельные обновления
	for i := 0; i < 5; i++ {
		go func(version int) {
			for j := 0; j < 10; j++ {
				err := pkm.UpdateKey(publicKeyPEM, version+j)
				if err != nil {
					errors <- err
				}
			}
			done <- true
		}(i * 10)
	}

	// Ждем завершения всех горутин
	for i := 0; i < 10; i++ {
		<-done
	}

	// Assert
	close(errors)
	errorCount := 0
	for err := range errors {
		if err != nil {
			errorCount++
		}
	}
	assert.Equal(t, 0, errorCount, "No errors should occur during concurrent access")
	assert.True(t, pkm.HasKey())
}
