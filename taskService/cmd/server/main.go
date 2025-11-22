package main

import (
	"common/config"
	"common/db"
	"common/kafka"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"strconv"
	"taskService/internal/controllers"
	"taskService/internal/handlers"
	"taskService/internal/repositories"
	"taskService/internal/routes"
	"taskService/internal/services"
)

// @title Task Service API
// @version 1.0
// @description API сервиса для работы с задачами
// @host localhost:8082
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

// @tag.name tasks
// @tag.description Операции с задачами

// @tag.name task-statuses
// @tag.description Операции со статусами задач

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
	config.ApplyAppEnvOverrides(cfg)

	//Init DB
	initDB, err := db.InitDB(cfg)
	if err != nil {
		log.Fatal("Can't init DB: " + err.Error())
		return
	}

	// Init Kafka notification service
	kafkaConfig := &kafka.ProducerConfig{
		Brokers: kafka.GetKafkaBrokers(),
		Topic:   kafka.GetNotificationsTopic(),
	}

	notificationService, err := services.NewNotificationService(kafkaConfig)
	if err != nil {
		log.Printf("Warning: Failed to initialize notification service: %v", err)
		notificationService = nil
	}

	//// Init repositories
	taskRepo := repositories.NewTaskRepository(initDB)
	taskFileRepo := repositories.NewTaskFileRepository(initDB)
	taskStatusRepo := repositories.NewTaskStatusRepository(initDB)

	//// Init controllers
	taskController := controllers.NewTaskController(taskRepo, taskStatusRepo, taskFileRepo, notificationService)
	taskStatusController := controllers.NewTaskStatusController(taskStatusRepo)

	//// Init handlers
	taskHandler := handlers.NewTaskHandler(taskController)
	taskStatusHandler := handlers.NewTaskStatusHandler(taskStatusController)

	r := gin.Default()

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	routes.RegisterTaskStatusRoutes(r, taskStatusHandler)
	routes.RegisterTaskRoutes(r, taskHandler)

	// Graceful shutdown для Kafka producer
	defer func() {
		if notificationService != nil {
			if err := notificationService.Close(); err != nil {
				log.Printf("Error closing notification service: %v", err)
			}
		}
	}()

	_ = r.Run(":" + strconv.Itoa(cfg.App.Port))
}
