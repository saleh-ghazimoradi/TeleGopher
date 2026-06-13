package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/saleh-ghazimoradi/TeleGopher/internal/domain"
	"gorm.io/gorm"
	"time"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *domain.User) error
	GetUserById(ctx context.Context, id uint) (*domain.User, error)
	GetUserByPhone(ctx context.Context, phone string) (*domain.User, error)
	GetUsersByName(ctx context.Context, name string) ([]*domain.User, error)
	UpdateUser(ctx context.Context, user *domain.User) error
	DeleteUser(ctx context.Context, id uint) error

	GetValidUserByRefreshToken(ctx context.Context, refreshToken string, platform string, maxAge time.Duration) (*domain.User, error)
	GetUserByRefreshToken(ctx context.Context, refreshToken string, platform string) (*domain.User, error)
	UpdateRefreshToken(ctx context.Context, userId uint, refreshToken string, platform string, oldVersion uint) error
	DeleteRefreshToken(ctx context.Context, userId uint, platform string) error
}

type userRepository struct {
	dbWrite *gorm.DB
	dbRead  *gorm.DB
}

func (u *userRepository) CreateUser(ctx context.Context, user *domain.User) error {
	return u.dbWrite.WithContext(ctx).Create(&user).Error
}

func (u *userRepository) GetUserById(ctx context.Context, id uint) (*domain.User, error) {
	var user domain.User
	if err := u.dbRead.WithContext(ctx).First(&user, id).Error; err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &user, nil
}

func (u *userRepository) GetUserByPhone(ctx context.Context, phone string) (*domain.User, error) {
	var user domain.User
	if err := u.dbRead.WithContext(ctx).Where("phone = ?", phone).First(&user).Error; err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &user, nil
}

func (u *userRepository) GetUsersByName(ctx context.Context, name string) ([]*domain.User, error) {
	var users []*domain.User

	if err := u.dbRead.WithContext(ctx).Where("name ILIKE ?", "%"+name+"%").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (u *userRepository) UpdateUser(ctx context.Context, user *domain.User) error {
	return u.dbWrite.WithContext(ctx).Save(&user).Error
}

func (u *userRepository) DeleteUser(ctx context.Context, id uint) error {
	return u.dbWrite.WithContext(ctx).Delete(&domain.User{}, id).Error
}

func (u *userRepository) GetUserByRefreshToken(ctx context.Context, refreshToken string, platform string) (*domain.User, error) {
	var user domain.User

	condition := u.buildRefreshTokenCondition(platform, refreshToken)
	if condition == nil {
		return nil, errors.New("invalid platform")
	}

	if err := u.dbRead.WithContext(ctx).Where(condition).First(&user).Error; err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

func (u *userRepository) GetValidUserByRefreshToken(ctx context.Context, refreshToken string, platform string, maxAge time.Duration) (*domain.User, error) {
	var user domain.User
	var field, atField string
	if platform == "web" {
		field = "refresh_token_web = ?"
		atField = "refresh_token_web_at"
	} else if platform == "mobile" {
		field = "refresh_token_mobile = ?"
		atField = "refresh_token_mobile_at"
	} else {
		return nil, errors.New("invalid platform")
	}

	err := u.dbRead.WithContext(ctx).
		Where(field, refreshToken).
		Where(atField+" IS NOT NULL AND "+atField+" > ?", time.Now().Add(-maxAge)).
		First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrInvalidOrExpiredToken
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// repository/user.go
func (u *userRepository) UpdateRefreshToken(ctx context.Context, userId uint, refreshToken string, platform string, oldVersion uint) error {
	updates := u.buildRefreshTokenUpdate(platform, refreshToken)
	if updates == nil {
		return errors.New("invalid platform")
	}
	updates["version"] = gorm.Expr("version + 1")

	// Only update if version hasn't changed since we read it
	res := u.dbWrite.WithContext(ctx).
		Model(&domain.User{}).
		Where("id = ? AND version = ?", userId, oldVersion).
		Updates(updates)

	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		// Another request already updated the token → refresh token reused or replay
		return ErrRefreshTokenReused
	}
	return nil
}

// Add a sentinel error
var ErrRefreshTokenReused = errors.New("refresh token already used")

func (u *userRepository) DeleteRefreshToken(ctx context.Context, userId uint, platform string) error {
	updates := u.buildRefreshTokenDeletion(platform)
	if updates == nil {
		return fmt.Errorf("invalid platform: %s", platform)
	}

	updates["version"] = gorm.Expr("version + 1")

	return u.dbWrite.WithContext(ctx).Model(&domain.User{}).Where("id = ?", userId).Updates(updates).Error
}

func (u *userRepository) buildRefreshTokenCondition(platform string, refreshToken string) map[string]any {
	if platform == "web" {
		return map[string]any{
			"refresh_token_web": refreshToken,
		}
	}
	if platform == "mobile" {
		return map[string]any{
			"refresh_token_mobile": refreshToken,
		}
	}
	return nil
}

func (u *userRepository) buildRefreshTokenUpdate(platform string, refreshToken string) map[string]any {
	if platform == "web" {
		return map[string]any{
			"refresh_token_web":    refreshToken,
			"refresh_token_web_at": time.Now(),
		}
	}
	if platform == "mobile" {
		return map[string]any{
			"refresh_token_mobile":    refreshToken,
			"refresh_token_mobile_at": time.Now(),
		}
	}

	return nil
}

func (u *userRepository) buildRefreshTokenDeletion(platform string) map[string]any {
	if platform == "web" {
		return map[string]any{
			"refresh_token_web":    nil,
			"refresh_token_web_at": nil,
		}
	}
	if platform == "mobile" {
		return map[string]any{
			"refresh_token_mobile":    nil,
			"refresh_token_mobile_at": nil,
		}
	}
	return nil
}

func NewUserRepository(dbWrite, dbRead *gorm.DB) UserRepository {
	return &userRepository{
		dbWrite: dbWrite,
		dbRead:  dbRead,
	}
}
