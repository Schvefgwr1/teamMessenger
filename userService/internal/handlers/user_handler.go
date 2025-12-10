package handlers

import (
	au "common/contracts/api-user"
	"errors"
	"fmt"
	"net/http"
	"userService/internal/controllers"
	"userService/internal/custom_errors"
	dto "userService/internal/handlers/dto" // для Swagger документации

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserHandler struct {
	userController controllers.UserControllerInterface
}

func NewUserHandler(userController controllers.UserControllerInterface) *UserHandler {
	return &UserHandler{userController: userController}
}

// GetProfile Получение информации о пользователе
// @Summary Получить профиль пользователя
// @Description Возвращает информацию о пользователе и его аватаре по ID
// @Tags users
// @Produce json
// @Param user_id path string true "UUID пользователя"
// @Success 200 {object} map[string]interface{} "Информация о пользователе и аватаре"
// @Failure 400 {object} map[string]interface{} "Неверный UUID"
// @Failure 404 {object} map[string]interface{} "Пользователь не найден"
// @Router /users/{user_id} [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	idParam := c.Param("user_id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, file, err := h.userController.GetUserProfile(id)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"user":  user,
			"file":  nil,
			"error": "User retrieved, but failed to load avatar",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
		"file": file,
	})
}

// UpdateProfile Обновление профиля
// @Summary Обновление информации профиля пользователя
// @Description Обновляет данные пользователя по ID
// @Tags users
// @Accept json
// @Produce json
// @Param user_id path string true "UUID пользователя"
// @Param profile body dto.UpdateUserRequestSwagger true "Новые данные профиля"
// @Success 200 {object} map[string]interface{} "Профиль успешно обновлен"
// @Failure 400 {object} map[string]interface{} "Некорректный запрос или файл"
// @Failure 401 {object} map[string]interface{} "Неверные учетные данные"
// @Failure 409 {object} map[string]interface{} "Логин уже используется"
// @Failure 502 {object} map[string]interface{} "Ошибка при обращении к файлу"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /users/{user_id} [put]
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	idParam := c.Param("user_id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req au.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}

	err = h.userController.UpdateUserProfile(&req, &id)
	if err != nil {
		errMsg := err.Error()

		var (
			usernameConflictErr *custom_errors.UserUsernameConflictError
			roleNotFoundErr     *custom_errors.RoleNotFoundError
			getFileHTTPErr      *custom_errors.GetFileHTTPError
			fileNotFoundErr     *custom_errors.FileNotFoundError
		)

		switch {
		case errors.Is(err, custom_errors.ErrInvalidCredentials):
			c.JSON(http.StatusUnauthorized, gin.H{"error": errMsg}) // 401
		case errors.As(err, &usernameConflictErr):
			c.JSON(http.StatusConflict, gin.H{"error": errMsg}) // 409
		case errors.As(err, &roleNotFoundErr):
			c.JSON(http.StatusBadRequest, gin.H{"error": errMsg}) // 400
		case errors.As(err, &getFileHTTPErr):
			c.JSON(http.StatusBadGateway, gin.H{"error": errMsg}) // 502
		case errors.As(err, &fileNotFoundErr):
			c.JSON(http.StatusBadRequest, gin.H{"error": errMsg}) // 400
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error: " + errMsg}) // 500
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile updated"})
}

// GetUserBrief Получение краткой информации о пользователе
// @Summary Получить краткую информацию о пользователе
// @Description Возвращает краткую информацию о пользователе: ник, почту, возраст, описание, аватар и роль в чате
// @Tags users
// @Produce json
// @Param user_id path string true "UUID пользователя"
// @Param chatId query string true "UUID чата для получения роли пользователя"
// @Param X-User-ID header string true "UUID запрашивающего пользователя"
// @Success 200 {object} dto.UserBriefResponse "Краткая информация о пользователе"
// @Failure 400 {object} map[string]interface{} "Неверный UUID"
// @Failure 401 {object} map[string]interface{} "Неверные учетные данные"
// @Failure 404 {object} map[string]interface{} "Пользователь не найден"
// @Router /users/{user_id}/brief [get]
func (h *UserHandler) GetUserBrief(c *gin.Context) {
	idParam := c.Param("user_id")
	userID, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	chatID := c.Query("chatId")
	if chatID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "chatId query parameter is required"})
		return
	}

	requesterID := c.GetHeader("X-User-ID")
	if requesterID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "X-User-ID header is required"})
		return
	}

	userBrief, err := h.userController.GetUserBrief(userID, chatID, requesterID)
	if err != nil {
		if errors.Is(err, custom_errors.ErrInvalidCredentials) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, userBrief)
}

// SearchUsers Поиск пользователей
// @Summary Поиск пользователей по имени или email
// @Description Ищет пользователей по частичному совпадению имени или email
// @Tags users
// @Produce json
// @Param q query string true "Поисковый запрос (имя или email)"
// @Param limit query int false "Максимальное количество результатов (по умолчанию 10, максимум 20)"
// @Success 200 {object} dto.UserSearchResponse "Список найденных пользователей"
// @Failure 400 {object} map[string]interface{} "Некорректный запрос"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /users/search [get]
func (h *UserHandler) SearchUsers(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Query parameter 'q' is required"})
		return
	}

	// Минимальная длина запроса
	if len(query) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Query must be at least 2 characters"})
		return
	}

	limit := 10
	if limitParam := c.Query("limit"); limitParam != "" {
		if _, err := fmt.Sscanf(limitParam, "%d", &limit); err != nil {
			limit = 10
		}
	}

	result, err := h.userController.SearchUsers(query, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// UpdateUserRole Изменение роли пользователя
// @Summary Изменить роль пользователя
// @Description Изменяет роль указанного пользователя
// @Tags users
// @Accept json
// @Produce json
// @Param user_id path string true "UUID пользователя"
// @Param request body dto.UpdateUserRoleRequest true "ID новой роли"
// @Success 200 {object} map[string]interface{} "Роль успешно изменена"
// @Failure 400 {object} map[string]interface{} "Некорректный запрос"
// @Failure 401 {object} map[string]interface{} "Пользователь не найден"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /users/{user_id}/role [patch]
func (h *UserHandler) UpdateUserRole(c *gin.Context) {
	userIDParam := c.Param("user_id")
	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req dto.UpdateUserRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}

	err = h.userController.UpdateUserRole(userID, req.RoleID)
	if err != nil {
		errMsg := err.Error()

		var roleNotFoundErr *custom_errors.RoleNotFoundError

		switch {
		case errors.Is(err, custom_errors.ErrInvalidCredentials):
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"}) // 404
		case errors.As(err, &roleNotFoundErr):
			c.JSON(http.StatusBadRequest, gin.H{"error": errMsg}) // 400
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error: " + errMsg}) // 500
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User role updated successfully"})
}
