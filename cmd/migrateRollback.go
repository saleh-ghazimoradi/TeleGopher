package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// migrateRollbackCmd represents the migrateRollback command
var migrateRollbackCmd = &cobra.Command{
	Use:   "migrateRollback",
	Short: "Migrating Rollback",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("migrateRollback called")
	},
}

func init() {
	rootCmd.AddCommand(migrateRollbackCmd)
}
