package helpers

import (
	"errors"
	"log"
	"os"
	"steamshark-api/utils"

	"github.com/joho/godotenv"
)

type Config struct {
	Env    string
	Host   string // "development" | "production" | "test"
	Port   string
	DBPath string
}

func LoadConfig() (*Config, error) {
	// Try to load .env only if it exists (so Coolify won't complain)
	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(".ENV"); err != nil {
			return nil, errors.New("error loading .env")
		}
	} else {
		// No .env file – normal in Docker/Coolify
		log.Println("No .env file found; using environment variables only.")
	}

	config := Config{
		Env:  utils.GetEnv("APP_ENV", "development"),
		Host: utils.GetEnv("HOST", "0.0.0.0"),
		Port: utils.GetEnv("PORT", "8800"),
		// default to databases/steamshark.db relative to the app WORKDIR
		DBPath: utils.GetEnv("DB_PATH", "databases/steamshark.db"),
	}

	return &config, nil
}
