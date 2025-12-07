// Package urlshortenerconfig define the config of urlshortener
package urlshortenerconfig

// Config for urlshortener
// it need to use github.com/sethvargo/go-envconfig package to read
type Config struct {
	Port string `env:"PORT"`
}
