package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"taskService/internal/controllers"
)

type TaskStatusHandler struct {
	Controller controllers.TaskStatusControllerInterface
}

func NewTaskStatusHandler(controller controllers.TaskStatusControllerInterface) *TaskStatusHandler {
	return &TaskStatusHandler{Controller: controller}
}

type CreateStatusDTO struct {
	Name string `json:"name" binding:"required"`
}

// Create Создание нового статуса задачи
// @Summary Создать новый статус задачи
// @Description Создает новый статус задачи с указанным названием
// @Tags task-statuses
// @Accept json
// @Produce json
// @Param status body CreateStatusDTO true "Данные для создания статуса"
// @Success 201 {object} models.TaskStatus "Статус успешно создан"
// @Failure 400 {object} map[string]interface{} "Некорректный запрос или статус уже существует"
// @Router /tasks/statuses [post]
func (h *TaskStatusHandler) Create(c *gin.Context) {
	var dto CreateStatusDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	status, err := h.Controller.Create(dto.Name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, status)
}

// GetByID Получение статуса задачи по ID
// @Summary Получить статус задачи по ID
// @Description Возвращает информацию о статусе задачи по его ID
// @Tags task-statuses
// @Produce json
// @Param id path int true "ID статуса"
// @Success 200 {object} models.TaskStatus "Информация о статусе"
// @Failure 400 {object} map[string]interface{} "Некорректный ID"
// @Failure 404 {object} map[string]interface{} "Статус не найден"
// @Router /tasks/statuses/{id} [get]
func (h *TaskStatusHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	status, err := h.Controller.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, status)
}

// DeleteByID Удаление статуса задачи
// @Summary Удалить статус задачи
// @Description Удаляет статус задачи по его ID
// @Tags task-statuses
// @Produce json
// @Param id path int true "ID статуса"
// @Success 204 "Статус успешно удален"
// @Failure 400 {object} map[string]interface{} "Некорректный ID"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /tasks/statuses/{id} [delete]
func (h *TaskStatusHandler) DeleteByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	err = h.Controller.DeleteByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not delete task status"})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetAll Получение всех статусов задач
// @Summary Получить все статусы задач
// @Description Возвращает список всех доступных статусов задач
// @Tags task-statuses
// @Produce json
// @Success 200 {array} models.TaskStatus "Список статусов задач"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /tasks/statuses [get]
func (h *TaskStatusHandler) GetAll(c *gin.Context) {
	statuses, err := h.Controller.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not get statuses"})
		return
	}

	c.JSON(http.StatusOK, statuses)
}
