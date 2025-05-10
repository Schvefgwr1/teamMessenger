package handlers

import (
	"chatService/internal/controllers"
	"chatService/internal/custom_errors"
	"chatService/internal/handlers/dto"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"strconv"
)

type MessageHandler struct {
	MessageController *controllers.MessageController
}

func NewMessageHandler(messageController *controllers.MessageController) *MessageHandler {
	return &MessageHandler{messageController}
}

// SendMessage POST /api/v1/chats/:chat_id/messages
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

// GetChatMessages GET /api/v1/chats/:chat_id/messages
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

// SearchMessages GET /api/v1/chats/:chat_id/messages/search
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

	messages, total, err := h.MessageController.SearchMessages(userID, chatID, query, limit, offset)
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

	c.JSON(http.StatusOK, gin.H{
		"messages": messages,
		"total":    total,
	})
}
