package route

import (
	"github.com/saleh-ghazimoradi/TeleGopher/internal/gateway/handler"
	"github.com/saleh-ghazimoradi/TeleGopher/internal/gateway/middleware"
	"net/http"
)

type MessageRoute struct {
	messageHandler *handler.MessageHandler
	middleware     *middleware.Middleware
}

func (m *MessageRoute) MessageRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /v1/messages",
		m.middleware.Authenticate(m.messageHandler.SendMessage))

	mux.HandleFunc("GET /v1/messages/{id}",
		m.middleware.Authenticate(m.messageHandler.GetMessage))

	mux.HandleFunc("GET /v1/conversations/privates/{private_id}/messages",
		m.middleware.Authenticate(m.messageHandler.GetPrivateMessages))

	mux.HandleFunc("PATCH /v1/messages/{id}/read",
		m.middleware.Authenticate(m.messageHandler.MarkMessageAsRead))

	mux.HandleFunc("PATCH /v1/messages/{id}/delivered",
		m.middleware.Authenticate(m.messageHandler.MarkMessageAsDelivered))
}

func NewMessageRoute(messageHandler *handler.MessageHandler, middleware *middleware.Middleware) *MessageRoute {
	return &MessageRoute{
		messageHandler: messageHandler,
		middleware:     middleware,
	}
}
