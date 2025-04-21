package main

import (
	"common/config"
	"common/db"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
	"userService/internal/controllers"
	"userService/internal/handlers"
	"userService/internal/repositories"
	"userService/internal/routes"
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

	// Init repositories
	permissionRepository := repositories.NewPermissionRepository(db)
	roleRepository := repositories.NewRoleRepository(db)
	userRepository := repositories.NewUserRepository(db)

	// Init controllers
	permissionController := controllers.NewPermissionController(permissionRepository)
	roleController := controllers.NewRoleController(roleRepository, permissionRepository)
	userController := controllers.NewUserController(userRepository, roleRepository)
	authController := controllers.NewAuthController(userRepository, roleRepository)

	// Init handlers
	permissionHandler := handlers.NewPermissionHandler(permissionController)
	roleHandler := handlers.NewRoleHandler(roleController)
	userHandler := handlers.NewUserHandler(userController)
	authHandler := handlers.NewAuthHandler(authController)

	r := gin.Default()

	routes.RegisterRoutes(r, authHandler, userHandler, roleHandler, permissionHandler)

	_ = r.Run(":" + strconv.Itoa(cfg.App.Port))
}
