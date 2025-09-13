package db

import (
	"log"
	"os"
	"steamshark-api/models"

	"github.com/glebarez/sqlite"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

func InitUsersDB() *gorm.DB {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		log.Fatal("DB_PATH not set in .env")
	}

	// Open DB using glebarez/sqlite
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to DB_PATH: %v", err)
	}

	if err := db.AutoMigrate(&models.Website{}, &models.Occurrence{}); err != nil {
		log.Fatalf("AutoMigrate failed: %v", err)
	}

	log.Printf("Connected to SQLite DB: %s", dbPath)

	return db
}
