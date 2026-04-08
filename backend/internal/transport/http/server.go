// Package http provides HTTP transport layer implementation
package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog"
)

// Server represents an HTTP server
type Server struct {
	server *http.Server
	host   string
	port   int
}

// NewServer creates a new HTTP server
func NewServer(host string, port int, handler http.Handler) (*Server, error) {
	if host == "" {
		return nil, fmt.Errorf("host cannot be empty")
	}
	if port < 0 || port > 65535 {
		return nil, fmt.Errorf("port must be between 0 and 65535")
	}

	if handler == nil {
		handler = http.NewServeMux()
	}

	return &Server{
		server: &http.Server{
			Addr:         fmt.Sprintf("%s:%d", host, port),
			Handler:      handler,
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
			IdleTimeout:  120 * time.Second,
		},
		host: host,
		port: port,
	}, nil
}

// Run starts the HTTP server
func (s *Server) Run() error {
	logger := zerolog.Nop()
	logger.Info().Str("addr", s.Address()).Msg("starting HTTP server")

	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("server error: %w", err)
	}
	return nil
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	logger := zerolog.Nop()
	logger.Info().Msg("shutting down HTTP server")

	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("shutdown error: %w", err)
	}
	return nil
}

// Address returns the server address
func (s *Server) Address() string {
	return fmt.Sprintf("%s:%d", s.host, s.port)
}
