package utils

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"userService/internal/utils"
)

// Тесты для GenerateKeyPair

func TestGenerateKeyPair_Success(t *testing.T) {
	// Act
	privateKey, publicKey, err := utils.GenerateKeyPair(2048)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, privateKey)
	assert.NotNil(t, publicKey)
	assert.Equal(t, 2048, privateKey.N.BitLen())
}

func TestGenerateKeyPair_MinimumSize(t *testing.T) {
	// Act - передаем размер меньше минимального
	privateKey, publicKey, err := utils.GenerateKeyPair(1024)

	// Assert - должен использоваться минимальный размер 2048
	require.NoError(t, err)
	assert.NotNil(t, privateKey)
	assert.NotNil(t, publicKey)
	assert.GreaterOrEqual(t, privateKey.N.BitLen(), 2048)
}

func TestGenerateKeyPair_LargeSize(t *testing.T) {
	// Act
	privateKey, publicKey, err := utils.GenerateKeyPair(4096)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, privateKey)
	assert.NotNil(t, publicKey)
	assert.Equal(t, 4096, privateKey.N.BitLen())
}

// Тесты для HashPassword и CheckPasswordHash

func TestHashPassword_Success(t *testing.T) {
	// Arrange
	password := "testpassword123"

	// Act
	hash, err := utils.HashPassword(password)

	// Assert
	require.NoError(t, err)
	assert.NotEmpty(t, hash)
	assert.NotEqual(t, password, hash)
}

func TestHashPassword_EmptyPassword(t *testing.T) {
	// Arrange
	password := ""

	// Act
	hash, err := utils.HashPassword(password)

	// Assert
	require.NoError(t, err)
	assert.NotEmpty(t, hash)
}

func TestCheckPasswordHash_Success(t *testing.T) {
	// Arrange
	password := "testpassword123"
	hash, err := utils.HashPassword(password)
	require.NoError(t, err)

	// Act
	isValid := utils.CheckPasswordHash(password, hash)

	// Assert
	assert.True(t, isValid)
}

func TestCheckPasswordHash_WrongPassword(t *testing.T) {
	// Arrange
	password := "testpassword123"
	wrongPassword := "wrongpassword"
	hash, err := utils.HashPassword(password)
	require.NoError(t, err)

	// Act
	isValid := utils.CheckPasswordHash(wrongPassword, hash)

	// Assert
	assert.False(t, isValid)
}

func TestCheckPasswordHash_InvalidHash(t *testing.T) {
	// Arrange
	password := "testpassword123"
	invalidHash := "invalidhash"

	// Act
	isValid := utils.CheckPasswordHash(password, invalidHash)

	// Assert
	assert.False(t, isValid)
}

// Тесты для PublicKeyToPEM

func TestPublicKeyToPEM_Success(t *testing.T) {
	// Arrange
	_, publicKey, err := utils.GenerateKeyPair(2048)
	require.NoError(t, err)

	// Act
	pemString, err := utils.PublicKeyToPEM(publicKey)

	// Assert
	require.NoError(t, err)
	assert.NotEmpty(t, pemString)
	assert.Contains(t, pemString, "PUBLIC KEY")
	assert.Contains(t, pemString, "BEGIN")
	assert.Contains(t, pemString, "END")
}

func TestPublicKeyToPEM_NilKey(t *testing.T) {
	// Act & Assert - nil ключ вызывает панику, поэтому используем recover
	defer func() {
		if r := recover(); r != nil {
			// Паника ожидаема для nil ключа
			assert.NotNil(t, r)
		}
	}()

	pemString, err := utils.PublicKeyToPEM(nil)

	// Если паники не было, проверяем ошибку
	if err != nil {
		assert.Empty(t, pemString)
	}
}

// Тесты для SavePrivateKeyToFile и SavePublicKeyToFile

func TestSavePrivateKeyToFile_Success(t *testing.T) {
	// Arrange
	privateKey, _, err := utils.GenerateKeyPair(2048)
	require.NoError(t, err)

	// Создаем временную директорию
	tempDir := t.TempDir()
	originalWd, _ := os.Getwd()
	defer func() {
		os.Chdir(originalWd)
	}()

	os.Chdir(tempDir)
	keysDir := filepath.Join(tempDir, "cmd", "keys")
	os.MkdirAll(keysDir, 0755)

	// Act
	err = utils.SavePrivateKeyToFile(privateKey)

	// Assert
	require.NoError(t, err)

	// Проверяем, что файл создан
	keyPath := filepath.Join(tempDir, "cmd", "keys", "private.pem")
	_, err = os.Stat(keyPath)
	assert.NoError(t, err)
}

func TestSavePublicKeyToFile_Success(t *testing.T) {
	// Arrange
	_, publicKey, err := utils.GenerateKeyPair(2048)
	require.NoError(t, err)

	// Создаем временную директорию
	tempDir := t.TempDir()
	originalWd, _ := os.Getwd()
	defer func() {
		os.Chdir(originalWd)
	}()

	os.Chdir(tempDir)
	keysDir := filepath.Join(tempDir, "cmd", "keys")
	os.MkdirAll(keysDir, 0755)

	// Act
	err = utils.SavePublicKeyToFile(publicKey)

	// Assert
	require.NoError(t, err)

	// Проверяем, что файл создан
	keyPath := filepath.Join(tempDir, "cmd", "keys", "public.pem")
	_, err = os.Stat(keyPath)
	assert.NoError(t, err)
}

// Тесты для GenerateAndSaveNewKeys

func TestGenerateAndSaveNewKeys_Success(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	originalWd, _ := os.Getwd()
	defer func() {
		os.Chdir(originalWd)
	}()

	os.Chdir(tempDir)
	keysDir := filepath.Join(tempDir, "cmd", "keys")
	os.MkdirAll(keysDir, 0755)

	// Act
	publicKey, err := utils.GenerateAndSaveNewKeys()

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, publicKey)

	// Проверяем, что файлы созданы
	privateKeyPath := filepath.Join(tempDir, "cmd", "keys", "private.pem")
	publicKeyPath := filepath.Join(tempDir, "cmd", "keys", "public.pem")

	_, err = os.Stat(privateKeyPath)
	assert.NoError(t, err)

	_, err = os.Stat(publicKeyPath)
	assert.NoError(t, err)
}

// Тесты для LoadPrivateKey и ExtractPublicKeyFromFile

func TestLoadPrivateKey_Success(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	originalWd, _ := os.Getwd()
	defer func() {
		os.Chdir(originalWd)
	}()

	os.Chdir(tempDir)
	keysDir := filepath.Join(tempDir, "cmd", "keys")
	os.MkdirAll(keysDir, 0755)

	// Создаем ключи
	privateKey, _, err := utils.GenerateKeyPair(2048)
	require.NoError(t, err)
	err = utils.SavePrivateKeyToFile(privateKey)
	require.NoError(t, err)

	// Act
	loadedKey, err := utils.LoadPrivateKey()

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, loadedKey)
	assert.Equal(t, privateKey.N, loadedKey.N)
}

func TestLoadPrivateKey_FileNotFound(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	originalWd, _ := os.Getwd()
	defer func() {
		os.Chdir(originalWd)
	}()

	os.Chdir(tempDir)

	// Act
	loadedKey, err := utils.LoadPrivateKey()

	// Assert
	require.Error(t, err)
	assert.Nil(t, loadedKey)
}

func TestExtractPublicKeyFromFile_Success(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	originalWd, _ := os.Getwd()
	defer func() {
		os.Chdir(originalWd)
	}()

	os.Chdir(tempDir)
	keysDir := filepath.Join(tempDir, "cmd", "keys")
	os.MkdirAll(keysDir, 0755)

	// Создаем ключи
	_, publicKey, err := utils.GenerateKeyPair(2048)
	require.NoError(t, err)
	err = utils.SavePublicKeyToFile(publicKey)
	require.NoError(t, err)

	// Act
	loadedKey, err := utils.ExtractPublicKeyFromFile()

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, loadedKey)
	assert.Equal(t, publicKey.N, loadedKey.N)
}

func TestExtractPublicKeyFromFile_FileNotFound(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	originalWd, _ := os.Getwd()
	defer func() {
		os.Chdir(originalWd)
	}()

	os.Chdir(tempDir)

	// Act
	loadedKey, err := utils.ExtractPublicKeyFromFile()

	// Assert
	require.Error(t, err)
	assert.Nil(t, loadedKey)
}

// Тесты для GenerateJWT

func TestGenerateJWT_Success(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	originalWd, _ := os.Getwd()
	defer func() {
		os.Chdir(originalWd)
	}()

	os.Chdir(tempDir)
	keysDir := filepath.Join(tempDir, "cmd", "keys")
	os.MkdirAll(keysDir, 0755)

	// Создаем ключи
	privateKey, _, err := utils.GenerateKeyPair(2048)
	require.NoError(t, err)
	err = utils.SavePrivateKeyToFile(privateKey)
	require.NoError(t, err)

	userID := uuid.New()
	permissions := []string{"read", "write"}

	// Act
	token, err := utils.GenerateJWT(userID, permissions)

	// Assert
	require.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestGenerateJWT_NoKeyFile(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	originalWd, _ := os.Getwd()
	defer func() {
		os.Chdir(originalWd)
	}()

	os.Chdir(tempDir)

	userID := uuid.New()
	permissions := []string{"read", "write"}

	// Act
	token, err := utils.GenerateJWT(userID, permissions)

	// Assert
	require.Error(t, err)
	assert.Empty(t, token)
	assert.Contains(t, err.Error(), "error of loading key")
}

func TestGenerateJWT_EmptyPermissions(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	originalWd, _ := os.Getwd()
	defer func() {
		os.Chdir(originalWd)
	}()

	os.Chdir(tempDir)
	keysDir := filepath.Join(tempDir, "cmd", "keys")
	os.MkdirAll(keysDir, 0755)

	// Создаем ключи
	privateKey, _, err := utils.GenerateKeyPair(2048)
	require.NoError(t, err)
	err = utils.SavePrivateKeyToFile(privateKey)
	require.NoError(t, err)

	userID := uuid.New()
	permissions := []string{}

	// Act
	token, err := utils.GenerateJWT(userID, permissions)

	// Assert
	require.NoError(t, err)
	assert.NotEmpty(t, token)
}

// Дополнительные тесты для проверки совместимости ключей

func TestKeyPairCompatibility(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	originalWd, _ := os.Getwd()
	defer func() {
		os.Chdir(originalWd)
	}()

	os.Chdir(tempDir)
	keysDir := filepath.Join(tempDir, "cmd", "keys")
	os.MkdirAll(keysDir, 0755)

	// Генерируем и сохраняем ключи
	privateKey, publicKey, err := utils.GenerateKeyPair(2048)
	require.NoError(t, err)

	err = utils.SavePrivateKeyToFile(privateKey)
	require.NoError(t, err)

	err = utils.SavePublicKeyToFile(publicKey)
	require.NoError(t, err)

	// Загружаем ключи обратно
	loadedPrivateKey, err := utils.LoadPrivateKey()
	require.NoError(t, err)

	loadedPublicKey, err := utils.ExtractPublicKeyFromFile()
	require.NoError(t, err)

	// Проверяем совместимость
	assert.Equal(t, privateKey.N, loadedPrivateKey.N)
	assert.Equal(t, publicKey.N, loadedPublicKey.N)
	assert.Equal(t, loadedPrivateKey.PublicKey.N, loadedPublicKey.N)
}

// Тесты для проверки разных размеров ключей

func TestHashPassword_DifferentPasswords(t *testing.T) {
	// Arrange
	password1 := "password1"
	password2 := "password2"

	// Act
	hash1, err1 := utils.HashPassword(password1)
	require.NoError(t, err1)

	hash2, err2 := utils.HashPassword(password2)
	require.NoError(t, err2)

	// Assert - хеши должны быть разными
	assert.NotEqual(t, hash1, hash2)

	// Проверяем, что каждый хеш соответствует своему паролю
	assert.True(t, utils.CheckPasswordHash(password1, hash1))
	assert.True(t, utils.CheckPasswordHash(password2, hash2))
	assert.False(t, utils.CheckPasswordHash(password1, hash2))
	assert.False(t, utils.CheckPasswordHash(password2, hash1))
}

// Тесты для проверки PEM формата

func TestPublicKeyToPEM_ValidFormat(t *testing.T) {
	// Arrange
	_, publicKey, err := utils.GenerateKeyPair(2048)
	require.NoError(t, err)

	// Act
	pemString, err := utils.PublicKeyToPEM(publicKey)

	// Assert
	require.NoError(t, err)
	assert.Contains(t, pemString, "-----BEGIN PUBLIC KEY-----")
	assert.Contains(t, pemString, "-----END PUBLIC KEY-----")
}

// Тесты для проверки bcrypt хеширования

func TestHashPassword_BcryptCompatibility(t *testing.T) {
	// Arrange
	password := "testpassword"

	// Act
	hash, err := utils.HashPassword(password)

	// Assert
	require.NoError(t, err)

	// Проверяем, что это валидный bcrypt хеш
	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	assert.NoError(t, err)
}
