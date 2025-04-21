package utils

import (
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
func ExtractPublicKeyFromFile(path string) (*rsa.PublicKey, error) {
	keyBytes, err := os.ReadFile(path)
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
