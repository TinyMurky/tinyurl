package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/joho/godotenv"

	"github.com/TinyMurky/tinyurl/internal/setup"
	"github.com/TinyMurky/tinyurl/internal/urlshortener"
	urlshortenerconfig "github.com/TinyMurky/tinyurl/internal/urlshortener/config"
	"github.com/TinyMurky/tinyurl/pkg/logging"
	"github.com/TinyMurky/tinyurl/pkg/server"
)

func main() {
	ctx, done := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	loadDotEnvIfNotLoaded()

	logger := logging.NewLoggerFromEnv()
	ctx = logging.WithLogger(ctx, logger)

	defer func() {
		done()
		if r := recover(); r != nil {
			logger.Fatalw("application panic", "panic", r)
		}
	}()

	err := realMain(ctx)
	done()

	if err != nil {
		logger.Fatal(err)
	}
	logger.Info("successful shutdown")
}

func realMain(ctx context.Context) error {
	logger := logging.FromContext(ctx)

	var config urlshortenerconfig.Config

	serverEnv, err := setup.Setup(ctx, &config)

	if err != nil {
		return fmt.Errorf("setup.Setup: %w", serverEnv)
	}
	defer serverEnv.Close(ctx)

	urlShortenerServer := urlshortener.NewServer(&config, serverEnv)

	srv, err := server.New(config.Port)
	if err != nil {
		return fmt.Errorf("server.New: %w", err)
	}

	logger.Infof("listening on :%s", config.Port)

	return srv.ServeHTTPHandler(ctx, urlShortenerServer.Routes(ctx))
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
