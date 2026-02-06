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
	protected := http.NewServeMux()

	protected.HandleFunc("GET /v1/conversations/privates/{private_id}", p.privateHandler.GetPrivate)
	protected.HandleFunc("POST /v1/conversations/privates", p.privateHandler.CreatePrivate)
	protected.HandleFunc("GET /v1/conversations", p.privateHandler.GetConversations)

	mux.Handle("/v1/conversations/", p.middleware.AuthMiddleware(protected))
}

func NewPrivateRoute(privateHandler *handler.PrivateHandler, middleware *middleware.Middleware) *PrivateRoute {
	return &PrivateRoute{
		privateHandler: privateHandler,
		middleware:     middleware,
	}
}
