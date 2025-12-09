package handlers

import (
	"fileService/internal/controllers"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type FileTypeHandler struct {
	controller controllers.FileTypeControllerInterface
}

func NewFileTypeHandler(controller controllers.FileTypeControllerInterface) *FileTypeHandler {
	return &FileTypeHandler{controller: controller}
}

// CreateFileTypeHandler обрабатывает создание нового типа файла
// @Summary Создание типа файла
// @Description Создает новый тип файла в системе
// @Tags file-types
// @Produce json
// @Param name query string true "Название типа файла"
// @Success 200 {object} models.FileType "Тип файла успешно создан"
// @Failure 400 {object} map[string]interface{} "Отсутствует название типа файла"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /file-types [post]
func (h *FileTypeHandler) CreateFileTypeHandler(c *gin.Context) {
	name := c.Query("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name must be exist"})
		return
	}

	fileType, err := h.controller.CreateFileType(name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, fileType)
}

// GetFileTypeByIDHandler обрабатывает получение типа файла по ID
// @Summary Получение типа файла по ID
// @Description Получает информацию о типе файла по его ID
// @Tags file-types
// @Produce json
// @Param id path int true "ID типа файла"
// @Success 200 {object} models.FileType "Информация о типе файла"
// @Failure 400 {object} map[string]interface{} "Некорректный ID"
// @Failure 404 {object} map[string]interface{} "Тип файла не найден"
// @Router /file-types/{id} [get]
func (h *FileTypeHandler) GetFileTypeByIDHandler(c *gin.Context) {
	idParam := c.Param("id")
	fileID, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect id"})
		return
	}

	fileType, err := h.controller.GetFileTypeByID(fileID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File type does not exist"})
		return
	}

	c.JSON(http.StatusOK, fileType)
}

// GetFileTypeByNameHandler обрабатывает получение типа файла по названию
// @Summary Получение типа файла по названию
// @Description Получает информацию о типе файла по его названию
// @Tags file-types
// @Produce json
// @Param name path string true "Название типа файла"
// @Success 200 {object} models.FileType "Информация о типе файла"
// @Failure 400 {object} map[string]interface{} "Отсутствует название типа файла"
// @Failure 404 {object} map[string]interface{} "Тип файла не найден"
// @Router /file-types/name/{name} [get]
func (h *FileTypeHandler) GetFileTypeByNameHandler(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name must be exist"})
		return
	}

	fileType, err := h.controller.GetFileTypeByName(name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File type does not exist"})
		return
	}

	c.JSON(http.StatusOK, fileType)
}

// DeleteFileTypeByIDHandler обрабатывает удаление типа файла по ID
// @Summary Удаление типа файла
// @Description Удаляет тип файла по его ID
// @Tags file-types
// @Produce json
// @Param id path int true "ID типа файла"
// @Success 200 {object} map[string]interface{} "Тип файла успешно удален"
// @Failure 400 {object} map[string]interface{} "Некорректный ID"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /file-types/{id} [delete]
func (h *FileTypeHandler) DeleteFileTypeByIDHandler(c *gin.Context) {
	idParam := c.Param("id")
	fileID, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect id"})
		return
	}

	err = h.controller.DeleteFileTypeByID(fileID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "FileType has successfully deleted"})
}
