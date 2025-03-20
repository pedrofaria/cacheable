package cacheable

import (
	"time"

	"github.com/pedrofaria/cacheable/serder"
	"github.com/pedrofaria/cacheable/serder/json"
)

type config struct {
	serder     serder.Serder
	keyPrefix  string
	defaultTtl time.Duration
}

type Option func(*config)

var defaultConfig = config{
	serder:     json.NewJsonSerde(),
	keyPrefix:  "",
	defaultTtl: 0,
}

func WithSerder(serder serder.Serder) Option {
	return func(c *config) {
		c.serder = serder
	}
}

func WithKeyPrefix(keyPrefix string) Option {
	return func(c *config) {
		c.keyPrefix = keyPrefix
	}
}

func WithTtl(defaultTtl time.Duration) Option {
	return func(c *config) {
		c.defaultTtl = defaultTtl
	}
}
