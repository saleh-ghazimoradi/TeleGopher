package handler

import (
	"fmt"
	"github.com/saleh-ghazimoradi/TeleGopher/internal/domain"
	"github.com/saleh-ghazimoradi/TeleGopher/internal/dto"
	"github.com/saleh-ghazimoradi/TeleGopher/internal/helper"
	"github.com/saleh-ghazimoradi/TeleGopher/internal/service"
	"github.com/saleh-ghazimoradi/TeleGopher/utils"
	"net/http"
	"strings"
)

type AuthenticationHandler struct {
	errResponse           *helper.ErrResponse
	validator             *helper.Validator
	authenticationService service.AuthenticationService
}

func (a *AuthenticationHandler) Register(w http.ResponseWriter, r *http.Request) {
	var payload dto.RegisterRequest

	if err := helper.ReadJSON(w, r, &payload); err != nil {
		a.errResponse.BadRequestResponse(w, r, err)
		return
	}

	dto.ValidateRegisterRequest(a.validator, &payload)
	if !a.validator.Valid() {
		a.errResponse.InvalidCredentialsResponse(w, r)
		return
	}

	response, err := a.authenticationService.Register(r.Context(), &payload)
	if err != nil {
		switch {
		case err.Error() == "email already in use":
			a.errResponse.EditConflictResponse(w, r)
		default:
			a.errResponse.ServerErrorResponse(w, r, err)
		}
		return
	}

	if err := helper.WriteJSON(w, http.StatusCreated, helper.Envelope{"data": response}, nil); err != nil {
		a.errResponse.ServerErrorResponse(w, r, err)
	}
}

func (a *AuthenticationHandler) Login(w http.ResponseWriter, r *http.Request) {
	platform, err := a.extractPlatform(r)
	if err != nil {
		a.errResponse.BadRequestResponse(w, r, err)
		return
	}

	var payload dto.LoginRequest
	if err := helper.ReadJSON(w, r, &payload); err != nil {
		a.errResponse.BadRequestResponse(w, r, err)
		return
	}

	dto.ValidateLoginRequest(a.validator, &payload)
	if !a.validator.Valid() {
		a.errResponse.InvalidCredentialsResponse(w, r)
		return
	}

	response, err := a.authenticationService.Login(r.Context(), &payload, platform)
	if err != nil {
		switch {
		case err.Error() == "invalid credentials":
			a.errResponse.InvalidCredentialsResponse(w, r)
		default:
			a.errResponse.ServerErrorResponse(w, r, err)
		}
		return
	}

	if err := helper.WriteJSON(w, http.StatusOK, helper.Envelope{"data": response}, nil); err != nil {
		a.errResponse.ServerErrorResponse(w, r, err)
	}
}

func (a *AuthenticationHandler) Logout(w http.ResponseWriter, r *http.Request) {
	userId, ok := utils.WithIdFromContext(r.Context())
	if !ok {
		a.errResponse.InvalidCredentialsResponse(w, r)
		return
	}

	platform, ok := utils.PlatformFromContext(r.Context())
	if !ok {
		a.errResponse.InvalidCredentialsResponse(w, r)
		return
	}

	if err := a.authenticationService.Logout(r.Context(), userId, domain.Platform(platform)); err != nil {
		a.errResponse.ServerErrorResponse(w, r, err)
		return
	}

	if err := helper.WriteJSON(w, http.StatusOK, helper.Envelope{"data": nil}, nil); err != nil {
		a.errResponse.ServerErrorResponse(w, r, err)
	}
}

func (a *AuthenticationHandler) extractPlatform(r *http.Request) (domain.Platform, error) {
	platform := strings.ToLower(strings.TrimSpace(r.Header.Get("X-Platform")))

	switch platform {
	case string(domain.PlatformWeb):
		return domain.PlatformWeb, nil
	case string(domain.PlatformMobile):
		return domain.PlatformMobile, nil
	default:
		return "", fmt.Errorf("invalid platform")
	}
}

func NewAuthenticationHandler(errResponse *helper.ErrResponse, validator *helper.Validator, authenticationService service.AuthenticationService) *AuthenticationHandler {
	return &AuthenticationHandler{
		errResponse:           errResponse,
		validator:             validator,
		authenticationService: authenticationService,
	}
}
