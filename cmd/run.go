package cmd

import (
	"fmt"
	"github.com/saleh-ghazimoradi/TeleGopher/config"
	"github.com/saleh-ghazimoradi/TeleGopher/infra/postgresql"
	"github.com/saleh-ghazimoradi/TeleGopher/internal/server"
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
			logger.Error("error connecting to database", "err", err.Error())
			os.Exit(1)
		}

		defer func() {
			if err := db.Close(); err != nil {
				logger.Error("error closing database connection", "err", err.Error())
				os.Exit(1)
			}
		}()

		s := server.NewServer(
			server.WithHost(cfg.Server.Host),
			server.WithPort(cfg.Server.Port),
			server.WithHandler(nil),
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
