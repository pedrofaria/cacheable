package cacheable

import (
	"context"
	"errors"
	"time"

	"github.com/pedrofaria/cacheable/driver"
	"github.com/pedrofaria/cacheable/serder"
)

type cacheable[T any] struct {
	driver     driver.Driver
	serder     serder.Serder
	keyPrefix  string
	defaultTtl time.Duration
}

func NewCacheable[T any](driver driver.Driver, serder serder.Serder, keyPrefix string, defaultTtl time.Duration) *cacheable[T] {
	return &cacheable[T]{
		driver:     driver,
		serder:     serder,
		keyPrefix:  keyPrefix,
		defaultTtl: defaultTtl,
	}
}

func (c *cacheable[T]) Get(ctx context.Context, key string, loadFn func(ctx context.Context, id string) (*T, error)) (*T, error) {
	data, err := c.driver.Get(ctx, c.keyPrefix+key)

	if err != nil && errors.Is(err, driver.ErrNotFound) {
		obj, err := loadFn(ctx, key)
		if err != nil {
			return nil, err
		}

		data, err := c.serder.Serialize(obj)
		if err != nil {
			return nil, err
		}

		if err := c.driver.Set(ctx, c.keyPrefix+key, data, c.defaultTtl); err != nil {
			return nil, err
		}

		return obj, nil
	}

	if err != nil {
		return nil, err
	}

	var obj T
	if err := c.serder.Deserialize(data, &obj); err != nil {
		return nil, err
	}

	return &obj, nil
}

func New[T any](driver driver.Driver, opts ...Option) *cacheable[T] {
	cfg := defaultConfig

	for _, opt := range opts {
		opt(&cfg)
	}

	c := &cacheable[T]{
		driver:     driver,
		serder:     cfg.serder,
		keyPrefix:  cfg.keyPrefix,
		defaultTtl: cfg.defaultTtl,
	}

	return c
}
