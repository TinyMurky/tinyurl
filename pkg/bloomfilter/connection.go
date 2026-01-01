package bloomfilter

import (
	"context"

	"github.com/redis/go-redis/v9"

	"github.com/TinyMurky/tinyurl/pkg/logging"
)

// BloomFilter is the pool of database Conn
type BloomFilter struct {
	RDB *redis.Client
}

// NewFromEnv sets up the redis connections using the configuration in the
// process's environment variables. This should be called just once per server
// instance.
func NewFromEnv(ctx context.Context, cfg *Config) (*BloomFilter, error) {
	logger := logging.FromContext(ctx)

	rdbOpt := redis.Options{
		Addr:     cfg.RedisConnectURL,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisBloomFilterDB,
		Protocol: cfg.RedisSerializationProtocol,
	}

	rdb := redis.NewClient(&rdbOpt)

	logger.Infof("Open redis cache at URL: %s", rdbOpt.Addr)

	return &BloomFilter{
		RDB: rdb,
	}, nil
}

// Close will close connection with redis
func (c *BloomFilter) Close() error {
	return c.RDB.Close()
}
