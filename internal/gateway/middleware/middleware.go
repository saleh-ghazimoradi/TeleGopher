package middleware

import (
	"errors"
	"fmt"
	"github.com/saleh-ghazimoradi/TeleGopher/config"
	"github.com/saleh-ghazimoradi/TeleGopher/internal/helper"
	"github.com/saleh-ghazimoradi/TeleGopher/utils"
	"net/http"
	"strings"
)

type Middleware struct {
	logger utils.LoggerStrategy
	cfg    *config.Config
}

func (m *Middleware) Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.logger.Info("Incoming request: ", "method", r.Method, "path", r.URL.Path, "protocol", r.Proto, "remote_addr", r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Platform")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				helper.InternalServerError(w, "panic recovery hit", fmt.Errorf("%v", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			helper.UnauthorizedResponse(w, "Authorization header missing")
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			helper.UnauthorizedResponse(w, "Invalid authorization header format")
			return
		}

		platform := r.Header.Get("X-Platform")
		if platform != "web" && platform != "mobile" {
			helper.BadRequestResponse(w, "invalid platform", errors.New("invalid platform"))
			return
		}

		claims, err := utils.ValidateToken(tokenParts[1], m.cfg.JWT.Secret)
		if err != nil {
			helper.UnauthorizedResponse(w, "Unauthorized")
			return
		}

		if claims.Platform != platform {
			helper.UnauthorizedResponse(w, "Unauthorized")
		}

		ctx := r.Context()
		ctx = utils.WithUserId(ctx, claims.UserId)
		ctx = utils.WithName(ctx, claims.Name)
		ctx = utils.WithPlatform(ctx, claims.Platform)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)

	})
}

func (m *Middleware) WrapAuth(handlerFunc http.HandlerFunc) http.Handler {
	return m.Authenticate(handlerFunc)
}

func NewMiddleware(logger utils.LoggerStrategy, cfg *config.Config) *Middleware {
	return &Middleware{
		logger: logger,
		cfg:    cfg,
	}
}
