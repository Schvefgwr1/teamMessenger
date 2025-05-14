package http_clients

import (
	"bytes"
	au "common/contracts/api-user"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type UserClient interface {
	RegisterUser(data au.RegisterUserRequest) (*au.RegisterUserResponse, error)
	Login(body *au.Login) (string, error)
	GetUserByID(s string) (*au.GetUserResponse, error)
	UpdateUser(userID string, req *au.UpdateUserRequest) (*au.UpdateUserResponse, error)
}

type userClient struct {
	host string
}

func NewUserClient(host string) UserClient {
	return &userClient{host: host}
}

func (c *userClient) RegisterUser(data au.RegisterUserRequest) (*au.RegisterUserResponse, error) {
	payload, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to encode request: %w", err)
	}

	resp, err := http.Post(c.host+"/api/v1/auth/register", "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("request to user service failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("%d: %s", resp.StatusCode, string(body))
	}

	var user au.RegisterUserResponse
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("failed to decode user service response: %w", err)
	}
	return &user, nil
}

func (c *userClient) Login(body *au.Login) (string, error) {
	payload, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("failed to encode login request: %w", err)
	}

	resp, err := http.Post(c.host+"/api/v1/auth/login", "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return "", fmt.Errorf("login request to user service failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("login failed: %s", string(bodyBytes))
	}

	var respBody struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return "", fmt.Errorf("failed to decode login response: %w", err)
	}

	return respBody.Token, nil
}

func (c *userClient) GetUserByID(userID string) (*au.GetUserResponse, error) {
	url := fmt.Sprintf("%s/users/%s", c.host, userID)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("user service returned error: %s", string(bodyBytes))
	}

	var user au.GetUserResponse
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("failed to decode user response: %w", err)
	}

	return &user, nil
}

func (c *userClient) UpdateUser(userID string, req *au.UpdateUserRequest) (*au.UpdateUserResponse, error) {
	url := fmt.Sprintf("%s/users/%s", c.host, userID)

	payload, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to encode update request: %w", err)
	}

	reqBody := bytes.NewBuffer(payload)
	httpReq, err := http.NewRequest(http.MethodPut, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create update request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("update request failed: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("update failed: %s", string(bodyBytes))
	}

	var response au.UpdateUserResponse
	if err := json.Unmarshal(bodyBytes, &response); err != nil {
		return nil, fmt.Errorf("failed to decode update response: %w", err)
	}

	return &response, nil
}
