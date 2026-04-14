// Package repository_test tests transaction manager
package repository_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/liang21/aitestos/internal/repository"
)

// mockTx is a mock implementation of Tx interface for testing
type mockTx struct {
	commitErr   error
	rollbackErr error
	committed   bool
	rolledBack  bool
}

func (m *mockTx) Commit() error {
	m.committed = true
	return m.commitErr
}

func (m *mockTx) Rollback() error {
	m.rolledBack = true
	return m.rollbackErr
}

func (m *mockTx) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return nil, nil
}

func (m *mockTx) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return nil, nil
}

func (m *mockTx) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return nil
}

func TestNewTxManager(t *testing.T) {
	beginFn := func(ctx context.Context) (repository.Tx, error) {
		return &mockTx{}, nil
	}

	tm := repository.NewTxManager(beginFn)
	if tm == nil {
		t.Error("NewTxManager() returned nil")
	}
}

func TestTxManager_Begin(t *testing.T) {
	tests := []struct {
		name    string
		beginFn func(ctx context.Context) (repository.Tx, error)
		wantErr bool
	}{
		{
			name: "successful begin",
			beginFn: func(ctx context.Context) (repository.Tx, error) {
				return &mockTx{}, nil
			},
			wantErr: false,
		},
		{
			name: "begin error",
			beginFn: func(ctx context.Context) (repository.Tx, error) {
				return nil, errors.New("connection error")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			tm := repository.NewTxManager(tt.beginFn)

			tx, err := tm.Begin(ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("TxManager.Begin() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tx == nil {
				t.Error("TxManager.Begin() returned nil tx")
			}
		})
	}
}

func TestTxManager_WithTransaction(t *testing.T) {
	tests := []struct {
		name         string
		fn           func(ctx context.Context) error
		commitErr    error
		wantErr      bool
		wantCommit   bool
		wantRollback bool
	}{
		{
			name: "successful transaction",
			fn: func(ctx context.Context) error {
				return nil
			},
			wantCommit: true,
		},
		{
			name: "function error causes rollback",
			fn: func(ctx context.Context) error {
				return errors.New("business error")
			},
			wantErr:      true,
			wantRollback: true,
		},
		{
			name: "commit error",
			fn: func(ctx context.Context) error {
				return nil
			},
			commitErr: errors.New("commit failed"),
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockTx := &mockTx{commitErr: tt.commitErr}

			beginFn := func(ctx context.Context) (repository.Tx, error) {
				return mockTx, nil
			}
			tm := repository.NewTxManager(beginFn)

			ctx := context.Background()
			err := tm.WithTransaction(ctx, tt.fn)
			if (err != nil) != tt.wantErr {
				t.Errorf("TxManager.WithTransaction() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantCommit && !mockTx.committed {
				t.Error("expected transaction to be committed")
			}
			if tt.wantRollback && !mockTx.rolledBack {
				t.Error("expected transaction to be rolled back")
			}
		})
	}
}

func TestTxFromContext(t *testing.T) {
	t.Run("no transaction in context", func(t *testing.T) {
		ctx := context.Background()
		tx := repository.TxFromContext(ctx)
		if tx != nil {
			t.Error("TxFromContext() should return nil for context without transaction")
		}
	})

	t.Run("with transaction in context", func(t *testing.T) {
		tx := &mockTx{}
		ctx := repository.ContextWithTx(context.Background(), tx)
		retrieved := repository.TxFromContext(ctx)
		if retrieved == nil {
			t.Error("TxFromContext() should return transaction from context")
		}
	})
}

func TestContextWithTx(t *testing.T) {
	tx := &mockTx{}
	ctx := repository.ContextWithTx(context.Background(), tx)

	// Verify transaction can be retrieved
	retrieved := repository.TxFromContext(ctx)
	if retrieved == nil {
		t.Error("ContextWithTx() should store transaction in context")
	}
}

func TestTxManager_BeginNested(t *testing.T) {
	mockTx := &mockTx{}

	beginFn := func(ctx context.Context) (repository.Tx, error) {
		return mockTx, nil
	}
	tm := repository.NewTxManager(beginFn)

	ctx := context.Background()

	// Start first transaction
	tx1, err := tm.Begin(ctx)
	if err != nil {
		t.Fatalf("Begin() error = %v", err)
	}
	ctx1 := repository.ContextWithTx(ctx, tx1)

	// Try to start nested transaction (should return same tx)
	tx2, err := tm.Begin(ctx1)
	if err != nil {
		t.Fatalf("Begin() nested error = %v", err)
	}
	if tx2 == nil {
		t.Error("Begin() nested should return a transaction")
	}

	// Should return the same transaction instance
	if tx1 != tx2 {
		t.Error("Nested Begin() should return the same transaction")
	}
}

func TestTxManager_WithTransaction_AlreadyInTx(t *testing.T) {
	existingTx := &mockTx{}

	beginFn := func(ctx context.Context) (repository.Tx, error) {
		return &mockTx{}, nil // Should not be called
	}
	tm := repository.NewTxManager(beginFn)

	// Create context with existing transaction
	ctx := repository.ContextWithTx(context.Background(), existingTx)

	executed := false
	err := tm.WithTransaction(ctx, func(ctx context.Context) error {
		executed = true
		return nil
	})

	if err != nil {
		t.Errorf("WithTransaction() error = %v", err)
	}
	if !executed {
		t.Error("function should have been executed")
	}
	// Original existingTx should not be committed because we're in existing transaction
	if existingTx.committed {
		t.Error("existing transaction should not be committed by nested WithTransaction")
	}
}

func TestMustTxFromContext(t *testing.T) {
	t.Run("panics without transaction", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("MustTxFromContext() should panic without transaction")
			}
		}()
		repository.MustTxFromContext(context.Background())
	})

	t.Run("returns transaction", func(t *testing.T) {
		tx := &mockTx{}
		ctx := repository.ContextWithTx(context.Background(), tx)
		got := repository.MustTxFromContext(ctx)
		if got == nil {
			t.Error("MustTxFromContext() should return transaction")
		}
	})
}
