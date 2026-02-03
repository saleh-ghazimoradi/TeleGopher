package service

import (
	"context"
	"fmt"
	"github.com/saleh-ghazimoradi/TeleGopher/internal/domain"
	"github.com/saleh-ghazimoradi/TeleGopher/internal/dto"
	"github.com/saleh-ghazimoradi/TeleGopher/internal/repository"
)

type UserService interface {
	GetUserById(ctx context.Context, userId int64) (*dto.RegisterResponse, error)
}

type userService struct {
	userRepository repository.UserRepository
}

func (u *userService) GetUserById(ctx context.Context, userId int64) (*dto.RegisterResponse, error) {
	user, err := u.userRepository.GetUserById(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return u.toRegisterResponse(user), nil
}

func (u *userService) toRegisterResponse(user *domain.User) *dto.RegisterResponse {
	return &dto.RegisterResponse{
		Id:        user.Id,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}
}

func NewUserService(repository repository.UserRepository) UserService {
	return &userService{
		userRepository: repository,
	}
}
