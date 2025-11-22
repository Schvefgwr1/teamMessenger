package handlers

import (
	"apiService/internal/controllers"
	"apiService/internal/dto"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type UserHandler struct {
	userController *controllers.UserController
}

func NewUserHandler(userController *controllers.UserController) *UserHandler {
	return &UserHandler{userController: userController}
}

// GetUser Получение информации о текущем пользователе
// @Summary Получить информацию о текущем пользователе
// @Description Возвращает информацию о текущем аутентифицированном пользователе
// @Tags users
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "Информация о пользователе"
// @Failure 400 {object} map[string]interface{} "Некорректный запрос"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /users/me [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	userIDValue, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID not found"})
		return
	}

	userID, ok := userIDValue.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	user, err := h.userController.GetUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateUser Обновление информации о текущем пользователе
// @Summary Обновить информацию о текущем пользователе
// @Description Обновляет данные текущего пользователя с возможностью загрузки нового аватара
// @Tags users
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param data formData string true "JSON данные для обновления" example({"username":"newusername","age":26})
// @Param file formData file false "Новый аватар пользователя"
// @Success 201 {object} map[string]interface{} "Пользователь успешно обновлен"
// @Failure 400 {object} map[string]interface{} "Некорректный запрос"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /users/me [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	var updateData dto.UpdateUserRequestGateway

	jsonData := c.PostForm("data")
	if err := json.Unmarshal([]byte(jsonData), &updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	file, err := c.FormFile("file")
	if err != nil && !errors.Is(err, http.ErrMissingFile) {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIDValue, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID not found"})
		return
	}

	userID, ok := userIDValue.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	userResponse, err := h.userController.UpdateUser(userID, &updateData, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if userResponse.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": userResponse.Error})
		return
	}

	c.JSON(http.StatusCreated, userResponse)
}

// GetAllPermissions Получение всех разрешений
// @Summary Получить все разрешения
// @Description Возвращает список всех доступных разрешений в системе
// @Tags users
// @Produce json
// @Security BearerAuth
// @Success 200 {array} map[string]interface{} "Список разрешений"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /permissions [get]
func (h *UserHandler) GetAllPermissions(c *gin.Context) {
	permissions, err := h.userController.GetAllPermissions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, permissions)
}

// GetAllRoles Получение всех ролей
// @Summary Получить все роли
// @Description Возвращает список всех доступных ролей в системе с их разрешениями
// @Tags users
// @Produce json
// @Security BearerAuth
// @Success 200 {array} map[string]interface{} "Список ролей"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /roles [get]
func (h *UserHandler) GetAllRoles(c *gin.Context) {
	roles, err := h.userController.GetAllRoles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, roles)
}

// CreateRole Создание новой роли
// @Summary Создать новую роль
// @Description Создает новую роль с указанными разрешениями
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateRoleRequestGateway true "Данные для создания роли"
// @Success 201 {object} map[string]interface{} "Роль успешно создана"
// @Failure 400 {object} map[string]interface{} "Некорректный запрос"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /roles [post]
func (h *UserHandler) CreateRole(c *gin.Context) {
	var req dto.CreateRoleRequestGateway
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Конвертируем Gateway DTO в контракт
	roleReq := req.ToCreateRoleRequest()

	role, err := h.userController.CreateRole(roleReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, role)
}

// GetUserProfileByID Получение профиля пользователя по ID
// @Summary Получить профиль пользователя по ID
// @Description Возвращает информацию о пользователе по указанному ID
// @Tags users
// @Produce json
// @Security BearerAuth
// @Param user_id path string true "UUID пользователя"
// @Success 200 {object} map[string]interface{} "Информация о пользователе"
// @Failure 400 {object} map[string]interface{} "Некорректный UUID пользователя"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /users/{user_id} [get]
func (h *UserHandler) GetUserProfileByID(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	userProfile, err := h.userController.GetUserProfileByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, userProfile)
}
