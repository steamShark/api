package model_config

type Config struct {
	Env  string
	Host string // "development" | "production" | "test"
	Port string

	//database
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
}
