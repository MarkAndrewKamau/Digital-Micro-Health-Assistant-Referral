package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	Client *redis.Client
}

func NewRedis(redisURL string) (*Redis, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse redis URL: %w", err)
	}

	client := redis.NewClient(opts)

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to ping redis: %w", err)
	}

	log.Println("Redis connection established successfully")

	return &Redis{Client: client}, nil
}

func (r *Redis) Close() error {
	if r.Client != nil {
		if err := r.Client.Close(); err != nil {
			return fmt.Errorf("failed to close redis connection: %w", err)
		}
		log.Println("Redis connection closed")
	}
	return nil
}

func (r *Redis) Health(ctx context.Context) error {
	return r.Client.Ping(ctx).Err()
}