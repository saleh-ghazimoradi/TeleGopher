package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// migrateDownCmd represents the migrateDown command
var migrateDownCmd = &cobra.Command{
	Use:   "migrateDown",
	Short: "Migrating Down",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("migrateDown called")
	},
}

func init() {
	rootCmd.AddCommand(migrateDownCmd)
}
