package http_clients

import (
	"apiService/internal/dto"
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
	UpdateUserRole(userID string, roleID int) error
	UpdateRolePermissions(roleID int, permissionIDs []int) error
	DeleteRole(roleID int) error
	GetUserBrief(userID, chatID, requesterID string) (*dto.UserBriefResponse, error)
	SearchUsers(query string, limit int) (*dto.UserSearchResponse, error)
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

// UpdateUserRole - изменить роль пользователя
func (c *userClient) UpdateUserRole(userID string, roleID int) error {
	url := fmt.Sprintf("%s/api/v1/users/%s/role", c.host, userID)

	reqBody, err := json.Marshal(map[string]int{"role_id": roleID})
	if err != nil {
		return fmt.Errorf("failed to encode request: %w", err)
	}

	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("request to user service failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("user service returned error: %s", string(bodyBytes))
	}

	return nil
}

// DeleteRole - удалить роль
func (c *userClient) DeleteRole(roleID int) error {
	url := fmt.Sprintf("%s/api/v1/roles/%d", c.host, roleID)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("request to user service failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("user service returned error: %s", string(bodyBytes))
	}

	return nil
}

// UpdateRolePermissions - обновить permissions роли
func (c *userClient) UpdateRolePermissions(roleID int, permissionIDs []int) error {
	url := fmt.Sprintf("%s/api/v1/roles/%d/permissions", c.host, roleID)

	reqBody, err := json.Marshal(map[string]interface{}{
		"permission_ids": permissionIDs,
	})
	if err != nil {
		return fmt.Errorf("failed to encode request: %w", err)
	}

	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("request to user service failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("user service returned error: %s", string(bodyBytes))
	}

	return nil
}

// GetUserBrief - получить краткую информацию о пользователе
func (c *userClient) GetUserBrief(userID, chatID, requesterID string) (*dto.UserBriefResponse, error) {
	url := fmt.Sprintf("%s/api/v1/users/%s/brief?chatId=%s", c.host, userID, chatID)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Передаем ID запрашивающего пользователя
	req.Header.Set("X-User-ID", requesterID)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get user brief: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("user service returned error: %s", string(bodyBytes))
	}

	// Используем промежуточную структуру для десериализации, т.к. AvatarFile может быть interface{}
	var userBrief dto.UserBriefResponse
	if err := json.NewDecoder(resp.Body).Decode(&userBrief); err != nil {
		return nil, fmt.Errorf("failed to decode user brief response: %w", err)
	}

	return &userBrief, nil
}

// SearchUsers - поиск пользователей по имени или email
func (c *userClient) SearchUsers(query string, limit int) (*dto.UserSearchResponse, error) {
	url := fmt.Sprintf("%s/api/v1/users/search?q=%s&limit=%d", c.host, query, limit)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to search users: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("user service returned error: %s", string(bodyBytes))
	}

	var result dto.UserSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode search response: %w", err)
	}

	return &result, nil
}
