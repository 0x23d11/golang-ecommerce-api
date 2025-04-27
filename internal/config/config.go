package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all configuration for our application
type Config struct {
	ServerPort  string
	PostgresURL string
	MongoURI    string
}

// LoadConfig loads configuration from environment variables
// It first checks for environment variables and then falls back to default values
func LoadConfig() *Config {
	// Attempt to load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, reading environment variables directly")
	}

	// Return the configuration
	return &Config{
		ServerPort:  os.Getenv("SERVER_PORT"),
		PostgresURL: os.Getenv("POSTGRES_URL"),
		MongoURI:    os.Getenv("MONGO_URI"),
	}
}

// Helper function to get environment variables with fallback
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	log.Printf("Environment variable %s not set, using default value: %s", key, fallback)
	return fallback
}

// Helper function to get environment variable as integer
func getEnvAsInt(key string, fallback int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	log.Printf("Environment variable %s is not a valid integer, using default value: %d", key, fallback)
	return fallback
}

// Helper function to get environment variable as boolean
func getEnvAsBool(key string, fallback bool) bool {
	valueStr := getEnv(key, "")
	if value, err := strconv.ParseBool(valueStr); err == nil {
		return value
	}
	log.Printf("Environment variable %s is not a valid boolean, using default value: %t", key, fallback)
	return fallback
}
