// Package http provides HTTP transport layer implementation
package http

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewServer(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		host    string
		port    int
		handler http.Handler
		wantErr bool
	}{
		{
			name:    "valid config",
			host:    "localhost",
			port:    8080,
			handler: http.NewServeMux(),
			wantErr: false,
		},
		{
			name:    "empty host",
			host:    "",
			port:    8080,
			handler: http.NewServeMux(),
			wantErr: true,
		},
		{
			name:    "invalid port negative",
			host:    "localhost",
			port:    -1,
			handler: http.NewServeMux(),
			wantErr: true,
		},
		{
			name:    "invalid port too high",
			host:    "localhost",
			port:    70000,
			handler: http.NewServeMux(),
			wantErr: true,
		},
		{
			name:    "nil handler defaults to empty mux",
			host:    "localhost",
			port:    8080,
			handler: nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			srv, err := NewServer(tt.host, tt.port, tt.handler)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, srv)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, srv)
			}
		})
	}
}

func TestServerShutdown(t *testing.T) {
	t.Parallel()

	srv, err := NewServer("localhost", 0, http.NewServeMux())
	require.NoError(t, err)
	require.NotNil(t, srv)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = srv.Shutdown(ctx)
	assert.NoError(t, err)
}

func TestServerAddress(t *testing.T) {
	t.Parallel()

	srv, err := NewServer("127.0.0.1", 8888, http.NewServeMux())
	require.NoError(t, err)

	assert.Equal(t, "127.0.0.1:8888", srv.Address())
}

func TestServerRun(t *testing.T) {
	t.Parallel()

	srv, err := NewServer("localhost", 0, http.NewServeMux())
	require.NoError(t, err)

	// Run server in goroutine
	errCh := make(chan error, 1)
	go func() {
		if err := srv.Run(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
		close(errCh)
	}()

	// Give server time to start
	time.Sleep(100 * time.Millisecond)

	// Shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err = srv.Shutdown(ctx)
	assert.NoError(t, err)

	// Check no error from Run
	select {
	case err := <-errCh:
		assert.NoError(t, err)
	case <-time.After(1 * time.Second):
		// Server shutdown cleanly
	}
}
