package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/saleh-ghazimoradi/TeleGopher/internal/domain"
	"github.com/saleh-ghazimoradi/TeleGopher/internal/dto"
	"github.com/saleh-ghazimoradi/TeleGopher/internal/repository"
)

type PrivateService interface {
	CreatePrivate(ctx context.Context, user1Id, user2Id uint) (*dto.PrivateResponse, error)
	GetPrivateById(ctx context.Context, privateId, userId uint) (*dto.PrivateResponse, error)
	GetPrivatesForUser(ctx context.Context, userId uint) ([]dto.PrivateResponse, error)
}

type privateService struct {
	privateRepository repository.PrivateRepository
	userRepository    repository.UserRepository
}

func (p *privateService) CreatePrivate(ctx context.Context, user1Id, user2Id uint) (*dto.PrivateResponse, error) {
	if user1Id == user2Id {
		return nil, repository.ErrSameUser
	}

	if err := p.verifyUsersExist(ctx, user1Id, user2Id); err != nil {
		return nil, err
	}

	exists, err := p.privateRepository.CheckPrivateExists(ctx, user1Id, user2Id)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing private: %w", err)
	}
	if exists {
		return nil, repository.ErrPrivateAlreadyExists
	}

	private := &domain.Private{
		User1Id: user1Id,
		User2Id: user2Id,
	}

	if err := p.privateRepository.CreatePrivate(ctx, private); err != nil {
		return nil, fmt.Errorf("failed to create private: %w", err)
	}

	return p.toPrivateResponse(private), nil
}

func (p *privateService) GetPrivateById(ctx context.Context, privateId, userId uint) (*dto.PrivateResponse, error) {
	private, err := p.privateRepository.GetPrivateById(ctx, privateId)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return nil, err
		}
		return nil, fmt.Errorf("failed to get private: %w", err)
	}

	if !p.canAccessPrivate(private, userId) {
		return nil, errors.New("unauthorized to access this private conversation")
	}

	return p.toPrivateResponse(private), nil
}

func (p *privateService) GetPrivatesForUser(ctx context.Context, userId uint) ([]dto.PrivateResponse, error) {
	// Verify user exists
	if _, err := p.userRepository.GetUserById(ctx, userId); err != nil {
		return nil, repository.ErrRecordNotFound
	}

	privates, err := p.privateRepository.GetPrivatesForUser(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to get privates: %w", err)
	}

	responses := make([]dto.PrivateResponse, len(privates))
	for i, private := range privates {
		responses[i] = *p.toPrivateResponse(&private)
	}

	return responses, nil
}

func (p *privateService) validateUsers(ctx context.Context, user1Id, user2Id uint) error {
	if user1Id == user2Id {
		return repository.ErrSameUser
	}

	return p.verifyUsersExist(ctx, user1Id, user2Id)
}

func (p *privateService) verifyUsersExist(ctx context.Context, userIds ...uint) error {
	for _, userId := range userIds {
		if _, err := p.userRepository.GetUserById(ctx, userId); err != nil {
			return fmt.Errorf("user %d not found: %w", userId, err)
		}
	}
	return nil
}

func (p *privateService) canAccessPrivate(private *domain.Private, userId uint) bool {
	return private.User1Id == userId || private.User2Id == userId
}

func (p *privateService) toPrivateResponse(private *domain.Private) *dto.PrivateResponse {
	return &dto.PrivateResponse{
		Id:        private.Id,
		User1Id:   private.User1Id,
		User2Id:   private.User2Id,
		CreatedAt: private.CreatedAt,
	}
}

func NewPrivateService(privateRepository repository.PrivateRepository, userRepository repository.UserRepository) PrivateService {
	return &privateService{
		privateRepository: privateRepository,
		userRepository:    userRepository,
	}
}
