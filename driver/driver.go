package driver

import (
	"context"
	"errors"
	"fmt"
	"time"
)

type Driver interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
	Del(ctx context.Context, key string) error
	Close() error
}

var (
	ErrDriverError = errors.New("driver error")
	ErrNotFound    = fmt.Errorf("%w: key not found", ErrDriverError)
	ErrKeyExists   = fmt.Errorf("%w: key exists", ErrDriverError)
)
