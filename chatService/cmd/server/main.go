package main

import (
	"chatService/internal/controllers"
	"chatService/internal/handlers"
	"chatService/internal/middlewares"
	"chatService/internal/repositories"
	"chatService/internal/routes"
	"chatService/internal/services"
	"common/config"
	"common/db"
	"common/kafka"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"strconv"
)

// @title Chat Service API
// @version 1.0
// @description API сервиса для работы с чатами
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

// @tag.name auth
// @tag.description Регистрация и аутентификация

// @tag.name users
// @tag.description Операции с пользователем

// @tag.name permissions
// @tag.description Операции с правами доступа

// @tag.name roles
// @tag.description Операции с ролями пользователей

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
	db, err := db.InitDB(cfg)
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

	// Init repositories
	messageRepository := repositories.NewMessageRepository(db)
	chatRoleRepository := repositories.NewChatRoleRepository(db)
	chatUserRepository := repositories.NewChatUserRepository(db)
	chatRepository := repositories.NewChatRepository(db)

	// Init controllers
	messageController := controllers.NewMessageController(messageRepository, chatRepository, chatUserRepository)
	chatController := controllers.NewChatController(chatRepository, chatUserRepository, chatRoleRepository, notificationService)

	// Init handlers
	messageHandler := handlers.NewMessageHandler(messageController)
	chatHandler := handlers.NewChatHandler(chatController)

	//Init services
	permissionsService := services.NewChatPermissionService(chatUserRepository)

	//Init middlewares
	permissionsMiddleware := middlewares.NewChatPermissionMiddleware(permissionsService)

	r := gin.Default()

	//r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	routes.RegisterChatRoutes(r, chatHandler, messageHandler, permissionsMiddleware)

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
