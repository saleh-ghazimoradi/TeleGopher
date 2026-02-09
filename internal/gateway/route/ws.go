package route

import (
	"github.com/saleh-ghazimoradi/TeleGopher/internal/gateway/handler"
	"net/http"
)

type WSRoute struct {
	wsHandler *handler.WSHandler
}

func (ws *WSRoute) WSRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /ws", ws.wsHandler.WebSocketHandler)
}

func NewWSRoute(wsHandler *handler.WSHandler) *WSRoute {
	return &WSRoute{
		wsHandler: wsHandler,
	}
}
