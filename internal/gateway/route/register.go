package route

import "net/http"

type RegisterRoutes struct {
	healthRoute *HealthCheckRoute
}

type Options func(*RegisterRoutes)

func WithHealthCheckRoute(route *HealthCheckRoute) Options {
	return func(r *RegisterRoutes) {
		r.healthRoute = route
	}
}

func (r *RegisterRoutes) Register() http.Handler {
	mux := http.NewServeMux()

	r.healthRoute.Healthcheck(mux)

	return mux
}

func NewRegisterRoutes(opts ...Options) *RegisterRoutes {
	r := &RegisterRoutes{}
	for _, opt := range opts {
		opt(r)
	}
	return r
}
