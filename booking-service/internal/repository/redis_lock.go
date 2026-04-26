package repository

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisLock interface {
	AcquireLock(ctx context.Context, key string, expiration time.Duration) (bool, error)
	ReleaseLock(ctx context.Context, key string) error
}

type redisLockWrapper struct {
	client *redis.Client
}

func NewRedisLock(client *redis.Client) RedisLock {
	return &redisLockWrapper{client: client}
}

func (r *redisLockWrapper) AcquireLock(ctx context.Context, key string, expiration time.Duration) (bool, error) {
	// SetNX executes the "SET if Not eXists" redis command.
	return r.client.SetNX(ctx, key, "locked", expiration).Result()
}

func (r *redisLockWrapper) ReleaseLock(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}
