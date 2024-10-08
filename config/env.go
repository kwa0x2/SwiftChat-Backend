package config

import (
	"github.com/joho/godotenv"
	"log"
)

// region "LoadEnv" loads environment variables from a .env file.
func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

// endregion
