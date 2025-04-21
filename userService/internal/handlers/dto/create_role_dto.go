package dto

type CreateRole struct {
	Name          string `json:"name"`
	Description   string `json:"description"`
	PermissionIds []int  `json:"permissionIds"`
}
