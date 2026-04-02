// Package app provides application lifecycle management
package app

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/liang21/aitestos/internal/config"
)

// HTTPServer interface for HTTP server operations
type HTTPServer interface {
	ListenAndServe() error
	Shutdown(ctx context.Context) error
	Addr() string
}

// App is the main application container
type App struct {
	config      *config.Config
	httpServer  HTTPServer
	shutdownMgr *ShutdownManager
}

// New creates a new application instance
func New(cfg *config.Config, httpServer HTTPServer, shutdownMgr *ShutdownManager) *App {
	return &App{
		config:      cfg,
		httpServer:  httpServer,
		shutdownMgr: shutdownMgr,
	}
}

// Run starts the HTTP server
func (a *App) Run() error {
	addr := a.httpServer.Addr()
	log.Info().Str("addr", addr).Msg("starting HTTP server")

	if err := a.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

// Shutdown gracefully stops the application
func (a *App) Shutdown(ctx context.Context) error {
	log.Info().Msg("shutting down application")

	// Shutdown HTTP server first
	if err := a.httpServer.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("HTTP server shutdown error")
		return err
	}

	// Shutdown other registered components
	if a.shutdownMgr != nil {
		if err := a.shutdownMgr.Shutdown(ctx); err != nil {
			log.Error().Err(err).Msg("component shutdown error")
		}
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
	mu      sync.Mutex
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
	m.mu.Lock()
	defer m.mu.Unlock()
	m.closers = append(m.closers, closer)
}

// Shutdown executes shutdown in reverse registration order
func (m *ShutdownManager) Shutdown(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, m.timeout)
	defer cancel()

	m.mu.Lock()
	closers := make([]Closer, len(m.closers))
	copy(closers, m.closers)
	m.mu.Unlock()

	for i := len(closers) - 1; i >= 0; i-- {
		closer := closers[i]
		log.Info().Str("component", closer.Name()).Msg("shutting down component")

		if err := closer.Close(ctx); err != nil {
			log.Error().Err(err).Str("component", closer.Name()).Msg("shutdown error")
		}
	}

	return nil
}
