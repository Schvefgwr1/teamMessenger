package controllers

import (
	"userService/internal/models"
	"userService/internal/repositories"
)

type PermissionController struct {
	permRepo *repositories.PermissionRepository
}

func NewPermissionController(permRepo *repositories.PermissionRepository) *PermissionController {
	return &PermissionController{permRepo: permRepo}
}

func (c *PermissionController) GetPermissions() ([]models.Permission, error) {
	return c.permRepo.GetAllPermissions()
}
