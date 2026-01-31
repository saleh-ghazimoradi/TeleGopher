package handler

import (
	"github.com/saleh-ghazimoradi/TeleGopher/config"
	"github.com/saleh-ghazimoradi/TeleGopher/internal/dto"
	"github.com/saleh-ghazimoradi/TeleGopher/internal/helper"
	"net/http"
)

type HealthCheckHandler struct {
	cfg *config.Config
}

func (h *HealthCheckHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	payload := &dto.Health{
		SystemInfo: map[string]any{
			"version":     h.cfg.Application.Version,
			"environment": h.cfg.Application.Environment,
		},
	}

	if err := helper.SuccessResponse(w, "I'm breathing", payload, nil); err != nil {
		return
	}
}

func NewHealthCheckHandler(cfg *config.Config) *HealthCheckHandler {
	return &HealthCheckHandler{
		cfg: cfg,
	}
}
