package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/saleh-ghazimoradi/TeleGopher/config"
	infra "github.com/saleh-ghazimoradi/TeleGopher/infra/TXManager"
	"github.com/saleh-ghazimoradi/TeleGopher/internal/domain"
	"github.com/saleh-ghazimoradi/TeleGopher/internal/dto"
	"github.com/saleh-ghazimoradi/TeleGopher/internal/repository"
	"github.com/saleh-ghazimoradi/TeleGopher/utils"
)

type AuthenticationService interface {
	Register(ctx context.Context, req *dto.RegisterRequest) (*dto.RegisterResponse, error)
	Login(ctx context.Context, req *dto.LoginRequest, platform domain.Platform) (*dto.LoginResponse, error)
	GetUserById(ctx context.Context, userId int64) (*dto.RegisterResponse, error)
	GetUserByRefreshToken(ctx context.Context, req *dto.RefreshTokenRequest, platform domain.Platform) (*dto.RegisterResponse, error)
	RefreshToken(ctx context.Context, req *dto.RefreshTokenRequest, platform domain.Platform) (*dto.RefreshTokenResponse, error)
	Logout(ctx context.Context, userId int64, platform domain.Platform) error
}

type authenticationService struct {
	cfg            *config.Config
	userRepository repository.UserRepository
	tx             infra.TxManager
}

func (a *authenticationService) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.RegisterResponse, error) {
	var user *domain.User

	err := a.tx.WithTransaction(ctx, func(tx *sql.Tx) error {

		existing, err := a.userRepository.WithTx(tx).GetUserByEmail(ctx, req.Email)
		if err != nil && !errors.Is(err, repository.ErrRecordNotFound) {
			return fmt.Errorf("checking email: %w", err)
		}

		if existing != nil {
			return repository.ErrDuplicateEmail
		}

		hashedPassword, err := utils.HashPassword(req.Password)
		if err != nil {
			return fmt.Errorf("hashing password: %w", err)
		}

		user = &domain.User{
			Name:     req.Name,
			Email:    req.Email,
			Password: hashedPassword,
		}

		if err := a.userRepository.WithTx(tx).CreateUser(ctx, user); err != nil {
			return fmt.Errorf("creating user: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return a.toRegisterResponse(user), nil
}

func (a *authenticationService) Login(ctx context.Context, req *dto.LoginRequest, platform domain.Platform) (*dto.LoginResponse, error) {
	var user *domain.User
	var accessToken, refreshToken string

	err := a.tx.WithTransaction(ctx, func(tx *sql.Tx) error {
		var err error

		user, err = a.userRepository.WithTx(tx).GetUserByEmail(ctx, req.Email)
		if err != nil {
			if errors.Is(err, repository.ErrRecordNotFound) {
				return fmt.Errorf("invalid credentials")
			}
			return err
		}

		if !utils.CheckPassword(user.Password, req.Password) {
			return fmt.Errorf("invalid credentials")
		}

		accessToken, err = utils.GenerateToken(a.cfg.JWT.Expire, a.cfg.JWT.Secret, user.Id, user.Name, string(platform))
		if err != nil {
			return err
		}

		refreshToken, err = utils.GenerateRefreshToken()
		if err != nil {
			return err
		}

		return a.userRepository.WithTx(tx).UpdateRefreshToken(ctx, user.Id, platform, refreshToken)
	})

	if err != nil {
		return nil, err
	}

	return a.toLoginResponse(user, accessToken, refreshToken), nil
}

func (a *authenticationService) RefreshToken(ctx context.Context, req *dto.RefreshTokenRequest, platform domain.Platform) (*dto.RefreshTokenResponse, error) {
	var user *domain.User
	var accessToken, newRefreshToken string

	err := a.tx.WithTransaction(ctx, func(tx *sql.Tx) error {
		var err error

		user, err = a.userRepository.WithTx(tx).GetUserByRefreshToken(ctx, req.RefreshToken, platform)
		if err != nil {
			if errors.Is(err, repository.ErrRecordNotFound) {
				return fmt.Errorf("invalid refresh token")
			}
			return err
		}

		accessToken, err = utils.GenerateToken(a.cfg.JWT.Expire, a.cfg.JWT.Secret, user.Id, user.Name, string(platform))
		if err != nil {
			return err
		}

		newRefreshToken, err = utils.GenerateRefreshToken()
		if err != nil {
			return err
		}

		return a.userRepository.WithTx(tx).UpdateRefreshToken(ctx, user.Id, platform, newRefreshToken)
	})

	if err != nil {
		return nil, err
	}

	return a.toRefreshTokenResponse(accessToken, newRefreshToken), nil
}

func (a *authenticationService) GetUserByRefreshToken(ctx context.Context, req *dto.RefreshTokenRequest, platform domain.Platform) (*dto.RegisterResponse, error) {
	user, err := a.userRepository.GetUserByRefreshToken(ctx, req.RefreshToken, platform)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return nil, fmt.Errorf("invalid refresh token")
		}
		return nil, err
	}
	return a.toRegisterResponse(user), nil
}

func (a *authenticationService) Logout(ctx context.Context, userId int64, platform domain.Platform) error {
	return a.tx.WithTransaction(ctx, func(tx *sql.Tx) error {
		return a.userRepository.WithTx(tx).DeleteRefreshToken(ctx, userId, platform)
	})
}

func (a *authenticationService) toRegisterResponse(user *domain.User) *dto.RegisterResponse {
	return &dto.RegisterResponse{
		Id:        user.Id,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}
}

func (a *authenticationService) toLoginResponse(user *domain.User, accessToken, refreshToken string) *dto.LoginResponse {
	return &dto.LoginResponse{
		User:         a.toRegisterResponse(user),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
}

func (a *authenticationService) toRefreshTokenResponse(accessToken, refreshToken string) *dto.RefreshTokenResponse {
	return &dto.RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
}

func NewAuthenticationService(cfg *config.Config, userRepository repository.UserRepository, tx infra.TxManager) AuthenticationService {
	return &authenticationService{
		cfg:            cfg,
		userRepository: userRepository,
		tx:             tx,
	}
}
