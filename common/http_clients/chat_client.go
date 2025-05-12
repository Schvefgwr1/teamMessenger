package http_clients

import (
	cc "common/contracts/chat-contracts"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// GetChatByID делает HTTP-запрос к чат-сервису и возвращает структуру Chat
func GetChatByID(chatID string) (*cc.Chat, error) {
	url := fmt.Sprintf("http://localhost:8083/api/v1/chats/%s", chatID)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error in request's processing: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("can't get chat: " + resp.Status)
	}

	var chat cc.Chat
	if err := json.NewDecoder(resp.Body).Decode(&chat); err != nil {
		return nil, fmt.Errorf("error of JSON encoding: %w", err)
	}

	return &chat, nil
}
