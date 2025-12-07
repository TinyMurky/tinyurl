// Package urlshortener provide http handler that provide api to:
//  1. create and store shorten url by providing longer url
//  2. transfer shorten url to longer url
//
// Package urlshortener provides the main entry point for the URL shortening service.
// It acts as the composer, wiring together the configuration, environment, and
// API handlers to provide a complete HTTP server implementation.
package urlshortener

import (
	"context"
	"net/http"

	"github.com/TinyMurky/tinyurl/internal/middleware"
	"github.com/TinyMurky/tinyurl/internal/serverenv"
	"github.com/TinyMurky/tinyurl/internal/urlshortener/api"
	urlshortenerconfig "github.com/TinyMurky/tinyurl/internal/urlshortener/config"
	"github.com/TinyMurky/tinyurl/pkg/logging"
)

// Server represents the HTTP server for the URL shortener application.
// It holds the necessary configuration and global server environment dependencies.
type Server struct {
	config *urlshortenerconfig.Config
	env    *serverenv.ServerEnv
}

// NewServer creates and returns a new Server instance.
func NewServer(cfg *urlshortenerconfig.Config, env *serverenv.ServerEnv) *Server {
	return &Server{
		config: cfg,
		env:    env,
	}
}

// Routes initializes the routing logic and registers all application endpoints.
// It returns the top-level http.Handler that can be used by the HTTP server.
func (s *Server) Routes(ctx context.Context) http.Handler {
	logger := logging.FromContext(ctx).Named("urlshortener")

	router := http.NewServeMux()
	// Initialize the API handler with dependencies
	apiHandler := api.NewAPIHandler(s.config, s.env)

	// Mount the API handler under the "/api/" path.
	// We use StripPrefix so the inner handler doesn't need to know about the "/api" prefix.
	// Note: The trailing slash in "/api/" ensures it matches all paths under /api.
	router.Handle("/api/", http.StripPrefix("/api", apiHandler.Handler()))

	// Wrap router with middlewares
	// request will perform middleware before it enter the route
	middlewareStack := middleware.CreateStack(
		middleware.PopulateLogger(logger),
	)

	return middlewareStack(router)
}
