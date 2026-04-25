package cmd

import (
	"fmt"
	"github.com/saleh-ghazimoradi/TeleGopher/config"
	"github.com/saleh-ghazimoradi/TeleGopher/infra/postgresql"
	"github.com/saleh-ghazimoradi/TeleGopher/internal/gateway/handler"
	"github.com/saleh-ghazimoradi/TeleGopher/internal/gateway/middleware"
	"github.com/saleh-ghazimoradi/TeleGopher/internal/gateway/route"
	"github.com/saleh-ghazimoradi/TeleGopher/internal/repository"
	"github.com/saleh-ghazimoradi/TeleGopher/internal/server"
	"github.com/saleh-ghazimoradi/TeleGopher/internal/service"
	"github.com/saleh-ghazimoradi/TeleGopher/internal/ws"
	"github.com/saleh-ghazimoradi/TeleGopher/utils"
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

		/*----------Slog Logger----------*/
		slogLogger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				if len(groups) == 0 && a.Key == slog.TimeKey {
					t := a.Value.Time()
					a.Value = slog.StringValue(t.Format("2006-01-02T15:04:05"))
				}
				return a
			},
		}))

		/*----------Strategy Logger----------*/
		logger := utils.NewLoggerContext(slogLogger)

		/*----------Config----------*/
		cfg, err := config.GetCfg()
		if err != nil {
			logger.Error("failed to get the config", "error", err)
			return
		}

		/*----------Postgresql----------*/
		postDB := postgresql.NewPostgresql(
			postgresql.WithHost(cfg.Postgresql.Host),
			postgresql.WithPort(cfg.Postgresql.Port),
			postgresql.WithUser(cfg.Postgresql.User),
			postgresql.WithPassword(cfg.Postgresql.Password),
			postgresql.WithName(cfg.Postgresql.Name),
			postgresql.WithMaxOpenConn(cfg.Postgresql.MaxOpenConn),
			postgresql.WithMaxIdleConn(cfg.Postgresql.MaxIdleConn),
			postgresql.WithMaxIdleTime(cfg.Postgresql.MaxIdleTime),
			postgresql.WithSSLMode(cfg.Postgresql.SSLMode),
			postgresql.WithTimeout(cfg.Postgresql.Timeout),
			postgresql.WithLogger(logger),
		)

		gormDB, _, err := postDB.Connect()
		if err != nil {
			logger.Error("failed to connect to the database", "error", err)
			return
		}

		/*----------Dependencies----------*/
		middlewares := middleware.NewMiddleware(logger, cfg)

		/*----------Repositories----------*/
		userRepository := repository.NewUserRepository(gormDB, gormDB)
		privateRepository := repository.NewPrivateRepository(gormDB, gormDB)
		messageRepository := repository.NewMessageRepository(gormDB, gormDB)

		/*----------Services----------*/
		authService := service.NewAuthService(userRepository, cfg)
		userService := service.NewUserService(userRepository)
		privateService := service.NewPrivateService(privateRepository, userRepository)
		messageService := service.NewMessageService(messageRepository, privateRepository)

		/*----------WS HUB----------*/
		wsHub := ws.NewHub(privateService, messageService, logger)

		/*----------Handlers----------*/
		healthCheck := handler.NewHealthCheckHandler(cfg)
		authHandler := handler.NewAuthHandler(authService)
		userHandler := handler.NewUserHandler(userService)
		privateHandler := handler.NewPrivateHandler(privateService)
		messageHandler := handler.NewMessageHandler(messageService)
		uploadFileHandler := handler.NewUploadFileHandler()
		wsHandler := handler.NewWebSocketHandler(userService, messageService, logger, wsHub, cfg)

		/*----------Routes----------*/
		healthRoute := route.NewHealthCheckRoute(healthCheck)
		authRoute := route.NewAuthRoute(middlewares, authHandler)
		userRoute := route.NewUserRoute(middlewares, userHandler)
		privateRoute := route.NewPrivateRoute(middlewares, privateHandler)
		messageRoute := route.NewMessageRoute(middlewares, messageHandler)
		uploadFileRoute := route.NewUploadFileRoute(middlewares, uploadFileHandler)
		wsRoute := route.NewWSRoute(wsHandler)

		/*----------Route Registery----------*/
		register := route.NewRegisterRoute(
			route.WithMiddleware(middlewares),
			route.WithHealthCheckRoute(healthRoute),
			route.WithAuthRoute(authRoute),
			route.WithUserRoute(userRoute),
			route.WithPrivateRoute(privateRoute),
			route.WithMessageRoute(messageRoute),
			route.WithUploadFileRoute(uploadFileRoute),
			route.WithWsRoute(wsRoute),
		)

		/*----------HTTP Server----------*/
		httpServer := server.NewServer(
			server.WithHost(cfg.Server.Host),
			server.WithPort(cfg.Server.Port),
			server.WithHandler(register.RegisterRoutes()),
			server.WithWriteTimeout(cfg.Server.WriteTimeout),
			server.WithReadTimeout(cfg.Server.ReadTimeout),
			server.WithIdleTimeout(cfg.Server.IdleTimeout),
			server.WithErrLog(slog.NewLogLogger(slogLogger.Handler(), slog.LevelError)),
			server.WithLogger(logger),
			server.WithHub(wsHub),
		)

		logger.Info("starting server", "addr", cfg.Server.Host+":"+cfg.Server.Port, "env", cfg.Application.Environment)
		if err := httpServer.Connect(); err != nil {
			logger.Error("failed to connect to the http server", "error", err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
