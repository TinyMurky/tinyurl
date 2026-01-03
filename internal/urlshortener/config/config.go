// Package urlshortenerconfig define the config of urlshortener
package urlshortenerconfig

import (
	"github.com/TinyMurky/tinyurl/pkg/bloomfilter"
	"github.com/TinyMurky/tinyurl/pkg/cache"
	"github.com/TinyMurky/tinyurl/pkg/database"
	"github.com/TinyMurky/tinyurl/pkg/singleflight"
)

// Config for urlshortener
// it need to use github.com/sethvargo/go-envconfig package to read
type Config struct {
	Database     database.Config
	Cache        cache.Config
	BloomFilter  bloomfilter.Config
	SingleFlight singleflight.Config

	IDGenerator            IDGeneratorConfig
	Port                   string `env:"PORT"`
	ShortURLPrefix         string `env:"SHORT_URL_PREFIX, default=http://localhost:3000"`
	RedisCacheTTLInMiliSec int    `env:"SHORT_URL_CACHE_TTL_IN_MILI_SEC, default=300000"`
}

// DatabaseConfig return Database config
func (c *Config) DatabaseConfig() *database.Config {
	return &c.Database
}

// CacheConfig return the config of cache
func (c *Config) CacheConfig() *cache.Config {
	return &c.Cache
}

// BloomFilterConfig return the config of cache
func (c *Config) BloomFilterConfig() *bloomfilter.Config {
	return &c.BloomFilter
}

// SingleFlightConfig return the config of singleflight
func (c *Config) SingleFlightConfig() *singleflight.Config {
	return &c.SingleFlight
}
