package cache

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

type Cache struct {
	client *redis.Client
}

func NewCache(addr string) *Cache {
	return &Cache{
		client: redis.NewClient(&redis.Options{Addr: addr}),
	}
}

func (c *Cache) Get(key string) (string, error) {
	return c.client.Get(ctx, key).Result()
}

func (c *Cache) Set(key string, value string, ttl time.Duration) error {
	return c.client.Set(ctx, key, value, ttl).Err()
}
