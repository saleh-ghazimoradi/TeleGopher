package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/saleh-ghazimoradi/TeleGopher/internal/domain"
	"gorm.io/gorm"
)

type MessageRepository interface {
	CreateMessage(ctx context.Context, message *domain.Message) error
	GetMessageById(ctx context.Context, id uint) (*domain.Message, error)
	GetMessageByPrivateId(ctx context.Context, privateId uint, page, limit int) ([]domain.Message, error)
	GetUndeliveredMessagesByPrivateId(ctx context.Context, privateId uint) ([]domain.Message, error)
	MarkMessageAsRead(ctx context.Context, id uint) error
	MarkMessageAsDelivered(ctx context.Context, id uint) error
}

type messageRepository struct {
	dbWrite *gorm.DB
	dbRead  *gorm.DB
}

func (m *messageRepository) CreateMessage(ctx context.Context, message *domain.Message) error {
	return m.dbWrite.WithContext(ctx).Create(&message).Error
}

func (m *messageRepository) GetMessageById(ctx context.Context, id uint) (*domain.Message, error) {
	var message domain.Message

	if err := m.dbRead.WithContext(ctx).Preload("From").Preload("Private").First(&message, id).Error; err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &message, nil
}

func (m *messageRepository) GetMessageByPrivateId(ctx context.Context, privateId uint, offset, limit int) ([]domain.Message, error) {
	var messages []domain.Message

	if err := m.dbRead.WithContext(ctx).Where("private_id = ?", privateId).Preload("From").Order("created_at DESC").Offset(offset).Limit(limit).Find(&messages).Error; err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return messages, nil
}

func (m *messageRepository) GetUndeliveredMessagesByPrivateId(ctx context.Context, privateId uint) ([]domain.Message, error) {
	var messages []domain.Message
	if err := m.dbRead.WithContext(ctx).
		Where("private_id = ? AND delivered = ?", privateId, false).
		Preload("From").
		Order("created_at ASC").
		Find(&messages).Error; err != nil {
		return nil, fmt.Errorf("failed to get undelivered messages: %w", err)
	}

	return messages, nil
}

func (m *messageRepository) MarkMessageAsRead(ctx context.Context, id uint) error {
	if err := m.dbWrite.WithContext(ctx).Model(&domain.Message{}).Where("id = ? AND read = ?", id, false).Update("read", true).Error; err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return ErrRecordNotFound
		default:
			return err
		}
	}
	return nil
}

func (m *messageRepository) MarkMessageAsDelivered(ctx context.Context, id uint) error {
	if err := m.dbWrite.WithContext(ctx).Model(&domain.Message{}).
		Where("id = ? AND delivered = ?", id, false).
		Update("delivered", true).Error; err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return ErrRecordNotFound
		default:
			return err
		}
	}
	return nil
}

func NewMessageRepository(dbWrite, dbRead *gorm.DB) MessageRepository {
	return &messageRepository{
		dbWrite: dbWrite,
		dbRead:  dbRead,
	}
}
