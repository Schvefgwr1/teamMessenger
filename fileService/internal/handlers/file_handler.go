package handlers

import (
	"fileService/internal/controllers"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

type FileHandler struct {
	controller *controllers.FileController
}

func NewFileHandler(controller *controllers.FileController) *FileHandler {
	return &FileHandler{controller: controller}
}

// UploadFileHandler обрабатывает загрузку файла
// @Summary Загрузка файла
// @Description Загружает файл в MinIO и сохраняет метаданные в БД
// @Tags files
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Файл для загрузки"
// @Success 201 {object} models.File "Файл успешно загружен"
// @Failure 400 {object} map[string]interface{} "Файл отсутствует в форме"
// @Failure 415 {object} map[string]interface{} "Неподдерживаемый тип файла"
// @Failure 409 {object} map[string]interface{} "Файл уже существует"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /files [post]
func (h *FileHandler) UploadFileHandler(c *gin.Context) {
	// Получаем файл из запроса
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File in form does not exist"})
		return
	}

	// Загружаем файл через контроллер
	uploadedFile, err := h.controller.UploadFile(file)
	if err != nil {
		if strings.Contains(err.Error(), "Unsupported file type") {
			c.JSON(http.StatusUnsupportedMediaType, gin.H{"error": err.Error()})
			return
		}
		if strings.Contains(err.Error(), "already exists in database") {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		if strings.Contains(err.Error(), "already exists in MinIO") {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, uploadedFile)
}

// GetFileHandler обрабатывает получение информации о файле
// @Summary Получение информации о файле
// @Description Получает информацию о файле по его ID
// @Tags files
// @Produce json
// @Param file_id path int true "ID файла"
// @Success 200 {object} models.File "Информация о файле"
// @Failure 400 {object} map[string]interface{} "Некорректный ID файла"
// @Failure 404 {object} map[string]interface{} "Файл не найден"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /files/{file_id} [get]
func (h *FileHandler) GetFileHandler(c *gin.Context) {
	// Получаем file_id из параметра URL
	idParam := c.Param("file_id")
	fileID, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect file_id"})
		return
	}

	// Получаем информацию о файле через контроллер
	file, err := h.controller.GetFile(fileID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, file)
}

// RenameFileHandler обрабатывает переименование файла
// @Summary Переименование файла
// @Description Переименовывает файл в MinIO и обновляет информацию в БД
// @Tags files
// @Produce json
// @Param file_id path int true "ID файла"
// @Param name query string true "Новое имя файла"
// @Success 200 {object} models.File "Файл успешно переименован"
// @Failure 400 {object} map[string]interface{} "Некорректные параметры запроса"
// @Failure 404 {object} map[string]interface{} "Файл не найден"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /files/{file_id}/rename [put]
func (h *FileHandler) RenameFileHandler(c *gin.Context) {
	// Получаем file_id из параметра URL
	idParam := c.Param("file_id")
	fileID, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect file_id"})
		return
	}

	newName := c.Query("name")
	if newName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "New name cannot be empty"})
		return
	}

	// Вызываем метод контроллера для переименования
	updatedFile, err := h.controller.RenameFile(fileID, newName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedFile)
}

// GetFileNamesHandler обрабатывает запрос списка ID и Name файлов с пагинацией
// @Summary Получение списка файлов
// @Description Получает список файлов с пагинацией
// @Tags files
// @Produce json
// @Param limit query int false "Количество файлов на странице" default(10)
// @Param offset query int false "Смещение для пагинации" default(0)
// @Success 200 {array} dto.FileInformation "Список файлов"
// @Failure 400 {object} map[string]interface{} "Некорректные параметры пагинации"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /files/names [get]
func (h *FileHandler) GetFileNamesHandler(c *gin.Context) {
	// Читаем limit и offset из запроса (по умолчанию limit=10, offset=0)
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Limit must be a natural number"})
		return
	}

	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil || offset < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Offset must be a natural number"})
		return
	}

	// Получаем данные через контроллер
	files, err := h.controller.GetFileNamesWithPagination(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, files)
}
