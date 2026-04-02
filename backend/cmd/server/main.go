// Package main is the entry point of the aitestos server
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/edy/zane/go/aitestos/internal/app"
	"github.com/edy/zane/go/aitestos/internal/config"
)

func main() {
	// Initialize logger
	zerolog.TimeFieldFormat = time.RFC3339
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

	// Load configuration
	cfg, err := config.Load("configs/config.yaml")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	// Create application
	application, cleanup, err := app.Initialize(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to initialize app")
	}
	defer cleanup()

	// Start server in goroutine
	errCh := make(chan error, 1)
	go func() {
		log.Info().Str("addr", fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)).Msg("starting server")
		if err := application.Run(); err != nil && err != http.ErrServerClosed {
			errCh <- fmt.Errorf("server error: %w", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errCh:
		log.Error().Err(err).Msg("server error")
	case sig := <-quit:
		log.Info().Str("signal", sig.String()).Msg("shutting down server")
	}

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := application.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("shutdown error")
	}

	log.Info().Msg("server stopped")
}
