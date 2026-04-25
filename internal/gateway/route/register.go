package route

import (
	"github.com/saleh-ghazimoradi/TeleGopher/internal/gateway/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
)

type RegisterRoute struct {
	Middleware       *middleware.Middleware
	HealthCheckRoute *HealthCheckRoute
	AuthRoute        *AuthRoute
	UserRoute        *UserRoute
	PrivateRoute     *PrivateRoute
	MessageRoute     *MessageRoute
	UploadFileRoute  *UploadFileRoute
	WsRoute          *WSRoute
}

type Options func(*RegisterRoute)

func WithMiddleware(middleware *middleware.Middleware) Options {
	return func(r *RegisterRoute) {
		r.Middleware = middleware
	}
}

func WithHealthCheckRoute(route *HealthCheckRoute) Options {
	return func(r *RegisterRoute) {
		r.HealthCheckRoute = route
	}
}

func WithAuthRoute(route *AuthRoute) Options {
	return func(r *RegisterRoute) {
		r.AuthRoute = route
	}
}

func WithUserRoute(route *UserRoute) Options {
	return func(r *RegisterRoute) {
		r.UserRoute = route
	}
}

func WithPrivateRoute(route *PrivateRoute) Options {
	return func(r *RegisterRoute) {
		r.PrivateRoute = route
	}
}

func WithMessageRoute(route *MessageRoute) Options {
	return func(r *RegisterRoute) {
		r.MessageRoute = route
	}
}

func WithUploadFileRoute(route *UploadFileRoute) Options {
	return func(r *RegisterRoute) {
		r.UploadFileRoute = route
	}
}

func WithWsRoute(route *WSRoute) Options {
	return func(r *RegisterRoute) {
		r.WsRoute = route
	}
}

func (r *RegisterRoute) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET /swagger/*any", httpSwagger.Handler(httpSwagger.URL("/docs/swagger.json")))

	// Serve static docs directory (for swagger.json and rapidoc.html)
	mux.Handle("GET /docs/", http.StripPrefix("/docs/", http.FileServer(http.Dir("./docs"))))

	// Redirect /api-docs to rapidoc
	mux.Handle("GET /api-docs", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./docs/rapidoc.html")
	}))

	r.HealthCheckRoute.HealthCheck(mux)
	r.AuthRoute.AuthRoutes(mux)
	r.UserRoute.UserRoutes(mux)
	r.PrivateRoute.PrivateRoutes(mux)
	r.MessageRoute.MessageRoutes(mux)
	r.UploadFileRoute.UploadFileRoutes(mux)
	r.WsRoute.WSRoutes(mux)
	return r.Middleware.Recover(r.Middleware.Logging(r.Middleware.CORS(mux)))
}

func NewRegisterRoute(opts ...Options) *RegisterRoute {
	r := &RegisterRoute{}
	for _, opt := range opts {
		opt(r)
	}

	return r
}
