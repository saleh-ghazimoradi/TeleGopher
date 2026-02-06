package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/saleh-ghazimoradi/TeleGopher/internal/domain"
)

type MessageRepository interface {
	CreateMessage(ctx context.Context, message *domain.Message) error
	GetMessageById(ctx context.Context, messageId int64) (*domain.Message, error)
	GetMessageByPrivateId(ctx context.Context, privateId int64, page, limit int) ([]domain.Message, error)
	GetUndeliveredMessagesByPrivateId(ctx context.Context, privateId int64) ([]domain.Message, error)
	MarkMessageRead(ctx context.Context, messageId int64) error
	MarkMessageDelivered(ctx context.Context, messageId int64) error
	WithTx(tx *sql.Tx) MessageRepository
}

type messageRepository struct {
	dbWrite *sql.DB
	dbRead  *sql.DB
	tx      *sql.Tx
}

func (m *messageRepository) CreateMessage(ctx context.Context, message *domain.Message) error {
	query := `
		INSERT INTO messages (from_id, private_id, message_type, content, delivered, read)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, version
	`

	args := []any{message.FromId, message.PrivateId, message.MessageType, message.Content, message.Delivered, message.Read}

	if err := querier(m.dbWrite, m.tx).QueryRowContext(ctx, query, args...).Scan(&message.Id, &message.CreatedAt, &message.Version); err != nil {
		return fmt.Errorf("failed to create message: %w", err)
	}
	return nil
}

func (m *messageRepository) GetMessageByPrivateId(ctx context.Context, privateId int64, page, limit int) ([]domain.Message, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	query := `
		SELECT id, from_id, private_id, message_type, content, delivered, read, created_at, version
		FROM messages
		WHERE private_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := querier(m.dbRead, m.tx).QueryContext(ctx, query, privateId, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query messages: %w", err)
	}
	defer rows.Close()

	var messages []domain.Message
	for rows.Next() {
		var msg domain.Message
		if err := rows.Scan(
			&msg.Id,
			&msg.FromId,
			&msg.PrivateId,
			&msg.MessageType,
			&msg.Content,
			&msg.Delivered,
			&msg.Read,
			&msg.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan message: %w", err)
		}
		messages = append(messages, msg)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return messages, nil
}

func (m *messageRepository) GetMessageById(ctx context.Context, messageId int64) (*domain.Message, error) {
	message := &domain.Message{}

	query := `
		SELECT id, from_id, private_id, message_type, content, delivered, read, created_at, version
		FROM messages
		WHERE id = $1
	`

	if err := querier(m.dbRead, m.tx).QueryRowContext(ctx, query, messageId).Scan(
		&message.Id,
		&message.FromId,
		&message.PrivateId,
		&message.MessageType,
		&message.Content,
		&message.Delivered,
		&message.Read,
		&message.CreatedAt,
		&message.Version,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
		return nil, fmt.Errorf("failed to get message by id: %w", err)
	}

	return message, nil
}

func (m *messageRepository) GetUndeliveredMessagesByPrivateId(ctx context.Context, privateId int64) ([]domain.Message, error) {
	query := `
		SELECT id, from_id, private_id, message_type, content, delivered, read, created_at, version
		FROM messages
		WHERE private_id = $1 AND delivered = false
		ORDER BY created_at ASC
	`

	rows, err := querier(m.dbRead, m.tx).QueryContext(ctx, query, privateId)
	if err != nil {
		return nil, fmt.Errorf("failed to query undelivered messages: %w", err)
	}
	defer rows.Close()

	var messages []domain.Message
	for rows.Next() {
		var msg domain.Message
		if err := rows.Scan(
			&msg.Id,
			&msg.FromId,
			&msg.PrivateId,
			&msg.MessageType,
			&msg.Content,
			&msg.Delivered,
			&msg.Read,
			&msg.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan message: %w", err)
		}
		messages = append(messages, msg)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return messages, nil
}

func (m *messageRepository) MarkMessageRead(ctx context.Context, messageId int64) error {
	query := `
		UPDATE messages
		SET read = true
		WHERE id = $1 AND read = false
	`

	result, err := querier(m.dbWrite, m.tx).ExecContext(ctx, query, messageId)
	if err != nil {
		return fmt.Errorf("failed to mark message as read: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("message already read or not found")
	}

	return nil
}

func (m *messageRepository) MarkMessageDelivered(ctx context.Context, messageId int64) error {
	query := `
		UPDATE messages
		SET delivered = true
		WHERE id = $1 AND delivered = false
	`

	result, err := querier(m.dbWrite, m.tx).ExecContext(ctx, query, messageId)
	if err != nil {
		return fmt.Errorf("failed to mark message as delivered: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("message already delivered or not found")
	}

	return nil
}

func (m *messageRepository) WithTx(tx *sql.Tx) MessageRepository {
	return &messageRepository{
		dbWrite: m.dbWrite,
		dbRead:  m.dbRead,
		tx:      tx,
	}
}

func NewMessageRepository(dbWrite *sql.DB, dbRead *sql.DB) MessageRepository {
	return &messageRepository{
		dbWrite: dbWrite,
		dbRead:  dbRead,
	}
}
