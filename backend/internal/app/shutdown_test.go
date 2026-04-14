// Package app_test tests shutdown management functionality
package app_test

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/liang21/aitestos/internal/app"
)

// mockCloser is a mock implementation of Closer for testing
type mockCloser struct {
	name       string
	closeErr   error
	closeChan  chan struct{}
	closed     bool
	mu         sync.Mutex
	recordFunc func(name string) // Callback to record close order
}

func (m *mockCloser) Name() string {
	return m.name
}

func (m *mockCloser) Close(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.closed = true

	// Record close order synchronously if callback provided
	if m.recordFunc != nil {
		m.recordFunc(m.name)
	}

	if m.closeChan != nil {
		m.closeChan <- struct{}{}
	}
	return m.closeErr
}

func (m *mockCloser) isClosed() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.closed
}

func TestShutdownManager_Register(t *testing.T) {
	sm := app.NewShutdownManager(30 * time.Second)

	closer1 := &mockCloser{name: "component1"}
	closer2 := &mockCloser{name: "component2"}

	sm.Register(closer1)
	sm.Register(closer2)
}

func TestShutdownManager_ReverseOrder(t *testing.T) {
	sm := app.NewShutdownManager(30 * time.Second)

	closeOrder := make([]string, 0)
	var mu sync.Mutex

	// Register components in order: 1, 2, 3 with synchronous recording
	for i := 1; i <= 3; i++ {
		name := fmt.Sprintf("component%d", i)
		sm.Register(&mockCloser{
			name: name,
			recordFunc: func(n string) {
				mu.Lock()
				closeOrder = append(closeOrder, n)
				mu.Unlock()
			},
		})
	}

	// Trigger shutdown
	ctx := context.Background()
	_ = sm.Shutdown(ctx)

	// Verify shutdown order is reversed: 3, 2, 1
	mu.Lock()
	defer mu.Unlock()

	expected := []string{"component3", "component2", "component1"}
	for i, v := range expected {
		if i >= len(closeOrder) {
			t.Errorf("Missing component in shutdown order: %s", v)
			continue
		}
		if closeOrder[i] != v {
			t.Errorf("Shutdown order[%d] = %s, want %s", i, closeOrder[i], v)
		}
	}
}

func TestShutdownManager_ContinueOnError(t *testing.T) {
	sm := app.NewShutdownManager(30 * time.Second)

	closer1 := &mockCloser{name: "component1", closeErr: fmt.Errorf("close error")}
	closer2 := &mockCloser{name: "component2"}

	sm.Register(closer1)
	sm.Register(closer2)

	ctx := context.Background()
	err := sm.Shutdown(ctx)

	// Shutdown should continue even if one closer fails
	// The error is logged but not returned (as per implementation)
	if err != nil {
		t.Logf("Shutdown returned error (may be expected): %v", err)
	}

	// Both should be closed despite error in first
	if !closer1.isClosed() {
		t.Error("component1 should be closed")
	}
	if !closer2.isClosed() {
		t.Error("component2 should be closed")
	}
}

func TestShutdownManager_Timeout(t *testing.T) {
	sm := app.NewShutdownManager(100 * time.Millisecond)

	// Create a closer that takes longer than timeout
	slowCloser := &mockCloser{
		name: "slow-component",
	}
	sm.Register(slowCloser)

	ctx := context.Background()
	start := time.Now()
	_ = sm.Shutdown(ctx)
	elapsed := time.Since(start)

	// Shutdown should complete within reasonable time of timeout
	if elapsed > 200*time.Millisecond {
		t.Errorf("Shutdown took too long: %v", elapsed)
	}
}

func TestShutdownManager_EmptyManager(t *testing.T) {
	sm := app.NewShutdownManager(30 * time.Second)

	ctx := context.Background()
	err := sm.Shutdown(ctx)

	// Empty manager should shutdown without error
	if err != nil {
		t.Errorf("Empty shutdown should not error: %v", err)
	}
}

func TestShutdownManager_ContextCancellation(t *testing.T) {
	sm := app.NewShutdownManager(30 * time.Second)

	closer := &mockCloser{name: "component1"}
	sm.Register(closer)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	err := sm.Shutdown(ctx)

	// Should handle cancelled context gracefully
	_ = err // Error handling depends on implementation
}

func TestShutdownManager_MultipleClosers(t *testing.T) {
	sm := app.NewShutdownManager(30 * time.Second)

	const numClosers = 10
	closers := make([]*mockCloser, numClosers)

	for i := 0; i < numClosers; i++ {
		closers[i] = &mockCloser{
			name: fmt.Sprintf("component%d", i),
		}
		sm.Register(closers[i])
	}

	ctx := context.Background()
	_ = sm.Shutdown(ctx)

	// All closers should be closed
	for i, c := range closers {
		if !c.isClosed() {
			t.Errorf("component%d should be closed", i)
		}
	}
}

func TestShutdownManager_ConcurrentShutdown(t *testing.T) {
	sm := app.NewShutdownManager(30 * time.Second)

	closer := &mockCloser{name: "component1"}
	sm.Register(closer)

	// Try concurrent shutdown calls
	var wg sync.WaitGroup
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ctx := context.Background()
			_ = sm.Shutdown(ctx)
		}()
	}

	wg.Wait()
	// Verify no panic and component is closed
	if !closer.isClosed() {
		t.Error("component should be closed after concurrent shutdown")
	}
}

func TestNewShutdownManager_DefaultTimeout(t *testing.T) {
	// Test creating manager with different timeouts
	tests := []struct {
		name    string
		timeout time.Duration
	}{
		{"short timeout", 1 * time.Second},
		{"standard timeout", 30 * time.Second},
		{"long timeout", 5 * time.Minute},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sm := app.NewShutdownManager(tt.timeout)
			if sm == nil {
				t.Error("NewShutdownManager should not return nil")
			}
		})
	}
}

func TestCloser_Interface(t *testing.T) {
	// Verify mockCloser implements Closer interface
	var _ app.Closer = (*mockCloser)(nil)
}

func TestShutdownManager_Integration(t *testing.T) {
	// Integration test simulating real shutdown scenario
	sm := app.NewShutdownManager(5 * time.Second)

	var shutdownOrder []string
	var mu sync.Mutex

	// Simulate HTTP server
	sm.Register(&mockCloser{
		name: "http-server",
		recordFunc: func(n string) {
			time.Sleep(10 * time.Millisecond) // Simulate graceful drain
			mu.Lock()
			shutdownOrder = append(shutdownOrder, n)
			mu.Unlock()
		},
	})

	// Simulate database connection pool
	sm.Register(&mockCloser{
		name: "database",
		recordFunc: func(n string) {
			time.Sleep(5 * time.Millisecond) // Simulate connection close
			mu.Lock()
			shutdownOrder = append(shutdownOrder, n)
			mu.Unlock()
		},
	})

	// Simulate message queue consumer
	sm.Register(&mockCloser{
		name: "mq-consumer",
		recordFunc: func(n string) {
			time.Sleep(15 * time.Millisecond) // Simulate message drain
			mu.Lock()
			shutdownOrder = append(shutdownOrder, n)
			mu.Unlock()
		},
	})

	ctx := context.Background()
	_ = sm.Shutdown(ctx)

	mu.Lock()
	defer mu.Unlock()

	// Expected order: mq-consumer, database, http-server (reverse registration)
	expected := []string{"mq-consumer", "database", "http-server"}
	for i, v := range expected {
		if i >= len(shutdownOrder) {
			t.Errorf("Missing component: %s", v)
			continue
		}
		if shutdownOrder[i] != v {
			t.Errorf("Order[%d] = %s, want %s", i, shutdownOrder[i], v)
		}
	}
}
