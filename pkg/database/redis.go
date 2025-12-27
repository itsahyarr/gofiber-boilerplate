package database

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"

	"github.com/itsahyarr/gofiber-boilerplate/pkg/logger"
	"go.uber.org/zap"
)

// Redis holds the Redis client
type Redis struct {
	Client *redis.Client
}

// NewRedis creates a new Redis connection
func NewRedis(host, port, password string, db int) (*Redis, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: password,
		DB:       db,
	})

	// Ping the database to verify connection
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		logger.Error("Failed to connect to Redis", zap.Error(err))
		return nil, err
	}

	logger.Info("Connected to Redis", zap.String("host", host), zap.String("port", port))

	return &Redis{
		Client: client,
	}, nil
}

// Close disconnects from Redis
func (r *Redis) Close() error {
	if err := r.Client.Close(); err != nil {
		logger.Error("Failed to disconnect from Redis", zap.Error(err))
		return err
	}
	logger.Info("Disconnected from Redis")
	return nil
}
