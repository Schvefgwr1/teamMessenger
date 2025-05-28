package dto

type UpdateUserRequestGateway struct {
	Username    *string `json:"username"`
	Description *string `json:"description"`
	Gender      *string `json:"gender"`
	Age         *int    `json:"age"`
	RoleID      *int    `json:"roleID"`
}
