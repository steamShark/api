package helpers

import (
	"errors"
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
	err := godotenv.Load(".env")
	if err != nil {
		return nil, errors.New("error loading .env")
	}

	config := Config{
		Env:  utils.GetEnv("ENV", "development"),
		Host: utils.GetEnv("HOST", "0.0.0.0"),
		Port: utils.GetEnv("PORT", "8800"),
		// default to databases/steamshark.db relative to the app WORKDIR
		DBPath: utils.GetEnv("DB_PATH", "databases/steamshark.db"),
	}

	if config.Host == "" || config.Port == "" || config.DBPath == "" {
		return nil, errors.New("some env variables were not defined")
	}

	return &config, nil
}
