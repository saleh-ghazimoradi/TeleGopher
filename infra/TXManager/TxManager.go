package infra

import (
	"context"
	"database/sql"
	"fmt"
)

type TxManager interface {
	WithTransaction(ctx context.Context, fn func(tx *sql.Tx) error) error
	WithTransactionOpts(ctx context.Context, opts *sql.TxOptions, fn func(tx *sql.Tx) error) error
}

type txManager struct {
	db *sql.DB
}

func (m *txManager) WithTransaction(ctx context.Context, fn func(tx *sql.Tx) error) error {
	return m.WithTransactionOpts(ctx, nil, fn)
}

func (m *txManager) WithTransactionOpts(ctx context.Context, opts *sql.TxOptions, fn func(tx *sql.Tx) error) error {
	tx, err := m.db.BeginTx(ctx, opts)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p) // re-panic after rollback
		}
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				err = fmt.Errorf("rollback failed: %v (original: %w)", rbErr, err)
			}
			return
		}
		if cmErr := tx.Commit(); cmErr != nil {
			err = fmt.Errorf("commit failed: %w", cmErr)
		}
	}()

	err = fn(tx)
	return err
}

func NewTxManager(db *sql.DB) TxManager {
	return &txManager{db: db}
}
