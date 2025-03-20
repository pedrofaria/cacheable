package redis

import (
	"context"
	"time"

	"github.com/pedrofaria/cacheable/driver"
	redisClient "github.com/redis/go-redis/v9"
)

type client interface {
	Get(ctx context.Context, key string) *redisClient.StringCmd
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redisClient.StatusCmd
	SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) *redisClient.BoolCmd
}

type RedisDriver struct {
	client    client
	useAtomic bool
}

func (r *RedisDriver) Get(ctx context.Context, key string) ([]byte, error) {
	res := r.client.Get(ctx, key)

	if res.Err() != nil {
		if res.Err() == redisClient.Nil {
			return nil, driver.ErrNotFound
		}

		return nil, res.Err()
	}

	return res.Bytes()
}

func (r *RedisDriver) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	if r.useAtomic {
		return r.setAtomic(ctx, key, value, ttl)
	}

	res := r.client.Set(ctx, key, value, ttl)
	if res.Err() != nil {
		return res.Err()
	}

	return nil
}

func (r *RedisDriver) setAtomic(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	res := r.client.SetNX(ctx, key, value, ttl)
	if res.Err() != nil {
		return res.Err()
	}

	if !res.Val() {
		return driver.ErrKeyExists
	}

	return nil
}
