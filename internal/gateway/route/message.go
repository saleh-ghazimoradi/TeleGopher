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
	protected := http.NewServeMux()

	protected.HandleFunc("POST /v1/messages", m.messageHandler.SendMessage)
	protected.HandleFunc("GET /v1/messages/{id}", m.messageHandler.GetMessage)
	protected.HandleFunc("GET /v1/conversations/privates/{private_id}/messages", m.messageHandler.GetPrivateMessages)
	protected.HandleFunc("PATCH /v1/messages/{id}/read", m.messageHandler.MarkMessageAsRead)
	protected.HandleFunc("PATCH /v1/messages/{id}/delivered", m.messageHandler.MarkMessageAsDelivered)

	mux.Handle("/v1/messages/", m.middleware.AuthMiddleware(protected))
	mux.Handle("/v1/conversations/privates/", m.middleware.AuthMiddleware(protected))
}

func NewMessageRoute(messageHandler *handler.MessageHandler, middleware *middleware.Middleware) *MessageRoute {
	return &MessageRoute{
		messageHandler: messageHandler,
		middleware:     middleware,
	}
}
