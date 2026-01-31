package cmd

import (
	"fmt"
	"github.com/saleh-ghazimoradi/TeleGopher/config"
	"github.com/saleh-ghazimoradi/TeleGopher/infra/migration"
	"github.com/saleh-ghazimoradi/TeleGopher/infra/postgresql"
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

		migrator, err := migration.NewMigrator(db, dbConfig.Name)
		if err != nil {
			logger.Error("migration init failed", "error", err.Error())
			os.Exit(1)
		}

		defer func() {
			if err := migrator.Close(); err != nil {
				logger.Error("failed to close migration", "error", err.Error())
				os.Exit(1)
			}
		}()

		if err := migrator.Up(); err != nil {
			logger.Error("migration failed", "error", err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(migrateUpCmd)
}
