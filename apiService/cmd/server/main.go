package main

import (
	"apiService/internal/controllers"
	"apiService/internal/handlers"
	"apiService/internal/http_clients"
	"apiService/internal/routes"
	"apiService/internal/services"
	"common/config"
	"crypto/rsa"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
	_ "userService/docs"
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
	//Upload config
	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		log.Fatal("Can't load config " + err.Error())
		return
	}

	// Init clients
	fileClient := http_clients.NewFileClient("http://localhost:8080")
	userClient := http_clients.NewUserClient("http://localhost:8082")
	chatClient := http_clients.NewChatClient("http://localhost:8083")
	taskClient := http_clients.NewTaskClient("http://localhost:8081")

	var publicKey *rsa.PublicKey

	errLoad := services.LoadPublicKeyFromService(userClient, &publicKey)
	if errLoad != nil {
		log.Fatalf("Could not load public key: %v", errLoad)
	}

	//Init controllers
	authController := controllers.NewAuthController(fileClient, userClient)
	userController := controllers.NewUserController(fileClient, userClient)
	chatController := controllers.NewChatController(chatClient, fileClient)
	taskController := controllers.NewTaskController(taskClient, fileClient)

	//Init handlers
	authHandler := handlers.NewAuthHandler(authController)
	userHandler := handlers.NewUserHandler(userController)
	chatHandler := handlers.NewChatHandler(chatController)
	taskHandler := handlers.NewTaskHandler(taskController)

	r := gin.Default()

	routes.RegisterAuthRoutes(r, authHandler)
	routes.RegisterUserRoutes(r, userHandler, publicKey)
	routes.RegisterChatRoutes(r, chatHandler, publicKey)
	routes.RegisterTaskRoutes(r, taskHandler, publicKey)

	_ = r.Run(":" + strconv.Itoa(cfg.App.Port))
}
