package cmd

import (
	"fmt"
	"imohamedsheta/gocrud/config"
	"imohamedsheta/gocrud/pkg/support"
	"os"

	"github.com/spf13/cobra"
)

// Root command
var rootCmd = &cobra.Command{
	Use:   "",
	Short: "Command runner",
	Long:  "A simple CLI tool for running application commands",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		support.PrintHowToUseApp()
	},
}

func init() {
	// Disable the "completion" command
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}

// registerCommands registers all commands
func registerCommands(root *cobra.Command) {
	for _, cmd := range config.GetRegisteredCommands() {
		root.AddCommand(cmd)
	}
}

// Execute runs the root command
func Execute() {
	registerCommands(rootCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
