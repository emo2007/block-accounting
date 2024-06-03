package cache

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/gob"
	"fmt"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache interface {
	// NOTE: dst MUST be a pointer
	Get(ctx context.Context, key any, dst any) error
	Cache(ctx context.Context, key any, val any, ttl time.Duration) error
}

type redisCache struct {
	log    *slog.Logger
	client *redis.Client
}

func NewRedisCache(
	log *slog.Logger,
	client *redis.Client,
) Cache {
	return &redisCache{
		log:    log,
		client: client,
	}
}

func (c *redisCache) Get(ctx context.Context, key any, dst any) error {
	res := c.client.Get(ctx, c.hashKeyStr(key))

	if res.Err() != nil {
		return fmt.Errorf("error fetch data from cache. %w", res.Err())
	}

	return res.Scan(dst)
}

func (c *redisCache) Cache(ctx context.Context, k any, v any, ttl time.Duration) error {
	res := c.client.Set(ctx, c.hashKeyStr(k), v, ttl)

	if res.Err() != nil {
		return fmt.Errorf("error add record to cache. %w", res.Err())
	}

	return nil
}

func (c *redisCache) hashKey(k any) []byte {
	var b bytes.Buffer

	gob.NewEncoder(&b).Encode(k)

	return b.Bytes()
}

func (c *redisCache) hashKeyStr(k any) string {
	return base64.StdEncoding.EncodeToString(c.hashKey(k))
}
