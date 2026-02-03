package cmd

import (
	"fmt"
	"github.com/saleh-ghazimoradi/TeleGopher/config"
	infra "github.com/saleh-ghazimoradi/TeleGopher/infra/TXManager"
	"github.com/saleh-ghazimoradi/TeleGopher/infra/postgresql"
	"github.com/saleh-ghazimoradi/TeleGopher/internal/gateway/handler"
	"github.com/saleh-ghazimoradi/TeleGopher/internal/gateway/middleware"
	"github.com/saleh-ghazimoradi/TeleGopher/internal/gateway/route"
	"github.com/saleh-ghazimoradi/TeleGopher/internal/helper"
	"github.com/saleh-ghazimoradi/TeleGopher/internal/repository"
	"github.com/saleh-ghazimoradi/TeleGopher/internal/server"
	"github.com/saleh-ghazimoradi/TeleGopher/internal/service"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Running the TeleGopher app",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("run called")

		logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

		cfg, err := config.GetConfigInstance()
		if err != nil {
			logger.Error("error getting config", "err", err.Error())
			os.Exit(1)
		}

		dbConfig := postgresql.NewPostgresql(
			postgresql.WithHost(cfg.Postgresql.Host),
			postgresql.WithPort(cfg.Postgresql.Port),
			postgresql.WithUser(cfg.Postgresql.User),
			postgresql.WithPassword(cfg.Postgresql.Password),
			postgresql.WithName(cfg.Postgresql.Name),
			postgresql.WithTimeZone(cfg.Postgresql.TimeZone),
			postgresql.WithMaxOpenConn(cfg.Postgresql.MaxOpenConn),
			postgresql.WithMaxIdleConn(cfg.Postgresql.MaxIdleConn),
			postgresql.WithMaxIdleTime(cfg.Postgresql.MaxIdleTime),
			postgresql.WithMaxLifetime(cfg.Postgresql.MaxLifetime),
			postgresql.WithSSLMode(cfg.Postgresql.SSLMode),
			postgresql.WithConnectTimeout(cfg.Postgresql.ConnectionTimeout),
		)

		db, err := dbConfig.Connect()
		if err != nil {
			logger.Error("error connecting", "err", err.Error())
			os.Exit(1)
		}

		_ = db

		defer func() {
			if err := db.Close(); err != nil {
				logger.Error("error closing db", "err", err.Error())
				os.Exit(1)
			}
		}()

		txManager := infra.NewTxManager(db)
		errResponse := helper.NewErrResponse(logger)
		middlewares := middleware.NewMiddleware(cfg, logger, errResponse)
		validator := helper.NewValidator()

		userRepository := repository.NewUserRepository(db, db)
		userService := service.NewUserService(userRepository)

		authService := service.NewAuthenticationService(cfg, userRepository, txManager)
		authHandler := handler.NewAuthenticationHandler(errResponse, validator, authService)

		healthcheckHandler := handler.NewHealthCheckHandler(cfg, errResponse)
		healthcheckRoute := route.NewHealthCheckRoute(healthcheckHandler)
		userHandler := handler.NewUserHandler(userService)
		authRoute := route.NewAuthenticationRoute(authHandler, middlewares)
		userRoute := route.NewUserRoute(userHandler, middlewares)
		registerRoutes := route.NewRegisterRoutes(
			route.WithHealthCheckRoute(healthcheckRoute),
			route.WithAuthenticationRoute(authRoute),
			route.WithUserRoute(userRoute),
			route.WithMiddleware(middlewares),
		)

		s := server.NewServer(
			server.WithHost(cfg.Server.Host),
			server.WithPort(cfg.Server.Port),
			server.WithHandler(registerRoutes.Register()),
			server.WithReadTimeout(cfg.Server.ReadTimeout),
			server.WithWriteTimeout(cfg.Server.WriteTimeout),
			server.WithIdleTimeout(cfg.Server.IdleTimeout),
			server.WithErrLog(slog.NewLogLogger(logger.Handler(), slog.LevelError)),
			server.WithLogger(logger),
		)

		logger.Info("Server is running on:", "port", cfg.Server.Port)
		if err := s.Connect(); err != nil {
			logger.Error("error connecting to server", "err", err.Error())
			os.Exit(1)
		}

	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
