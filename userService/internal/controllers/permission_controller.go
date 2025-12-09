package controllers

import (
	"userService/internal/models"
	"userService/internal/repositories"
)

type PermissionController struct {
	permRepo repositories.PermissionRepositoryInterface
}

func NewPermissionController(permRepo repositories.PermissionRepositoryInterface) *PermissionController {
	return &PermissionController{permRepo: permRepo}
}

func (c *PermissionController) GetPermissions() ([]models.Permission, error) {
	return c.permRepo.GetAllPermissions()
}
