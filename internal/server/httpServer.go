package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	Host         string
	Port         string
	Handler      http.Handler
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
	ErrLog       *log.Logger
	Logger       *slog.Logger
}

type Options func(*Server)

func WithHost(host string) Options {
	return func(s *Server) {
		s.Host = host
	}
}

func WithPort(port string) Options {
	return func(s *Server) {
		s.Port = port
	}
}

func WithHandler(h http.Handler) Options {
	return func(s *Server) {
		s.Handler = h
	}
}

func WithReadTimeout(timeout time.Duration) Options {
	return func(s *Server) {
		s.ReadTimeout = timeout
	}
}

func WithWriteTimeout(timeout time.Duration) Options {
	return func(s *Server) {
		s.WriteTimeout = timeout
	}
}

func WithIdleTimeout(timeout time.Duration) Options {
	return func(s *Server) {
		s.IdleTimeout = timeout
	}
}

func WithErrLog(logger *log.Logger) Options {
	return func(s *Server) {
		s.ErrLog = logger
	}
}

func WithLogger(logger *slog.Logger) Options {
	return func(server *Server) {
		server.Logger = logger
	}
}

func (s *Server) addr() string {
	return fmt.Sprintf("%s:%s", s.Host, s.Port)
}

func (s *Server) Connect() error {
	server := &http.Server{
		Addr:         s.addr(),
		Handler:      s.Handler,
		IdleTimeout:  s.IdleTimeout,
		ReadTimeout:  s.ReadTimeout,
		WriteTimeout: s.WriteTimeout,
		ErrorLog:     s.ErrLog,
	}

	shutdownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		se := <-quit

		s.Logger.Info("caught shutdown signal", "signal", se.String())

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		err := server.Shutdown(ctx)
		if err != nil {
			shutdownError <- err
		}

		s.Logger.Info("completing background tasks", "addr", server.Addr)
		shutdownError <- nil
	}()

	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	if err := <-shutdownError; err != nil {
		return err
	}

	s.Logger.Info("Stopped server", "addr", server.Addr)
	return nil
}

func NewServer(opts ...Options) *Server {
	s := &Server{}
	for _, opt := range opts {
		opt(s)
	}
	return s
}
