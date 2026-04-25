package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/saleh-ghazimoradi/TeleGopher/config"
	"github.com/saleh-ghazimoradi/TeleGopher/internal/domain"
	"github.com/saleh-ghazimoradi/TeleGopher/internal/dto"
	"github.com/saleh-ghazimoradi/TeleGopher/internal/repository"
	"github.com/saleh-ghazimoradi/TeleGopher/utils"
)

type AuthService interface {
	Register(ctx context.Context, input *dto.RegisterRequest) (*dto.RegisterResponse, error)
	Login(ctx context.Context, input *dto.LoginRequest, platform string) (*dto.LoginResponse, error)
	Logout(ctx context.Context, userId uint, platform string) error
	GetUserByRefreshToken(ctx context.Context, input *dto.RefreshTokenRequest, platform string) (*dto.UserResponse, error)
	RefreshToken(ctx context.Context, input *dto.RefreshTokenRequest, platform string) (*dto.RefreshTokenResponse, error)
}

type authService struct {
	userRepository repository.UserRepository
	cfg            *config.Config
}

func (a *authService) Register(ctx context.Context, input *dto.RegisterRequest) (*dto.RegisterResponse, error) {
	if _, err := a.userRepository.GetUserByEmail(ctx, input.Email); err == nil {
		return nil, repository.ErrEmailExists
	}

	user, err := a.toUserDomain(input)
	if err != nil {
		return nil, err
	}

	if err := a.userRepository.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	return a.toRegisterResponse(user), nil
}

func (a *authService) Login(ctx context.Context, input *dto.LoginRequest, platform string) (*dto.LoginResponse, error) {
	user, err := a.userRepository.GetUserByEmail(ctx, input.Email)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrRecordNotFound):
			return nil, errors.New("invalid credentials")
		default:
			return nil, err
		}
	}

	if !utils.CheckPasswordHash(user.Password, input.Password) {
		return nil, errors.New("invalid credentials")
	}

	accessToken, err := utils.GenerateToken(a.cfg, user.Id, user.Name, platform)
	if err != nil {
		return nil, err
	}

	refreshToken, err := utils.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	if err := a.userRepository.UpdateRefreshToken(ctx, user.Id, refreshToken, platform); err != nil {
		return nil, err
	}

	return a.toLoginResponse(user, accessToken, refreshToken), nil
}

func (a *authService) Logout(ctx context.Context, userId uint, platform string) error {
	return a.userRepository.DeleteRefreshToken(ctx, userId, platform)
}

func (a *authService) RefreshToken(ctx context.Context, input *dto.RefreshTokenRequest, platform string) (*dto.RefreshTokenResponse, error) {
	user, err := a.userRepository.GetUserByRefreshToken(ctx, input.RefreshToken, platform)
	if err != nil {
		return nil, err
	}

	accessToken, err := utils.GenerateToken(a.cfg, user.Id, user.Name, platform)
	if err != nil {
		return nil, err
	}

	refreshToken, err := utils.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	if err := a.userRepository.UpdateRefreshToken(ctx, user.Id, refreshToken, platform); err != nil {
		return nil, err
	}

	return a.toRefreshToken(accessToken, refreshToken), nil
}

func (a *authService) GetUserByRefreshToken(ctx context.Context, input *dto.RefreshTokenRequest, platform string) (*dto.UserResponse, error) {
	user, err := a.userRepository.GetUserByRefreshToken(ctx, input.RefreshToken, platform)
	if err != nil {
		return nil, err
	}

	return a.toUserResponse(user), nil
}

func (a *authService) toUserDomain(input *dto.RegisterRequest) (*domain.User, error) {
	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}
	return &domain.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: hashedPassword,
	}, nil
}

func (a *authService) toRegisterResponse(user *domain.User) *dto.RegisterResponse {
	return &dto.RegisterResponse{
		Id:        user.Id,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}
}

func (a *authService) toLoginResponse(user *domain.User, accessToken, refreshToken string) *dto.LoginResponse {
	return &dto.LoginResponse{
		User:         a.toRegisterResponse(user),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
}

func (a *authService) toRefreshToken(accessToken, refreshToken string) *dto.RefreshTokenResponse {
	return &dto.RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
}

func (a *authService) toUserResponse(user *domain.User) *dto.UserResponse {
	return &dto.UserResponse{
		Id:        user.Id,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}
}

func NewAuthService(userRepository repository.UserRepository, cfg *config.Config) AuthService {
	return &authService{
		userRepository: userRepository,
		cfg:            cfg,
	}
}
