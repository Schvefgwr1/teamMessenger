package dto

// ==================== Requests ====================

type CreateRoleRequest struct {
	Name          string `json:"name" binding:"required"`
	PermissionIDs []int  `json:"permissionIds"`
}

type UpdateRolePermissionsRequest struct {
	PermissionIDs []int `json:"permissionIds" binding:"required"`
}

type CreatePermissionRequest struct {
	Name string `json:"name" binding:"required"`
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
