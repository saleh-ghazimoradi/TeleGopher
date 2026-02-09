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
	if user1Id == user2Id {
		return nil, ErrSameUser
	}

	if user1Id > user2Id {
		user1Id, user2Id = user2Id, user1Id
	}

	_, err := p.GetPrivateByUsers(ctx, user1Id, user2Id)
	if err == nil {
		return nil, ErrPrivateAlreadyExists
	}
	if !errors.Is(err, ErrRecordNotFound) {
		return nil, fmt.Errorf("failed to check for existing private conversation: %w", err)
	}

	priv := &domain.Private{}

	query := `
        INSERT INTO privates (user1_id, user2_id)
        VALUES ($1, $2)
        RETURNING id, user1_id, user2_id, created_at, version
    `

	err = querier(p.dbWrite, p.tx).QueryRowContext(ctx, query, user1Id, user2Id).
		Scan(
			&priv.Id,
			&priv.User1Id,
			&priv.User2Id,
			&priv.CreatedAt,
			&priv.Version,
		)
	if err != nil {
		return nil, fmt.Errorf("failed to create private conversation: %w", err)
	}

	return priv, nil
}

func (p *privateRepository) GetPrivateById(ctx context.Context, id int64) (*domain.Private, error) {
	var pvt domain.Private

	query := `
		SELECT id, user1_id, user2_id, created_at, version 
		FROM privates 
		WHERE id = $1
	`

	err := querier(p.dbRead, p.tx).QueryRowContext(ctx, query, id).
		Scan(
			&pvt.Id,
			&pvt.User1Id,
			&pvt.User2Id,
			&pvt.CreatedAt,
			&pvt.Version,
		)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
		return nil, fmt.Errorf("failed to get private by id: %w", err)
	}

	return &pvt, nil
}

func (p *privateRepository) GetPrivateByUsers(ctx context.Context, user1Id, user2Id int64) (*domain.Private, error) {
	var pvt domain.Private

	// Same canonical order as in Create
	if user1Id > user2Id {
		user1Id, user2Id = user2Id, user1Id
	}

	query := `
		SELECT id, user1_id, user2_id, created_at, version 
		FROM privates 
		WHERE user1_id = $1 AND user2_id = $2
	`

	err := querier(p.dbRead, p.tx).QueryRowContext(ctx, query, user1Id, user2Id).
		Scan(
			&pvt.Id,
			&pvt.User1Id,
			&pvt.User2Id,
			&pvt.CreatedAt,
			&pvt.Version,
		)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
		return nil, fmt.Errorf("failed to get private by users: %w", err)
	}

	return &pvt, nil
}

func (p *privateRepository) GetPrivateForUser(ctx context.Context, userId int64) ([]domain.Private, error) {
	query := `
		SELECT id, user1_id, user2_id, created_at, version 
		FROM privates 
		WHERE user1_id = $1 OR user2_id = $1
		ORDER BY created_at DESC
	`

	rows, err := querier(p.dbRead, p.tx).QueryContext(ctx, query, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to query private conversations for user: %w", err)
	}
	defer rows.Close()

	privates := make([]domain.Private, 0, 8) // initial capacity hint

	for rows.Next() {
		var pvt domain.Private
		if err := rows.Scan(
			&pvt.Id,
			&pvt.User1Id,
			&pvt.User2Id,
			&pvt.CreatedAt,
			&pvt.Version,
		); err != nil {
			return nil, fmt.Errorf("failed to scan private conversation: %w", err)
		}
		privates = append(privates, pvt)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating private conversations: %w", err)
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

func NewPrivateRepository(dbWrite, dbRead *sql.DB) PrivateRepository {
	return &privateRepository{
		dbWrite: dbWrite,
		dbRead:  dbRead,
	}
}
