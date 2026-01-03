package singleflight

import (
	"context"

	"golang.org/x/sync/singleflight"
)

// DefaultGroup is a wrapper around singleflight.Group.
// We redefine it here because it was referenced in config.go but that file was just config.
// Actually, let's keep the struct definition here.
type Wrapper struct {
	group singleflight.Group
}

// New creates a new singleflight group wrapper.
func New(_ context.Context, _ *Config) *Wrapper {
	return &Wrapper{}
}

// Do executes and returns the results of the given function, making sure that
// only one execution is in-flight for a given key at a time.
func (g *Wrapper) Do(key string, fn func() (any, error)) (v any, err error, shared bool) {
	return g.group.Do(key, fn)
}

// Ensure Wrapper implements Group
var _ Group = (*Wrapper)(nil)
