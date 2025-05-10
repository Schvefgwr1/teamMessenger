package middlewares

import (
	"chatService/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type ChatPermissionMiddleware struct {
	permissionService *services.ChatPermissionService
}

func NewChatPermissionMiddleware(permissionService *services.ChatPermissionService) *ChatPermissionMiddleware {
	return &ChatPermissionMiddleware{permissionService: permissionService}
}

func (m *ChatPermissionMiddleware) RequireChatPermission(permissionName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Извлекаем userID и chatID
		userIDStr := c.GetHeader("X-User-ID") // или откуда у тебя userID (например из middleware авторизации)
		chatIDStr := c.Param("chat_id")

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid User ID"})
			return
		}

		chatID, err := uuid.Parse(chatIDStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid Chat ID"})
			return
		}

		hasPermission, err := m.permissionService.HasPermission(userID, chatID, permissionName)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to check permissions"})
			return
		}

		if !hasPermission {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden: insufficient permissions"})
			return
		}

		c.Next()
	}
}
