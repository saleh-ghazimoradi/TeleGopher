package handler

import (
	"github.com/saleh-ghazimoradi/TeleGopher/config"
	"github.com/saleh-ghazimoradi/TeleGopher/internal/helper"
	"net/http"
)

type HealthCheckHandler struct {
	cfg         *config.Config
	errResponse *helper.ErrResponse
}

func (h *HealthCheckHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	env := helper.Envelope{
		"status": "available",
		"system_info": map[string]string{
			"environment": h.cfg.Application.Environment,
			"version":     h.cfg.Application.Version,
		},
	}

	if err := helper.WriteJSON(w, http.StatusOK, env, nil); err != nil {
		h.errResponse.ServerErrorResponse(w, r, err)
	}
}

func NewHealthCheckHandler(cfg *config.Config, errResponse *helper.ErrResponse) *HealthCheckHandler {
	return &HealthCheckHandler{
		cfg:         cfg,
		errResponse: errResponse,
	}
}
