// Copyright 2020 the Exposure Notifications Server authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	"go.uber.org/zap"

	"github.com/TinyMurky/tinyurl/internal/setup"
	"github.com/TinyMurky/tinyurl/pkg/database"
	"github.com/TinyMurky/tinyurl/pkg/logging"
)

var (
	pathFlag = flag.String("path", "migrations", "path to migrations folder")
)

func main() {
	flag.Parse()

	ctx, done := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	loadDotEnvIfNotLoaded()

	logger := logging.NewLoggerFromEnv()

	ctx = logging.WithLogger(ctx, logger)

	defer func() {
		done()
		if r := recover(); r != nil {
			logger.Fatalw("migration panic", "panic", r)
		}
	}()

	err := realMain(ctx)

	done()
	if err != nil {
		log.Fatalf("migration failed: %s", err.Error())
	}

	logger.Info("migration complete succeessfully.")
}

func realMain(ctx context.Context) error {
	logger := logging.FromContext(ctx)

	var config database.Config

	env, err := setup.Setup(ctx, &config)

	if err != nil {
		return fmt.Errorf("setup.Setup: %w", err)
	}
	defer env.Close(ctx)

	db := env.Database()
	defer db.Close(ctx)

	// ./migrate.out -path="/path/to/migrations" (need / at the front of path)
	migrateSQLDir := fmt.Sprintf("file://%s", *pathFlag)

	m, err := migrate.New(migrateSQLDir, config.ToSQLiteDSN())

	if err != nil {
		return fmt.Errorf("failed create migrate: %w, migration path:%q", err, migrateSQLDir)
	}

	m.Log = newLogger(logger)

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed run migrate: %w", err)
	}
	srcErr, dbErr := m.Close()
	if srcErr != nil {
		return fmt.Errorf("migrate source error: %w", srcErr)
	}
	if dbErr != nil {
		return fmt.Errorf("migrate database error: %w", dbErr)
	}

	logger.Debugw("finished running migrations")
	return nil
}

func loadDotEnvIfNotLoaded() {
	mode := strings.TrimSpace(strings.ToLower(os.Getenv("RUN_MODE")))
	isEnvLoaded := mode != ""

	if !isEnvLoaded {
		// it will be where the binary is located
		exePath, err := os.Executable()
		if err != nil {
			panic(err)
		}

		exeDir := filepath.Dir(exePath)

		envPath := filepath.Join(exeDir, "../../.env")
		// load from .env
		if err := godotenv.Load(envPath); err != nil {
			panicMsg := fmt.Sprintf("Warning: failed to load .env file from path %q: %v\n", envPath, err)
			panic(panicMsg)
		}
	}
}

type logger struct {
	logger *zap.SugaredLogger
}

func newLogger(zapLogger *zap.SugaredLogger) migrate.Logger {
	return &logger{
		logger: zapLogger,
	}
}

func (l *logger) Printf(format string, v ...any) {
	l.logger.Infof(format, v...)
}

func (l *logger) Verbose() bool {
	return true
}
