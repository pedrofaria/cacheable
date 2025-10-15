package cacheable

import (
	"context"
	"testing"

	"github.com/pedrofaria/cacheable/driver"
	"github.com/pedrofaria/cacheable/driver/test"
	"github.com/pedrofaria/cacheable/serder/binary"
	"github.com/pedrofaria/cacheable/serder/msgpack"
	"github.com/stretchr/testify/assert"
)

type Data struct {
	Name string
	Age  int
}

func Test_integration_Cacheable_Load_BinarySerder(t *testing.T) {
	ctx := context.Background()

	cache := New[Data](test.New(), WithKeyPrefix("user:"), WithSerder(binary.New()))
	defer cache.Close()

	d, err := cache.Load(ctx, "1", func(ctx context.Context) (*Data, error) {
		return &Data{Name: "John", Age: 30}, nil
	})

	assert.NoError(t, err)
	assert.Equal(t, &Data{Name: "John", Age: 30}, d)

	d2, err := cache.Load(ctx, "1", func(ctx context.Context) (*Data, error) {
		return &Data{Name: "Anna", Age: 10}, nil
	})

	assert.NoError(t, err)
	assert.Equal(t, &Data{Name: "John", Age: 30}, d2)
}

func Test_integration_Cacheable_Remove_BinarySerder(t *testing.T) {
	ctx := context.Background()

	cache := New[Data](test.New(), WithKeyPrefix("user:"), WithSerder(binary.New()))
	defer cache.Close()

	d, err := cache.Load(ctx, "1", func(ctx context.Context) (*Data, error) {
		return &Data{Name: "John", Age: 30}, nil
	})

	assert.NoError(t, err)
	assert.Equal(t, &Data{Name: "John", Age: 30}, d)

	err = cache.Remove(ctx, "1")
	assert.NoError(t, err)

	d2, err := cache.Load(ctx, "1", func(ctx context.Context) (*Data, error) {
		return &Data{Name: "Anna", Age: 10}, nil
	})

	assert.NoError(t, err)
	assert.Equal(t, &Data{Name: "Anna", Age: 10}, d2)
}

func Test_integration_Cacheable_Remove_BinarySerder_ErrNotFound(t *testing.T) {
	ctx := context.Background()

	cache := New[Data](test.New(), WithKeyPrefix("user:"), WithSerder(binary.New()))
	defer cache.Close()

	err := cache.Remove(ctx, "1")
	assert.ErrorIs(t, err, driver.ErrNotFound)
}

func Test_integration_Cacheable_Load_JsonSerder(t *testing.T) {
	ctx := context.Background()

	cache := New[Data](test.New(), WithKeyPrefix("user:"))
	defer cache.Close()

	d, err := cache.Load(ctx, "1", func(ctx context.Context) (*Data, error) {
		return &Data{Name: "John", Age: 30}, nil
	})

	assert.NoError(t, err)
	assert.Equal(t, &Data{Name: "John", Age: 30}, d)

	d2, err := cache.Load(ctx, "1", func(ctx context.Context) (*Data, error) {
		return &Data{Name: "Anna", Age: 10}, nil
	})

	assert.NoError(t, err)
	assert.Equal(t, &Data{Name: "John", Age: 30}, d2)
}

func Test_integration_Cacheable_Load_MsgpackSerder(t *testing.T) {
	ctx := context.Background()

	cache := New[Data](test.New(), WithKeyPrefix("user:"), WithSerder(msgpack.New()))
	defer cache.Close()

	d, err := cache.Load(ctx, "1", func(ctx context.Context) (*Data, error) {
		return &Data{Name: "John", Age: 30}, nil
	})

	assert.NoError(t, err)
	assert.Equal(t, &Data{Name: "John", Age: 30}, d)

	d2, err := cache.Load(ctx, "1", func(ctx context.Context) (*Data, error) {
		return &Data{Name: "Anna", Age: 10}, nil
	})

	assert.NoError(t, err)
	assert.Equal(t, &Data{Name: "John", Age: 30}, d2)
}

func Benchmark_Cacheable_Load_Binary(b *testing.B) {
	ctx := context.Background()

	cache := New[Data](test.New(), WithKeyPrefix("user:"), WithSerder(binary.New()))
	defer cache.Close()

	for i := 0; i < b.N; i++ {
		_, err := cache.Load(ctx, "1", func(ctx context.Context) (*Data, error) {
			return &Data{Name: "John", Age: 30}, nil
		})
		if err != nil {
			b.Fatal(err)
		}
	}
}

func Benchmark_Cacheable_Load_Json(b *testing.B) {
	ctx := context.Background()

	cache := New[Data](test.New(), WithKeyPrefix("user:"))
	defer cache.Close()

	for i := 0; i < b.N; i++ {
		_, err := cache.Load(ctx, "1", func(ctx context.Context) (*Data, error) {
			return &Data{Name: "John", Age: 30}, nil
		})
		if err != nil {
			b.Fatal(err)
		}
	}
}

func Benchmark_Cacheable_Load_Msgpack(b *testing.B) {
	ctx := context.Background()

	cache := New[Data](test.New(), WithKeyPrefix("user:"), WithSerder(msgpack.New()))
	defer cache.Close()

	for i := 0; i < b.N; i++ {
		_, err := cache.Load(ctx, "1", func(ctx context.Context) (*Data, error) {
			return &Data{Name: "John", Age: 30}, nil
		})
		if err != nil {
			b.Fatal(err)
		}
	}
}
