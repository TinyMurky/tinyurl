// Package model is a model abstraction of url shorterner
package model

import (
	"fmt"
	"time"

	"github.com/TinyMurky/snowflake"
)

// URL represents the mapping between a SnowflakeID (ShortURL) and a LongURL.
type URL struct {
	ID        snowflake.SID `json:"id"`
	LongURL   string        `json:"long_url"`
	CreatedAt time.Time     `json:"created_at"`
}

// GetIDBase62 returns the snowflake id in base62 format.
func (u *URL) GetIDBase62() string {
	return u.ID.Base62()
}

// IsZero will return that if URL is zero value
func (u *URL) IsZero() bool {
	if u == nil {
		return false
	}
	isIDZero := u.ID == 0
	isLongURLZero := u.LongURL == ""
	isCreatedAtZero := u.CreatedAt.IsZero()

	return isIDZero && isLongURLZero && isCreatedAtZero
}

// NewURL create a new URL item
func NewURL(id snowflake.SID, longURL string) *URL {
	return &URL{
		ID:        id,
		LongURL:   longURL,
		CreatedAt: time.Now(),
	}
}

// NewURLFromBase62 can create URL from base62 for search query
func NewURLFromBase62(base62 string) (*URL, error) {
	sid, err := snowflake.ParseBase62(base62)
	if err != nil {
		return nil, fmt.Errorf("invalid base62 string: %w", err)
	}

	return &URL{
		ID: sid,
	}, nil
}
