# Cacheable

Golang simplified implentation of the java spring cacheable annotation.

## Installation

`go get github.com/pedrofaria/cacheable`

## Usage

main.go
```golang
package main

import (
    "context"
    "log"

    "github.com/pedrofaria/cachable"
    "github.com/pedrofaria/cachable/driver/redisdb"
)

func main() {
    cacheNotification := cacheable.New[model.Notification](
		redisdb.New(redisConn, false),
		cacheable.WithKeyPrefix("notification:communication_type_id:"),
		cacheable.WithTtl(24*time.Hour),
	)
	defer cacheNotification.Close()

    repo := repository.NewNotificationRepository(dbConn)
    storage := service.NewCostumerNotificationStorage(repo, cacheNotification)

    ctx := context.Background()

    n, err := storage.FetchNotification(ctx, "2")
    if err != nil {
        log.Fatalf("failed during fetch notification: %s", err.Error())
    }
}
```

internal/service/storage.go
```golang
package service

import (
	"context"

	"github.com/pedrofaria/notifications/internal/model"
	"github.com/pedrofaria/cacheable"
)

type CostumerNotificationRepositorier interface {
	Fetch(ctx context.Context, id string) (*model.Notification, error)
}

type CostumerNotificationStorage struct {
	cache cacheable.Cacheler[model.Notification]
	repo  CostumerNotificationRepositorier
}

func NewCostumerNotificationStorage(repo CostumerNotificationRepositorier, cache cacheable.Cacheler[model.Notification]) *CostumerNotificationStorage {
	return &CostumerNotificationStorage{
		cache: cache,
		repo:  repo,
	}
}

func (s *CostumerNotificationStorage) FetchNotification(ctx context.Context, id string) (*model.Notification, error) {
    // Try to fetch from cache. If it's not present in cache, will call repo.Fetch and will store on cache.
	return s.cache.Load(ctx, id, func(ctx context.Context) (*model.Notification, error) {
		return s.repo.Fetch(ctx, id)
	})
}
```

## Drivers

* Redis
* Ristretto (in memory)

### Redis

https://github.com/pedrofaria/cacheable/tree/main/driver/redisdb

```golang
import "github.com/pedrofaria/cachable/driver/redisdb"
```

Example:

```golang
package main

import (
    "context"
    "log"

	"github.com/redis/go-redis/v9"
    "github.com/pedrofaria/cachable"
    "github.com/pedrofaria/cachable/driver/redisdb"
)

func main() {
	redisConn := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:          []string{"localhost:6379"},
		RouteByLatency: true,
		Username:       "default",
		Password:       "default",
	})

	cache := cacheable.New[model.Notification](
		redisdb.New(redisConn, false),
		cacheable.WithKeyPrefix("notification:communication_type_id:"),
		cacheable.WithTtl(24*time.Hour),
	)
	defer cache.Close()

	// ...
}
```

### Ristretto

https://github.com/pedrofaria/cacheable/tree/main/driver/ristretto

Library: https://github.com/dgraph-io/ristretto

```golang
import "github.com/pedrofaria/cachable/driver/ristretto"
```

Example:

```golang
package main

import (
    "context"
    "log"

    "github.com/pedrofaria/cachable"
    "github.com/pedrofaria/cachable/driver/ristretto"
)

func main() {
	ristrettoCache, _ := ristretto.NewCache(&ristretto.Config[string, []byte]{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 30, // maximum cost of cache (1GB).
		BufferItems: 64,      // number of keys per Get buffer.
	})

	cache := cacheable.New[model.Notification](
		ristretto..New(ristrettoCache),
		cacheable.WithKeyPrefix("notification:communication_type_id:"),
		cacheable.WithTtl(24*time.Hour),
	)
	defer cache.Close()

	// ...
}
```