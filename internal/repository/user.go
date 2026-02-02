package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/saleh-ghazimoradi/TeleGopher/internal/domain"
	"time"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *domain.User) error
	GetUserById(ctx context.Context, id int64) (*domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	GetUserByRefreshToken(ctx context.Context, refreshToken string, platform domain.Platform) (*domain.User, error)
	UpdateRefreshToken(ctx context.Context, userId int64, platform domain.Platform, refreshToken string) error
	DeleteRefreshToken(ctx context.Context, userId int64, platform domain.Platform) error
	WithTx(tx *sql.Tx) UserRepository
}

type userRepository struct {
	dbWrite *sql.DB
	dbRead  *sql.DB
	tx      *sql.Tx
}

func (u *userRepository) CreateUser(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (name, email, password)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, version`

	args := []any{user.Name, user.Email, user.Password}

	if err := querier(u.dbWrite, u.tx).QueryRowContext(ctx, query, args...).Scan(&user.Id, &user.CreatedAt, &user.Version); err != nil {
		if err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"` {
			return ErrDuplicateEmail
		}
		return fmt.Errorf("create user: %w", err)
	}

	return nil
}

func (u *userRepository) GetUserById(ctx context.Context, id int64) (*domain.User, error) {
	query := `SELECT id, name, email, password, refresh_token_web, refresh_token_web_at, refresh_token_mobile, refresh_token_mobile_at, created_at, version FROM users WHERE id = $1`

	user := &domain.User{}

	if err := querier(u.dbRead, u.tx).QueryRowContext(ctx, query, id).Scan(
		&user.Id,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.RefreshTokenWeb,
		&user.RefreshTokenWebAt,
		&user.RefreshTokenMobile,
		&user.RefreshTokenMobileAt,
		&user.CreatedAt,
		&user.Version,
	); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return user, nil
}

func (u *userRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `SELECT id, name, email, password, refresh_token_web, refresh_token_web_at, refresh_token_mobile, refresh_token_mobile_at, created_at, version FROM users WHERE email = $1`

	user := &domain.User{}

	if err := querier(u.dbRead, u.tx).QueryRowContext(ctx, query, email).Scan(
		&user.Id,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.RefreshTokenWeb,
		&user.RefreshTokenWebAt,
		&user.RefreshTokenMobile,
		&user.RefreshTokenMobileAt,
		&user.CreatedAt,
		&user.Version,
	); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return user, nil
}

func (u *userRepository) GetUserByRefreshToken(ctx context.Context, refreshToken string, platform domain.Platform) (*domain.User, error) {
	user := &domain.User{}
	var query string
	switch platform {
	case domain.PlatformWeb:
		query = `SELECT id, name, email, password,
		                refresh_token_web, refresh_token_web_at,
		                refresh_token_mobile, refresh_token_mobile_at,
		                created_at, version
		         FROM users 
		         WHERE refresh_token_web = $1`

	case domain.PlatformMobile:
		query = `SELECT id, name, email, password,
		                refresh_token_web, refresh_token_web_at,
		                refresh_token_mobile, refresh_token_mobile_at,
		                created_at, version
		         FROM users 
		         WHERE refresh_token_mobile = $1`
	default:
		return nil, fmt.Errorf("invalid platform: %s", platform)
	}

	if err := querier(u.dbRead, u.tx).QueryRowContext(ctx, query, refreshToken).Scan(
		&user.Id, &user.Name, &user.Email, &user.Password,
		&user.RefreshTokenWeb, &user.RefreshTokenWebAt,
		&user.RefreshTokenMobile, &user.RefreshTokenMobileAt,
		&user.CreatedAt, &user.Version); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		}
		return nil, err
	}

	return user, nil
}

func (u *userRepository) UpdateRefreshToken(ctx context.Context, userId int64, platform domain.Platform, refreshToken string) error {
	var query string

	switch platform {
	case domain.PlatformWeb:
		query = `
			UPDATE users 
			SET refresh_token_web = $1, 
			    refresh_token_web_at = $2,
			    version = version + 1
			WHERE id = $3
		`
		_, err := querier(u.dbWrite, u.tx).ExecContext(ctx, query, refreshToken, time.Now(), userId)
		return err

	case domain.PlatformMobile:
		query = `
			UPDATE users 
			SET refresh_token_mobile = $1, 
			    refresh_token_mobile_at = $2,
			    version = version + 1
			WHERE id = $3
		`
		_, err := querier(u.dbWrite, u.tx).ExecContext(ctx, query, refreshToken, time.Now(), userId)
		return err
	default:
		return fmt.Errorf("invalid platform: %s", platform)
	}
}

func (u *userRepository) DeleteRefreshToken(ctx context.Context, userId int64, platform domain.Platform) error {
	var query string

	switch platform {
	case domain.PlatformWeb:
		query = `
			UPDATE users 
			SET refresh_token_web = NULL, 
			    refresh_token_web_at = NULL,
			    version = version + 1
			WHERE id = $1
		`
	case domain.PlatformMobile:
		query = `
			UPDATE users 
			SET refresh_token_mobile = NULL, 
			    refresh_token_mobile_at = NULL,
			    version = version + 1
			WHERE id = $1
		`
	default:
		return fmt.Errorf("invalid platform: %s", platform)
	}

	_, err := querier(u.dbWrite, u.tx).ExecContext(ctx, query, userId)
	return err
}

func (u *userRepository) WithTx(tx *sql.Tx) UserRepository {
	return &userRepository{
		dbWrite: u.dbWrite,
		dbRead:  u.dbRead,
		tx:      tx,
	}
}

func NewUserRepository(dbWrite, dbRead *sql.DB) UserRepository {
	return &userRepository{
		dbWrite: dbWrite,
		dbRead:  dbRead,
	}
}
