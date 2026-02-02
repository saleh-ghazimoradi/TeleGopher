package repository

import (
	"context"
	"database/sql"
)

type DBTX interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

func querier(readWrite *sql.DB, tx *sql.Tx) DBTX {
	if tx != nil {
		return tx
	}
	return readWrite
}
