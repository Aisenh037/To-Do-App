package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port           string
	DBPath         string // New field for SQLite
	JWTSecret      string
	JWTExpiryHours int
}

var AppConfig *Config

func LoadConfig() {
	// Load .env file (ignore error if not present)
	_ = godotenv.Load()

	jwtExpiry, err := strconv.Atoi(getEnv("JWT_EXPIRY_HOURS", "24"))
	if err != nil {
		jwtExpiry = 24
	}

	AppConfig = &Config{
		Port:           getEnv("PORT", "8080"),
		DBPath:         getEnv("DB_PATH", "todo.db"),
		JWTSecret:      getEnv("JWT_SECRET", "default-secret-change-me"),
		JWTExpiryHours: jwtExpiry,
	}

	log.Printf("Configuration loaded: Port=%s, DBPath=%s", AppConfig.Port, AppConfig.DBPath)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
