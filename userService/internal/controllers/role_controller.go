package controllers

import (
	"fmt"
	"log"
	"userService/internal/handlers/dto"
	"userService/internal/models"
	"userService/internal/repositories"
)

type RoleController struct {
	roleRepo       repositories.RoleRepositoryInterface
	permissionRepo repositories.PermissionRepositoryInterface
}

func NewRoleController(
	roleRepo repositories.RoleRepositoryInterface,
	permissionRepo repositories.PermissionRepositoryInterface,
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

func (c *RoleController) DeleteRole(roleID int) error {
	// Проверяем, существует ли роль
	_, err := c.roleRepo.GetRoleByID(roleID)
	if err != nil {
		return err
	}

	return c.roleRepo.DeleteRole(roleID)
}

func (c *RoleController) UpdateRolePermissions(roleID int, permissionIDs []int) error {
	// Проверяем, существует ли роль
	_, err := c.roleRepo.GetRoleByID(roleID)
	if err != nil {
		return err
	}

	// Проверяем, что все permissions существуют
	for _, permID := range permissionIDs {
		_, err := c.permissionRepo.GetPermissionById(permID)
		if err != nil {
			log.Default().Printf("Can't find permission with id: %d, error: %v", permID, err)
			return fmt.Errorf("permission with id %d not found", permID)
		}
	}

	return c.roleRepo.UpdateRolePermissions(roleID, permissionIDs)
}
