package main

import (
	"common/config"
	"common/db"
	"common/kafka"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	_ "userService/docs"
	"userService/internal/controllers"
	"userService/internal/handlers"
	"userService/internal/repositories"
	"userService/internal/routes"
	"userService/internal/services"
	"userService/internal/utils"
)

// @title User Service API
// @version 1.0
// @description API сервиса для работы с пользователями
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

// @tag.name keys
// @tag.description Операции с ключами шифрования

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
	config.ApplyKeysEnvOverrides(cfg)

	// Initialize keys if they don't exist
	if err := initializeKeysIfNeeded(); err != nil {
		log.Printf("Warning: Failed to initialize keys: %v", err)
	}

	//Init DB
	db, err := db.InitDB(cfg)
	if err != nil {
		log.Fatal("Can't init DB: " + err.Error())
		return
	}

	// Init Kafka notification service
	notificationKafkaConfig := &kafka.ProducerConfig{
		Brokers: kafka.GetKafkaBrokers(),
		Topic:   kafka.GetNotificationsTopic(),
	}

	notificationService, err := services.NewNotificationService(notificationKafkaConfig)
	if err != nil {
		log.Printf("Warning: Failed to initialize notification service: %v", err)
		notificationService = nil
	}

	// Init Kafka key management service
	keyKafkaConfig := &kafka.ProducerConfig{
		Brokers: kafka.GetKafkaBrokers(),
		Topic:   kafka.GetKeyUpdatesTopic(),
	}

	keyManagementService, err := services.NewKeyManagementService(keyKafkaConfig)
	if err != nil {
		log.Printf("Warning: Failed to initialize key management service: %v", err)
		keyManagementService = nil
	}

	// Init key scheduler
	var keySchedulerService *services.KeySchedulerService
	if keyManagementService != nil {
		keyRotationInterval, err := cfg.GetKeyRotationInterval()
		if err != nil {
			log.Printf("Warning: Invalid key rotation interval in config, using default: %v", err)
			keyRotationInterval = 24 * 60 * 60 * 1000000000 // 24 часа в наносекундах
		}

		keySchedulerService = services.NewKeySchedulerService(keyManagementService, keyRotationInterval)

		// Запускаем scheduler в отдельной горутине
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		keySchedulerService.Start(ctx)

		log.Printf("Key rotation scheduler started with interval: %v", keyRotationInterval)
	}

	// Init repositories
	permissionRepository := repositories.NewPermissionRepository(db)
	roleRepository := repositories.NewRoleRepository(db)
	userRepository := repositories.NewUserRepository(db)

	// Init controllers
	permissionController := controllers.NewPermissionController(permissionRepository)
	roleController := controllers.NewRoleController(roleRepository, permissionRepository)
	userController := controllers.NewUserController(userRepository, roleRepository)
	authController := controllers.NewAuthController(userRepository, roleRepository, notificationService)

	// Init handlers
	permissionHandler := handlers.NewPermissionHandler(permissionController)
	roleHandler := handlers.NewRoleHandler(roleController)
	userHandler := handlers.NewUserHandler(userController)
	authHandler := handlers.NewAuthHandler(authController)
	keyHandler := handlers.NewKeyHandler(keyManagementService)

	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	routes.RegisterRoutes(r, authHandler, userHandler, roleHandler, permissionHandler, keyHandler)

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
	if keySchedulerService != nil {
		keySchedulerService.Stop()
	}

	if keyManagementService != nil {
		if err := keyManagementService.Close(); err != nil {
			log.Printf("Error closing key management service: %v", err)
		}
	}

	if notificationService != nil {
		if err := notificationService.Close(); err != nil {
			log.Printf("Error closing notification service: %v", err)
		}
	}

	log.Println("Server exited")
}

// initializeKeysIfNeeded проверяет наличие ключей и создает их если их нет
func initializeKeysIfNeeded() error {
	// Пытаемся загрузить существующий публичный ключ
	_, err := utils.ExtractPublicKeyFromFile()
	if err == nil {
		// Ключи уже существуют
		log.Println("RSA keys already exist")
		return nil
	}

	// Ключей нет, генерируем новые
	log.Println("RSA keys not found, generating new keys...")
	_, err = utils.GenerateAndSaveNewKeys()
	if err != nil {
		return err
	}

	log.Println("RSA keys generated successfully")
	return nil
}
