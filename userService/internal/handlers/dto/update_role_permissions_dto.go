package dto

// UpdateRolePermissionsRequest - запрос на обновление permissions роли
type UpdateRolePermissionsRequest struct {
	PermissionIds []int `json:"permission_ids" binding:"required"`
}
