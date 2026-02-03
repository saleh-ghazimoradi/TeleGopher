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
	mux.HandleFunc("POST /v1/auth/register", a.authenticationHandler.Register)
	mux.HandleFunc("POST /v1/auth/login", a.authenticationHandler.Login)
	mux.HandleFunc("POST /v1/auth/refresh-token", a.authenticationHandler.RefreshToken)

	protectedMux := http.NewServeMux()
	protectedMux.HandleFunc("POST /v1/auth/logout", a.authenticationHandler.Logout)
	protectedMux.HandleFunc("POST /v1/auth/me", a.authenticationHandler.GetCurrentUser)
	mux.Handle("/v1/auth/", a.middleware.AuthMiddleware(protectedMux))
}

func NewAuthenticationRoute(authenticationHandler *handler.AuthenticationHandler, middleware *middleware.Middleware) *AuthenticationRoute {
	return &AuthenticationRoute{
		authenticationHandler: authenticationHandler,
		middleware:            middleware,
	}
}
