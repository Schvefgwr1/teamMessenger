package handlers

import (
	"apiService/internal/controllers"
	"apiService/internal/dto"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ChatRolePermissionHandler struct {
	controller controllers.ChatRolePermissionControllerInterface
}

func NewRolePermissionHandler(controller controllers.ChatRolePermissionControllerInterface) *ChatRolePermissionHandler {
	return &ChatRolePermissionHandler{controller: controller}
}

// ==================== Roles ====================

// GetAllRoles Получение списка всех ролей чатов
// @Summary Получить список всех ролей чата
// @Description Возвращает список всех доступных ролей с их permissions
// @Tags chat-roles
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dto.ChatRoleResponseGateway "Список ролей"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /chat-roles [get]
func (h *ChatRolePermissionHandler) GetAllRoles(c *gin.Context) {
	roles, err := h.controller.GetAllRoles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, roles)
}

// GetRoleByID Получение роли по ID
// @Summary Получить роль по ID
// @Description Возвращает информацию о роли с её permissions
// @Tags chat-roles
// @Produce json
// @Security BearerAuth
// @Param role_id path int true "ID роли"
// @Success 200 {object} dto.ChatRoleResponseGateway "Информация о роли"
// @Failure 400 {object} map[string]interface{} "Некорректный ID роли"
// @Failure 404 {object} map[string]interface{} "Роль не найдена"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /chat-roles/{role_id} [get]
func (h *ChatRolePermissionHandler) GetRoleByID(c *gin.Context) {
	roleID, err := strconv.Atoi(c.Param("role_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role ID"})
		return
	}

	role, err := h.controller.GetRoleByID(roleID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "role not found"})
		return
	}
	c.JSON(http.StatusOK, role)
}

// CreateRole Создание новой роли
// @Summary Создать новую роль
// @Description Создает новую роль с указанными permissions
// @Tags chat-roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param role body dto.CreateChatRoleRequestGateway true "Данные для создания роли"
// @Success 201 {object} dto.ChatRoleResponseGateway "Созданная роль"
// @Failure 400 {object} map[string]interface{} "Некорректный запрос"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /chat-roles [post]
func (h *ChatRolePermissionHandler) CreateRole(c *gin.Context) {
	var req dto.CreateChatRoleRequestGateway
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	role, err := h.controller.CreateRole(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, role)
}

// DeleteRole Удаление роли
// @Summary Удалить роль
// @Description Удаляет роль по указанному ID
// @Tags chat-roles
// @Produce json
// @Security BearerAuth
// @Param role_id path int true "ID роли"
// @Success 204 "Роль успешно удалена"
// @Failure 400 {object} map[string]interface{} "Некорректный ID роли"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /chat-roles/{role_id} [delete]
func (h *ChatRolePermissionHandler) DeleteRole(c *gin.Context) {
	roleID, err := strconv.Atoi(c.Param("role_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role ID"})
		return
	}

	if err := h.controller.DeleteRole(roleID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// UpdateRolePermissions Обновление permissions роли
// @Summary Обновить permissions роли
// @Description Полностью заменяет список permissions у роли
// @Tags chat-roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param role_id path int true "ID роли"
// @Param permissions body dto.UpdateChatRolePermissionsRequestGateway true "Новый список permission IDs"
// @Success 200 {object} dto.ChatRoleResponseGateway "Обновленная роль"
// @Failure 400 {object} map[string]interface{} "Некорректный запрос"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /chat-roles/{role_id}/permissions [patch]
func (h *ChatRolePermissionHandler) UpdateRolePermissions(c *gin.Context) {
	roleID, err := strconv.Atoi(c.Param("role_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role ID"})
		return
	}

	var req dto.UpdateChatRolePermissionsRequestGateway
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	role, err := h.controller.UpdateRolePermissions(roleID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, role)
}

// ==================== Permissions ====================

// GetAllPermissions Получение списка всех permissions
// @Summary Получить список всех permissions
// @Description Возвращает список всех доступных permissions для чатов
// @Tags chat-permissions
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dto.ChatPermissionResponseGateway "Список permissions"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /chat-permissions [get]
func (h *ChatRolePermissionHandler) GetAllPermissions(c *gin.Context) {
	permissions, err := h.controller.GetAllPermissions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, permissions)
}

// CreatePermission Создание нового permission
// @Summary Создать новый permission
// @Description Создает новый permission для чатов
// @Tags chat-permissions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param permission body dto.CreateChatPermissionRequestGateway true "Данные для создания permission"
// @Success 201 {object} dto.ChatPermissionResponseGateway "Созданный permission"
// @Failure 400 {object} map[string]interface{} "Некорректный запрос"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /chat-permissions [post]
func (h *ChatRolePermissionHandler) CreatePermission(c *gin.Context) {
	var req dto.CreateChatPermissionRequestGateway
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	permission, err := h.controller.CreatePermission(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, permission)
}

// DeletePermission Удаление permission
// @Summary Удалить permission
// @Description Удаляет permission по указанному ID
// @Tags chat-permissions
// @Produce json
// @Security BearerAuth
// @Param permission_id path int true "ID permission"
// @Success 204 "Permission успешно удален"
// @Failure 400 {object} map[string]interface{} "Некорректный ID permission"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /chat-permissions/{permission_id} [delete]
func (h *ChatRolePermissionHandler) DeletePermission(c *gin.Context) {
	permissionID, err := strconv.Atoi(c.Param("permission_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid permission ID"})
		return
	}

	if err := h.controller.DeletePermission(permissionID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
