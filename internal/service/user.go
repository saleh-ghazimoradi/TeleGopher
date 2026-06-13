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
	GetUserByPhone(ctx context.Context, phone string) (*dto.UserResponse, error)
	GetUsersByName(ctx context.Context, name string) ([]*dto.UserResponse, error)
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

func (u *userService) GetUserByPhone(ctx context.Context, phone string) (*dto.UserResponse, error) {
	user, err := u.userRepository.GetUserByPhone(ctx, phone)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return u.toUserResponse(user), nil
}

func (u *userService) GetUsersByName(ctx context.Context, name string) ([]*dto.UserResponse, error) {
	users, err := u.userRepository.GetUsersByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}
	return u.toUsersResponse(users), nil
}

func (u *userService) toUserResponse(user *domain.User) *dto.UserResponse {
	return &dto.UserResponse{
		Id:        user.Id,
		Name:      user.Name,
		Phone:     user.Phone,
		CreatedAt: user.CreatedAt,
	}
}

func (u *userService) toUsersResponse(users []*domain.User) []*dto.UserResponse {
	userResponses := make([]*dto.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = u.toUserResponse(user)
	}
	return userResponses
}

func NewUserService(userRepository repository.UserRepository) UserService {
	return &userService{
		userRepository: userRepository,
	}
}
