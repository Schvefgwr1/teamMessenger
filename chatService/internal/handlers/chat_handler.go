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
	ChatController controllers.ChatControllerInterface
}

func NewChatHandler(chatController controllers.ChatControllerInterface) *ChatHandler {
	return &ChatHandler{chatController}
}

// ChangeUserRole Изменение роли пользователя в чате
// @Summary Изменение роли пользователя в чате
// @Description Изменяет роль пользователя в указанном чате
// @Tags chats
// @Accept json
// @Produce json
// @Param chat_id path string true "UUID чата"
// @Param request body object true "Данные для изменения роли" example({"user_id":"00000000-0000-0000-0000-000000000000","role_id":1})
// @Success 200 "Роль успешно изменена"
// @Failure 400 {object} map[string]interface{} "Некорректный запрос"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /chats/{chat_id}/roles/change [patch]
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

// GetUserChats Получение списка чатов пользователя
// @Summary Получить список чатов пользователя
// @Description Возвращает список всех чатов, в которых участвует указанный пользователь
// @Tags chats
// @Produce json
// @Param user_id path string true "UUID пользователя"
// @Success 200 {array} dto.ChatResponse "Список чатов пользователя"
// @Failure 400 {object} map[string]interface{} "Некорректный UUID пользователя"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /chats/user/{user_id} [get]
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

// CreateChat Создание нового чата
// @Summary Создать новый чат
// @Description Создает новый чат с указанными параметрами и участниками
// @Tags chats
// @Accept json
// @Produce json
// @Param chat body dto.CreateChatDTO true "Данные для создания чата"
// @Success 201 {object} map[string]interface{} "Чат успешно создан"
// @Failure 400 {object} map[string]interface{} "Некорректный запрос или неверные учетные данные"
// @Failure 404 {object} map[string]interface{} "Файл не найден"
// @Failure 502 {object} map[string]interface{} "Ошибка при обращении к внешнему сервису"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /chats [post]
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

// UpdateChat Обновление информации о чате
// @Summary Обновить информацию о чате
// @Description Обновляет данные чата: название, описание, аватар, список участников
// @Tags chats
// @Accept json
// @Produce json
// @Param chat_id path string true "UUID чата"
// @Param chat body dto.UpdateChatDTO true "Данные для обновления чата"
// @Success 200 {object} dto.UpdateChatResponse "Чат успешно обновлен"
// @Failure 400 {object} map[string]interface{} "Некорректный запрос или неверные учетные данные"
// @Failure 404 {object} map[string]interface{} "Файл не найден"
// @Failure 502 {object} map[string]interface{} "Ошибка при обращении к внешнему сервису"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /chats/{chat_id} [put]
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

// DeleteChat Удаление чата
// @Summary Удалить чат
// @Description Удаляет чат по указанному ID
// @Tags chats
// @Produce json
// @Param chat_id path string true "UUID чата"
// @Success 204 "Чат успешно удален"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /chats/{chat_id} [delete]
func (h *ChatHandler) DeleteChat(c *gin.Context) {
	chatID, _ := uuid.Parse(c.Param("chat_id"))
	if err := h.ChatController.DeleteChat(chatID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// BanUser Блокировка пользователя в чате
// @Summary Заблокировать пользователя в чате
// @Description Блокирует указанного пользователя в чате
// @Tags chats
// @Produce json
// @Param chat_id path string true "UUID чата"
// @Param user_id path string true "UUID пользователя"
// @Success 200 "Пользователь успешно заблокирован"
// @Failure 400 {object} map[string]interface{} "Некорректный UUID чата или пользователя"
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

	if err := h.ChatController.BanUser(chatID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

// GetUserRoleInChat Получение роли пользователя в чате
// @Summary Получить роль пользователя в чате
// @Description Возвращает название роли указанного пользователя в чате
// @Tags chats
// @Produce json
// @Param chat_id path string true "UUID чата"
// @Param user_id path string true "UUID пользователя"
// @Param X-User-ID header string true "UUID запрашивающего пользователя"
// @Success 200 {object} dto.UserRoleResponse "Роль пользователя в чате"
// @Failure 400 {object} map[string]interface{} "Некорректный UUID"
// @Failure 403 {object} map[string]interface{} "Нет доступа к чату"
// @Failure 404 {object} map[string]interface{} "Пользователь не найден в чате"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /chats/{chat_id}/user-roles/{user_id} [get]
func (h *ChatHandler) GetUserRoleInChat(c *gin.Context) {
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

	requesterIDStr := c.GetHeader("X-User-ID")
	requesterID, err := uuid.Parse(requesterIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid requester ID"})
		return
	}

	roleName, err := h.ChatController.GetUserRoleInChat(chatID, userID, requesterID)
	if err != nil {
		switch {
		case errors.Is(err, custom_errors.ErrUnauthorizedChat):
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		case errors.Is(err, custom_errors.ErrUserNotInChat):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, dto.UserRoleResponse{RoleName: roleName})
}

// GetMyRoleInChat Получение своей роли в чате с permissions
// @Summary Получить свою роль в чате с permissions
// @Description Возвращает роль текущего пользователя в чате с полным списком permissions
// @Tags chats
// @Produce json
// @Param chat_id path string true "UUID чата"
// @Param X-User-ID header string true "UUID пользователя"
// @Success 200 {object} dto.UserRoleWithPermissionsResponse "Роль с permissions"
// @Failure 400 {object} map[string]interface{} "Некорректный UUID"
// @Failure 404 {object} map[string]interface{} "Пользователь не найден в чате"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /chats/{chat_id}/me/role [get]
func (h *ChatHandler) GetMyRoleInChat(c *gin.Context) {
	chatID, err := uuid.Parse(c.Param("chat_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid chat ID"})
		return
	}

	userIDStr := c.GetHeader("X-User-ID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	role, err := h.ChatController.GetMyRoleWithPermissions(chatID, userID)
	if err != nil {
		switch {
		case errors.Is(err, custom_errors.ErrUserNotInChat):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	// Конвертируем permissions в DTO
	permissions := make([]dto.ChatPermissionResponse, len(role.Permissions))
	for i, p := range role.Permissions {
		permissions[i] = dto.ChatPermissionResponse{
			ID:   p.ID,
			Name: p.Name,
		}
	}

	c.JSON(http.StatusOK, dto.UserRoleWithPermissionsResponse{
		RoleID:      role.ID,
		RoleName:    role.Name,
		Permissions: permissions,
	})
}

// GetChatByID Получение чата по ID
// @Summary Получить чат по ID
// @Description Возвращает информацию о чате по указанному ID (без участников и сообщений)
// @Tags chats
// @Produce json
// @Param chat_id path string true "UUID чата"
// @Success 200 {object} dto.ChatResponse "Информация о чате"
// @Failure 400 {object} map[string]interface{} "Некорректный UUID"
// @Failure 404 {object} map[string]interface{} "Чат не найден"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /chats/{chat_id} [get]
func (h *ChatHandler) GetChatByID(c *gin.Context) {
	chatID, err := uuid.Parse(c.Param("chat_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid chat ID"})
		return
	}

	chat, err := h.ChatController.GetChatByID(chatID)
	if err != nil {
		if errors.Is(err, custom_errors.ErrChatNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, chat)
}

// GetChatMembers Получение списка участников чата
// @Summary Получить список участников чата
// @Description Возвращает список всех участников чата с их ролями
// @Tags chats
// @Produce json
// @Param chat_id path string true "UUID чата"
// @Success 200 {array} dto.ChatMemberResponse "Список участников"
// @Failure 400 {object} map[string]interface{} "Некорректный UUID"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /chats/{chat_id}/members [get]
func (h *ChatHandler) GetChatMembers(c *gin.Context) {
	chatID, err := uuid.Parse(c.Param("chat_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid chat ID"})
		return
	}

	chatUsers, err := h.ChatController.GetChatMembers(chatID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Конвертируем в DTO
	members := make([]dto.ChatMemberResponse, len(chatUsers))
	for i, cu := range chatUsers {
		members[i] = dto.ChatMemberResponse{
			UserID:   cu.UserID.String(),
			RoleID:   cu.RoleID,
			RoleName: cu.Role.Name,
		}
	}

	c.JSON(http.StatusOK, members)
}
