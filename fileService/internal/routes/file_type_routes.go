package routes

import (
	"fileService/internal/handlers"
	"github.com/gin-gonic/gin"
)

func SetupFileTypeRoutes(router *gin.Engine, fileTypeHandler *handlers.FileTypeHandler) {
	fileTypeGroup := router.Group("/api/v1/file-types")
	{
		fileTypeGroup.POST("/", fileTypeHandler.CreateFileTypeHandler)             // Создание типа файла
		fileTypeGroup.GET("/:id", fileTypeHandler.GetFileTypeByIDHandler)          // Получение типа файла по ID
		fileTypeGroup.GET("/name/:name", fileTypeHandler.GetFileTypeByNameHandler) // Получение типа файла по названию
		fileTypeGroup.DELETE("/:id", fileTypeHandler.DeleteFileTypeByIDHandler)    // Удаление типа файла по ID
	}
}
