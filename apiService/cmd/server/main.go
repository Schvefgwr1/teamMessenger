package main

import (
	"apiService/internal/controllers"
	"apiService/internal/handlers"
	"apiService/internal/http_clients"
	"apiService/internal/routes"
	"apiService/internal/services"
	common "common/config"
	"common/kafka"
	commonRedis "common/redis"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

// @title Service API
// @version 1.0
// @description API сервиса
// @host localhost:8084
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.url    http://www.swagger.io/support
// @contact.email support@swagger.io
// @license.name Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html
// @BasePath /api/v1
// @schemes http
// @externalDocs.description OpenAPI
// @externalDocs.url https://swagger.io/resources/open-api/

func main() {
	// Загружаем переменные окружения из .env файла (если существует)
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found or error loading .env file: %v", err)
	}

	//Upload config
	cfg, err := common.LoadConfig("config/config.yaml")
	if err != nil {
		log.Fatal("Can't load config " + err.Error())
		return
	}

	// Apply environment variable overrides
	common.ApplyRedisEnvOverrides(cfg)
	common.ApplyAppEnvOverrides(cfg)
	common.ApplyKafkaEnvOverrides(cfg)

	redisClient := commonRedis.NewRedisClient(&cfg.Redis)

	// Init Redis services
	sessionService := services.NewSessionService(redisClient)
	cacheService := services.NewCacheService(redisClient)

	// Init clients
	fileClient := http_clients.NewFileClient(common.GetEnvOrDefault("FILE_SERVICE_URL", "http://localhost:8080"))
	userClient := http_clients.NewUserClient(common.GetEnvOrDefault("USER_SERVICE_URL", "http://localhost:8082"))
	chatClient := http_clients.NewChatClient(common.GetEnvOrDefault("CHAT_SERVICE_URL", "http://localhost:8083"))
	taskClient := http_clients.NewTaskClient(common.GetEnvOrDefault("TASK_SERVICE_URL", "http://localhost:8081"))

	// Init PublicKeyManager
	publicKeyManager := services.NewPublicKeyManager()

	// Load initial public key from userService
	errLoad := services.LoadPublicKeyFromService(userClient, publicKeyManager)
	if errLoad != nil {
		log.Fatalf("Could not load initial public key: %v", errLoad)
	}
	log.Printf("Initial public key loaded (version %d)", publicKeyManager.GetKeyVersion())

	// Init Kafka key update consumer
	keyConsumerConfig := &services.KeyUpdateConsumerConfig{
		Brokers: kafka.GetKafkaBrokers(),
		Topic:   kafka.GetKeyUpdatesTopic(),
		GroupID: cfg.Kafka.GroupID,
	}

	keyUpdateConsumer, err := services.NewKeyUpdateConsumer(
		keyConsumerConfig,
		publicKeyManager,
		sessionService,
		redisClient,
	)
	if err != nil {
		log.Printf("Warning: Failed to initialize key update consumer: %v", err)
		keyUpdateConsumer = nil
	}

	// Start key update consumer
	var consumerCtx context.Context
	var consumerCancel context.CancelFunc
	if keyUpdateConsumer != nil {
		consumerCtx, consumerCancel = context.WithCancel(context.Background())
		go func() {
			if err := keyUpdateConsumer.Start(consumerCtx); err != nil {
				log.Printf("Key update consumer error: %v", err)
			}
		}()
		log.Println("Key update consumer started")
	}

	//Init controllers with cache service
	authController := controllers.NewAuthController(fileClient, userClient)
	userController := controllers.NewUserController(fileClient, userClient, cacheService)
	chatController := controllers.NewChatController(chatClient, fileClient, cacheService)
	taskController := controllers.NewTaskController(taskClient, fileClient)

	//Init handlers with session service
	authHandler := handlers.NewAuthHandler(authController, sessionService)
	userHandler := handlers.NewUserHandler(userController)
	chatHandler := handlers.NewChatHandler(chatController)
	taskHandler := handlers.NewTaskHandler(taskController)

	r := gin.Default()

	// Use new middleware with PublicKeyManager
	routes.RegisterAuthRoutes(r, authHandler, publicKeyManager, sessionService)
	routes.RegisterUserRoutes(r, userHandler, publicKeyManager, sessionService)
	routes.RegisterChatRoutes(r, chatHandler, publicKeyManager, sessionService)
	routes.RegisterTaskRoutes(r, taskHandler, publicKeyManager, sessionService)

	// Graceful shutdown
	go func() {
		_ = r.Run(":" + strconv.Itoa(cfg.App.Port))
	}()

	// Ожидание сигнала завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Остановка сервисов
	if consumerCancel != nil {
		consumerCancel()
	}

	if keyUpdateConsumer != nil {
		if err := keyUpdateConsumer.Close(); err != nil {
			log.Printf("Error closing key update consumer: %v", err)
		}
	}

	log.Println("Server exited")
}
