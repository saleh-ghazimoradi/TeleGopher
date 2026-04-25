package route

import (
	"github.com/saleh-ghazimoradi/TeleGopher/internal/gateway/handler"
	"github.com/saleh-ghazimoradi/TeleGopher/internal/gateway/middleware"
	"net/http"
)

type PrivateRoute struct {
	middleware     *middleware.Middleware
	privateHandler *handler.PrivateHandler
}

func (p *PrivateRoute) PrivateRoutes(mux *http.ServeMux) {
	mux.Handle("POST /v1/conversations/privates", p.middleware.WrapAuth(p.privateHandler.CreatePrivate))
	mux.Handle("GET /v1/conversations/privates/{id}", p.middleware.WrapAuth(p.privateHandler.GetPrivateById))
	mux.Handle("GET /v1/conversations", p.middleware.WrapAuth(p.privateHandler.GetConversations))
}

func NewPrivateRoute(middleware *middleware.Middleware, privateHandler *handler.PrivateHandler) *PrivateRoute {
	return &PrivateRoute{
		middleware:     middleware,
		privateHandler: privateHandler,
	}
}
