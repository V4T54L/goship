package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisCache wraps a Redis client for caching operations.
type RedisCache struct {
	Client *redis.Client
}

// NewRedisCache creates a new RedisCache instance using an existing Redis client.
//
// Example:
//
//	client, _ := ConnectToRedisDb(ctx, "redis://localhost:6379/0")
//	cache := NewRedisCache(client)
func NewRedisCache(client *redis.Client) *RedisCache {
	return &RedisCache{Client: client}
}

// Set stores a value with a TTL (time-to-live in seconds).
//
// Example:
//
//	err := cache.Set(ctx, "user:1", "john", 60)
func (r *RedisCache) Set(ctx context.Context, key, value string, ttlSeconds int) error {
	return r.Client.Set(ctx, key, value, time.Duration(ttlSeconds)*time.Second).Err()
}

// Get retrieves a value by key.
//
// Example:
//
//	val, err := cache.Get(ctx, "user:1")
func (r *RedisCache) Get(ctx context.Context, key string) (string, error) {
	return r.Client.Get(ctx, key).Result()
}

// Delete removes a key from the cache.
//
// Example:
//
//	err := cache.Delete(ctx, "user:1")
func (r *RedisCache) Delete(ctx context.Context, key string) error {
	return r.Client.Del(ctx, key).Err()
}
