package route

import (
	"github.com/saleh-ghazimoradi/TeleGopher/internal/gateway/handler"
	"github.com/saleh-ghazimoradi/TeleGopher/internal/gateway/middleware"
	"net/http"
)

type UserRoute struct {
	middleware  *middleware.Middleware
	userHandler *handler.UserHandler
}

func (u *UserRoute) UserRoutes(mux *http.ServeMux) {
	mux.Handle("GET /v1/users/{id}", u.middleware.WrapAuth(u.userHandler.GetUserById))
}

func NewUserRoute(middleware *middleware.Middleware, userHandler *handler.UserHandler) *UserRoute {
	return &UserRoute{
		middleware:  middleware,
		userHandler: userHandler,
	}
}
