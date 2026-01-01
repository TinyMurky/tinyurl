// Package cache can set and get of url
package cache

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/TinyMurky/tinyurl/internal/urlshortener/model"
	"github.com/TinyMurky/tinyurl/pkg/cache"
)

type URLShortenerCache struct {
	cache *cache.Cache
}

func New(c *cache.Cache) *URLShortenerCache {
	return &URLShortenerCache{
		cache: c,
	}
}

// SetLongURL set longURL from model.URL into
func (uc *URLShortenerCache) SetLongURL(
	ctx context.Context,
	u model.URL,
	expiration time.Duration,
) error {
	// check validation
	if u.IsZero() || u.ID == 0 {
		return errors.New("SetLongURL: invalid model.URL or ID")
	}

	if u.LongURL == "" {
		return errors.New("SetLongURL: longURL is empty")
	}

	key := genURLKey(u)
	return uc.set(ctx, key, u.LongURL, expiration)
}

// GetLongURL get url model (with longURL) from cache
func (uc *URLShortenerCache) GetLongURL(
	ctx context.Context,
	u model.URL,
) (model.URL, error) {
	if u.IsZero() || u.ID == 0 {
		return u, errors.New("GetLongURL: invalid model.URL or ID")
	}

	key := genURLKey(u)
	longURL, err := uc.get(ctx, key)

	if err != nil {
		if errors.Is(err, redis.Nil) {
			return u, redis.Nil
		}
		return u, fmt.Errorf("GetLongURL failed for ID %s: %w", u.GetIDBase62(), err)
	}

	u.LongURL = longURL
	return u, nil
}

func (uc *URLShortenerCache) set(
	ctx context.Context,
	key string,
	value any,
	expiration time.Duration,
) error {
	return uc.cache.RDB.Set(
		ctx, key, value, expiration,
	).Err()
}

func (uc *URLShortenerCache) get(ctx context.Context, key string) (string, error) {
	return uc.cache.RDB.Get(ctx, key).Result()
}
