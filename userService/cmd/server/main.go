package main

import (
	"common/config"
	"common/db"
	"common/kafka"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"strconv"
	_ "userService/docs"
	"userService/internal/controllers"
	"userService/internal/handlers"
	"userService/internal/repositories"
	"userService/internal/routes"
	"userService/internal/services"
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

func main() {
	//Upload config
	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		log.Fatal("Can't load config " + err.Error())
		return
	}

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
	keyHandler := handlers.NewKeyHandler()

	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	routes.RegisterRoutes(r, authHandler, userHandler, roleHandler, permissionHandler, keyHandler)

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
