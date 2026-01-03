package singleflight

import (
	"context"
)

// Config holds configuration for the singleflight package.
type Config struct {
	// Add fields here if we need to configure timeouts, metrics, etc.
	// For now, it's a placeholder for future extensibility as per design.
}

// NewFromEnv creates a new Config from environment variables.
func NewFromEnv(_ context.Context) *Config {
	return &Config{}
}

// Group is the interface that wraps the singleflight Do method.
type Group interface {
	Do(key string, fn func() (any, error)) (v any, err error, shared bool)
}

// MockGroup can be used for testing.
type MockGroup struct {
	DoFunc func(key string, fn func() (any, error)) (v any, err error, shared bool)
}

func (m *MockGroup) Do(key string, fn func() (any, error)) (v any, err error, shared bool) {
	if m.DoFunc != nil {
		return m.DoFunc(key, fn)
	}
	// Default behavior: just run it, not shared
	v, err = fn()
	return v, err, false
}

// ConfigProvider ensures that the environment config can provide a singleflight config.
type ConfigProvider interface {
	SingleFlightConfig() *Config
}