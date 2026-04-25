package service

import (
	"context"
	"fmt"
	"github.com/saleh-ghazimoradi/TeleGopher/internal/domain"
	"github.com/saleh-ghazimoradi/TeleGopher/internal/dto"
	"github.com/saleh-ghazimoradi/TeleGopher/internal/repository"
)

type UserService interface {
	GetUserById(ctx context.Context, id uint) (*dto.UserResponse, error)
}

type userService struct {
	userRepository repository.UserRepository
}

func (u *userService) GetUserById(ctx context.Context, id uint) (*dto.UserResponse, error) {
	user, err := u.userRepository.GetUserById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return u.toUserResponse(user), nil
}

func (u *userService) toUserResponse(user *domain.User) *dto.UserResponse {
	return &dto.UserResponse{
		Id:        user.Id,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}
}

func NewUserService(userRepository repository.UserRepository) UserService {
	return &userService{
		userRepository: userRepository,
	}
}
