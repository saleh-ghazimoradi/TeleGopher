package route

import (
	"github.com/saleh-ghazimoradi/TeleGopher/internal/gateway/handler"
	"github.com/saleh-ghazimoradi/TeleGopher/internal/gateway/middleware"
	"net/http"
)

type MessageRoute struct {
	middleware     *middleware.Middleware
	messageHandler *handler.MessageHandler
}

func (m *MessageRoute) MessageRoutes(mux *http.ServeMux) {
	mux.Handle("POST /v1/messages", m.middleware.WrapAuth(m.messageHandler.SendMessage))
	mux.Handle("GET /v1/messages/{id}", m.middleware.WrapAuth(m.messageHandler.GetMessage))
	mux.Handle("GET /v1/conversations/privates/{id}/messages", m.middleware.WrapAuth(m.messageHandler.GetPrivateMessages))
	mux.Handle("PATCH /v1/messages/{id}/read", m.middleware.WrapAuth(m.messageHandler.MarkMessageAsRead))
	mux.Handle("PATCH /v1/messages/{id}/delivered", m.middleware.WrapAuth(m.messageHandler.MarkMessageAsDelivered))
}

func NewMessageRoute(middleware *middleware.Middleware, messageHandler *handler.MessageHandler) *MessageRoute {
	return &MessageRoute{
		middleware:     middleware,
		messageHandler: messageHandler,
	}
}
