package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"taskService/internal/controllers"
	"taskService/internal/custom_errors"
	"taskService/internal/handlers/dto"
)

type TaskHandler struct {
	TaskController *controllers.TaskController
}

func NewTaskHandler(controller *controllers.TaskController) *TaskHandler {
	return &TaskHandler{TaskController: controller}
}

// CreateTask Создание новой задачи
// @Summary Создать новую задачу
// @Description Создает новую задачу с указанными параметрами, исполнителем и прикрепленными файлами
// @Tags tasks
// @Accept json
// @Produce json
// @Param task body dto.CreateTaskDTO true "Данные для создания задачи"
// @Success 200 {object} models.Task "Задача успешно создана"
// @Failure 400 {object} map[string]interface{} "Некорректный запрос или статус не найден"
// @Failure 502 {object} map[string]interface{} "Ошибка при обращении к внешнему сервису"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /tasks [post]
func (h *TaskHandler) CreateTask(c *gin.Context) {
	var taskDTO dto.CreateTaskDTO
	if err := c.ShouldBindJSON(&taskDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	task, err := h.TaskController.Create(&taskDTO)
	if err != nil {
		var userErr *custom_errors.GetUserHTTPError
		var chatErr *custom_errors.GetChatHTTPError
		var fileErr *custom_errors.GetFileHTTPError
		var statusErr *custom_errors.TaskStatusNotFoundError

		switch {
		case errors.As(err, &userErr),
			errors.As(err, &chatErr),
			errors.As(err, &fileErr):
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		case errors.As(err, &statusErr):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, task)
}

// UpdateTaskStatus Обновление статуса задачи
// @Summary Обновить статус задачи
// @Description Изменяет статус задачи на указанный
// @Tags tasks
// @Produce json
// @Param task_id path int true "ID задачи"
// @Param status_id path int true "ID статуса"
// @Success 200 "Статус задачи успешно обновлен"
// @Failure 400 {object} map[string]interface{} "Некорректный ID задачи или статуса, или статус/задача не найдены"
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

	err = h.TaskController.UpdateStatus(taskID, statusID)
	if err != nil {
		var statusErr *custom_errors.TaskStatusNotFoundError
		var taskErr *custom_errors.TaskNotFoundError
		if errors.As(err, &statusErr) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			if errors.As(err, &taskErr) {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			}
		}
		return
	}

	c.Status(http.StatusOK)
}

// GetTaskByID Получение задачи по ID
// @Summary Получить задачу по ID
// @Description Возвращает информацию о задаче и прикрепленных файлах по ID
// @Tags tasks
// @Produce json
// @Param task_id path int true "ID задачи"
// @Success 200 {object} dto.TaskResponse "Информация о задаче"
// @Failure 400 {object} map[string]interface{} "Некорректный ID задачи"
// @Failure 502 {object} map[string]interface{} "Ошибка при обращении к сервису файлов"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /tasks/{task_id} [get]
func (h *TaskHandler) GetTaskByID(c *gin.Context) {
	taskID, err := strconv.Atoi(c.Param("task_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task ID"})
		return
	}

	task, err := h.TaskController.GetByID(taskID)
	if err != nil {
		var httpClientError *custom_errors.GetFileHTTPError
		if errors.As(err, &httpClientError) {
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, task)
}

// GetUserTasks Получение списка задач пользователя
// @Summary Получить список задач пользователя
// @Description Возвращает список задач указанного пользователя с пагинацией
// @Tags tasks
// @Produce json
// @Param user_id path string true "UUID пользователя"
// @Param limit query int false "Количество задач на странице" default(20)
// @Param offset query int false "Смещение для пагинации" default(0)
// @Success 200 {array} dto.TaskToList "Список задач пользователя"
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
	if err != nil || limit <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit"})
		return
	}

	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil || offset < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid offset"})
		return
	}

	tasks, err := h.TaskController.GetUserTasks(userID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, tasks)
}
