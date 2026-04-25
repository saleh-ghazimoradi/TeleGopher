package cmd

import (
	"fmt"
	"github.com/saleh-ghazimoradi/TeleGopher/config"
	"github.com/saleh-ghazimoradi/TeleGopher/infra/postgresql"
	"github.com/saleh-ghazimoradi/TeleGopher/internal/domain"
	"github.com/saleh-ghazimoradi/TeleGopher/utils"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

// migrateUpCmd represents the migrateUp command
var migrateUpCmd = &cobra.Command{
	Use:   "migrateUp",
	Short: "Migrating UP",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("migrateUp called")

		slogLogger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
		logger := utils.NewLoggerContext(slogLogger)

		cfg, err := config.GetCfg()
		if err != nil {
			logger.Error("failed to get the config", "error", err)
			return
		}

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

		if err := gormDB.Migrator().AutoMigrate(&domain.User{}, &domain.Private{}, &domain.Message{}); err != nil {
			logger.Error("failed to migrate up", "error", err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(migrateUpCmd)
}
