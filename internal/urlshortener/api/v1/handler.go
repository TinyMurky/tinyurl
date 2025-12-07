// Package v1 provides the HTTP handlers for the V1 version of the URL shortener service.
// It defines the routing logic and endpoints for the /api/v1 layer.
package v1

import (
	"net/http"

	"github.com/TinyMurky/tinyurl/internal/serverenv"
	handlegetshorturl "github.com/TinyMurky/tinyurl/internal/urlshortener/api/v1/handle_get_shorturl"
	urlshortenerconfig "github.com/TinyMurky/tinyurl/internal/urlshortener/config"
)

// Handler encapsulates the dependencies required for handling V1 version of the URL shortener requests.
// It holds references to the configuration and server environment.
type Handler struct {
	config *urlshortenerconfig.Config
	env    *serverenv.ServerEnv
}

// NewV1Handler creates and returns a new instance of Handler with the provided
// configuration and server environment.
func NewV1Handler(cfg *urlshortenerconfig.Config, env *serverenv.ServerEnv) *Handler {
	return &Handler{
		config: cfg,
		env:    env,
	}
}

// Handler constructs and returns an http.Handler with all V1 routes registered.
// This handler serves as the entry point for V1 traffic and can be mounted
// onto a parent router.
func (a *Handler) Handler() http.Handler {
	mux := http.NewServeMux()

	getShortURLHandler := handlegetshorturl.New(a.config, a.env)
	mux.Handle("GET /shortUrl/{id}", getShortURLHandler)

	return mux
}
