package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Migration command
var MigrateCommand = &cobra.Command{
	Use:   "migrate",
	Short: "Run database migrations",
	Run:   runMigrations,
}

func runMigrations(cmd *cobra.Command, args []string) {
	fmt.Println("Running migrations...")
	// Call your migration logic here
}
