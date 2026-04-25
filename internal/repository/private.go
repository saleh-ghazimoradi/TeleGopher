package repository

import (
	"context"
	"errors"
	"github.com/saleh-ghazimoradi/TeleGopher/internal/domain"
	"gorm.io/gorm"
)

type PrivateRepository interface {
	CreatePrivate(ctx context.Context, private *domain.Private) error
	GetPrivateById(ctx context.Context, id uint) (*domain.Private, error)
	GetPrivateByUsers(ctx context.Context, user1Id, user2Id uint) (*domain.Private, error)
	GetPrivatesForUser(ctx context.Context, userId uint) ([]domain.Private, error)
	CheckPrivateExists(ctx context.Context, user1Id, user2Id uint) (bool, error)
}

type privateRepository struct {
	dbWrite *gorm.DB
	dbRead  *gorm.DB
}

func (p *privateRepository) CreatePrivate(ctx context.Context, private *domain.Private) error {
	if private.User1Id > private.User2Id {
		private.User1Id, private.User2Id = private.User2Id, private.User1Id
	}
	return p.dbWrite.WithContext(ctx).Create(&private).Error
}

func (p *privateRepository) GetPrivateById(ctx context.Context, id uint) (*domain.Private, error) {
	var private domain.Private

	if err := p.dbRead.WithContext(ctx).
		Preload("User1").
		Preload("User2").
		First(&private, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}
	return &private, nil
}

func (p *privateRepository) GetPrivateByUsers(ctx context.Context, user1Id, user2Id uint) (*domain.Private, error) {
	var private domain.Private

	if user1Id > user2Id {
		user1Id, user2Id = user2Id, user1Id
	}

	if err := p.dbRead.WithContext(ctx).Where("user1_id = ? AND user2_id = ?", user1Id, user2Id).First(&private).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}
	return &private, nil
}

func (p *privateRepository) GetPrivatesForUser(ctx context.Context, userId uint) ([]domain.Private, error) {
	var privates []domain.Private

	if err := p.dbRead.WithContext(ctx).
		Where("user1_id = ? OR user2_id = ?", userId, userId).
		Order("created_at DESC").
		Find(&privates).Error; err != nil {
		return nil, err
	}
	return privates, nil
}

func (p *privateRepository) CheckPrivateExists(ctx context.Context, user1Id, user2Id uint) (bool, error) {
	var count int64

	// Ensure canonical ordering
	if user1Id > user2Id {
		user1Id, user2Id = user2Id, user1Id
	}

	if err := p.dbRead.WithContext(ctx).
		Model(&domain.Private{}).
		Where("user1_id = ? AND user2_id = ?", user1Id, user2Id).
		Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func NewPrivateRepository(dbWrite, dbRead *gorm.DB) PrivateRepository {
	return &privateRepository{
		dbWrite: dbWrite,
		dbRead:  dbRead,
	}
}
