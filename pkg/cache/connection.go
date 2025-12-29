package cache

import (
	"context"

	"github.com/redis/go-redis/v9"

	"github.com/TinyMurky/tinyurl/pkg/logging"
)

// Cache is the pool of database Conn
type Cache struct {
	RDB *redis.Client
}

// NewFromEnv sets up the redis cache connections using the configuration in the
// process's environment variables. This should be called just once per server
// instance.
func NewFromEnv(ctx context.Context, cfg *Config) (*Cache, error) {
	logger := logging.FromContext(ctx)

	rdbOpt := redis.Options{
		Addr:     cfg.RedisConnectURL,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisCacheDB,
		Protocol: cfg.RedisSerializationProtocol,
	}

	rdb := redis.NewClient(&rdbOpt)

	logger.Infof("Open redis cache at URL: %s", rdbOpt.Addr)

	return &Cache{
		RDB: rdb,
	}, nil
}

// Close will close connection with redis
func (c *Cache) Close() error {
	return c.RDB.Close()
}
