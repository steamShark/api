package helpers

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type EnvConfig struct {
	Env    string
	Host   string // "development" | "production" | "test"
	Port   string
	DBPath string
}

func LoadEnvConfig() (*EnvConfig, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return nil, errors.New("error loading .env")
	}

	envConfig := EnvConfig{
		Env:    os.Getenv("APP_ENV"),
		Host:   os.Getenv("HOST"),
		Port:   os.Getenv("PORT"),
		DBPath: os.Getenv("DB_PATH"),
	}

	if envConfig.Host == "" || envConfig.Port == "" || envConfig.DBPath == "" {
		return nil, errors.New("some env variables were not defined")
	}

	return &envConfig, nil
}
