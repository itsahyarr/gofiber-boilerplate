package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/itsahyarr/gofiber-boilerplate/pkg/database"
)

const refreshTokenPrefix = "refresh_token:"

var (
	ErrTokenNotFound = errors.New("refresh token not found")
)

// TokenRepository defines the interface for refresh token management
type TokenRepository interface {
	Store(ctx context.Context, userID string, token string, expiration time.Duration) error
	Get(ctx context.Context, userID string) (string, error)
	Delete(ctx context.Context, userID string) error
	Exists(ctx context.Context, userID string) (bool, error)
}

type tokenRepositoryRedis struct {
	redis *database.Redis
}

// NewTokenRepository creates a new Redis token repository
func NewTokenRepository(redis *database.Redis) TokenRepository {
	return &tokenRepositoryRedis{
		redis: redis,
	}
}

func (r *tokenRepositoryRedis) Store(ctx context.Context, userID string, token string, expiration time.Duration) error {
	key := fmt.Sprintf("%s%s", refreshTokenPrefix, userID)
	return r.redis.Client.Set(ctx, key, token, expiration).Err()
}

func (r *tokenRepositoryRedis) Get(ctx context.Context, userID string) (string, error) {
	key := fmt.Sprintf("%s%s", refreshTokenPrefix, userID)
	token, err := r.redis.Client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", ErrTokenNotFound
		}
		return "", err
	}
	return token, nil
}

func (r *tokenRepositoryRedis) Delete(ctx context.Context, userID string) error {
	key := fmt.Sprintf("%s%s", refreshTokenPrefix, userID)
	return r.redis.Client.Del(ctx, key).Err()
}

func (r *tokenRepositoryRedis) Exists(ctx context.Context, userID string) (bool, error) {
	key := fmt.Sprintf("%s%s", refreshTokenPrefix, userID)
	n, err := r.redis.Client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return n > 0, nil
}
