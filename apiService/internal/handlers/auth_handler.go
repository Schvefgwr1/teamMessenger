package handlers

import (
	"apiService/internal/controllers"
	"apiService/internal/dto"
	"apiService/internal/services"
	au "common/contracts/api-user"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthHandler struct {
	authController *controllers.AuthController
	sessionService *services.SessionService
}

func NewAuthHandler(authController *controllers.AuthController, sessionService *services.SessionService) *AuthHandler {
	return &AuthHandler{
		authController: authController,
		sessionService: sessionService,
	}
}

// Register Регистрация нового пользователя
// @Summary Зарегистрировать нового пользователя
// @Description Регистрирует нового пользователя с возможностью загрузки аватара
// @Tags auth
// @Accept multipart/form-data
// @Produce json
// @Param data formData string true "JSON данные пользователя" example({"username":"user","email":"user@example.com","password":"password123","age":25,"roleID":1})
// @Param file formData file false "Аватар пользователя"
// @Success 201 {object} map[string]interface{} "Пользователь успешно зарегистрирован"
// @Failure 400 {object} map[string]interface{} "Некорректный запрос"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var registerData dto.RegisterUserRequestGateway

	// Parse JSON part from the multipart form
	jsonData := c.PostForm("data")
	if err := json.Unmarshal([]byte(jsonData), &registerData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Handle optional file upload
	file, err := c.FormFile("file")
	if err != nil && !errors.Is(err, http.ErrMissingFile) {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userResponse := h.authController.Register(&registerData, file)
	if userResponse.Error != nil {
		c.JSON(http.StatusInternalServerError, userResponse)
		return
	}

	c.JSON(http.StatusCreated, userResponse)
}

// Login Аутентификация пользователя
// @Summary Войти в систему
// @Description Выполняет аутентификацию пользователя и возвращает JWT токен
// @Tags auth
// @Accept json
// @Produce json
// @Param login body dto.LoginSwag true "Данные для входа"
// @Success 200 {object} map[string]interface{} "Успешная аутентификация" example({"token":"eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9...","userID":"00000000-0000-0000-0000-000000000000"})
// @Failure 400 {object} map[string]interface{} "Некорректный запрос"
// @Failure 401 {object} map[string]interface{} "Неверные учетные данные"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var loginData au.Login

	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}

	token, userID, errReq := h.authController.Login(&loginData)
	if errReq != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": errReq})
		return
	}

	// Создаем сессию в Redis если есть sessionService
	if h.sessionService != nil && token != "" {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Извлекаем userID из ответа (предполагаем что он есть в response)
		if userID != uuid.Nil {
			expiresAt := time.Now().Add(24 * time.Hour)
			if err := h.sessionService.CreateSession(ctx, userID, token, expiresAt); err != nil {
				fmt.Printf("Failed to create session in Redis: %v\n", err)
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{"token": token, "userID": userID})
}

// Logout Выход из системы
// @Summary Выйти из системы
// @Description Отзывает текущую сессию пользователя
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "Успешный выход" example({"message":"Logged out successfully"})
// @Failure 401 {object} map[string]interface{} "Пользователь не аутентифицирован"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	if h.sessionService == nil {
		c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	token, exists := c.Get("token")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token not found"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := h.sessionService.RevokeSession(ctx, userID.(uuid.UUID), token.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to revoke session"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}
