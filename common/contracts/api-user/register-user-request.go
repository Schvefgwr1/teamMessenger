package api_user

type RegisterUserRequest struct {
	Username    string  `json:"username" binding:"required"`
	Email       string  `json:"email" binding:"required,email"`
	Password    string  `json:"password" binding:"required,min=6"`
	Description *string `json:"description" binding:""`
	Gender      string  `json:"gender"`
	Age         int     `json:"age" binding:"required"`
	AvatarID    *int    `json:"avatar"`
	RoleID      int     `json:"roleID" binding:"required"`
}
