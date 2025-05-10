package handlers

import (
	"chatService/internal/controllers"
	"chatService/internal/custom_errors"
	"chatService/internal/handlers/dto"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type ChatHandler struct {
	ChatController *controllers.ChatController
}

func NewChatHandler(chatController *controllers.ChatController) *ChatHandler {
	return &ChatHandler{chatController}
}

// ChangeUserRole POST /api/v1/chats/:chat_id/roles/change
func (h *ChatHandler) ChangeUserRole(c *gin.Context) {
	var request struct {
		UserID uuid.UUID `json:"user_id" binding:"required"`
		RoleID int       `json:"role_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	chatID, err := uuid.Parse(c.Param("chat_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.ChatController.ChangeUserRole(chatID, request.UserID, request.RoleID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

// GetUserChats GET /api/v1/chats
func (h *ChatHandler) GetUserChats(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	chats, err := h.ChatController.GetUserChats(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, chats)
}

// CreateChat POST /api/v1/chats
func (h *ChatHandler) CreateChat(c *gin.Context) {
	var createChatDTO dto.CreateChatDTO
	if err := c.ShouldBindJSON(&createChatDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input: " + err.Error()})
		return
	}

	chatID, err := h.ChatController.CreateChat(&createChatDTO)
	if err != nil {
		var getFileErr *custom_errors.GetFileHTTPError
		var fileNotFoundErr *custom_errors.FileNotFoundError
		var userClientErr *custom_errors.UserClientError
		var dbErr *custom_errors.DatabaseError

		switch {
		case errors.Is(err, custom_errors.ErrInvalidCredentials):
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid credentials"})
		case errors.As(err, &getFileErr):
			c.JSON(http.StatusBadGateway, gin.H{"error": getFileErr.Error()})
		case errors.As(err, &fileNotFoundErr):
			c.JSON(http.StatusNotFound, gin.H{"error": fileNotFoundErr.Error()})
		case errors.As(err, &userClientErr):
			c.JSON(http.StatusBadGateway, gin.H{"error": userClientErr.Error()})
		case errors.As(err, &dbErr):
			c.JSON(http.StatusInternalServerError, gin.H{"error": dbErr.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "unknown error: " + err.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"chat_id": chatID})
}

// UpdateChat PATCH /api/v1/chats/:chat_id
func (h *ChatHandler) UpdateChat(c *gin.Context) {
	chatIDParam := c.Param("chat_id")
	chatID, err := uuid.Parse(chatIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid chat ID"})
		return
	}

	var updateChatDTO dto.UpdateChatDTO
	if err := c.ShouldBindJSON(&updateChatDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input: " + err.Error()})
		return
	}

	response, err := h.ChatController.UpdateChat(chatID, &updateChatDTO)
	if err != nil {
		var getFileErr *custom_errors.GetFileHTTPError
		var fileNotFoundErr *custom_errors.FileNotFoundError
		var userClientErr *custom_errors.UserClientError
		var dbErr *custom_errors.DatabaseError

		switch {
		case errors.Is(err, custom_errors.ErrInvalidCredentials):
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid credentials"})
		case errors.As(err, &getFileErr):
			c.JSON(http.StatusBadGateway, gin.H{"error": getFileErr.Error()})
		case errors.As(err, &fileNotFoundErr):
			c.JSON(http.StatusNotFound, gin.H{"error": fileNotFoundErr.Error()})
		case errors.As(err, &userClientErr):
			c.JSON(http.StatusBadGateway, gin.H{"error": userClientErr.Error()})
		case errors.As(err, &dbErr):
			c.JSON(http.StatusInternalServerError, gin.H{"error": dbErr.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "unknown error: " + err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, response)
}

// DeleteChat DELETE /api/v1/chats/:chat_id
func (h *ChatHandler) DeleteChat(c *gin.Context) {
	chatID, _ := uuid.Parse(c.Param("chat_id"))
	if err := h.ChatController.DeleteChat(chatID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// BanUser POST /api/v1/chats/:chat_id/ban/:user_id
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

	if err := h.ChatController.BanUser(chatID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}
