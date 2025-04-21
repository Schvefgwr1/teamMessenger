package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"userService/internal/controllers"
)

type PermissionHandler struct {
	permController *controllers.PermissionController
}

func NewPermissionHandler(permController *controllers.PermissionController) *PermissionHandler {
	return &PermissionHandler{permController: permController}
}

// GetPermissions Получение списка прав
// @Summary Получение всех прав доступа
// @Description Возвращает список всех прав (permissions)
// @Tags permissions
// @Produce json
// @Success 200 {array} models.Permission "Список прав"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /permissions/ [get]
func (h *PermissionHandler) GetPermissions(c *gin.Context) {
	permissions, err := h.permController.GetPermissions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, permissions)
}
