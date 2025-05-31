package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"os"
	"path/filepath"
	"time"
)

// Claims структура для JWT
type Claims struct {
	UserID      uuid.UUID `json:"user_id"`
	Permissions []string  `json:"permissions"`
	jwt.RegisteredClaims
}

// LoadPrivateKey загружает RSA private key из PEM-файла
func LoadPrivateKey() (*rsa.PrivateKey, error) {
	wd, err := os.Getwd() // Получаем текущую директорию
	if err != nil {
		return nil, err
	}

	keyPath := filepath.Join(wd, "cmd", "keys", "private.pem")

	keyBytes, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(keyBytes)
	if block == nil || block.Type != "PRIVATE KEY" {
		return nil, errors.New("invalid PEM block type or empty private key")
	}

	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	privateKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("parsed key is not an RSA private key")
	}

	return privateKey, nil
}

// ExtractPublicKeyFromFile загружает RSA публичный ключ из PEM-файла
func ExtractPublicKeyFromFile() (*rsa.PublicKey, error) {
	wd, err := os.Getwd()
	keyPath := filepath.Join(wd, "cmd", "keys", "public.pem")
	keyBytes, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(keyBytes)
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, errors.New("invalid PEM block type or empty public key")
	}

	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	pubKey, ok := pubInterface.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("not RSA public key")
	}

	return pubKey, nil
}

// GenerateKeyPair генерирует новую пару RSA ключей
func GenerateKeyPair(bitSize int) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	if bitSize < 2048 {
		bitSize = 2048 // Минимальный размер для безопасности
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		return nil, nil, err
	}

	return privateKey, &privateKey.PublicKey, nil
}

// SavePrivateKeyToFile сохраняет приватный ключ в PEM-файл
func SavePrivateKeyToFile(privateKey *rsa.PrivateKey) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	keyPath := filepath.Join(wd, "cmd", "keys", "private.pem")

	// Создаем директорию если её нет
	if err := os.MkdirAll(filepath.Dir(keyPath), 0755); err != nil {
		return err
	}

	// Кодируем ключ в PKCS8 формат
	privateKeyBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return err
	}

	// Создаем PEM блок
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privateKeyBytes,
	})

	// Записываем в файл
	return os.WriteFile(keyPath, privateKeyPEM, 0600)
}

// SavePublicKeyToFile сохраняет публичный ключ в PEM-файл
func SavePublicKeyToFile(publicKey *rsa.PublicKey) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	keyPath := filepath.Join(wd, "cmd", "keys", "public.pem")

	// Создаем директорию если её нет
	if err := os.MkdirAll(filepath.Dir(keyPath), 0755); err != nil {
		return err
	}

	// Кодируем ключ в PKIX формат
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return err
	}

	// Создаем PEM блок
	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	})

	// Записываем в файл
	return os.WriteFile(keyPath, publicKeyPEM, 0644)
}

// GenerateAndSaveNewKeys генерирует новую пару ключей и сохраняет их в файлы
func GenerateAndSaveNewKeys() (*rsa.PublicKey, error) {
	// Генерируем новую пару ключей
	privateKey, publicKey, err := GenerateKeyPair(2048)
	if err != nil {
		return nil, err
	}

	// Сохраняем приватный ключ
	if err := SavePrivateKeyToFile(privateKey); err != nil {
		return nil, err
	}

	// Сохраняем публичный ключ
	if err := SavePublicKeyToFile(publicKey); err != nil {
		return nil, err
	}

	return publicKey, nil
}

// PublicKeyToPEM конвертирует публичный ключ в PEM строку
func PublicKeyToPEM(publicKey *rsa.PublicKey) (string, error) {
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

// GenerateJWT создает JWT токен по приватному ключу
func GenerateJWT(userID uuid.UUID, permissions []string) (string, error) {
	claims := Claims{
		UserID:      userID,
		Permissions: permissions,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(72 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	privateKey, err := LoadPrivateKey()
	if err != nil {
		return "", errors.New("error of loading key")
	}
	return token.SignedString(privateKey)
}

// HashPassword хеширует пароль
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword), err
}

// CheckPasswordHash проверяет соответствие пароля хешу
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
