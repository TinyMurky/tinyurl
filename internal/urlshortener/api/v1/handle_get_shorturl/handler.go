// Package handlegetshorturl will get snowflake ID and return
// original longer url
package handlegetshorturl

import (
	"fmt"
	"net/http"

	"github.com/TinyMurky/tinyurl/internal/serverenv"
	urlshortenerconfig "github.com/TinyMurky/tinyurl/internal/urlshortener/config"
	"github.com/TinyMurky/tinyurl/pkg/logging"
)

// Handler encapsulates the dependencies required for handling V1 version of
// looking up original URL from id provided
// It holds references to the configuration and server environment.
type Handler struct {
	config *urlshortenerconfig.Config
	env    *serverenv.ServerEnv
}

var _ http.Handler = (*Handler)(nil)

// New will return http.Handler that can
// get snowflake ID and return original longer url
func New(cfg *urlshortenerconfig.Config, env *serverenv.ServerEnv) *Handler {
	return &Handler{
		config: cfg,
		env:    env,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logging.FromContext(ctx).Named("handel_get_shorturl")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	id := r.PathValue("id")

	if len(id) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	format := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Document</title>
</head>
<body>
    <h1> %s </h1>
</body>
</html>
`

	responseHTML := fmt.Sprintf(format, id)

	w.Header().Set("Content-Type", "text/html; charset=utf-8") // normal header
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, responseHTML)
	logger.Debug("method", r.Method, "id", id)
}
