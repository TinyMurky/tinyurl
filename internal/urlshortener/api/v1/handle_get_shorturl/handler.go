// Package handlegetshorturl will get snowflake ID and return
// original longer url
package handlegetshorturl

import (
	"context"
	"fmt"
	"net/http"

	"github.com/TinyMurky/tinyurl/internal/serverenv"
	urlshortenerconfig "github.com/TinyMurky/tinyurl/internal/urlshortener/config"
	"github.com/TinyMurky/tinyurl/internal/urlshortener/database"
	"github.com/TinyMurky/tinyurl/internal/urlshortener/model"
	"github.com/TinyMurky/tinyurl/pkg/logging"
)

// Handler encapsulates the dependencies required for handling V1 version of
// looking up original URL from id provided
// It holds references to the configuration and server environment.
type Handler struct {
	config *urlshortenerconfig.Config
	env    *serverenv.ServerEnv
	db     *database.URLShortenerDB
}

var _ http.Handler = (*Handler)(nil)

// New will return http.Handler that can
// get snowflake ID and return original longer url
func New(cfg *urlshortenerconfig.Config, env *serverenv.ServerEnv) *Handler {
	db := database.New(env.Database())

	return &Handler{
		config: cfg,
		env:    env,
		db:     db,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logging.FromContext(ctx).Named("handel_get_shorturl")

	if r.Method != http.MethodGet {
		http.Error(w, "Method Not allow", http.StatusMethodNotAllowed)
		return
	}

	id := r.PathValue("id")

	if len(id) == 0 {
		http.Error(w, "Bad Request: id not provided", http.StatusBadRequest)
		return
	}

	url, err := model.NewURLFromBase62(id)

	if err != nil {
		http.Error(w, "invalid idL not base62", http.StatusBadRequest)
		return
	}

	longURL, err := h.findURL(ctx, url)

	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		logger.Errorf("findURL: %s", err.Error())
		return
	}

	http.Redirect(w, r, longURL, http.StatusMovedPermanently)
	logger.Debug("method=", r.Method, "id=", id, "tinyURL=", longURL)
}

func (h *Handler) findURL(ctx context.Context, url *model.URL) (string, error) {
	urlFromDB, err := h.db.GetFirstByID(ctx, url.ID)

	if err != nil {
		return "", fmt.Errorf("findURL: %w", err)
	}

	if urlFromDB.IsZero() {
		return h.genDomain(), nil
	}

	return urlFromDB.LongURL, nil
}

func (h *Handler) genDomain() string {
	return h.config.ShortURLPrefix
}
