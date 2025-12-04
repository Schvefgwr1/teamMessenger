package dto

// UserRoleResponse - ответ с ролью пользователя в чате (только имя)
type UserRoleResponse struct {
	RoleName string `json:"roleName"`
}

// ChatPermissionResponse - permission в ответе
type ChatPermissionResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// UserRoleWithPermissionsResponse - ответ с ролью и permissions пользователя в чате
type UserRoleWithPermissionsResponse struct {
	RoleID      int                      `json:"roleId"`
	RoleName    string                   `json:"roleName"`
	Permissions []ChatPermissionResponse `json:"permissions"`
}

// ChatMemberResponse - участник чата
type ChatMemberResponse struct {
	UserID   string `json:"userId"`
	RoleID   int    `json:"roleId"`
	RoleName string `json:"roleName"`
}
