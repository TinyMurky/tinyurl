// Package bloomfilter is the config that can use redis as bloomfilter
package bloomfilter

// Config is the config of cache
type Config struct {
	RedisConnectURL            string  `env:"REDIS_CONNECT_URL, default=localhost:6379"`
	RedisPassword              string  `env:"REDIS_PASSWORD"`
	RedisBloomFilterDB         int     `env:"REDIS_BLOOM_FILTER_DB, default=1"`
	RedisBloomFilterErrorRate  float64 `env:"REDIS_BLOOM_FILTER_ERROR_RATE, default=0.001"`
	RedisBloomFilterCapacity   int64   `env:"REDIS_BLOOM_FILTER_CAPACITY, default=1000"`
	RedisSerializationProtocol int     `env:"REDIS_SERIALIZATION_PROTOCAL, default=2"`
}

// BloomFilterConfig return the config of cache
func (c *Config) BloomFilterConfig() *Config {
	return c
}
