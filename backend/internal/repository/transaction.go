// Package repository provides transaction management
package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

// ctxKey is the context key for storing transaction
type ctxKey struct{}

// Tx is the interface for database transaction
type Tx interface {
	Commit() error
	Rollback() error
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

// sqlxTxWrapper wraps sqlx.Tx to implement Tx interface
type sqlxTxWrapper struct {
	*sqlx.Tx
}

func (w *sqlxTxWrapper) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return w.Tx.ExecContext(ctx, query, args...)
}

func (w *sqlxTxWrapper) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return w.Tx.QueryContext(ctx, query, args...)
}

func (w *sqlxTxWrapper) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return w.Tx.QueryRowContext(ctx, query, args...)
}

// WrapTx wraps sqlx.Tx to implement Tx interface
func WrapTx(tx *sqlx.Tx) Tx {
	return &sqlxTxWrapper{Tx: tx}
}

// TxManager manages database transactions
type TxManager struct {
	beginFn func(ctx context.Context) (Tx, error)
}

// NewTxManager creates a new transaction manager
func NewTxManager(beginFn func(ctx context.Context) (Tx, error)) *TxManager {
	return &TxManager{
		beginFn: beginFn,
	}
}

// Begin starts a new transaction
func (tm *TxManager) Begin(ctx context.Context) (Tx, error) {
	// Check if there's already a transaction in context
	if tx := TxFromContext(ctx); tx != nil {
		// Return existing transaction for nested transaction support
		return tx, nil
	}
	return tm.beginFn(ctx)
}

// WithTransaction executes a function within a transaction
func (tm *TxManager) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	// Check if already in a transaction
	if TxFromContext(ctx) != nil {
		return fn(ctx)
	}

	tx, err := tm.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	// Store transaction in context
	txCtx := ContextWithTx(ctx, tx)

	// Defer rollback in case of panic
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	if err := fn(txCtx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("rollback failed: %v, original error: %w", rbErr, err)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

// TxFromContext retrieves the transaction from context
func TxFromContext(ctx context.Context) Tx {
	if tx, ok := ctx.Value(ctxKey{}).(Tx); ok {
		return tx
	}
	return nil
}

// ContextWithTx stores the transaction in context
func ContextWithTx(ctx context.Context, tx Tx) context.Context {
	return context.WithValue(ctx, ctxKey{}, tx)
}

// MustTxFromContext retrieves the transaction from context or panics
func MustTxFromContext(ctx context.Context) Tx {
	tx := TxFromContext(ctx)
	if tx == nil {
		panic(errors.New("no transaction in context"))
	}
	return tx
}
