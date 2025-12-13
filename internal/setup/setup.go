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

// Package setup provides common logic for configuring the various services.
package setup

import (
	"context"
	"fmt"

	"github.com/sethvargo/go-envconfig"

	"github.com/TinyMurky/tinyurl/internal/serverenv"
	"github.com/TinyMurky/tinyurl/pkg/database"
	"github.com/TinyMurky/tinyurl/pkg/logging"
)

// DatabaseConfigProvider ensures that the environment config can provide a DB config.
// All binaries in this application connect to the database via the same method.
type DatabaseConfigProvider interface {
	DatabaseConfig() *database.Config
}

// Setup runs common initialization code for all servers. See SetupWith.
func Setup(ctx context.Context, config any) (*serverenv.ServerEnv, error) {
	//logger := logging.FromContext(ctx)
	//osLookuper := NewOSLookuper(logger)
	//return SetupWith(ctx, config, osLookuper)
	return SetupWith(ctx, config, envconfig.OsLookuper())
}

// SetupWith processes the given configuration using envconfig. It is
// responsible for establishing database connections and
// accessing app configs. The provided interface must implement the various
// interfaces.
func SetupWith(ctx context.Context, config any, l envconfig.Lookuper) (*serverenv.ServerEnv, error) {
	logger := logging.FromContext(ctx)

	// Build a list of mutators (can use envconfig.MutatorFunc).
	// This list will grow as we initialize more of the
	// configuration, such as the secret manager.
	var mutatorFuncs []envconfig.Mutator

	// Build a list of options to pass to the server env.
	var serverEnvOpts []serverenv.Option

	if err := envconfig.ProcessWith(ctx, &envconfig.Config{
		Target:   config,
		Lookuper: l,
		Mutators: mutatorFuncs,
	}); err != nil {
		return nil, fmt.Errorf("error loading environment variables: %w", err)
	}

	if provider, ok := config.(DatabaseConfigProvider); ok {
		logger.Info("configuring database")
		dbConfig := provider.DatabaseConfig()

		db, err := database.NewFromEnv(ctx, dbConfig)

		if err != nil {
			return nil, fmt.Errorf("unable to connect to database: %w", err)
		}

		serverEnvOpt := serverenv.WithDatabase(db)
		serverEnvOpts = append(serverEnvOpts, serverEnvOpt)
	}

	return serverenv.New(ctx, serverEnvOpts...), nil
}
