// Package app provides application lifecycle management
package app

import (
	"context"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/edy/zane/go/aitestos/internal/config"
)

// App is the main application container
type App struct {
	config     *config.Config
	httpServer *http.Server
}

// New creates a new application instance
func New(cfg *config.Config, httpServer *http.Server) *App {
	return &App{
		config:     cfg,
		httpServer: httpServer,
	}
}

// Run starts the HTTP server
func (a *App) Run() error {
	addr := a.httpServer.Addr
	log.Info().Str("addr", addr).Msg("starting HTTP server")

	if err := a.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

// Shutdown gracefully stops the application
func (a *App) Shutdown(ctx context.Context) error {
	log.Info().Msg("shutting down application")

	// Shutdown HTTP server
	if err := a.httpServer.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("HTTP server shutdown error")
		return err
	}

	log.Info().Msg("application shutdown complete")
	return nil
}

// Closer interface for shutdown management
type Closer interface {
	Name() string
	Close(ctx context.Context) error
}

// ShutdownManager manages graceful shutdown of multiple components
type ShutdownManager struct {
	closers []Closer
	timeout time.Duration
}

// NewShutdownManager creates a new shutdown manager
func NewShutdownManager(timeout time.Duration) *ShutdownManager {
	return &ShutdownManager{
		closers: make([]Closer, 0),
		timeout: timeout,
	}
}

// Register adds a closer to be managed
func (m *ShutdownManager) Register(closer Closer) {
	m.closers = append(m.closers, closer)
}

// Shutdown executes shutdown in reverse registration order
func (m *ShutdownManager) Shutdown(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, m.timeout)
	defer cancel()

	for i := len(m.closers) - 1; i >= 0; i-- {
		closer := m.closers[i]
		log.Info().Str("component", closer.Name()).Msg("shutting down component")

		if err := closer.Close(ctx); err != nil {
			log.Error().Err(err).Str("component", closer.Name()).Msg("shutdown error")
		}
	}

	return nil
}
