package services

import (
	"apiService/internal/http_clients"
	"crypto/rsa"
)

func LoadPublicKeyFromService(client http_clients.UserClient, publicKey **rsa.PublicKey) error {
	key, err := client.GetPublicKey()
	if err != nil {
		return err
	}
	*publicKey = key
	return nil
}
