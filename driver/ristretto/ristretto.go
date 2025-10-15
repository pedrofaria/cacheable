package ristretto

import (
	"context"
	"time"

	rClt "github.com/dgraph-io/ristretto/v2"
	"github.com/pedrofaria/cacheable/driver"
)

type client interface {
	Get(key string) (value []byte, ok bool)
	SetWithTTL(key string, value []byte, cost int64, ttl time.Duration) bool
	Del(key string)
	Close()
	Wait()
}

type Config[K string, V any] = rClt.Config[K, V]
type Cache[K string, V any] = rClt.Cache[K, V]

// NewCache creates a new Ristretto cache instance with the provided configuration.
// It returns a cache that stores string keys and byte slice values, or an error
// if the cache creation fails due to invalid configuration parameters.
//
// The config parameter must contain valid Ristretto cache settings including
// buffer sizes, cost calculations, and other cache-specific options.
//
// Returns an error if the underlying Ristretto cache cannot be initialized
// with the given configuration.
func NewCache(config *Config[string, []byte]) (*Cache[string, []byte], error) {
	return rClt.NewCache(config)
}

type ristrettoDriver struct {
	client client
}

// New creates a new ristretto driver instance with the provided client.
// It returns a pointer to ristrettoDriver that wraps the given client for caching operations.
func New(clt client) *ristrettoDriver {
	return &ristrettoDriver{
		client: clt,
	}
}

func (r *ristrettoDriver) Get(ctx context.Context, key string) ([]byte, error) {
	v, ok := r.client.Get(key)
	if !ok {
		return nil, driver.ErrNotFound
	}

	return v, nil
}

func (r *ristrettoDriver) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	ok := r.client.SetWithTTL(key, value, 1, ttl)

	if !ok {
		return driver.ErrKeyExists
	}

	r.client.Wait()

	return nil
}

func (r *ristrettoDriver) Del(ctx context.Context, key string) error {
	r.client.Del(key)

	return nil
}

func (r *ristrettoDriver) Close() error {
	r.client.Close()

	return nil
}
