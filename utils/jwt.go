package utils

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type Claims struct {
	UserId   int64  `json:"user_id"`
	Name     string `json:"name"`
	Platform string `json:"X-Platform"`
	jwt.RegisteredClaims
}

func GenerateToken(expiry time.Duration, secret string, userId int64, name string, platform string) (string, error) {
	if platform != "web" && platform != "mobile" {
		return "", errors.New("invalid platform")
	}
	claim := &Claims{
		UserId:   userId,
		Name:     name,
		Platform: platform,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
			Subject:   fmt.Sprint(userId),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateToken(tokenString, secret string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claim, ok := token.Claims.(*Claims); ok && token.Valid {
		return claim, nil
	}

	return nil, errors.New("invalid token")
}
