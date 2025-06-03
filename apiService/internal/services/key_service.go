package services

import (
	"apiService/internal/http_clients"
	"fmt"
	"log"
	"time"
)

// LoadPublicKeyFromService загружает начальный публичный ключ из userService с retry логикой
func LoadPublicKeyFromService(client http_clients.UserClient, publicKeyManager *PublicKeyManager) error {
	maxRetries := 10
	retryDelay := 2 * time.Second

	for attempt := 1; attempt <= maxRetries; attempt++ {
		log.Printf("Attempting to load public key from userService (attempt %d/%d)", attempt, maxRetries)

		key, err := client.GetPublicKey()
		if err != nil {
			log.Printf("Failed to load public key (attempt %d): %v", attempt, err)

			if attempt == maxRetries {
				return fmt.Errorf("failed to load public key from userService after %d attempts: %w", maxRetries, err)
			}

			log.Printf("Retrying in %v...", retryDelay)
			time.Sleep(retryDelay)
			continue
		}

		publicKeyManager.SetInitialKey(key)
		log.Printf("Successfully loaded public key on attempt %d", attempt)
		return nil
	}

	return fmt.Errorf("failed to load public key after %d attempts", maxRetries)
}
