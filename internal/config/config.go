package config

import (
	"os"
)

type Config struct {
	Environment string
	Port        string

	// Database
	DatabaseURL string

	// Redis
	RedisURL string

	// RabbitMQ
	RabbitMQURL string

	// MinIO/S3
	MinIOEndpoint  string
	MinIOAccessKey string
	MinIOSecretKey string
	MinIOBucket    string

	// Africa's Talking
	AfricasTalkingAPIKey  string
	AfricasTalkingUsername string

	// Safaricom Daraja (M-Pesa)
	MpesaConsumerKey    string
	MpesaConsumerSecret string
	MpesaPasskey        string
	MpesaShortcode      string

	// OpenAI (optional)
	OpenAIAPIKey string

	// JWT
	JWTSecret string
}

func Load() *Config {
	return &Config{
		Environment: getEnv("ENVIRONMENT", "development"),
		Port:        getEnv("PORT", "8080"),

		// Database
		DatabaseURL: getEnv("DATABASE_URL", "postgres://healthuser:healthpass@localhost:5432/healthdb?sslmode=disable"),

		// Redis
		RedisURL: getEnv("REDIS_URL", "redis://localhost:6379/0"),

		// RabbitMQ
		RabbitMQURL: getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),

		// MinIO
		MinIOEndpoint:  getEnv("MINIO_ENDPOINT", "localhost:9000"),
		MinIOAccessKey: getEnv("MINIO_ACCESS_KEY", "minioadmin"),
		MinIOSecretKey: getEnv("MINIO_SECRET_KEY", "minioadmin"),
		MinIOBucket:    getEnv("MINIO_BUCKET", "health-attachments"),

		// Africa's Talking
		AfricasTalkingAPIKey:  getEnv("AFRICASTALKING_API_KEY", ""),
		AfricasTalkingUsername: getEnv("AFRICASTALKING_USERNAME", "sandbox"),

		// M-Pesa
		MpesaConsumerKey:    getEnv("MPESA_CONSUMER_KEY", ""),
		MpesaConsumerSecret: getEnv("MPESA_CONSUMER_SECRET", ""),
		MpesaPasskey:        getEnv("MPESA_PASSKEY", ""),
		MpesaShortcode:      getEnv("MPESA_SHORTCODE", ""),

		// OpenAI
		OpenAIAPIKey: getEnv("OPENAI_API_KEY", ""),

		// JWT
		JWTSecret: getEnv("JWT_SECRET", "mark1234"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}