package route

import (
	"github.com/saleh-ghazimoradi/TeleGopher/internal/gateway/handler"
	"github.com/saleh-ghazimoradi/TeleGopher/internal/gateway/middleware"
	"net/http"
)

type AuthRoute struct {
	middleware  *middleware.Middleware
	authHandler *handler.AuthHandler
}

func (a *AuthRoute) AuthRoutes(mux *http.ServeMux) {
	/*----------Public Routes----------*/
	mux.HandleFunc("POST /v1/auth/signup", a.authHandler.Signup)
	mux.HandleFunc("POST /v1/auth/login", a.authHandler.Login)
	mux.HandleFunc("POST /v1/auth/refresh-token", a.authHandler.RefreshToken)

	/*----------Private Routes----------*/
	mux.Handle("POST /v1/auth/logout", a.middleware.WrapAuth(a.authHandler.Logout))
	mux.Handle("GET /v1/auth/me", a.middleware.WrapAuth(a.authHandler.Me))
}

func NewAuthRoute(middleware *middleware.Middleware, authHandler *handler.AuthHandler) *AuthRoute {
	return &AuthRoute{
		middleware:  middleware,
		authHandler: authHandler,
	}
}
