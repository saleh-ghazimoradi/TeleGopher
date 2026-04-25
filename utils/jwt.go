package utils

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/saleh-ghazimoradi/TeleGopher/config"
	"time"
)

type Claims struct {
	UserId   uint   `json:"user_id"`
	Name     string `json:"name"`
	Platform string `json:"X-Platform"`
	jwt.RegisteredClaims
}

func GenerateToken(cfg *config.Config, userId uint, name, platform string) (string, error) {
	if platform != "web" && platform != "mobile" {
		return "", errors.New("invalid platform for token")
	}

	accessClaims := &Claims{
		UserId:   userId,
		Name:     name,
		Platform: platform,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(cfg.JWT.ExpiresIn)),
			Subject:   fmt.Sprint(userId),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessToken, err := at.SignedString([]byte(cfg.JWT.Secret))
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func ValidateToken(tokenString, secret string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
