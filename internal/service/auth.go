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
	RefreshToken(ctx context.Context, input *dto.RefreshTokenRequest, platform string) (*dto.RefreshTokenResponse, error)
	GetUserById(ctx context.Context, id uint) (*dto.UserResponse, error)
}

type authService struct {
	userRepository repository.UserRepository
	cfg            *config.Config
}

func (a *authService) Register(ctx context.Context, input *dto.RegisterRequest) (*dto.RegisterResponse, error) {

	if _, err := a.userRepository.GetUserByPhone(ctx, input.Phone); err == nil {
		return nil, repository.ErrPhoneNumberExists
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

	user, err := a.userRepository.GetUserByPhone(ctx, input.Phone)
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
	if err := a.userRepository.UpdateRefreshToken(ctx, user.Id, refreshToken, platform, user.Version); err != nil {
		return nil, err
	}
	return a.toLoginResponse(user, accessToken, refreshToken), nil
}

func (a *authService) Logout(ctx context.Context, userId uint, platform string) error {
	return a.userRepository.DeleteRefreshToken(ctx, userId, platform)
}

func (a *authService) RefreshToken(ctx context.Context, input *dto.RefreshTokenRequest, platform string) (*dto.RefreshTokenResponse, error) {
	user, err := a.userRepository.GetValidUserByRefreshToken(ctx, input.RefreshToken, platform, a.cfg.JWT.RefreshTokenExpires)
	if err != nil {
		return nil, errors.New("invalid or expired refresh token")
	}

	accessToken, err := utils.GenerateToken(a.cfg, user.Id, user.Name, platform)
	if err != nil {
		return nil, err
	}
	newRefreshToken, err := utils.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	// Use the version from the fetched user to prevent concurrent overwrites
	if err := a.userRepository.UpdateRefreshToken(ctx, user.Id, newRefreshToken, platform, user.Version); err != nil {
		if errors.Is(err, repository.ErrRefreshTokenReused) {
			// Optional: delete all tokens to force re-login if replay detected
			// a.userRepository.DeleteRefreshToken(ctx, user.Id, platform)
			return nil, errors.New("refresh token already used – possible replay attack")
		}
		return nil, err
	}
	return a.toRefreshToken(accessToken, newRefreshToken), nil
}

func (a *authService) GetUserById(ctx context.Context, id uint) (*dto.UserResponse, error) {
	user, err := a.userRepository.GetUserById(ctx, id)
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
		Phone:    input.Phone,
		Password: hashedPassword,
	}, nil
}

func (a *authService) toRegisterResponse(user *domain.User) *dto.RegisterResponse {
	return &dto.RegisterResponse{
		Id:        user.Id,
		Name:      user.Name,
		Phone:     user.Phone,
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
		Phone:     user.Phone,
		CreatedAt: user.CreatedAt,
	}
}

func NewAuthService(userRepository repository.UserRepository, cfg *config.Config) AuthService {
	return &authService{
		userRepository: userRepository,
		cfg:            cfg,
	}
}
