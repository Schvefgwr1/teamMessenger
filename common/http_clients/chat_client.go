package http_clients

import (
	"common/config"
	cc "common/contracts/chat-contracts"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// UserRoleInChatResponse - ответ с ролью пользователя в чате
type UserRoleInChatResponse struct {
	RoleName string `json:"roleName"`
}

// GetChatByID делает HTTP-запрос к чат-сервису и возвращает структуру Chat
func GetChatByID(chatID string) (*cc.Chat, error) {
	baseURL := config.GetEnvOrDefault("CHAT_SERVICE_URL", "http://localhost:8083")
	url := fmt.Sprintf("%s/api/v1/chats/%s", baseURL, chatID)

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

// GetUserRoleInChat получает роль пользователя в чате
// chatID - ID чата
// userID - ID пользователя, чью роль запрашиваем
// requesterID - ID пользователя, который делает запрос (для проверки прав доступа)
func GetUserRoleInChat(chatID, userID, requesterID string) (*UserRoleInChatResponse, error) {
	baseURL := config.GetEnvOrDefault("CHAT_SERVICE_URL", "http://localhost:8083")
	url := fmt.Sprintf("%s/api/v1/chats/%s/user-roles/%s", baseURL, chatID, userID)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	// Передаем ID запрашивающего пользователя для проверки прав
	req.Header.Set("X-User-ID", requesterID)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error in request's processing: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusForbidden {
		return nil, errors.New("access denied: user is not a member of this chat")
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, errors.New("user not found in chat")
	}

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("can't get user role: %s - %s", resp.Status, string(bodyBytes))
	}

	var roleResponse UserRoleInChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&roleResponse); err != nil {
		return nil, fmt.Errorf("error of JSON decoding: %w", err)
	}

	return &roleResponse, nil
}
