package route

import (
	"github.com/saleh-ghazimoradi/TeleGopher/internal/gateway/handler"
	"net/http"
)

type HealthCheckRoute struct {
	healthHandler *handler.HealthCheckHandler
}

func (h *HealthCheckRoute) Healthcheck(mux *http.ServeMux) {
	mux.HandleFunc("GET /v1/healthcheck", h.healthHandler.HealthCheck)
}

func NewHealthCheckRoute(healthHandler *handler.HealthCheckHandler) *HealthCheckRoute {
	return &HealthCheckRoute{healthHandler: healthHandler}
}
