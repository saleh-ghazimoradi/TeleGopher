package route

import (
	"github.com/saleh-ghazimoradi/TeleGopher/internal/gateway/handler"
	"github.com/saleh-ghazimoradi/TeleGopher/internal/gateway/middleware"
	"net/http"
)

type PrivateRoute struct {
	privateHandler *handler.PrivateHandler
	middleware     *middleware.Middleware
}

func (p *PrivateRoute) PrivateRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /v1/conversations/privates",
		p.middleware.Authenticate(p.privateHandler.CreatePrivate))

	mux.HandleFunc("GET /v1/conversations/privates/{private_id}",
		p.middleware.Authenticate(p.privateHandler.GetPrivate))

	mux.HandleFunc("GET /v1/conversations",
		p.middleware.Authenticate(p.privateHandler.GetConversations))
}

func NewPrivateRoute(privateHandler *handler.PrivateHandler, middleware *middleware.Middleware) *PrivateRoute {
	return &PrivateRoute{
		privateHandler: privateHandler,
		middleware:     middleware,
	}
}
