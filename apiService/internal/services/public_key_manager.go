package services

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"sync"
)

type PublicKeyManager struct {
	currentKey *rsa.PublicKey
	keyVersion int
	mutex      sync.RWMutex
}

func NewPublicKeyManager() *PublicKeyManager {
	return &PublicKeyManager{
		keyVersion: 0,
	}
}

// GetCurrentKey возвращает текущий публичный ключ (thread-safe)
func (pkm *PublicKeyManager) GetCurrentKey() *rsa.PublicKey {
	pkm.mutex.RLock()
	defer pkm.mutex.RUnlock()
	return pkm.currentKey
}

// GetKeyVersion возвращает текущую версию ключа
func (pkm *PublicKeyManager) GetKeyVersion() int {
	pkm.mutex.RLock()
	defer pkm.mutex.RUnlock()
	return pkm.keyVersion
}

// UpdateKey обновляет публичный ключ из PEM строки (thread-safe)
func (pkm *PublicKeyManager) UpdateKey(publicKeyPEM string, version int) error {
	// Парсим PEM строку
	publicKey, err := pkm.parsePublicKeyFromPEM(publicKeyPEM)
	if err != nil {
		return fmt.Errorf("failed to parse public key from PEM: %w", err)
	}

	pkm.mutex.Lock()
	defer pkm.mutex.Unlock()

	pkm.currentKey = publicKey
	pkm.keyVersion = version

	log.Printf("Public key updated to version %d", version)
	return nil
}

// SetInitialKey устанавливает начальный ключ (из userService при старте)
func (pkm *PublicKeyManager) SetInitialKey(key *rsa.PublicKey) {
	pkm.mutex.Lock()
	defer pkm.mutex.Unlock()

	pkm.currentKey = key
	pkm.keyVersion = 0 // Начальная версия
	log.Println("Initial public key set")
}

// HasKey проверяет, установлен ли ключ
func (pkm *PublicKeyManager) HasKey() bool {
	pkm.mutex.RLock()
	defer pkm.mutex.RUnlock()
	return pkm.currentKey != nil
}

// parsePublicKeyFromPEM парсит RSA публичный ключ из PEM строки
func (pkm *PublicKeyManager) parsePublicKeyFromPEM(publicKeyPEM string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(publicKeyPEM))
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, fmt.Errorf("invalid PEM block type or empty public key")
	}

	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}

	pubKey, ok := pubInterface.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("not RSA public key")
	}

	return pubKey, nil
}
