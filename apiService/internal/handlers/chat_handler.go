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

type ChatHandler struct {
	chatController *controllers.ChatController
}

func NewChatHandler(chatController *controllers.ChatController) *ChatHandler {
	return &ChatHandler{chatController: chatController}
}

func getUserIDFromContext(c *gin.Context) (uuid.UUID, error) {
	userID, exists := c.Get("userID")
	if !exists {
		return uuid.Nil, fmt.Errorf("user ID not found in context")
	}
	return userID.(uuid.UUID), nil
}

// GetUserChats Получение списка чатов пользователя
// @Summary Получить список чатов пользователя
// @Description Возвращает список всех чатов, в которых участвует указанный пользователь
// @Tags chats
// @Produce json
// @Security BearerAuth
// @Param user_id path string true "UUID пользователя"
// @Success 200 {array} map[string]interface{} "Список чатов пользователя"
// @Failure 400 {object} map[string]interface{} "Некорректный UUID пользователя"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /chats/{user_id} [get]
func (h *ChatHandler) GetUserChats(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	chats, err := h.chatController.GetUserChats(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, chats)
}

// CreateChat Создание нового чата
// @Summary Создать новый чат
// @Description Создает новый чат с указанными параметрами и участниками
// @Tags chats
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param name formData string true "Название чата"
// @Param description formData string false "Описание чата"
// @Param ownerID formData string true "UUID владельца чата"
// @Param userIDs formData []string true "UUID участников чата"
// @Param avatar formData file false "Аватар чата"
// @Success 201 {object} dto.CreateChatResponse "Чат успешно создан"
// @Failure 400 {object} map[string]interface{} "Некорректный запрос или неверный формат UUID"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /chats [post]
func (h *ChatHandler) CreateChat(c *gin.Context) {
	var req dto.CreateChatRequestGateway
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Парсим UUID из строк
	ownerID, userIDs, err := req.ParseUUIDs()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid UUID format: " + err.Error()})
		return
	}

	// Создаем новый request с преобразованными UUID
	createReq := dto.CreateChatRequestGateway{
		Name:        req.Name,
		Description: req.Description,
		Avatar:      req.Avatar,
	}

	chat, err := h.chatController.CreateChat(&createReq, ownerID, userIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, chat)
}

// GetChatMessages Получение сообщений чата
// @Summary Получить сообщения чата
// @Description Возвращает список сообщений из указанного чата с пагинацией
// @Tags chats
// @Produce json
// @Security BearerAuth
// @Param chat_id path string true "UUID чата"
// @Param offset query int false "Смещение для пагинации" default(0)
// @Param limit query int false "Количество сообщений на странице" default(20) maximum(100)
// @Success 200 {array} map[string]interface{} "Список сообщений"
// @Failure 400 {object} map[string]interface{} "Некорректный UUID чата или параметры пагинации"
// @Failure 401 {object} map[string]interface{} "Пользователь не аутентифицирован"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /chats/messages/{chat_id} [get]
func (h *ChatHandler) GetChatMessages(c *gin.Context) {
	chatID, err := uuid.Parse(c.Param("chat_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid chat ID"})
		return
	}

	userID, err := getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil || offset < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid offset"})
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if err != nil || limit <= 0 || limit > 100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit"})
		return
	}

	messages, err := h.chatController.GetChatMessages(chatID, userID, offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, messages)
}

// SendMessage Отправка сообщения в чат
// @Summary Отправить сообщение в чат
// @Description Отправляет новое сообщение в указанный чат с возможностью прикрепления файлов
// @Tags chats
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param chat_id path string true "UUID чата"
// @Param content formData string true "Текст сообщения"
// @Param files formData []file false "Прикрепленные файлы"
// @Success 201 {object} map[string]interface{} "Сообщение успешно отправлено"
// @Failure 400 {object} map[string]interface{} "Некорректный запрос или неверный UUID чата"
// @Failure 401 {object} map[string]interface{} "Пользователь не аутентифицирован"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /chats/messages/{chat_id} [post]
func (h *ChatHandler) SendMessage(c *gin.Context) {
	chatID, err := uuid.Parse(c.Param("chat_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid chat ID"})
		return
	}

	userID, err := getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var req dto.SendMessageRequestGateway
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	message, err := h.chatController.SendMessage(chatID, userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, message)
}

// SearchMessages Поиск сообщений в чате
// @Summary Поиск сообщений в чате
// @Description Выполняет поиск сообщений по тексту в указанном чате с пагинацией
// @Tags chats
// @Produce json
// @Security BearerAuth
// @Param chat_id path string true "UUID чата"
// @Param query query string true "Поисковый запрос"
// @Param offset query int false "Смещение для пагинации" default(0)
// @Param limit query int false "Количество результатов на странице" default(20) maximum(100)
// @Success 200 {array} map[string]interface{} "Результаты поиска"
// @Failure 400 {object} map[string]interface{} "Некорректный запрос или пустой поисковый запрос"
// @Failure 401 {object} map[string]interface{} "Пользователь не аутентифицирован"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /chats/search/{chat_id} [get]
func (h *ChatHandler) SearchMessages(c *gin.Context) {
	chatID, err := uuid.Parse(c.Param("chat_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid chat ID"})
		return
	}

	userID, err := getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	query := c.Query("query")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query parameter is required"})
		return
	}

	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil || offset < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid offset"})
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if err != nil || limit <= 0 || limit > 100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit"})
		return
	}

	messages, err := h.chatController.SearchMessages(userID, chatID, query, offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, messages)
}

// UpdateChat Обновление чата
// @Summary Обновить чат
// @Description Обновляет параметры чата, включая добавление/удаление участников
// @Tags chats
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param chat_id path string true "UUID чата"
// @Param request body dto.UpdateChatRequestGateway true "Параметры обновления"
// @Success 200 {object} map[string]interface{} "Чат успешно обновлен"
// @Failure 400 {object} map[string]interface{} "Некорректный запрос или неверный UUID"
// @Failure 401 {object} map[string]interface{} "Пользователь не аутентифицирован"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /chats/{chat_id} [put]
func (h *ChatHandler) UpdateChat(c *gin.Context) {
	chatID, err := uuid.Parse(c.Param("chat_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid chat ID"})
		return
	}

	var req dto.UpdateChatRequestGateway
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Парсим UUID для добавления/удаления пользователей
	updateReq, err := req.ToUpdateChatRequest()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid UUID in request: " + err.Error()})
		return
	}

	result, err := h.chatController.UpdateChat(chatID, updateReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// DeleteChat Удаление чата
// @Summary Удалить чат
// @Description Удаляет чат и все связанные данные
// @Tags chats
// @Produce json
// @Security BearerAuth
// @Param chat_id path string true "UUID чата"
// @Success 200 {object} map[string]interface{} "Чат успешно удален"
// @Failure 400 {object} map[string]interface{} "Некорректный UUID чата"
// @Failure 401 {object} map[string]interface{} "Пользователь не аутентифицирован"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /chats/{chat_id} [delete]
func (h *ChatHandler) DeleteChat(c *gin.Context) {
	chatID, err := uuid.Parse(c.Param("chat_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid chat ID"})
		return
	}

	err = h.chatController.DeleteChat(chatID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "chat deleted successfully"})
}

// BanUser Блокировка пользователя в чате
// @Summary Заблокировать пользователя в чате
// @Description Блокирует пользователя, запрещая ему доступ к чату
// @Tags chats
// @Produce json
// @Security BearerAuth
// @Param chat_id path string true "UUID чата"
// @Param user_id path string true "UUID пользователя"
// @Success 200 {object} map[string]interface{} "Пользователь успешно заблокирован"
// @Failure 400 {object} map[string]interface{} "Некорректный UUID чата или пользователя"
// @Failure 401 {object} map[string]interface{} "Пользователь не аутентифицирован"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /chats/{chat_id}/ban/{user_id} [post]
func (h *ChatHandler) BanUser(c *gin.Context) {
	chatID, err := uuid.Parse(c.Param("chat_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid chat ID"})
		return
	}

	userID, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	err = h.chatController.BanUser(chatID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user banned successfully"})
}

// ChangeUserRole Изменение роли пользователя в чате
// @Summary Изменить роль пользователя в чате
// @Description Назначает новую роль пользователю в указанном чате
// @Tags chats
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param chat_id path string true "UUID чата"
// @Param request body dto.ChangeRoleRequestGateway true "Данные для изменения роли"
// @Success 200 {object} map[string]interface{} "Роль пользователя успешно изменена"
// @Failure 400 {object} map[string]interface{} "Некорректный запрос или неверный UUID"
// @Failure 401 {object} map[string]interface{} "Пользователь не аутентифицирован"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /chats/{chat_id}/roles/change [patch]
func (h *ChatHandler) ChangeUserRole(c *gin.Context) {
	chatID, err := uuid.Parse(c.Param("chat_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid chat ID"})
		return
	}

	var req dto.ChangeRoleRequestGateway
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Конвертируем DTO
	changeRoleReq, err := req.ToChangeRoleRequest()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid UUID in request: " + err.Error()})
		return
	}

	err = h.chatController.ChangeUserRole(chatID, changeRoleReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user role changed successfully"})
}
