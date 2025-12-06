package dto

// UpdateRolePermissionsRequestGateway - запрос на обновление permissions роли через API Gateway
type UpdateRolePermissionsRequestGateway struct {
	PermissionIds []int `json:"permission_ids" binding:"required"`
}
