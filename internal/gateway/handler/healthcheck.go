package handler

import (
	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/saleh-ghazimoradi/TeleGopher/config"
	"github.com/saleh-ghazimoradi/TeleGopher/internal/helper"
	"net/http"
)

type HealthCheckHandler struct {
	cfg *config.Config
}

// HealthCheck godoc
// @Summary      Health check
// @Description  Check if the server is running and get system information
// @Tags         health
// @Produce      json
// @Success      200 {object} map[string]interface{} "Server is running"
// @Failure      500 {object} map[string]interface{} "Internal server error"
// @Router       /healthcheck [get]
func (h *HealthCheckHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	data := map[string]any{
		"status": "available",
		"system_info": map[string]any{
			"environment": h.cfg.Application.Environment,
			"version":     h.cfg.Application.Version,
		},
	}
	helper.SuccessResponse(w, "I'm breathing", data)
}

func (h *HealthCheckHandler) WsHealthCheck(w http.ResponseWriter, r *http.Request) {
	opts := &websocket.AcceptOptions{
		OriginPatterns: []string{"*"},
	}

	conn, err := websocket.Accept(w, r, opts)
	if err != nil {
		return
	}
	defer conn.Close(websocket.StatusNormalClosure, "Connection closed")

	ctx := r.Context()

	for {
		var message string
		err := wsjson.Read(ctx, conn, &message)
		if err != nil {
			break
		}

		response := map[string]any{
			"data":    message,
			"from":    "server",
			"success": true,
		}

		err = wsjson.Write(ctx, conn, response)
		if err != nil {
			break
		}
	}
}

func NewHealthCheckHandler(cfg *config.Config) *HealthCheckHandler {
	return &HealthCheckHandler{
		cfg: cfg,
	}
}
