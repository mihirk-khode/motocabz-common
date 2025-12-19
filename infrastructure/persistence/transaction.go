package persistence

import (
	"context"
	"fmt"
)

// WithTransaction executes a function within a transaction - simple and clean
// This is a generic helper that works with any transaction type
// For Ent, use: WithEntTransaction
// For other ORMs, create similar wrappers
func WithTransaction(ctx context.Context, fn func(context.Context) error) error {
	// This is a placeholder - actual implementation depends on the ORM
	// See WithEntTransaction for Ent-specific implementation
	return fn(ctx)
}

// TransactionFunc represents a function that operates within a transaction
type TransactionFunc func(context.Context) error

// TransactionManager provides transaction management interface
type TransactionManager interface {
	Begin(ctx context.Context) (context.Context, error)
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

// ExecuteTransaction executes a function within a managed transaction
func ExecuteTransaction(ctx context.Context, manager TransactionManager, fn TransactionFunc) error {
	txCtx, err := manager.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			manager.Rollback(txCtx)
			panic(p)
		}
	}()

	if err := fn(txCtx); err != nil {
		if rbErr := manager.Rollback(txCtx); rbErr != nil {
			return fmt.Errorf("transaction failed: %w, rollback failed: %v", err, rbErr)
		}
		return err
	}

	if err := manager.Commit(txCtx); err != nil {
		return fmt.Errorf("commit failed: %w", err)
	}

	return nil
}

