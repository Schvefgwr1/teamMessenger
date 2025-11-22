package handlers

import (
	"apiService/internal/controllers"
	"apiService/internal/dto"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"strconv"
)

type TaskHandler struct {
	taskController *controllers.TaskController
}

func NewTaskHandler(taskController *controllers.TaskController) *TaskHandler {
	return &TaskHandler{taskController: taskController}
}

func getUserIDFromTaskContext(c *gin.Context) (uuid.UUID, error) {
	userID, exists := c.Get("userID")
	if !exists {
		return uuid.Nil, fmt.Errorf("user ID not found in context")
	}
	return userID.(uuid.UUID), nil
}

// CreateTask Создание новой задачи
// @Summary Создать новую задачу
// @Description Создает новую задачу с указанными параметрами, исполнителем и прикрепленными файлами
// @Tags tasks
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param title formData string true "Название задачи"
// @Param description formData string false "Описание задачи"
// @Param executor_id formData string false "UUID исполнителя задачи"
// @Param chat_id formData string false "UUID чата, связанного с задачей"
// @Param files formData []file false "Прикрепленные файлы"
// @Success 201 {object} map[string]interface{} "Задача успешно создана"
// @Failure 400 {object} map[string]interface{} "Некорректный запрос"
// @Failure 401 {object} map[string]interface{} "Пользователь не аутентифицирован"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /tasks [post]
func (h *TaskHandler) CreateTask(c *gin.Context) {
	userID, err := getUserIDFromTaskContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var req dto.CreateTaskRequestGateway
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task, err := h.taskController.CreateTask(&req, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, task)
}

// UpdateTaskStatus Обновление статуса задачи
// @Summary Обновить статус задачи
// @Description Изменяет статус задачи на указанный
// @Tags tasks
// @Produce json
// @Security BearerAuth
// @Param task_id path int true "ID задачи"
// @Param status_id path int true "ID статуса"
// @Success 200 "Статус задачи успешно обновлен"
// @Failure 400 {object} map[string]interface{} "Некорректный ID задачи или статуса"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /tasks/{task_id}/status/{status_id} [patch]
func (h *TaskHandler) UpdateTaskStatus(c *gin.Context) {
	taskID, err := strconv.Atoi(c.Param("task_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task ID"})
		return
	}

	statusID, err := strconv.Atoi(c.Param("status_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status ID"})
		return
	}

	err = h.taskController.UpdateTaskStatus(taskID, statusID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

// GetTaskByID Получение задачи по ID
// @Summary Получить задачу по ID
// @Description Возвращает информацию о задаче и прикрепленных файлах по ID
// @Tags tasks
// @Produce json
// @Security BearerAuth
// @Param task_id path int true "ID задачи"
// @Success 200 {object} map[string]interface{} "Информация о задаче"
// @Failure 400 {object} map[string]interface{} "Некорректный ID задачи"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /tasks/{task_id} [get]
func (h *TaskHandler) GetTaskByID(c *gin.Context) {
	taskID, err := strconv.Atoi(c.Param("task_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task ID"})
		return
	}

	task, err := h.taskController.GetTaskByID(taskID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, task)
}

// GetUserTasks Получение списка задач пользователя
// @Summary Получить список задач пользователя
// @Description Возвращает список задач указанного пользователя с пагинацией
// @Tags tasks
// @Produce json
// @Security BearerAuth
// @Param user_id path string true "UUID пользователя"
// @Param limit query int false "Количество задач на странице" default(20) maximum(100)
// @Param offset query int false "Смещение для пагинации" default(0)
// @Success 200 {array} map[string]interface{} "Список задач пользователя"
// @Failure 400 {object} map[string]interface{} "Некорректный UUID пользователя или параметры пагинации"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /users/{user_id}/tasks [get]
func (h *TaskHandler) GetUserTasks(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user ID is required"})
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if err != nil || limit <= 0 || limit > 100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit"})
		return
	}

	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil || offset < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid offset"})
		return
	}

	tasks, err := h.taskController.GetUserTasks(userID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

// GetAllStatuses Получение всех статусов задач
// @Summary Получить все статусы задач
// @Description Возвращает список всех доступных статусов задач
// @Tags tasks
// @Produce json
// @Security BearerAuth
// @Success 200 {array} map[string]interface{} "Список статусов задач"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /tasks/statuses [get]
func (h *TaskHandler) GetAllStatuses(c *gin.Context) {
	statuses, err := h.taskController.GetAllStatuses()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, statuses)
}

// CreateStatus Создание нового статуса задачи
// @Summary Создать новый статус задачи
// @Description Создает новый статус задачи с указанным названием
// @Tags tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateStatusRequestGateway true "Данные для создания статуса"
// @Success 201 {object} map[string]interface{} "Статус успешно создан"
// @Failure 400 {object} map[string]interface{} "Некорректный запрос"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /tasks/statuses [post]
func (h *TaskHandler) CreateStatus(c *gin.Context) {
	var req dto.CreateStatusRequestGateway
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	status, err := h.taskController.CreateStatus(req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, status)
}

// GetStatusByID Получение статуса задачи по ID
// @Summary Получить статус задачи по ID
// @Description Возвращает информацию о статусе задачи по ID
// @Tags tasks
// @Produce json
// @Security BearerAuth
// @Param status_id path int true "ID статуса"
// @Success 200 {object} map[string]interface{} "Информация о статусе"
// @Failure 400 {object} map[string]interface{} "Некорректный ID статуса"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /tasks/statuses/{status_id} [get]
func (h *TaskHandler) GetStatusByID(c *gin.Context) {
	statusID, err := strconv.Atoi(c.Param("status_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status ID"})
		return
	}

	status, err := h.taskController.GetStatusByID(statusID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, status)
}

// DeleteStatus Удаление статуса задачи
// @Summary Удалить статус задачи
// @Description Удаляет статус задачи по ID
// @Tags tasks
// @Produce json
// @Security BearerAuth
// @Param status_id path int true "ID статуса"
// @Success 200 {object} map[string]interface{} "Статус успешно удален"
// @Failure 400 {object} map[string]interface{} "Некорректный ID статуса"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /tasks/statuses/{status_id} [delete]
func (h *TaskHandler) DeleteStatus(c *gin.Context) {
	statusID, err := strconv.Atoi(c.Param("status_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status ID"})
		return
	}

	err = h.taskController.DeleteStatus(statusID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "status deleted successfully"})
}
