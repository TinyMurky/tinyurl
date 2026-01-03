// Package handlegetshorturl will get snowflake ID and return
// original longer url
package handlegetshorturl

import (
	"errors"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/TinyMurky/tinyurl/internal/serverenv"
	"github.com/TinyMurky/tinyurl/internal/urlshortener/bloomfilter"
	"github.com/TinyMurky/tinyurl/internal/urlshortener/cache"
	urlshortenerconfig "github.com/TinyMurky/tinyurl/internal/urlshortener/config"
	"github.com/TinyMurky/tinyurl/internal/urlshortener/database"
	"github.com/TinyMurky/tinyurl/internal/urlshortener/model"
	"github.com/TinyMurky/tinyurl/internal/urlshortener/singleflight"
	"github.com/TinyMurky/tinyurl/pkg/logging"
)

// Handler encapsulates the dependencies required for handling V1 version of
// looking up original URL from id provided
// It holds references to the configuration and server environment.
type Handler struct {
	config       *urlshortenerconfig.Config
	env          *serverenv.ServerEnv
	cache        *cache.URLShortenerCache
	bloomFilter  *bloomfilter.URLShortenerBloomFilter
	db           *database.URLShortenerDB
	singleFlight *singleflight.Group
}

var _ http.Handler = (*Handler)(nil)

// New will return http.Handler that can
// get snowflake ID and return original longer url
func New(cfg *urlshortenerconfig.Config, env *serverenv.ServerEnv) *Handler {

	cache := cache.New(env.Cache())
	bloomFilter := bloomfilter.New(env.BloomFilter(), cfg.BloomFilterConfig())
	db := database.New(env.Database())
	sf := singleflight.New(env.SingleFlight())

	return &Handler{
		config:       cfg,
		env:          env,
		cache:        cache,
		db:           db,
		bloomFilter:  bloomFilter,
		singleFlight: sf,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logging.FromContext(ctx).Named("handel_get_shorturl")
	cacheTTL := time.Millisecond * time.Duration(h.config.RedisCacheTTLInMiliSec)

	if r.Method != http.MethodGet {
		http.Error(w, "Method Not allow", http.StatusMethodNotAllowed)
		return
	}

	id := r.PathValue("id")

	if len(id) == 0 {
		http.Error(w, "Bad Request: id not provided", http.StatusBadRequest)
		return
	}

	u, err := model.NewURLFromBase62(id)

	if err != nil {
		http.Error(w, "invalid id: not base62", http.StatusBadRequest)
		return
	}

	// check bloom filter first
	isURLExists, err := h.bloomFilter.IsURLBase62IDExist(ctx, u)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		logger.Errorf("bloom filter IsURLBase62IDExist: %s", err.Error())
		return
	}

	if !isURLExists {
		http.Error(w, "not found", http.StatusNotFound)
		logger.Errorf("ID not found: %s", u.GetIDBase62())
		return
	}

	u, err = h.cache.GetLongURL(ctx, u)

	if err != nil && !errors.Is(err, redis.Nil) {
		http.Error(w, "internal error", http.StatusInternalServerError)
		logger.Errorf("cache GetLongURL: %s", err.Error())
		return
	}

	if !u.IsEmptyLongURL() {
		// 找到 cache 的資料
		http.Redirect(w, r, u.LongURL, http.StatusMovedPermanently)
		return
	}

	v, err, _ := h.singleFlight.Do(u.GetIDBase62(), func() (any, error) {
		u, err := h.db.GetFirstByID(ctx, u.ID)
		if err != nil {
			return model.URL{}, err
		}

		if u.IsEmptyLongURL() {
			return u, nil
		}

		err = h.cache.SetLongURL(ctx, u, cacheTTL)
		if err != nil {
			// If cache fails, we still return the URL but log the error.
			// Returning error here would fail the request, which might be too harsh if DB worked.
			// But sticking to original logic, let's return error if SetLongURL fails?
			// The original logic returned 500 on SetLongURL error. Let's keep it consistency.
			return model.URL{}, err
		}

		return u, nil
	})

	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		logger.Errorf("singleFlight error: %s", err.Error())
		return
	}

	u = v.(model.URL)

	if u.IsEmptyLongURL() {
		http.Error(w, "not found", http.StatusNotFound)
		logger.Errorf("ID not found: %s", u.GetIDBase62())
		return
	}

	http.Redirect(w, r, u.LongURL, http.StatusMovedPermanently)
	logger.Debug("method=", r.Method, "id=", id, "tinyURL=", u.LongURL)
}
