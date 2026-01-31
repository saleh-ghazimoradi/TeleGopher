package cmd

import (
	"fmt"
	"github.com/saleh-ghazimoradi/TeleGopher/config"
	"github.com/saleh-ghazimoradi/TeleGopher/infra/postgresql"
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

		log := slog.New(slog.NewJSONHandler(os.Stdout, nil))

		cfg, err := config.GetConfigInstance()
		if err != nil {
			log.Error("error getting config", err)
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
			log.Error("error connecting to database", err)
			os.Exit(1)
		}

		defer func() {
			if err := db.Close(); err != nil {
				log.Error("error closing database connection", err)
				os.Exit(1)
			}
		}()
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
