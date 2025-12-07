// Deprecated
package setup

import (
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
)

// Logger defines the minimal interface required for logging within the lookuper.
// This allows for dependency injection of any compatible logging implementation.
type Logger interface {
	Warnf(template string, args ...any)
}

// osLookuper implements the envconfig.Lookuper interface.
// It wraps the standard OS lookuper with additional logic to load .env files
// when running in development mode.
type osLookuper struct {
	loadDotEnvOnce sync.Once
	delegate       envconfig.Lookuper
	logger         Logger
}

// Ensure osLookuper implements the envconfig.Lookuper interface at compile time.
var _ envconfig.Lookuper = (*osLookuper)(nil)

// NewOSLookuper creates and returns a new instance of envconfig.Lookuper.
// It initializes the underlying OS lookuper and injects the provided logger.
func NewOSLookuper(logger Logger) envconfig.Lookuper {
	return &osLookuper{
		delegate: envconfig.OsLookuper(),
		logger:   logger,
	}
}

// loadDotEnv checks the "RUN_MODE" environment variable.
// If the mode is set to "development", it attempts to load environment variables
// from a .env file using godotenv. This operation is thread-safe and performed only once.
func (o *osLookuper) loadDotEnv() {
	o.loadDotEnvOnce.Do(func() {
		mode := strings.TrimSpace(strings.ToLower(os.Getenv("RUN_MODE")))
		isDevelopment := mode == "development" || mode == ""

		if isDevelopment {
			// it will be where the binary is located
			exePath, err := os.Executable()
			if err != nil {
				panic(err)
			}

			exeDir := filepath.Dir(exePath)

			envPath := filepath.Join(exeDir, "../../.env")
			// load from .env
			if err := godotenv.Load(envPath); err != nil {
				o.logger.Warnf("Warning: failed to load .env file from path %q: %v\n", envPath, err)
			}
		}
	})
}

// Lookup satisfies the envconfig.Lookuper interface.
// It ensures that the .env file is loaded (if in development mode) before
// delegating the lookup to the standard OS lookuper.
func (o *osLookuper) Lookup(key string) (string, bool) {
	o.loadDotEnv()
	return o.delegate.Lookup(key)
}
