package dto

// UpdateUserRoleRequest - запрос на изменение роли пользователя
type UpdateUserRoleRequest struct {
	RoleID int `json:"role_id" binding:"required"`
}
