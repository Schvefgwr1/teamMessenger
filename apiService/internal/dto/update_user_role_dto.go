package dto

// UpdateUserRoleRequestGateway - запрос на изменение роли пользователя через API Gateway
type UpdateUserRoleRequestGateway struct {
	RoleID int `json:"role_id" binding:"required"`
}
