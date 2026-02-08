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
	mux.HandleFunc("GET /v1/users/{id}",
		u.middleware.Authenticate(u.userHandler.GetUserById))
}

func NewUserRoute(userHandler *handler.UserHandler, middleware *middleware.Middleware) *UserRoute {
	return &UserRoute{
		userHandler: userHandler,
		middleware:  middleware,
	}
}
