package routes

import (
	"fileService/internal/handlers"
	"github.com/gin-gonic/gin"
)

func SetupFileRoutes(router *gin.Engine, fileHandler *handlers.FileHandler) {
	fileGroup := router.Group("/api/v1/files")
	{
		fileGroup.POST("/upload", fileHandler.UploadFileHandler)  // Загрузка файла
		fileGroup.GET("/:file_id", fileHandler.GetFileHandler)    // Получение информации о файле
		fileGroup.PUT("/:file_id", fileHandler.RenameFileHandler) // Переименование файла
		fileGroup.GET("/names", fileHandler.GetFileNamesHandler)  // Получение списка ID + Name
	}
}
