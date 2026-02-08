package middleware

import (
	"fmt"
	"github.com/saleh-ghazimoradi/TeleGopher/config"
	"github.com/saleh-ghazimoradi/TeleGopher/internal/domain"
	"github.com/saleh-ghazimoradi/TeleGopher/internal/helper"
	"github.com/saleh-ghazimoradi/TeleGopher/utils"
	"log/slog"
	"net/http"
	"strings"
)

type Middleware struct {
	cfg         *config.Config
	logger      *slog.Logger
	errResponse *helper.ErrResponse
}

func (m *Middleware) Logging(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m.logger.Info("Incoming request: ", "method", r.Method, "path", r.URL.Path, "protocol", r.Proto, "remote_addr", r.RemoteAddr)
		next(w, r)
	}
}

func (m *Middleware) CORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Platform")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next(w, r)
	}
}

func (m *Middleware) Recover(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				m.errResponse.ServerErrorResponse(w, r, fmt.Errorf("%v", err))
			}
		}()
		next(w, r)
	}
}

func (m *Middleware) Authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			m.errResponse.InvalidCredentialsResponse(w, r)
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			m.errResponse.InvalidAuthenticationTokenResponse(w, r)
			return
		}

		platform := r.Header.Get("X-Platform")
		if platform != string(domain.PlatformWeb) && platform != string(domain.PlatformMobile) {
			m.errResponse.InvalidCredentialsResponse(w, r)
			return
		}

		claims, err := utils.ValidateToken(tokenParts[1], m.cfg.JWT.Secret)
		if err != nil {
			m.errResponse.InvalidAuthenticationTokenResponse(w, r)
			return
		}

		if claims.Platform != platform {
			m.errResponse.InvalidAuthenticationTokenResponse(w, r)
			return
		}

		ctx := r.Context()
		ctx = utils.WithUserId(ctx, claims.UserId)
		ctx = utils.WithUserName(ctx, claims.Name)
		ctx = utils.WithPlatform(ctx, claims.Platform)
		next(w, r.WithContext(ctx))
	}
}

func NewMiddleware(cfg *config.Config, logger *slog.Logger, errResponse *helper.ErrResponse) *Middleware {
	return &Middleware{
		cfg:         cfg,
		logger:      logger,
		errResponse: errResponse,
	}
}
