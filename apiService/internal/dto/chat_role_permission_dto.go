package dto

// ==================== Requests ====================

type CreateChatRoleRequestGateway struct {
	Name          string `json:"name" binding:"required"`
	PermissionIDs []int  `json:"permissionIds"`
}

type UpdateChatRolePermissionsRequestGateway struct {
	PermissionIDs []int `json:"permissionIds" binding:"required"`
}

type CreateChatPermissionRequestGateway struct {
	Name string `json:"name" binding:"required"`
}

// ==================== Responses ====================

type ChatPermissionResponseGateway struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type ChatRoleResponseGateway struct {
	ID          int                             `json:"id"`
	Name        string                          `json:"name"`
	Permissions []ChatPermissionResponseGateway `json:"permissions"`
}
