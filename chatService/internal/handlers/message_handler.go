package handlers

import (
	"chatService/internal/controllers"
	"chatService/internal/custom_errors"
	"chatService/internal/handlers/dto"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type MessageHandler struct {
	MessageController *controllers.MessageController
}

func NewMessageHandler(messageController *controllers.MessageController) *MessageHandler {
	return &MessageHandler{messageController}
}

// SendMessage Отправка сообщения в чат
// @Summary Отправить сообщение в чат
// @Description Отправляет новое сообщение в указанный чат с возможностью прикрепления файлов
// @Tags messages
// @Accept json
// @Produce json
// @Param chat_id path string true "UUID чата"
// @Param X-User-ID header string true "UUID отправителя"
// @Param message body dto.CreateMessageDTO true "Данные сообщения"
// @Success 201 {object} models.Message "Сообщение успешно отправлено"
// @Failure 400 {object} map[string]interface{} "Некорректный запрос или неверный UUID"
// @Failure 401 {object} map[string]interface{} "Неверные учетные данные"
// @Failure 502 {object} map[string]interface{} "Ошибка при обращении к внешнему сервису"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /chats/messages/{chat_id} [post]
func (h *MessageHandler) SendMessage(c *gin.Context) {
	userIDStr := c.GetHeader("X-User-ID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid User ID"})
		return
	}

	chatID, err := uuid.Parse(c.Param("chat_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid chat ID"})
		return
	}

	var messageDTO dto.CreateMessageDTO
	if err := c.ShouldBindJSON(&messageDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	msg, err := h.MessageController.SendMessage(userID, chatID, &messageDTO)
	if err != nil {
		var userErr *custom_errors.UserClientError
		var fileNotFoundErr *custom_errors.FileNotFoundError
		var getFileHTTPError *custom_errors.GetFileHTTPError
		var dbErr *custom_errors.DatabaseError

		switch {
		case errors.Is(err, custom_errors.ErrInvalidCredentials):
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		case errors.As(err, &userErr):
			c.JSON(http.StatusBadRequest, gin.H{"error": userErr.Error()})
		case errors.As(err, &fileNotFoundErr):
			c.JSON(http.StatusBadRequest, gin.H{"error": fileNotFoundErr.Error()})
		case errors.As(err, &getFileHTTPError):
			c.JSON(http.StatusBadGateway, gin.H{"error": getFileHTTPError.Error()})
		case errors.As(err, &dbErr):
			c.JSON(http.StatusInternalServerError, gin.H{"error": dbErr.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": custom_errors.ErrInternalServerError.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, msg)
}

// GetChatMessages Получение сообщений чата
// @Summary Получить сообщения чата
// @Description Возвращает список сообщений из указанного чата с пагинацией
// @Tags messages
// @Produce json
// @Param chat_id path string true "UUID чата"
// @Param offset query int false "Смещение для пагинации" default(0)
// @Param limit query int false "Количество сообщений на странице" default(20)
// @Success 200 {array} dto.GetChatMessage "Список сообщений"
// @Failure 400 {object} map[string]interface{} "Некорректный UUID чата или параметры пагинации"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /chats/messages/{chat_id} [get]
func (h *MessageHandler) GetChatMessages(c *gin.Context) {
	chatIDStr := c.Param("chat_id")
	chatID, err := uuid.Parse(chatIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid chat ID"})
		return
	}

	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil || offset < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid offset"})
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if err != nil || limit <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit"})
		return
	}

	messages, err := h.MessageController.GetChatMessages(chatID, offset, limit)
	if err != nil {
		var dbErr *custom_errors.DatabaseError

		if errors.As(err, &dbErr) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": dbErr.Error()})
			return
		} else if errors.Is(err, custom_errors.ErrInvalidCredentials) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "chat not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": custom_errors.ErrInternalServerError.Error()})
		return
	}

	c.JSON(http.StatusOK, messages)
}

// SearchMessages Поиск сообщений в чате
// @Summary Поиск сообщений в чате
// @Description Выполняет поиск сообщений по тексту в указанном чате с пагинацией
// @Tags messages
// @Produce json
// @Param chat_id path string true "UUID чата"
// @Param X-User-ID header string true "UUID пользователя"
// @Param query query string true "Поисковый запрос"
// @Param limit query int false "Количество результатов на странице" default(20)
// @Param offset query int false "Смещение для пагинации" default(0)
// @Success 200 {object} dto.GetSearchResponse "Результаты поиска"
// @Failure 400 {object} map[string]interface{} "Некорректный запрос или пустой поисковый запрос"
// @Failure 403 {object} map[string]interface{} "Нет доступа к чату"
// @Failure 404 {object} map[string]interface{} "Чат не найден"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /chats/search/{chat_id} [get]
func (h *MessageHandler) SearchMessages(c *gin.Context) {
	chatID, err := uuid.Parse(c.Param("chat_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid chat_id"})
		return
	}

	userIDStr := c.GetHeader("X-User-ID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid User ID"})
		return
	}

	query := c.Query("query")
	limitStr := c.DefaultQuery("limit", "20")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit"})
		return
	}
	if limit > 50 {
		limit = 50
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid offset"})
		return
	}

	messages, err := h.MessageController.SearchMessages(userID, chatID, query, limit, offset)
	if err != nil {
		switch {
		case errors.Is(err, custom_errors.ErrEmptyQuery):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case errors.Is(err, custom_errors.ErrChatNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, custom_errors.ErrUnauthorizedChat):
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		}
		return
	}

	c.JSON(http.StatusOK, messages)
}
