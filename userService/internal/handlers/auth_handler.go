package handlers

import (
	au "common/contracts/api-user"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"userService/internal/controllers"
	"userService/internal/custom_errors"
)

type AuthHandler struct {
	authController *controllers.AuthController
}

func NewAuthHandler(authController *controllers.AuthController) *AuthHandler {
	return &AuthHandler{authController: authController}
}

// Register Регистрация пользователя
// @Summary Регистрация нового пользователя
// @Description Регистрирует нового пользователя в системе
// @Tags auth
// @Accept json
// @Produce json
// @Param user body au.RegisterUserRequest true "Данные для регистрации"
// @Success 201 {object} models.User "Пользователь успешно зарегистрирован"
// @Failure 400 {object} map[string]interface{} "Некорректный запрос"
// @Failure 409 {object} map[string]interface{} "Почта или логин уже заняты"
// @Failure 502 {object} map[string]interface{} "Ошибка при получении файла аватара"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req au.RegisterUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}

	user, err := h.authController.Register(&req)
	if err != nil {
		errMsg := err.Error()

		var userEmailConflictError *custom_errors.UserEmailConflictError
		var userUsernameConflictError *custom_errors.UserUsernameConflictError
		var roleNotFoundError *custom_errors.RoleNotFoundError
		var getFileHTTPError *custom_errors.GetFileHTTPError
		var fileNotFoundError *custom_errors.FileNotFoundError
		switch {
		case errors.As(err, &userEmailConflictError), errors.As(err, &userUsernameConflictError):
			c.JSON(http.StatusConflict, gin.H{"error": errMsg}) // 409
		case errors.As(err, &roleNotFoundError):
			c.JSON(http.StatusBadRequest, gin.H{"error": errMsg}) // 400
		case errors.As(err, &getFileHTTPError):
			c.JSON(http.StatusBadGateway, gin.H{"error": errMsg}) // 502
		case errors.As(err, &fileNotFoundError):
			c.JSON(http.StatusBadRequest, gin.H{"error": errMsg}) // 400
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error: " + errMsg}) // 500
		}
		return
	}

	c.JSON(http.StatusCreated, user)
}

// Login Логин пользователя
// @Summary Авторизация пользователя
// @Description Выполняет вход пользователя по логину и паролю
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body dto.Login true "Данные для входа"
// @Success 200 {object} map[string]string "Токен доступа"
// @Failure 400 {object} map[string]interface{} "Некорректный запрос"
// @Failure 401 {object} map[string]interface{} "Неверные учетные данные"
// @Failure 500 {object} map[string]interface{} "Ошибка генерации токена или сервера"
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req au.Login

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}

	token, userID, err := h.authController.Login(&req)
	if err != nil {
		switch {
		case errors.Is(err, custom_errors.ErrInvalidCredentials):
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		case errors.Is(err, custom_errors.ErrTokenGeneration):
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error: " + err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token, "userID": userID})
}
