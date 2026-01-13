package database

import (
	"log"

	"github.com/glebarez/sqlite"
	"github.com/user/go-todo-api/internal/config"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect() {
	var err error
	DB, err = gorm.Open(sqlite.Open(config.AppConfig.DBPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Printf("Database connection established (SQLite: %s)", config.AppConfig.DBPath)
}

func GetDB() *gorm.DB {
	return DB
}
