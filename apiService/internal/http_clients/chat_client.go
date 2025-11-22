package http_clients

import (
	"bytes"
	ac "common/contracts/api-chat"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io"
	"net/http"
)

type ChatClient interface {
	GetUserChats(userID uuid.UUID) ([]*ac.ChatResponse, error)
	CreateChat(req *ac.CreateChatRequest) (*ac.CreateChatServiceResponse, error)
	SendMessage(chatID uuid.UUID, senderID uuid.UUID, req *ac.CreateMessageRequest) (*ac.MessageResponse, error)
	GetChatMessages(chatID uuid.UUID, userID uuid.UUID, offset, limit int) ([]*ac.GetChatMessage, error)
	SearchMessages(userID uuid.UUID, chatID uuid.UUID, query string, offset, limit int) (*ac.GetSearchResponse, error)
	UpdateChat(chatID uuid.UUID, updateReq *ac.UpdateChatRequest) (*ac.UpdateChatResponse, error)
	DeleteChat(chatID uuid.UUID) error
	BanUser(chatID, userID uuid.UUID) error
	ChangeUserRole(chatID uuid.UUID, changeRoleReq *ac.ChangeRoleRequest) error
}

type chatClient struct {
	host string
}

func NewChatClient(host string) ChatClient {
	return &chatClient{host: host}
}

func (c *chatClient) GetUserChats(userID uuid.UUID) ([]*ac.ChatResponse, error) {
	url := fmt.Sprintf("%s/api/v1/chats/%s", c.host, userID.String())

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get user chats: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("chat service returned error: %s", string(bodyBytes))
	}

	var chats []*ac.ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&chats); err != nil {
		return nil, fmt.Errorf("failed to decode chats response: %w", err)
	}

	return chats, nil
}

func (c *chatClient) CreateChat(req *ac.CreateChatRequest) (*ac.CreateChatServiceResponse, error) {
	payload, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to encode request: %w", err)
	}

	resp, err := http.Post(c.host+"/api/v1/chats", "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("request to chat service failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("chat service returned error: %s", string(bodyBytes))
	}

	var serviceResp ac.CreateChatServiceResponse
	if err := json.NewDecoder(resp.Body).Decode(&serviceResp); err != nil {
		return nil, fmt.Errorf("failed to decode chat service response: %w", err)
	}

	return &serviceResp, nil
}

func (c *chatClient) SendMessage(chatID uuid.UUID, senderID uuid.UUID, req *ac.CreateMessageRequest) (*ac.MessageResponse, error) {
	payload, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to encode request: %w", err)
	}

	url := fmt.Sprintf("%s/api/v1/chats/messages/%s", c.host, chatID.String())
	httpReq, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-User-ID", senderID.String())

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("request to chat service failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("chat service returned error: %s", string(bodyBytes))
	}

	var message ac.MessageResponse
	if err := json.NewDecoder(resp.Body).Decode(&message); err != nil {
		return nil, fmt.Errorf("failed to decode chat service response: %w", err)
	}

	return &message, nil
}

func (c *chatClient) GetChatMessages(chatID uuid.UUID, userID uuid.UUID, offset, limit int) ([]*ac.GetChatMessage, error) {
	url := fmt.Sprintf("%s/api/v1/chats/messages/%s?offset=%d&limit=%d", c.host, chatID.String(), offset, limit)

	httpReq, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("X-User-ID", userID.String())

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to get chat messages: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("chat service returned error: %s", string(bodyBytes))
	}

	var messages []*ac.GetChatMessage
	if err := json.NewDecoder(resp.Body).Decode(&messages); err != nil {
		return nil, fmt.Errorf("failed to decode messages response: %w", err)
	}

	return messages, nil
}

func (c *chatClient) SearchMessages(userID uuid.UUID, chatID uuid.UUID, query string, offset, limit int) (*ac.GetSearchResponse, error) {
	url := fmt.Sprintf("%s/api/v1/chats/search/%s?query=%s&offset=%d&limit=%d", c.host, chatID.String(), query, offset, limit)

	httpReq, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("X-User-ID", userID.String())

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to search messages: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("chat service returned error: %s", string(bodyBytes))
	}

	var messages *ac.GetSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&messages); err != nil {
		return nil, fmt.Errorf("failed to decode messages response: %w", err)
	}

	return messages, nil
}

// UpdateChat - обновление чата
func (c *chatClient) UpdateChat(chatID uuid.UUID, updateReq *ac.UpdateChatRequest) (*ac.UpdateChatResponse, error) {
	url := fmt.Sprintf("%s/api/v1/chats/%s", c.host, chatID.String())

	reqBody, err := json.Marshal(updateReq)
	if err != nil {
		return nil, fmt.Errorf("failed to encode request: %w", err)
	}

	httpReq, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("request to chat service failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("chat service returned error: %s", string(bodyBytes))
	}

	var result ac.UpdateChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// DeleteChat - удаление чата
func (c *chatClient) DeleteChat(chatID uuid.UUID) error {
	url := fmt.Sprintf("%s/api/v1/chats/%s", c.host, chatID.String())

	httpReq, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("request to chat service failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("chat service returned error: status %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	return nil
}

// BanUser - блокировка пользователя в чате
func (c *chatClient) BanUser(chatID, userID uuid.UUID) error {
	url := fmt.Sprintf("%s/api/v1/chats/%s/ban/%s", c.host, chatID.String(), userID.String())

	httpReq, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("request to chat service failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("chat service returned error: status %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	return nil
}

// ChangeUserRole - изменение роли пользователя в чате
func (c *chatClient) ChangeUserRole(chatID uuid.UUID, changeRoleReq *ac.ChangeRoleRequest) error {
	url := fmt.Sprintf("%s/api/v1/chats/%s/roles/change", c.host, chatID.String())

	reqBody, err := json.Marshal(changeRoleReq)
	if err != nil {
		return fmt.Errorf("failed to encode request: %w", err)
	}

	httpReq, err := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("request to chat service failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("chat service returned error: status %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	return nil
}
