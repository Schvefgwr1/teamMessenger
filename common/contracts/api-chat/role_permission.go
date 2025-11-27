package api_chat

// ==================== Requests ====================

type CreateRoleRequest struct {
	Name          string `json:"name"`
	PermissionIDs []int  `json:"permissionIds"`
}

type UpdateRolePermissionsRequest struct {
	PermissionIDs []int `json:"permissionIds"`
}

type CreatePermissionRequest struct {
	Name string `json:"name"`
}

// ==================== Responses ====================

type PermissionResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type RoleResponse struct {
	ID          int                  `json:"id"`
	Name        string               `json:"name"`
	Permissions []PermissionResponse `json:"permissions"`
}
