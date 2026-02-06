package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/saleh-ghazimoradi/TeleGopher/internal/domain"
)

type PrivateRepository interface {
	GetPrivateById(ctx context.Context, id int64) (*domain.Private, error)
	GetPrivateByUsers(ctx context.Context, user1Id, user2Id int64) (*domain.Private, error)
	GetPrivateForUser(ctx context.Context, userId int64) ([]domain.Private, error)
	CreatePrivate(ctx context.Context, user1Id, user2Id int64) (*domain.Private, error)
	WithTx(tx *sql.Tx) PrivateRepository
}

type privateRepository struct {
	dbWrite *sql.DB
	dbRead  *sql.DB
	tx      *sql.Tx
}

func (p *privateRepository) CreatePrivate(ctx context.Context, user1Id, user2Id int64) (*domain.Private, error) {
	private := &domain.Private{}

	if user1Id == user2Id {
		return nil, ErrSameUser
	}

	if user1Id > user2Id {
		user1Id, user2Id = user2Id, user1Id
	}

	existing, err := p.GetPrivateByUsers(ctx, user1Id, user2Id)
	if err == nil && existing != nil {
		return nil, ErrPrivateAlreadyExists
	}
	if err != nil && !errors.Is(err, ErrPrivateAlreadyExists) {
		return nil, fmt.Errorf("failed to check existing private: %w", err)
	}

	query := `
		INSERT INTO privates (user1_id, user2_id)
		VALUES ($1, $2)
		RETURNING id, user1_id, user2_id, created_at, version
	`

	if err := querier(p.dbWrite, p.tx).QueryRowContext(ctx, query, user1Id, user2Id).Scan(
		&private.Id,
		&private.User1Id,
		&private.User2Id,
		&private.CreatedAt,
		&private.Version,
	); err != nil {
		return nil, fmt.Errorf("failed to create private: %w", err)
	}

	return private, nil
}

func (p *privateRepository) GetPrivateById(ctx context.Context, id int64) (*domain.Private, error) {
	var private domain.Private
	query := `SELECT id, user1_id, user2_id, created_at, version FROM privates WHERE id = $1`

	if err := querier(p.dbRead, p.tx).QueryRowContext(ctx, query, id).Scan(
		&private.Id,
		&private.User1Id,
		&private.User2Id,
		&private.CreatedAt,
		&private.Version,
	); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, fmt.Errorf("failed to get private by id: %w", err)
		}
	}

	return &private, nil
}

func (p *privateRepository) GetPrivateByUsers(ctx context.Context, user1Id, user2Id int64) (*domain.Private, error) {
	private := &domain.Private{}

	if user1Id > user2Id {
		user1Id, user2Id = user2Id, user1Id
	}

	query := `SELECT id, user1_id, user2_id, created_at, version FROM privates WHERE user1_id = $1 AND user2_id = $2`

	if err := querier(p.dbRead, p.tx).QueryRowContext(ctx, query, user1Id, user2Id).Scan(
		&private.Id,
		&private.User1Id,
		&private.User2Id,
		&private.CreatedAt,
		&private.Version,
	); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, fmt.Errorf("failed to get private by user: %w", err)
		}
	}
	return private, nil
}

func (p *privateRepository) GetPrivateForUser(ctx context.Context, userId int64) ([]domain.Private, error) {
	privates := make([]domain.Private, 0)

	query := `SELECT id, user1_id, user2_id, created_at, version 
          FROM privates 
          WHERE user1_id = $1 OR user2_id = $2
          ORDER BY created_at DESC`

	rows, err := querier(p.dbRead, p.tx).QueryContext(ctx, query, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to get privates for user: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		private := domain.Private{}
		if err := rows.Scan(&private.Id, &private.User1Id, &private.User2Id, &private.CreatedAt, &private.Version); err != nil {
			return nil, fmt.Errorf("failed to scan private for user: %w", err)
		}
		privates = append(privates, private)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to get privates for user: %w", err)
	}

	return privates, nil
}

func (p *privateRepository) WithTx(tx *sql.Tx) PrivateRepository {
	return &privateRepository{
		dbWrite: p.dbWrite,
		dbRead:  p.dbRead,
		tx:      tx,
	}
}

func NewPrivateRepository(dbWrite *sql.DB, dbRead *sql.DB) PrivateRepository {
	return &privateRepository{
		dbWrite: dbWrite,
		dbRead:  dbRead,
	}
}
