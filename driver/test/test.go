package test

import (
	"context"
	"sync"
	"time"

	"github.com/pedrofaria/cacheable/driver"
)

type testDriver struct {
	data sync.Map
}

func New() *testDriver {
	return &testDriver{
		data: sync.Map{},
	}
}

func (r *testDriver) Get(ctx context.Context, key string) ([]byte, error) {
	v, ok := r.data.Load(key)
	if !ok {
		return nil, driver.ErrNotFound
	}

	return v.([]byte), nil
}

func (r *testDriver) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	r.data.Store(key, value)
	return nil
}

func (r *testDriver) Del(ctx context.Context, key string) error {
	_, ok := r.data.LoadAndDelete(key)
	if !ok {
		return driver.ErrNotFound
	}

	return nil
}

func (r *testDriver) Close() error {
	r.data.Clear()
	return nil
}
