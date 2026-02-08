package route

import (
	"github.com/saleh-ghazimoradi/TeleGopher/internal/gateway/handler"
	"github.com/saleh-ghazimoradi/TeleGopher/internal/gateway/middleware"
	"net/http"
)

type AuthenticationRoute struct {
	authenticationHandler *handler.AuthenticationHandler
	middleware            *middleware.Middleware
}

func (a *AuthenticationRoute) AuthenticationRoutes(mux *http.ServeMux) {
	// public routes — no auth
	mux.HandleFunc("POST /v1/auth/register", a.authenticationHandler.Register)
	mux.HandleFunc("POST /v1/auth/login", a.authenticationHandler.Login)
	mux.HandleFunc("POST /v1/auth/refresh-token", a.authenticationHandler.RefreshToken)

	// protected routes — one by one (cleaner than sub-mux in most cases)
	mux.HandleFunc("POST /v1/auth/logout", a.middleware.Authenticate(a.authenticationHandler.Logout))
	mux.HandleFunc("POST /v1/auth/me", a.middleware.Authenticate(a.authenticationHandler.GetCurrentUser))
}

func NewAuthenticationRoute(authenticationHandler *handler.AuthenticationHandler, middleware *middleware.Middleware) *AuthenticationRoute {
	return &AuthenticationRoute{
		authenticationHandler: authenticationHandler,
		middleware:            middleware,
	}
}
