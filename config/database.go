package config

import (
	"fmt"
	"os"
	"time"

	"github.com/Chamanthra/TaskManager/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	var db *gorm.DB

	// Retry logic for Docker container startup
	for i := 0; i < 5; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		panic(fmt.Sprintf("Failed to connect to database after 5 attempts: %v", err))
	}

	DB = db

	// Enable auto-migration
	err = DB.AutoMigrate(
		&models.User{},
		&models.Task{},
		// Add all your models here
	)
	if err != nil {
		panic(fmt.Sprintf("AutoMigrate failed: %v", err))
	}
}
