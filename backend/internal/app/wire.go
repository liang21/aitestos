//go:build wireinject
// +build wireinject

// Package app provides Wire dependency injection setup
// This file is only compiled when using -tags=wireinject
// For normal builds, use initialize.go which has manual dependency injection
package app

import (
	"github.com/google/wire"

	"github.com/liang21/aitestos/internal/config"
)

// InfrastructureSet provides infrastructure dependencies (database, cache, etc.)
var InfrastructureSet = wire.NewSet(
	ProvideDB,
	ProvideShutdownManager,
)

// HTTPServerSet provides HTTP server dependencies
var HTTPServerSet = wire.NewSet(
	ProvideHTTPServer,
)

// AppSet provides application dependencies
var AppSet = wire.NewSet(
	ProvideApp,
)

// SuperSet combines all provider sets
var SuperSet = wire.NewSet(
	InfrastructureSet,
	HTTPServerSet,
	AppSet,
)

// Initialize creates a fully initialized application
// This function signature is used by Wire to generate wire_gen.go
func Initialize(cfg *config.Config) (*App, func(), error) {
	wire.Build(SuperSet)
	return nil, nil, nil
}

// ProvideDB provides database connection (placeholder)
func ProvideDB(cfg *config.Config) (interface{ Close() error }, error) {
	// Placeholder - will be replaced with actual DB provider
	return &noopCloser{}, nil
}

// ProvideShutdownManager provides shutdown manager
func ProvideShutdownManager(cfg *config.Config) *ShutdownManager {
	return NewShutdownManagerFromConfig(cfg)
}

// ProvideHTTPServer provides HTTP server
func ProvideHTTPServer(cfg *config.Config) HTTPServer {
	return NewHTTPServerFromConfig(cfg)
}

// ProvideApp provides the application
func ProvideApp(cfg *config.Config, httpServer HTTPServer, shutdownMgr *ShutdownManager) *App {
	return New(cfg, httpServer, shutdownMgr)
}

// noopCloser is a no-op closer for placeholder
type noopCloser struct{}

func (n *noopCloser) Close() error { return nil }
