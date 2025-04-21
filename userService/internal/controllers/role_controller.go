package controllers

import (
	"log"
	"userService/internal/handlers/dto"
	"userService/internal/models"
	"userService/internal/repositories"
)

type RoleController struct {
	roleRepo       *repositories.RoleRepository
	permissionRepo *repositories.PermissionRepository
}

func NewRoleController(
	roleRepo *repositories.RoleRepository,
	permissionRepo *repositories.PermissionRepository,
) *RoleController {
	return &RoleController{roleRepo: roleRepo, permissionRepo: permissionRepo}
}

func (c *RoleController) GetRoles() ([]models.Role, error) {
	return c.roleRepo.GetAllRoles()
}

func (c *RoleController) CreateRole(roleDTO *dto.CreateRole) error {
	role := &models.Role{
		Name:        roleDTO.Name,
		Description: roleDTO.Description,
		Permissions: make([]models.Permission, 0),
	}
	for _, permId := range roleDTO.PermissionIds {
		var perm *models.Permission
		perm, err := c.permissionRepo.GetPermissionById(permId)
		if err != nil {
			log.Default().Printf("Can't find permission with id: %d, error: %v", permId, err)
			continue
		}
		role.Permissions = append(role.Permissions, *perm)
	}
	return c.roleRepo.CreateRole(role)
}
