// Package handlepostdatashorten will shorten the url with snowflakeID
// original longer url
package handlepostdatashorten

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"go.uber.org/zap"

	"github.com/TinyMurky/tinyurl/internal/cache"
	"github.com/TinyMurky/tinyurl/internal/serverenv"
	urlshortenerconfig "github.com/TinyMurky/tinyurl/internal/urlshortener/config"
	"github.com/TinyMurky/tinyurl/internal/urlshortener/database"
	idgenerator "github.com/TinyMurky/tinyurl/internal/urlshortener/id_generator"
	"github.com/TinyMurky/tinyurl/internal/urlshortener/model"
	"github.com/TinyMurky/tinyurl/pkg/logging"
)

type response struct {
	Success  bool   `json:"success"`
	Message  string `json:"message,omitempty"`
	ShortURL string `json:"short_url,omitempty"`
}

// Handler encapsulates the dependencies required for handling V1 version of
// looking up original URL from id provided
// It holds references to the configuration and server environment.
type Handler struct {
	config      *urlshortenerconfig.Config
	env         *serverenv.ServerEnv
	cache       *cache.URLShortenerCache
	db          *database.URLShortenerDB
	idGenerator *idgenerator.Generator
}

var _ http.Handler = (*Handler)(nil)

// New will return http.Handler that can
// get snowflake ID and return original longer url
func New(cfg *urlshortenerconfig.Config, env *serverenv.ServerEnv) *Handler {
	cache := cache.New(env.Cache())
	db := database.New(env.Database())
	idGenerator, err := idgenerator.NewGenerator(cfg)

	if err != nil {
		log.Fatalf("New handle_post_data_shorten new id generator: %s", err.Error())
	}

	return &Handler{
		config:      cfg,
		env:         env,
		cache:       cache,
		db:          db,
		idGenerator: idGenerator,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logging.FromContext(ctx).Named("handel_post_data_shorten")

	var contentType = "application/x-www-form-urlencoded"

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if r.Header.Get("Content-Type") != contentType {
		msg := fmt.Sprintf("Content-Type need to be %s", contentType)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	// parse form will parse query and form
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// PostFormValue only got body, use FormValue to get all
	longURL := r.PostFormValue("long_url")

	if longURL == "" {
		msg := "long_url is required"
		sendBadRequest(w, msg, logger)
		return
	}

	if !isValidURL(longURL) {
		msg := fmt.Sprintf("long_url %q is invalid", longURL)
		sendBadRequest(w, msg, logger)
		return
	}

	u := model.URL{
		LongURL: longURL,
	}

	u, err := h.createURL(ctx, u)

	if err != nil {
		msg := fmt.Sprintf("create url error: %s", err.Error())
		sendInternalError(w, msg, logger)
		return
	}

	shortURL, err := h.genTinyURL(u)

	if err != nil {
		msg := fmt.Sprintf("gen tiny url error: %s", err.Error())
		sendInternalError(w, msg, logger)
		return
	}

	res := response{
		Success:  true,
		ShortURL: shortURL,
	}

	sendJSONResponse(w, http.StatusOK, res, logger)

	logger.Debug("method", r.Method, "response", res)
}

func (h *Handler) createURL(ctx context.Context, urlModel model.URL) (model.URL, error) {
	cacheTTL := time.Millisecond * time.Duration(h.config.RedisCacheTTLInMiliSec)

	if urlModel.LongURL == "" {
		return model.URL{}, errors.New("longURL is empty")
	}

	dbURLModel, err := h.db.GetFirstByLongURL(ctx, urlModel.LongURL)

	if err != nil {
		return model.URL{}, fmt.Errorf("database GetFirstByLongURL: %w", err)
	}

	// If exist just return
	if !dbURLModel.IsZero() {
		if err := h.cache.SetLongURL(ctx, dbURLModel, cacheTTL); err != nil {
			return model.URL{}, fmt.Errorf("cache SetLongURL: %w", err)
		}
		return dbURLModel, nil
	}

	newID, err := h.idGenerator.NextID()

	if err != nil {
		return model.URL{}, fmt.Errorf("idGenerator nextID: %w", err)
	}

	urlModel.ID = newID

	if err := h.db.CreateURL(ctx, urlModel); err != nil {
		return model.URL{}, fmt.Errorf("database CreateURL: %w", err)
	}

	if err := h.cache.SetLongURL(ctx, urlModel, cacheTTL); err != nil {
		return model.URL{}, fmt.Errorf("cache SetLongURL: %w", err)
	}

	return urlModel, nil
}

func (h *Handler) genTinyURL(u model.URL) (string, error) {
	urlPath := h.config.ShortURLPrefix
	if !isValidURL(urlPath) {
		return "", fmt.Errorf("config.ShortURLPrefix %q is not valid", urlPath)
	}

	shortURL, err := url.JoinPath(urlPath, "api", "v1", "shortUrl", u.GetIDBase62())

	if err != nil {
		return "", fmt.Errorf("url join path err: %w", err)
	}

	return shortURL, nil
}

func sendBadRequest(w http.ResponseWriter, msg string, logger *zap.SugaredLogger) {
	res := response{
		Success: false,
		Message: msg,
	}
	sendJSONResponse(w, http.StatusBadRequest, msg, logger)
	logger.Debug("response", res)
}

func sendInternalError(w http.ResponseWriter, msg string, logger *zap.SugaredLogger) {
	res := response{
		Success: false,
		Message: msg,
	}
	sendJSONResponse(w, http.StatusInternalServerError, msg, logger)
	logger.Debug("response", res)
}

func sendJSONResponse(w http.ResponseWriter, status int, data any, logger *zap.SugaredLogger) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		logger.Errorf("JSON encode err: %s", err.Error())
	}
}

func isValidURL(u string) bool {
	// need to start with http...
	_, err := url.ParseRequestURI(u)

	return err == nil
}
