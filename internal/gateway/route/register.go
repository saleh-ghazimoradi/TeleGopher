package route

import (
	"github.com/saleh-ghazimoradi/TeleGopher/internal/gateway/middleware"
	"net/http"
)

type RegisterRoutes struct {
	healthRoute         *HealthCheckRoute
	authenticationRoute *AuthenticationRoute
	userRoute           *UserRoute
	middleware          *middleware.Middleware
}

type Options func(*RegisterRoutes)

func WithHealthCheckRoute(route *HealthCheckRoute) Options {
	return func(r *RegisterRoutes) {
		r.healthRoute = route
	}
}

func WithAuthenticationRoute(route *AuthenticationRoute) Options {
	return func(r *RegisterRoutes) {
		r.authenticationRoute = route
	}
}

func WithUserRoute(route *UserRoute) Options {
	return func(r *RegisterRoutes) {
		r.userRoute = route
	}
}

func WithMiddleware(middleware *middleware.Middleware) Options {
	return func(r *RegisterRoutes) {
		r.middleware = middleware
	}
}

func (r *RegisterRoutes) Register() http.Handler {
	mux := http.NewServeMux()
	r.healthRoute.Healthcheck(mux)
	r.authenticationRoute.AuthenticationRoutes(mux)
	r.userRoute.UserRoutes(mux)
	return r.middleware.RecoverPanic(r.middleware.LoggingMiddleware(r.middleware.CORSMiddleware(mux)))
}

func NewRegisterRoutes(opts ...Options) *RegisterRoutes {
	r := &RegisterRoutes{}
	for _, opt := range opts {
		opt(r)
	}
	return r
}
