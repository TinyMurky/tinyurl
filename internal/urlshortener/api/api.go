// Package api provides the HTTP handlers for the URL shortener service.
// It defines the routing logic and endpoints for the API layer.
package api

import (
	"net/http"

	"github.com/TinyMurky/tinyurl/internal/serverenv"
	v1 "github.com/TinyMurky/tinyurl/internal/urlshortener/api/v1"
	urlshortenerconfig "github.com/TinyMurky/tinyurl/internal/urlshortener/config"
)

// Handler encapsulates the dependencies required for handling API requests.
// It holds references to the configuration and server environment.
type Handler struct {
	config *urlshortenerconfig.Config
	env    *serverenv.ServerEnv
}

// NewAPIHandler creates and returns a new instance of APIHandler with the provided
// configuration and server environment.
func NewAPIHandler(cfg *urlshortenerconfig.Config, env *serverenv.ServerEnv) *Handler {
	return &Handler{
		config: cfg,
		env:    env,
	}
}

// Handler constructs and returns an http.Handler with all API routes registered.
// This handler serves as the entry point for API traffic and can be mounted
// onto a parent router.
func (a *Handler) Handler() http.Handler {
	router := http.NewServeMux()

	v1Router := v1.NewV1Handler(a.config, a.env)

	router.Handle("/v1/", http.StripPrefix("/v1", v1Router.Handler()))

	return router
}
