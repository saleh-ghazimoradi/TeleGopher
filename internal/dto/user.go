package dto

import (
	"time"
)

type RegisterRequest struct {
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	Id        uint      `json:"id"`
	Name      string    `json:"name"`
	Phone     string    `json:"phone"`
	CreatedAt time.Time `json:"created_at"`
}

type LoginRequest struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type LoginResponse struct {
	User         *RegisterResponse
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type UserResponse struct {
	Id        uint      `json:"id"`
	Name      string    `json:"name"`
	Phone     string    `json:"phone"`
	CreatedAt time.Time `json:"created_at"`
}

type UserByPhoneRequest struct {
	Phone string `json:"phone"`
}
