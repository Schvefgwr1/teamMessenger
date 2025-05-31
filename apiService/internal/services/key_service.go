package services

import (
	"apiService/internal/http_clients"
	"fmt"
)

// LoadPublicKeyFromService загружает начальный публичный ключ из userService
func LoadPublicKeyFromService(client http_clients.UserClient, publicKeyManager *PublicKeyManager) error {
	key, err := client.GetPublicKey()
	if err != nil {
		return fmt.Errorf("failed to load public key from userService: %w", err)
	}

	publicKeyManager.SetInitialKey(key)
	return nil
}
