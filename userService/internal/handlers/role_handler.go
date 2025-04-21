package handlers

import (
	"github.com/gin-gonic/gin"
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
