package service

import (
	"context"
	"fmt"
	infra "github.com/saleh-ghazimoradi/TeleGopher/infra/TXManager"
	"github.com/saleh-ghazimoradi/TeleGopher/internal/domain"
	"github.com/saleh-ghazimoradi/TeleGopher/internal/dto"
	"github.com/saleh-ghazimoradi/TeleGopher/internal/repository"
)

type PrivateService interface {
	CreatePrivate(ctx context.Context, user1Id, user2Id int64) (*dto.PrivateResponse, error)
	GetPrivate(ctx context.Context, privateId, userId int64) (*dto.PrivateResponse, error)
	GetPrivateByUsers(ctx context.Context, user1Id, user2Id int64) (*dto.PrivateResponse, error)
	GetUserPrivates(ctx context.Context, userId int64) ([]dto.PrivateResponse, error)
}

type privateService struct {
	privateRepository repository.PrivateRepository
	tx                infra.TxManager
}

func (p *privateService) CreatePrivate(ctx context.Context, user1Id, user2Id int64) (*dto.PrivateResponse, error) {
	private, err := p.privateRepository.CreatePrivate(ctx, user1Id, user2Id)
	if err != nil {
		switch err {
		case repository.ErrSameUser:
			return nil, fmt.Errorf("cannot create private chat with same user")
		case repository.ErrPrivateAlreadyExists:
			return nil, fmt.Errorf("private chat already exists")
		default:
			return nil, fmt.Errorf("failed to create private: %w", err)
		}
	}

	return p.toPrivateResponse(private), nil
}

func (p *privateService) GetPrivate(ctx context.Context, privateId, userId int64) (*dto.PrivateResponse, error) {
	private, err := p.privateRepository.GetPrivateById(ctx, privateId)
	if err != nil {
		return nil, fmt.Errorf("failed to get private: %w", err)
	}

	if private.User1Id != userId && private.User2Id != userId {
		return nil, fmt.Errorf("unauthorized access to private chat")
	}

	return p.toPrivateResponse(private), nil
}

func (p *privateService) GetPrivateByUsers(ctx context.Context, user1Id, user2Id int64) (*dto.PrivateResponse, error) {
	private, err := p.privateRepository.GetPrivateByUsers(ctx, user1Id, user2Id)
	if err != nil {
		return nil, fmt.Errorf("failed to get private by users: %w", err)
	}

	if private.User1Id != user1Id && private.User2Id != user2Id {
		return nil, fmt.Errorf("unauthorized access to private chat")
	}

	return p.toPrivateResponse(private), nil
}

func (p *privateService) GetUserPrivates(ctx context.Context, userId int64) ([]dto.PrivateResponse, error) {
	privates, err := p.privateRepository.GetPrivateForUser(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to get privates: %w", err)
	}

	response := make([]dto.PrivateResponse, len(privates))
	for i, private := range privates {
		response[i] = *p.toPrivateResponse(&private)
	}

	return response, nil
}

func (p *privateService) toPrivateResponse(private *domain.Private) *dto.PrivateResponse {
	return &dto.PrivateResponse{
		Id:        private.Id,
		User1Id:   private.User1Id,
		User2Id:   private.User2Id,
		CreatedAt: private.CreatedAt,
	}
}

func NewPrivateService(privateRepository repository.PrivateRepository, tx infra.TxManager) PrivateService {
	return &privateService{
		privateRepository: privateRepository,
		tx:                tx,
	}
}
