package dto

import (
	"github.com/saleh-ghazimoradi/TeleGopher/internal/domain"
	"github.com/saleh-ghazimoradi/TeleGopher/internal/helper"
)

func validateName(v *helper.Validator, name string) {
	v.Check(helper.NotBlank(name), "name", "Name must be provided")
	v.Check(helper.MaxChars(name, 100), "name", "Name must be less than 100 characters")
}

func validateEmail(v *helper.Validator, email string) {
	v.Check(helper.NotBlank(email), "email", "Email must be provided")
	v.Check(helper.Matches(email, helper.EmailRX), "email", "Must be a valid email")
	v.Check(helper.MaxChars(email, 100), "email", "Email must be less than 100 characters")
}

func validatePassword(v *helper.Validator, password string) {
	v.Check(helper.NotBlank(password), "password", "Password must be provided")
	v.Check(helper.MinChars(password, 8), "password", "Password must be at least 8 characters")
	v.Check(helper.MaxChars(password, 72), "password", "Password must be less than 72 characters")
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

func ValidateRefreshToken(v *helper.Validator, req *RefreshTokenRequest) {
	v.Check(helper.NotBlank(req.RefreshToken), "refreshToken", "refreshToken must be provided")
}

func validatePrivateId(v *helper.Validator, privateId uint) {
	v.Check(privateId > 0, "private_id", "privateId must be provided")
}

func validateMessageType(v *helper.Validator, messageType string) {
	v.Check(helper.NotBlank(messageType), "message_type", "message type must be provided")
	v.Check(messageType == string(domain.MessageTypeText) || messageType == string(domain.MessageTypeFile) || messageType == string(domain.MessageTypeImage), "message_type", "Only Text, Image, and file are permitted")
}

func validateContent(v *helper.Validator, content string) {
	v.Check(helper.NotBlank(content), "content", "content must be provided")
	v.Check(helper.MaxChars("content", 5000), "content", "content must be less than 5000 characters")
}

func ValidateMessageRequest(v *helper.Validator, req *MessageRequest) {
	validatePrivateId(v, req.PrivateId)
	validateMessageType(v, req.MessageType)
	validateContent(v, req.Content)
}
