// Package cache is the config that can use redis as cache
package cache

// Config is the config of cache
type Config struct {
	RedisConnectURL            string `env:"REDIS_CONNECT_URL, default=localhost:6379"`
	RedisPassword              string `env:"REDIS_PASSWORD"`
	RedisCacheDB               int    `env:"REDIS_CACHE_DB, default=0"`
	RedisSerializationProtocol int    `env:"REDIS_SERIALIZATION_PROTOCAL, default=2"`
}

// CacheConfig return the config of cache
func (c *Config) CacheConfig() *Config {
	return c
}
