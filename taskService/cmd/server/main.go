package main

import (
	"common/config"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
)

func main() {
	//Upload config
	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		log.Fatal("Can't load config " + err.Error())
		return
	}

	////Init DB
	//db, err := db.InitDB(cfg)
	//if err != nil {
	//	log.Fatal("Can't init DB: " + err.Error())
	//	return
	//}

	//// Init repositories
	//fileRepo := repositories.NewFileRepository(db)
	//fileTypeRepo := repositories.NewFileTypeRepository(db)
	//
	//// Init controllers
	//fileController := controllers.NewFileController(fileRepo, fileTypeRepo, minio, &cfg.MinIO)
	//fileTypeController := controllers.NewFileTypeController(fileTypeRepo)
	//
	//// Init handlers
	//fileHandler := handlers.NewFileHandler(fileController)
	//fileTypeHandler := handlers.NewFileTypeHandler(fileTypeController)

	r := gin.Default()

	_ = r.Run(":" + strconv.Itoa(cfg.App.Port))
}
