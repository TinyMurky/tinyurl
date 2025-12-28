// Package urlshortenerconfig define the config of urlshortener
package urlshortenerconfig

import "github.com/TinyMurky/tinyurl/pkg/database"

// Config for urlshortener
// it need to use github.com/sethvargo/go-envconfig package to read
type Config struct {
	Database database.Config

	IDGenerator    IDGeneratorConfig
	Port           string `env:"PORT"`
	ShortURLPrefix string `env:"SHORT_URL_PREFIX, default=http://localhost:3000"`
}

// DatabaseConfig return Database config
func (c *Config) DatabaseConfig() *database.Config {
	return &c.Database
}
