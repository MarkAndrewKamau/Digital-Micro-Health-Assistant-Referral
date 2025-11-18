package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Environment string
	Port        string

	// Database
	DatabaseURL string

	// Redis
	RedisURL string

	// Session
	SessionDuration time.Duration
}

func Load() *Config {
	sessionDurationHours, _ := strconv.Atoi(getEnv("SESSION_DURATION_HOURS", "720")) // 30 days default

	return &Config{
		Environment: getEnv("ENVIRONMENT", "development"),
		Port:        getEnv("PORT", "8080"),

		// Database
		DatabaseURL: getEnv("DATABASE_URL", "postgres://healthuser:healthpass@localhost:5432/healthdb?sslmode=disable"),

		// Redis
		RedisURL: getEnv("REDIS_URL", "redis://localhost:6379/0"),

		// Session
		SessionDuration: time.Duration(sessionDurationHours) * time.Hour,
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}