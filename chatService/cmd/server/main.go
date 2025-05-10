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
	"github.com/gin-gonic/gin"
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

	// Init repositories
	messageRepository := repositories.NewMessageRepository(db)
	chatRoleRepository := repositories.NewChatRoleRepository(db)
	chatUserRepository := repositories.NewChatUserRepository(db)
	chatRepository := repositories.NewChatRepository(db)

	// Init controllers
	messageController := controllers.NewMessageController(messageRepository, chatRepository, chatUserRepository)
	chatController := controllers.NewChatController(chatRepository, chatUserRepository, chatRoleRepository)

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

	_ = r.Run(":" + strconv.Itoa(cfg.App.Port))
}
