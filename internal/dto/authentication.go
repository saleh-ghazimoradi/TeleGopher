package dto

import (
	"github.com/saleh-ghazimoradi/TeleGopher/internal/helper"
	"time"
)

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	Id        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type LoginResponse struct {
	User         *RegisterResponse `json:"user"`
	AccessToken  string            `json:"access_token"`
	RefreshToken string            `json:"refresh_token"`
}

func validateEmail(v *helper.Validator, email string) {
	v.Check(email != "", "email", "must be provided")
	v.Check(helper.Matches(email, helper.EmailRX), "email", "must be a valid email address")
}

func validatePassword(v *helper.Validator, password string) {
	v.Check(password != "", "password", "must be provided")
	v.Check(len(password) >= 8, "password", "must be at least 8 characters")
	v.Check(len(password) <= 72, "password", "must not be more than 72 characters")
}

func validateName(v *helper.Validator, name string) {
	v.Check(name != "", "name", "must be provided")
	v.Check(len(name) <= 500, "name", "must not be more than 500 characters")
}

func ValidateRegisterRequest(v *helper.Validator, req *RegisterRequest) {
	validateName(v, req.Name)
	validateEmail(v, req.Email)
	validatePassword(v, req.Password)
}

func ValidateLoginRequest(v *helper.Validator, req *LoginRequest) {
	validateEmail(v, req.Email)
	validatePassword(v, req.Password)
}
