package handlers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"userService/internal/controllers"
	"userService/internal/handlers/dto"
)

type RoleHandler struct {
	roleController *controllers.RoleController
}

func NewRoleHandler(roleController *controllers.RoleController) *RoleHandler {
	return &RoleHandler{roleController: roleController}
}

// GetRoles Получение списка ролей
// @Summary Получение всех ролей
// @Description Возвращает список всех ролей в системе
// @Tags roles
// @Produce json
// @Success 200 {array} models.Role "Список ролей"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /roles/ [get]
func (h *RoleHandler) GetRoles(c *gin.Context) {
	roles, err := h.roleController.GetRoles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, roles)
}

// CreateRole Создание новой роли
// @Summary Создание новой роли
// @Description Добавляет новую роль в систему
// @Tags roles
// @Accept json
// @Produce json
// @Param role body dto.CreateRole true "Данные новой роли"
// @Success 201 {object} dto.CreateRole "Роль создана"
// @Failure 400 {object} map[string]interface{} "Некорректный запрос"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /roles/ [post]
func (h *RoleHandler) CreateRole(c *gin.Context) {
	var role dto.CreateRole
	if err := c.ShouldBindJSON(&role); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.roleController.CreateRole(&role); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, role)
}

// DeleteRole Удаление роли
// @Summary Удалить роль
// @Description Удаляет роль по указанному ID
// @Tags roles
// @Produce json
// @Param role_id path int true "ID роли"
// @Success 200 {object} map[string]interface{} "Роль успешно удалена"
// @Failure 400 {object} map[string]interface{} "Некорректный ID роли"
// @Failure 404 {object} map[string]interface{} "Роль не найдена"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /roles/{role_id} [delete]
func (h *RoleHandler) DeleteRole(c *gin.Context) {
	roleIDParam := c.Param("role_id")
	var roleID int
	if _, err := fmt.Sscanf(roleIDParam, "%d", &roleID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID"})
		return
	}

	if err := h.roleController.DeleteRole(roleID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Role deleted successfully"})
}

// UpdateRolePermissions Обновление permissions роли
// @Summary Обновить permissions роли
// @Description Обновляет список permissions для указанной роли
// @Tags roles
// @Accept json
// @Produce json
// @Param role_id path int true "ID роли"
// @Param request body dto.UpdateRolePermissionsRequest true "Список ID permissions"
// @Success 200 {object} map[string]interface{} "Permissions успешно обновлены"
// @Failure 400 {object} map[string]interface{} "Некорректный запрос"
// @Failure 404 {object} map[string]interface{} "Роль не найдена"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /roles/{role_id}/permissions [patch]
func (h *RoleHandler) UpdateRolePermissions(c *gin.Context) {
	roleIDParam := c.Param("role_id")
	var roleID int
	if _, err := fmt.Sscanf(roleIDParam, "%d", &roleID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID"})
		return
	}

	var req dto.UpdateRolePermissionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}

	if err := h.roleController.UpdateRolePermissions(roleID, req.PermissionIds); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Role permissions updated successfully"})
}
