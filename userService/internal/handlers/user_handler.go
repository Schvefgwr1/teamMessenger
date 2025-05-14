package handlers

import (
	au "common/contracts/api-user"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"net/http"
	"userService/internal/controllers"
	"userService/internal/custom_errors"
)

type UserHandler struct {
	userController *controllers.UserController
}

func NewUserHandler(userController *controllers.UserController) *UserHandler {
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
// @Param profile body au.UpdateUserRequest true "Новые данные профиля"
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
