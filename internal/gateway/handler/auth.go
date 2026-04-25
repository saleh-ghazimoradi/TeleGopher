package handler

import (
	"fmt"
	"github.com/saleh-ghazimoradi/TeleGopher/internal/dto"
	"github.com/saleh-ghazimoradi/TeleGopher/internal/helper"
	"github.com/saleh-ghazimoradi/TeleGopher/internal/service"
	"github.com/saleh-ghazimoradi/TeleGopher/utils"
	"net/http"
)

type AuthHandler struct {
	authService service.AuthService
}

// Signup godoc
// @Summary      Signup a new user
// @Description  Create a new user account with email and password
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        request body dto.RegisterRequest true "User registration data"
// @Success      201 {object} helper.Response{data=dto.RegisterResponse} "User successfully created"
// @Failure      400 {object} helper.Response "Invalid request data or user already exists"
// @Failure      500 {object} helper.Response "Internal server error"
// @Router       /auth/signup [post]
func (a *AuthHandler) Signup(w http.ResponseWriter, r *http.Request) {
	var payload dto.RegisterRequest
	if err := helper.ReadJSON(w, r, &payload); err != nil {
		helper.BadRequestResponse(w, "Invalid given payload", err)
		return
	}

	v := helper.NewValidator()
	dto.ValidateRegisterRequest(v, &payload)
	if !v.Valid() {
		helper.FailedValidationResponse(w, "input's not valid")
		return
	}

	user, err := a.authService.Register(r.Context(), &payload)
	if err != nil {
		helper.InternalServerError(w, "failed to signup a user", err)
		return
	}

	helper.CreatedResponse(w, "user's successfully created", user)
}

// Login godoc
// @Summary      User login
// @Description  Authenticate user with email and password
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        request body dto.LoginRequest true "User login credentials"
// @Param        X-Platform header string true "Platform type (web or mobile)" Enums(web, mobile)
// @Success      200 {object} helper.Response{data=dto.LoginResponse} "Login successfully"
// @Failure      400 {object} helper.Response "Invalid platform or request data"
// @Failure      401 {object} helper.Response "Invalid credentials"
// @Failure      500 {object} helper.Response "Internal server error"
// @Router       /auth/login [post]
func (a *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	platform := r.Header.Get("X-Platform")
	if platform != "web" && platform != "mobile" {
		helper.BadRequestResponse(w, "Invalid platform", fmt.Errorf("platform must be web or mobile"))
		return
	}

	var payload dto.LoginRequest
	if err := helper.ReadJSON(w, r, &payload); err != nil {
		helper.BadRequestResponse(w, "Invalid given payload", err)
		return
	}

	v := helper.NewValidator()
	dto.ValidateLoginRequest(v, &payload)
	if !v.Valid() {
		helper.FailedValidationResponse(w, "input's not valid")
		return
	}

	login, err := a.authService.Login(r.Context(), &payload, platform)
	if err != nil {
		helper.InternalServerError(w, "failed to login", err)
		return
	}

	helper.SuccessResponse(w, "User Successfully logged in", login)
}

// Logout godoc
// @Summary      User logout
// @Description  Invalidate refresh token and logout user
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        X-Platform header string true "Platform type (web or mobile)" Enums(web, mobile)
// @Success      200 {object} helper.Response "Logout successful"
// @Failure      400 {object} helper.Response "Invalid platform"
// @Failure      401 {object} helper.Response "Unauthorized"
// @Failure      500 {object} helper.Response "Internal server error"
// @Router       /auth/logout [post]
func (a *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	userid, ok := utils.UserIdFromContext()
	userId, ok := utils.UserIdFromContext(r.Context())
	if !ok {
		helper.UnauthorizedResponse(w, "Unauthorized")
		return
	}

	platform, ok := utils.PlatformFromContext(r.Context())
	if !ok {
		helper.UnauthorizedResponse(w, "Unauthorized")
		return
	}

	if err := a.authService.Logout(r.Context(), userId, platform); err != nil {
		helper.InternalServerError(w, "failed to logout", err)
		return
	}

	helper.SuccessResponse(w, "User Successfully logged out", nil)
}

// RefreshToken godoc
// @Summary      Refresh access token
// @Description  Get a new access token using refresh token
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        X-Platform header string true "Platform type (web or mobile)" Enums(web, mobile)
// @Param        request body dto.RefreshTokenRequest true "Refresh token"
// @Success      200 {object} helper.Response{data=dto.RefreshTokenResponse} "Token refreshed successfully"
// @Failure      400 {object} helper.Response "Invalid request data"
// @Failure      401 {object} helper.Response "Invalid refresh token"
// @Failure      500 {object} helper.Response "Internal server error"
// @Router       /auth/refresh-token [post]
func (a *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	platform := r.Header.Get("X-Platform")
	if platform != "web" && platform != "mobile" {
		helper.BadRequestResponse(w, "Invalid platform", fmt.Errorf("platform must be web or mobile"))
		return
	}

	var payload dto.RefreshTokenRequest
	if err := helper.ReadJSON(w, r, &payload); err != nil {
		helper.BadRequestResponse(w, "Invalid given payload", err)
		return
	}

	v := helper.NewValidator()
	dto.ValidateRefreshToken(v, &payload)
	if !v.Valid() {
		helper.FailedValidationResponse(w, "input's not valid")
		return
	}

	refreshToken, err := a.authService.RefreshToken(r.Context(), &payload, platform)
	if err != nil {
		helper.InternalServerError(w, "failed to refresh token", err)
		return
	}

	helper.SuccessResponse(w, "Refresh Token Successfully", refreshToken)
}

// Me godoc
// @Summary      Get current authenticated user
// @Description  Get the currently authenticated user's information
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        X-Platform header string true "Platform type (web or mobile)" Enums(web, mobile)
// @Success      200 {object} helper.Response{data=dto.UserResponse} "User successfully fetched"
// @Failure      400 {object} helper.Response "Invalid platform"
// @Failure      401 {object} helper.Response "Unauthorized - Invalid or missing token"
// @Failure      404 {object} helper.Response "User not found"
// @Failure      500 {object} helper.Response "Internal server error"
// @Router       /auth/me [get]
func (a *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	platform := r.Header.Get("X-Platform")
	if platform != "web" && platform != "mobile" {
		helper.BadRequestResponse(w, "Invalid platform", fmt.Errorf("platform must be web or mobile"))
		return
	}

	var payload dto.RefreshTokenRequest
	if err := helper.ReadJSON(w, r, &payload); err != nil {
		helper.BadRequestResponse(w, "Invalid given payload", err)
		return
	}

	v := helper.NewValidator()
	dto.ValidateRefreshToken(v, &payload)
	if !v.Valid() {
		helper.FailedValidationResponse(w, "input's not valid")
		return
	}

	user, err := a.authService.GetUserByRefreshToken(r.Context(), &payload, platform)
	if err != nil {
		helper.InternalServerError(w, "failed to fetch user", err)
		return
	}

	helper.SuccessResponse(w, "User Successfully fetched", user)
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}
