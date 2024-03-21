package db

import (
	"log"
	"os"

	"gorm.io/driver/postgres"

	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	dsn := os.Getenv("DB_URL")
	var err error

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to initialize Db", err)
	}
}
