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
