package db

import (
	"errors"
	"fmt"
	"steamshark-api/models"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func InitDB(dbPath string) (*gorm.DB, error) {
	// Open DB using glebarez/sqlite
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		fmt.Printf("Failed to connect to DB_PATH: %v", err)
		return nil, errors.New("Failed to connect to DB_PATH")
	}

	if err := db.AutoMigrate(&models.Website{}, &models.Occurrence{}); err != nil {
		return nil, errors.Join(errors.New("Auto migrate failed: "), err)
	}

	return db, nil
}
