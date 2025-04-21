package db

import (
	"common/config"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// InitDB подключается к PostgreSQL
func InitDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		cfg.Database.Host, cfg.Database.User, cfg.Database.Password, cfg.Database.Name, cfg.Database.Port,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	return db, err
}
