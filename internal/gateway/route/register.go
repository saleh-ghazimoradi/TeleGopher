package route

import (
	"github.com/saleh-ghazimoradi/TeleGopher/internal/gateway/middleware"
	"net/http"
)

type RegisterRoutes struct {
	healthRoute         *HealthCheckRoute
	authenticationRoute *AuthenticationRoute
	userRoute           *UserRoute
	privateRoute        *PrivateRoute
	messageRoute        *MessageRoute
	uploadFileRoute     *UploadFileRoute
	wsRoute             *WSRoute
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

func WithPrivateRoute(privateRoute *PrivateRoute) Options {
	return func(r *RegisterRoutes) {
		r.privateRoute = privateRoute
	}
}

func WithMessageRoute(messageRoute *MessageRoute) Options {
	return func(r *RegisterRoutes) {
		r.messageRoute = messageRoute
	}
}

func WithUploadFileRoute(uploadFileRoute *UploadFileRoute) Options {
	return func(r *RegisterRoutes) {
		r.uploadFileRoute = uploadFileRoute
	}
}

func WithWSRoute(wsRoute *WSRoute) Options {
	return func(r *RegisterRoutes) {
		r.wsRoute = wsRoute
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
	r.privateRoute.PrivateRoutes(mux)
	r.messageRoute.MessageRoutes(mux)
	r.uploadFileRoute.UploadFileRoutes(mux)
	r.wsRoute.WSRoutes(mux)
	return r.middleware.Recover(r.middleware.Logging(r.middleware.CORS(mux.ServeHTTP)))
}

func NewRegisterRoutes(opts ...Options) *RegisterRoutes {
	r := &RegisterRoutes{}
	for _, opt := range opts {
		opt(r)
	}
	return r
}
