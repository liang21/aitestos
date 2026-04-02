// Package app provides application initialization
package app

import (
	"fmt"
	"net/http"
	"time"

	"github.com/liang21/aitestos/internal/config"
)

// Initialize creates a fully initialized application
// This is a manual implementation that can be replaced with Wire-generated code later
func Initialize(cfg *config.Config) (*App, func(), error) {
	// Create HTTP server
	httpServer := NewHTTPServerFromConfig(cfg)

	// Create shutdown manager with configured timeout
	shutdownMgr := NewShutdownManagerWithTimeout(cfg.Server.ShutdownTimeout)

	// Create application
	app := New(cfg, httpServer, shutdownMgr)

	// Return cleanup function
	cleanup := func() {
		// Cleanup resources if needed
	}

	return app, cleanup, nil
}

// NewShutdownManagerFromConfig creates a shutdown manager from config
func NewShutdownManagerFromConfig(cfg *config.Config) *ShutdownManager {
	return NewShutdownManagerWithTimeout(cfg.Server.ShutdownTimeout)
}

// NewShutdownManagerWithTimeout creates a shutdown manager with specific timeout
func NewShutdownManagerWithTimeout(timeout time.Duration) *ShutdownManager {
	return &ShutdownManager{
		closers: make([]Closer, 0),
		timeout: timeout,
	}
}

// NewHTTPServerFromConfig creates an HTTP server from config
func NewHTTPServerFromConfig(cfg *config.Config) HTTPServer {
	return &httpServerWrapper{
		Server: &http.Server{
			Addr:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
	}
}

// httpServerWrapper wraps http.Server to implement HTTPServer interface
type httpServerWrapper struct {
	*http.Server
}

// Addr returns the server address
func (s *httpServerWrapper) Addr() string {
	return s.Server.Addr
}
