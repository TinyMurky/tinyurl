// Package urlshortenerconfig define the config of urlshortener
package urlshortenerconfig

import "github.com/TinyMurky/tinyurl/pkg/database"

// Config for urlshortener
// it need to use github.com/sethvargo/go-envconfig package to read
type Config struct {
	Database database.Config

	Port string `env:"PORT"`
}

// DatabaseConfig return Database config
func (c *Config) DatabaseConfig() *database.Config {
	return &c.Database
}
