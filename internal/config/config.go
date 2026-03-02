package config

import (
	"errors"
	"log"
	"os"

	"steamshark-api/internal/models"
	"steamshark-api/internal/utils"

	"github.com/joho/godotenv"
)

func LoadConfig() (*models.Config, error) {
	// Try to load .env only if it exists (so Coolify won't complain)
	if _, err := os.Stat("../../.env"); err == nil {
		if err := godotenv.Load(".env"); err != nil {
			return nil, errors.New("error loading .env")
		}
	} else {
		// No .env file – normal in Docker/Coolify
		log.Println("No .env file found; using environment variables only.")
	}

	config := models.Config{
		Env:  utils.GetEnv("APP_ENV", "development"),
		Host: utils.GetEnv("HOST", "0.0.0.0"),
		Port: utils.GetEnv("PORT", "8800"),
		// default to databases/steamshark.db relative to the app WORKDIR
		DBHost:     utils.GetEnv("DB_HOST", "steamshark-db"), // docker service name
		DBPort:     utils.GetEnv("DB_PORT", "5432"),
		DBUser:     utils.GetEnv("POSTGRES_USER", "postgres"),
		DBPassword: utils.GetEnv("POSTGRES_PASSWORD", "postgres"),
		DBName:     utils.GetEnv("POSTGRES_DB", "steamshark"),
		DBSSLMode:  utils.GetEnv("DB_SSLMODE", "disable"),
	}

	return &config, nil
}
