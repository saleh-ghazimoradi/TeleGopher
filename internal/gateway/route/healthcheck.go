package route

import (
	"github.com/saleh-ghazimoradi/TeleGopher/internal/gateway/handler"
	"net/http"
)

type HealthCheckRoute struct {
	healthCheckHandler *handler.HealthCheckHandler
}

func (h *HealthCheckRoute) HealthCheck(mux *http.ServeMux) {
	mux.HandleFunc("GET /v1/healthcheck", h.healthCheckHandler.HealthCheck)
	mux.HandleFunc("GET /v1/ws-healthcheck", h.healthCheckHandler.WsHealthCheck)
}

func NewHealthCheckRoute(healthCheckHandler *handler.HealthCheckHandler) *HealthCheckRoute {
	return &HealthCheckRoute{
		healthCheckHandler: healthCheckHandler,
	}
}
