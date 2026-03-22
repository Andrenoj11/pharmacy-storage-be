package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort         string
	StorageProvider string
}

func Load() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found")
	}

	cfg := &Config{
		AppPort:         getEnv("APP_PORT", "8080"),
		StorageProvider: getEnv("STORAGE_PROVIDER", "spreadsheet"),
	}

	return cfg, nil
}

func getEnv(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
