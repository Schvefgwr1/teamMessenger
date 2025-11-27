package http_clients

import (
	"bytes"
	ac "common/contracts/api-chat"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type ChatRolePermissionClient interface {
	// Roles
	GetAllRoles() ([]ac.RoleResponse, error)
	GetRoleByID(roleID int) (*ac.RoleResponse, error)
	CreateRole(req *ac.CreateRoleRequest) (*ac.RoleResponse, error)
	DeleteRole(roleID int) error
	UpdateRolePermissions(roleID int, req *ac.UpdateRolePermissionsRequest) (*ac.RoleResponse, error)

	// Permissions
	GetAllPermissions() ([]ac.PermissionResponse, error)
	CreatePermission(req *ac.CreatePermissionRequest) (*ac.PermissionResponse, error)
	DeletePermission(permissionID int) error
}

type rolePermissionClient struct {
	host string
}

func NewRolePermissionClient(host string) ChatRolePermissionClient {
	return &rolePermissionClient{host: host}
}

// ==================== Roles ====================

func (c *rolePermissionClient) GetAllRoles() ([]ac.RoleResponse, error) {
	url := fmt.Sprintf("%s/api/v1/chat-roles", c.host)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get roles: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("chat service returned error: %s", string(bodyBytes))
	}

	var roles []ac.RoleResponse
	if err := json.NewDecoder(resp.Body).Decode(&roles); err != nil {
		return nil, fmt.Errorf("failed to decode roles response: %w", err)
	}

	return roles, nil
}

func (c *rolePermissionClient) GetRoleByID(roleID int) (*ac.RoleResponse, error) {
	url := fmt.Sprintf("%s/api/v1/chat-roles/%d", c.host, roleID)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get role: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("chat service returned error: %s", string(bodyBytes))
	}

	var role ac.RoleResponse
	if err := json.NewDecoder(resp.Body).Decode(&role); err != nil {
		return nil, fmt.Errorf("failed to decode role response: %w", err)
	}

	return &role, nil
}

func (c *rolePermissionClient) CreateRole(req *ac.CreateRoleRequest) (*ac.RoleResponse, error) {
	url := fmt.Sprintf("%s/api/v1/chat-roles", c.host)

	payload, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to encode request: %w", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("request to chat service failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("chat service returned error: %s", string(bodyBytes))
	}

	var role ac.RoleResponse
	if err := json.NewDecoder(resp.Body).Decode(&role); err != nil {
		return nil, fmt.Errorf("failed to decode role response: %w", err)
	}

	return &role, nil
}

func (c *rolePermissionClient) DeleteRole(roleID int) error {
	url := fmt.Sprintf("%s/api/v1/chat-roles/%d", c.host, roleID)

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

func (c *rolePermissionClient) UpdateRolePermissions(roleID int, req *ac.UpdateRolePermissionsRequest) (*ac.RoleResponse, error) {
	url := fmt.Sprintf("%s/api/v1/chat-roles/%d/permissions", c.host, roleID)

	payload, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to encode request: %w", err)
	}

	httpReq, err := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(payload))
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

	var role ac.RoleResponse
	if err := json.NewDecoder(resp.Body).Decode(&role); err != nil {
		return nil, fmt.Errorf("failed to decode role response: %w", err)
	}

	return &role, nil
}

// ==================== Permissions ====================

func (c *rolePermissionClient) GetAllPermissions() ([]ac.PermissionResponse, error) {
	url := fmt.Sprintf("%s/api/v1/chat-permissions", c.host)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get permissions: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("chat service returned error: %s", string(bodyBytes))
	}

	var permissions []ac.PermissionResponse
	if err := json.NewDecoder(resp.Body).Decode(&permissions); err != nil {
		return nil, fmt.Errorf("failed to decode permissions response: %w", err)
	}

	return permissions, nil
}

func (c *rolePermissionClient) CreatePermission(req *ac.CreatePermissionRequest) (*ac.PermissionResponse, error) {
	url := fmt.Sprintf("%s/api/v1/chat-permissions", c.host)

	payload, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to encode request: %w", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("request to chat service failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("chat service returned error: %s", string(bodyBytes))
	}

	var permission ac.PermissionResponse
	if err := json.NewDecoder(resp.Body).Decode(&permission); err != nil {
		return nil, fmt.Errorf("failed to decode permission response: %w", err)
	}

	return &permission, nil
}

func (c *rolePermissionClient) DeletePermission(permissionID int) error {
	url := fmt.Sprintf("%s/api/v1/chat-permissions/%d", c.host, permissionID)

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
