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
	ignoreErr  bool
}

type Option func(*config)

var defaultConfig = config{
	serder:     json.New(),
	keyPrefix:  "",
	defaultTtl: 0,
	ignoreErr:  false,
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

func WithIgnoreErr(ignoreErr bool) Option {
	return func(c *config) {
		c.ignoreErr = ignoreErr
	}
}
