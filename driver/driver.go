package driver

import (
	"context"
	"errors"
	"time"
)

type Driver interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
	Del(ctx context.Context, key string) error
}

var (
	ErrNotFound  = errors.New("not found")
	ErrKeyExists = errors.New("key exists")
)
