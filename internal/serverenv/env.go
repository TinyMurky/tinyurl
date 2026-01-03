// Copyright 2020 the Exposure Notifications Server authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package serverenv defines common parameters for the sever environment.
package serverenv

import (
	"context"

	"github.com/TinyMurky/tinyurl/pkg/bloomfilter"
	"github.com/TinyMurky/tinyurl/pkg/cache"
	"github.com/TinyMurky/tinyurl/pkg/database"
	"github.com/TinyMurky/tinyurl/pkg/singleflight"
)

// ServerEnv represents latent environment configuration for servers in this application.
type ServerEnv struct {
	database     *database.DB
	cache        *cache.Cache
	bloomFilter  *bloomfilter.BloomFilter
	singleFlight singleflight.Group
}

// Option defines function types to modify the ServerEnv on creation.
type Option func(*ServerEnv) *ServerEnv

// New creates a new ServerEnv with the requested options.
func New(_ context.Context, opts ...Option) *ServerEnv {
	env := new(ServerEnv)

	for _, f := range opts {
		env = f(env)
	}

	return env
}

// Close shuts down the server env, closing database connections, etc.
func (s *ServerEnv) Close(ctx context.Context) error {
	if s == nil {
		return nil
	}

	if s.database != nil {
		s.database.Close(ctx)
	}

	return nil
}

// WithDatabase add database to serverEnv
func WithDatabase(db *database.DB) Option {
	return func(s *ServerEnv) *ServerEnv {
		s.database = db
		return s
	}
}

// Database get database
func (s *ServerEnv) Database() *database.DB {
	return s.database
}

// WithCache add cache to serverEnv
func WithCache(c *cache.Cache) Option {
	return func(s *ServerEnv) *ServerEnv {
		s.cache = c
		return s
	}
}

// Cache get cache
func (s *ServerEnv) Cache() *cache.Cache {
	return s.cache
}

// WithBloomFilter add bloom filter to serverEnv
func WithBloomFilter(bf *bloomfilter.BloomFilter) Option {
	return func(s *ServerEnv) *ServerEnv {
		s.bloomFilter = bf
		return s
	}
}

// BloomFilter get Bloom Filter
func (s *ServerEnv) BloomFilter() *bloomfilter.BloomFilter {
	return s.bloomFilter
}

// WithSingleFlight add single flight to serverEnv
func WithSingleFlight(sf singleflight.Group) Option {
	return func(s *ServerEnv) *ServerEnv {
		s.singleFlight = sf
		return s
	}
}

// SingleFlight get Single Flight
func (s *ServerEnv) SingleFlight() singleflight.Group {
	return s.singleFlight
}
