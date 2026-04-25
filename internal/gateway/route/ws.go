package route

import (
	"github.com/saleh-ghazimoradi/TeleGopher/internal/gateway/handler"
	"net/http"
)

type WSRoute struct {
	wsHandler *handler.WebSocketHandler
}

func (ws *WSRoute) WSRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /ws", ws.wsHandler.HandleWebsocket)
}

func NewWSRoute(wsHandler *handler.WebSocketHandler) *WSRoute {
	return &WSRoute{
		wsHandler: wsHandler,
	}
}
