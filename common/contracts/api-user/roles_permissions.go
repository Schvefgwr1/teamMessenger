package api_user

// CreateRoleRequest - запрос на создание роли
type CreateRoleRequest struct {
	Name          string `json:"name" binding:"required"`
	Description   string `json:"description"`
	PermissionIds []int  `json:"permissionIds"`
}
