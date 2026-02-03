package route

import (
	"github.com/saleh-ghazimoradi/TeleGopher/internal/gateway/handler"
	"github.com/saleh-ghazimoradi/TeleGopher/internal/gateway/middleware"
	"net/http"
)

type UserRoute struct {
	userHandler *handler.UserHandler
	middleware  *middleware.Middleware
}

func (u *UserRoute) UserRoutes(mux *http.ServeMux) {
	protected := http.NewServeMux()
	protected.HandleFunc("GET /v1/users/{id}", u.userHandler.GetUserById)
	mux.Handle("/v1/users/", u.middleware.AuthMiddleware(protected))
}

func NewUserRoute(userHandler *handler.UserHandler, middleware *middleware.Middleware) *UserRoute {
	return &UserRoute{
		userHandler: userHandler,
		middleware:  middleware,
	}
}
