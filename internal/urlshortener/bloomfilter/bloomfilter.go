// Package bloomfilter can add/check exist of longURL or ID in bloom filter
package bloomfilter

import (
	"context"
	"fmt"

	"github.com/TinyMurky/tinyurl/internal/urlshortener/model"
	"github.com/TinyMurky/tinyurl/pkg/bloomfilter"
)

// URLShortenerBloomFilter is to create bloom filter for url shortener
type URLShortenerBloomFilter struct {
	filter *bloomfilter.BloomFilter
	cfg    *bloomfilter.Config
}

// New URLShortenerBloomFilter
func New(filter *bloomfilter.BloomFilter, cfg *bloomfilter.Config) *URLShortenerBloomFilter {
	bf := &URLShortenerBloomFilter{
		filter: filter,
		cfg:    cfg,
	}

	if err := bf.reserveBase62ID(context.Background()); err != nil {
		panic(fmt.Errorf("reserveBase62ID: %w", err))
	}

	return bf
}

func (bf *URLShortenerBloomFilter) reserveBase62ID(ctx context.Context) error {
	key := genBase62IDKey()
	errorRate := bf.cfg.RedisBloomFilterErrorRate
	capacity := bf.cfg.RedisBloomFilterCapacity
	_, err := bf.reserve(ctx, key, errorRate, capacity)
	return err
}

// AddURLBase62ID add base62 ID of url to bloom filter
func (bf *URLShortenerBloomFilter) AddURLBase62ID(ctx context.Context, u model.URL) error {
	key := genBase62IDKey()
	return bf.add(ctx, key, u.GetIDBase62())
}

// IsURLBase62IDExist check if base62 ID in bloom filter
func (bf *URLShortenerBloomFilter) IsURLBase62IDExist(ctx context.Context, u model.URL) (bool, error) {
	key := genBase62IDKey()
	return bf.exists(ctx, key, u.GetIDBase62())
}

// https://redis.io/docs/latest/commands/bf.reserve/
func (bf *URLShortenerBloomFilter) reserve(
	ctx context.Context, key string, errorRate float64, capacity int64,
) (bool, error) {
	err := bf.filter.RDB.BFReserve(ctx, key, errorRate, capacity).Err()

	if err != nil && err.Error() != "ERR item exists" {
		return false, err
	}

	if err != nil {
		return false, nil
	}

	return true, nil
}

// https://redis.io/docs/latest/commands/bf.add/
func (bf *URLShortenerBloomFilter) add(
	ctx context.Context, key string, element any,
) error {
	return bf.filter.RDB.BFAdd(ctx, key, element).Err()
}

// https://redis.io/docs/latest/commands/bf.exists/
func (bf *URLShortenerBloomFilter) exists(
	ctx context.Context, key string, element any,
) (bool, error) {
	return bf.filter.RDB.BFExists(ctx, key, element).Result()
}
