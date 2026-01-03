package singleflight

import (
	"fmt"

	"github.com/TinyMurky/tinyurl/pkg/singleflight"
)

// Group wraps pkg/singleflight.Group to provide domain-specific key namespacing.
type Group struct {
	group singleflight.Group
}

// New creates a new domain-specific singleflight group.
func New(g singleflight.Group) *Group {
	return &Group{group: g}
}

// Do executes the function with the key prefixed by "shorturl:".
// This ensures that the key is namespaced correctly for the URL shortener domain.
func (g *Group) Do(id string, fn func() (any, error)) (v any, err error, shared bool) {
	key := fmt.Sprintf("shorturl:%s", id)
	return g.group.Do(key, fn)
}
