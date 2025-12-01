package http_clients

import (
	"bytes"
	au "common/contracts/api-user"
	uc "common/contracts/user-contracts"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"

	"github.com/google/uuid"
)

type UserClient interface {
	RegisterUser(data au.RegisterUserRequest) (*au.RegisterUserResponse, error)
	Login(body *au.Login) (string, uuid.UUID, error)
	GetUserByID(s string) (*au.GetUserResponse, error)
	UpdateUser(userID string, req *au.UpdateUserRequest) (*au.UpdateUserResponse, error)
	GetPublicKey() (*rsa.PublicKey, error)
	GetAllPermissions() ([]*uc.Permission, error)
	GetAllRoles() ([]*uc.Role, error)
	CreateRole(req *au.CreateRoleRequest) (*uc.Role, error)
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

func (c *userClient) Login(body *au.Login) (string, uuid.UUID, error) {
	payload, err := json.Marshal(body)
	if err != nil {
		return "", uuid.Nil, fmt.Errorf("failed to encode login request: %w", err)
	}

	resp, err := http.Post(c.host+"/api/v1/auth/login", "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return "", uuid.Nil, fmt.Errorf("login request to user service failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", uuid.Nil, fmt.Errorf("login failed: %s", string(bodyBytes))
	}

	var respBody struct {
		Token  string    `json:"token"`
		UserID uuid.UUID `json:"userID"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return "", uuid.Nil, fmt.Errorf("failed to decode login response: %w", err)
	}

	return respBody.Token, respBody.UserID, nil
}

func (c *userClient) GetUserByID(userID string) (*au.GetUserResponse, error) {
	url := fmt.Sprintf("%s/api/v1/users/%s", c.host, userID)

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
	url := fmt.Sprintf("%s/api/v1/users/%s", c.host, userID)

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

func (c *userClient) GetPublicKey() (*rsa.PublicKey, error) {
	resp, err := http.Get(c.host + "/api/v1/keys/public")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch public key: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to fetch public key: %s", string(body))
	}

	var result struct {
		Key struct {
			N *big.Int `json:"N"`
			E int      `json:"E"`
		} `json:"key"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("invalid response body: %w", err)
	}

	return &rsa.PublicKey{
		N: result.Key.N,
		E: result.Key.E,
	}, nil
}

// GetAllPermissions - получить все права
func (c *userClient) GetAllPermissions() ([]*uc.Permission, error) {
	url := fmt.Sprintf("%s/api/v1/permissions/", c.host)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get permissions: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("user service returned error: %s", string(bodyBytes))
	}

	var permissions []*uc.Permission
	if err := json.NewDecoder(resp.Body).Decode(&permissions); err != nil {
		return nil, fmt.Errorf("failed to decode permissions response: %w", err)
	}

	return permissions, nil
}

// GetAllRoles - получить все роли
func (c *userClient) GetAllRoles() ([]*uc.Role, error) {
	url := fmt.Sprintf("%s/api/v1/roles/", c.host)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get roles: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("user service returned error: %s", string(bodyBytes))
	}

	var roles []*uc.Role
	if err := json.NewDecoder(resp.Body).Decode(&roles); err != nil {
		return nil, fmt.Errorf("failed to decode roles response: %w", err)
	}

	return roles, nil
}

// CreateRole - создать новую роль
func (c *userClient) CreateRole(req *au.CreateRoleRequest) (*uc.Role, error) {
	url := fmt.Sprintf("%s/api/v1/roles/", c.host)

	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to encode request: %w", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("request to user service failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("user service returned error: %s", string(bodyBytes))
	}

	var role uc.Role
	if err := json.NewDecoder(resp.Body).Decode(&role); err != nil {
		return nil, fmt.Errorf("failed to decode role response: %w", err)
	}

	return &role, nil
}
