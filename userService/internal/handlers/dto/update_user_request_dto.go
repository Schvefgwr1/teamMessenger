package dto

// UpdateUserRequestSwagger - структура для Swagger документации (аналог au.UpdateUserRequest)
type UpdateUserRequestSwagger struct {
	Username     *string `json:"username"`
	Description  *string `json:"description"`
	Gender       *string `json:"gender"`
	Age          *int    `json:"age"`
	AvatarFileID *int    `json:"avatar"`
	RoleID       *int    `json:"roleID"`
}
