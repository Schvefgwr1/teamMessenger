package main

import (
	"common/config"
	"common/db"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
	"taskService/internal/controllers"
	"taskService/internal/handlers"
	"taskService/internal/repositories"
	"taskService/internal/routes"
)

func main() {
	//Upload config
	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		log.Fatal("Can't load config " + err.Error())
		return
	}

	//Init DB
	initDB, err := db.InitDB(cfg)
	if err != nil {
		log.Fatal("Can't init DB: " + err.Error())
		return
	}

	//// Init repositories
	taskRepo := repositories.NewTaskRepository(initDB)
	taskFileRepo := repositories.NewTaskFileRepository(initDB)
	taskStatusRepo := repositories.NewTaskStatusRepository(initDB)

	//// Init controllers
	taskController := controllers.NewTaskController(taskRepo, taskStatusRepo, taskFileRepo)
	taskStatusController := controllers.NewTaskStatusController(taskStatusRepo)

	//// Init handlers
	taskHandler := handlers.NewTaskHandler(taskController)
	taskStatusHandler := handlers.NewTaskStatusHandler(taskStatusController)

	r := gin.Default()
	routes.RegisterTaskStatusRoutes(r, taskStatusHandler)
	routes.RegisterTaskRoutes(r, taskHandler)

	_ = r.Run(":" + strconv.Itoa(cfg.App.Port))
}
