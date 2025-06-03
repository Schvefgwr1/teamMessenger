package main

import (
	"common/config"
	"common/db"
	"common/minio"
	_ "fileService/docs"
	"fileService/internal/controllers"
	"fileService/internal/handlers"
	"fileService/internal/repositories"
	"fileService/internal/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"strconv"
)

// @title File Service API
// @version 1.0
// @description API сервиса для работы с файлами
// @host localhost:8080
// @termsOfService  http://swagger.io/terms/
// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io
// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html
// @BasePath /api/v1
// @schemes http
// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/

// @tag.name files
// @tag.description Операции с файлами

// @tag.name file-types
// @tag.description Операции с типами файлов

func main() {
	// Загружаем переменные окружения из .env файла (если существует)
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found or error loading .env file: %v", err)
	}

	//Upload config
	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		log.Fatal("Can't load config " + err.Error())
		return
	}

	// Apply environment variable overrides
	config.ApplyDatabaseEnvOverrides(cfg)
	config.ApplyMinIOEnvOverrides(cfg)
	config.ApplyAppEnvOverrides(cfg)

	//Init DB
	dbClient, err := db.InitDB(cfg)
	if err != nil {
		log.Fatal("Can't init DB: " + err.Error())
		return
	}

	//Init MinIO
	minioClient, err := minio.InitMinIO(&cfg.MinIO)
	if err != nil {
		log.Fatalf(err.Error())
		return
	}

	// Init repositories
	fileRepo := repositories.NewFileRepository(dbClient)
	fileTypeRepo := repositories.NewFileTypeRepository(dbClient)

	// Init controllers
	fileController := controllers.NewFileController(fileRepo, fileTypeRepo, minioClient, &cfg.MinIO)
	fileTypeController := controllers.NewFileTypeController(fileTypeRepo)

	// Init handlers
	fileHandler := handlers.NewFileHandler(fileController)
	fileTypeHandler := handlers.NewFileTypeHandler(fileTypeController)

	r := gin.Default()

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	//Init routes
	routes.SetupFileRoutes(r, fileHandler)
	routes.SetupFileTypeRoutes(r, fileTypeHandler)

	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	_ = r.Run(":" + strconv.Itoa(cfg.App.Port))
}
