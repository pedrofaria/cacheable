package cacheable

import (
	"context"
	"errors"
	"time"

	"github.com/pedrofaria/cacheable/driver"
	"github.com/pedrofaria/cacheable/serder"
)

// Cacheler is a generic interface for maintaining data in a cache.
// It defines a Load method that retrieves data from the cache if available,
// or calls loadFn to fetch and potentially store data otherwise.
type Cacheler[T any] interface {
	// Load method takes a context for handling cancellation and timeouts,
	// a key to identify the cached item, and a loadFn which is called when
	// the item is not found in the cache or needs to be refreshed.
	Load(ctx context.Context, key string, loadFn func(ctx context.Context) (*T, error)) (*T, error)
}

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

// Load returns a value from the cache by key. If the key does not exist in the cache,
// it will call the provided loadFn to fetch the value, cache it, and then return it.
// If the key exists, it will deserialize and return the cached value.
//
// Parameters:
//   - ctx: The context for the operation.
//   - key: The cache key without the prefix (prefix will be added automatically).
//   - loadFn: A function to load the value if it's not found in the cache.
//
// Returns:
//   - *T: The retrieved or loaded value.
//   - error: Any error that occurred during the operation.
func (c *cacheable[T]) Load(ctx context.Context, key string, loadFn func(ctx context.Context) (*T, error)) (*T, error) {
	data, err := c.driver.Get(ctx, c.keyPrefix+key)

	if err != nil && errors.Is(err, driver.ErrNotFound) {
		obj, err := loadFn(ctx)
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

// New creates a new cacheable instance with the given driver and options.
// It returns a pointer to a cacheable instance.
//
// The driver parameter is required and will be used to store and retrieve cache entries.
// The opts parameter is optional and can be used to customize the cacheable instance.
//
// Example:
//
//	cache := cachable.New[MyType](redisDriver,
//		cachable.WithKeyPrefix("myapp:"),
//		cachable.WithTTL(time.Hour),
//	)
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
