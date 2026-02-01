package dto

import "github.com/saleh-ghazimoradi/TeleGopher/internal/helper"

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
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
