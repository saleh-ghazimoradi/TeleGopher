package TXManager

import (
	"context"
	"database/sql"
	"fmt"
)

type TxManager interface {
	WithTransaction(ctx context.Context, fn TransactionFunc) error
}

type txManager struct {
	db *sql.DB
}

type TransactionFunc func(*sql.Tx) error

func (tm *txManager) WithTransaction(ctx context.Context, fn TransactionFunc) error {
	tx, err := tm.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
	})
	if err != nil {
		return fmt.Errorf("beginning transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}
	}()

	if err := fn(tx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("rolling back transaction: %v (original error: %w)", rbErr, err)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}

	return nil
}

func NewTxManager(db *sql.DB) TxManager {
	return &txManager{db: db}
}
