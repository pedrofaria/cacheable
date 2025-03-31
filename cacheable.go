package cacheable

import (
	"context"
	"errors"
	"sync"
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

	Remove(ctx context.Context, key string) error

	GetStats() Stats
}

var mtxStats = &sync.Mutex{}

type Stats struct {
	Hits       uint64
	Miss       uint64
	SetSuccess uint64
	SetError   uint64
	DelSuccess uint64
	DelError   uint64
}

type cacheable[T any] struct {
	driver     driver.Driver
	serder     serder.Serder
	keyPrefix  string
	defaultTtl time.Duration
	ignoreErr  bool
	stats      Stats
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
		atomicInc(&c.stats.Miss)

		obj, err := loadFn(ctx)
		if err != nil {
			return nil, err
		}

		data, err := c.serder.Serialize(obj)
		if err != nil {
			return nil, err
		}

		if err := c.driver.Set(ctx, c.keyPrefix+key, data, c.defaultTtl); err != nil {
			atomicInc(&c.stats.SetError)
			return nil, err
		}

		atomicInc(&c.stats.SetSuccess)

		return obj, nil
	}

	if err != nil {
		if c.ignoreErr {
			return loadFn(ctx)
		}

		return nil, err
	}

	var obj T
	if err := c.serder.Deserialize(data, &obj); err != nil {
		return nil, err
	}

	return &obj, nil
}

func (c *cacheable[T]) Remove(ctx context.Context, key string) error {
	if err := c.driver.Del(ctx, c.keyPrefix+key); err != nil {
		atomicInc(&c.stats.DelError)
		return err
	}

	atomicInc(&c.stats.DelSuccess)

	return nil
}

func (c *cacheable[T]) GetStats() Stats {
	return c.stats
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

	return &cacheable[T]{
		driver:     driver,
		serder:     cfg.serder,
		keyPrefix:  cfg.keyPrefix,
		defaultTtl: cfg.defaultTtl,
		ignoreErr:  cfg.ignoreErr,
	}
}

func atomicInc(v *uint64) {
	mtxStats.Lock()
	defer mtxStats.Unlock()

	if v == nil {
		return
	}

	*v += 1
}
