package service

import (
	"context"
	"errors"
	"fmt"
	infra "github.com/saleh-ghazimoradi/TeleGopher/infra/TXManager"
	"github.com/saleh-ghazimoradi/TeleGopher/internal/domain"
	"github.com/saleh-ghazimoradi/TeleGopher/internal/repository"
)

type PrivateService interface {
	CreatePrivate(ctx context.Context, user1Id, user2Id int64) (*domain.Private, error)
	GetPrivateById(ctx context.Context, id int64) (*domain.Private, error)
	GetPrivateByUsers(ctx context.Context, user1Id, user2Id int64) (*domain.Private, error)
	GetPrivatesForUser(ctx context.Context, userId int64) ([]domain.Private, error)
}

type privateService struct {
	privateRepository repository.PrivateRepository
	tx                infra.TxManager
}

func (p *privateService) CreatePrivate(ctx context.Context, user1Id, user2Id int64) (*domain.Private, error) {
	existing, err := p.privateRepository.GetPrivateByUsers(ctx, user1Id, user2Id)
	if err != nil && !errors.Is(err, repository.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed to check existing private: %w", err)
	}

	if existing != nil {
		return existing, nil
	}

	return p.privateRepository.CreatePrivate(ctx, user1Id, user2Id)
}

func (p *privateService) GetPrivateById(ctx context.Context, id int64) (*domain.Private, error) {
	return p.privateRepository.GetPrivateById(ctx, id)
}
func (p *privateService) GetPrivateByUsers(ctx context.Context, user1Id, user2Id int64) (*domain.Private, error) {
	return p.privateRepository.GetPrivateByUsers(ctx, user1Id, user2Id)
}

func (p *privateService) GetPrivatesForUser(ctx context.Context, userId int64) ([]domain.Private, error) {
	return p.privateRepository.GetPrivateForUser(ctx, userId)
}

func NewPrivateService(privateRepository repository.PrivateRepository, tx infra.TxManager) PrivateService {
	return &privateService{
		privateRepository: privateRepository,
		tx:                tx,
	}
}
