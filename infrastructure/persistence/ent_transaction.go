package persistence

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// EntTx represents an Ent transaction interface
// This allows the transaction helper to work with Ent without direct dependency
type EntTx interface {
	Commit() error
	Rollback() error
}

// EntClient represents an Ent client interface
type EntClient interface {
	Tx(context.Context) (EntTx, error)
}

// WithEntTransaction executes a function within an Ent transaction with tracing
func WithEntTransaction(ctx context.Context, db EntClient, fn func(EntTx) error) error {
	// Start span for transaction
	ctx, span := trace.SpanFromContext(ctx).TracerProvider().
		Tracer("motocabz-common/persistence").
		Start(ctx, "database.transaction")
	defer span.End()

	tx, err := db.Tx(ctx)
	if err != nil {
		if span.IsRecording() {
			span.SetStatus(codes.Error, "Failed to start transaction")
			span.RecordError(err)
		}
		return fmt.Errorf("start transaction: %w", err)
	}

	if span.IsRecording() {
		span.SetAttributes(attribute.String("db.operation", "transaction"))
	}

	defer func() {
		if p := recover(); p != nil {
			if span.IsRecording() {
				span.SetStatus(codes.Error, "Transaction panic")
				span.RecordError(fmt.Errorf("panic: %v", p))
			}
			tx.Rollback()
			panic(p)
		}
	}()

	if err := fn(tx); err != nil {
		if span.IsRecording() {
			span.SetStatus(codes.Error, "Transaction failed")
			span.RecordError(err)
		}
		if rbErr := tx.Rollback(); rbErr != nil {
			if span.IsRecording() {
				span.RecordError(rbErr)
			}
			return fmt.Errorf("transaction failed: %w, rollback failed: %v", err, rbErr)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		if span.IsRecording() {
			span.SetStatus(codes.Error, "Commit failed")
			span.RecordError(err)
		}
		return fmt.Errorf("commit failed: %w", err)
	}

	if span.IsRecording() {
		span.SetStatus(codes.Ok, "Transaction committed")
	}
	return nil
}
