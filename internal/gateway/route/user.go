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
	mux.Handle("POST /v1/users/phone", u.middleware.WrapAuth(u.userHandler.GetUserByPhone))
	mux.Handle("GET /v1/users", u.middleware.WrapAuth(u.userHandler.GetUsersByName))
}

func NewUserRoute(middleware *middleware.Middleware, userHandler *handler.UserHandler) *UserRoute {
	return &UserRoute{
		middleware:  middleware,
		userHandler: userHandler,
	}
}
