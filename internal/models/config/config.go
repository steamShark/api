package model_config

type Config struct {
	Env      string
	Host     string // "development" | "production" | "test"
	Port     string
	AdminKey string // secret key required on admin-only endpoints

	//database
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
}
