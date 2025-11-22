package dto

import (
	au "common/contracts/api-user"
)

// CreateRoleRequestGateway - запрос на создание роли через API Gateway
type CreateRoleRequestGateway struct {
	Name          string `json:"name" binding:"required"`
	Description   string `json:"description"`
	PermissionIds []int  `json:"permissionIds"`
}

// ToCreateRoleRequest преобразует Gateway DTO в контракт
func (r *CreateRoleRequestGateway) ToCreateRoleRequest() *au.CreateRoleRequest {
	return &au.CreateRoleRequest{
		Name:          r.Name,
		Description:   r.Description,
		PermissionIds: r.PermissionIds,
	}
}
