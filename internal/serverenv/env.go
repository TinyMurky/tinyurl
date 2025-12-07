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

import "context"

// ServerEnv represents latent environment configuration for servers in this application.
type ServerEnv struct{}

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

	return nil
}
