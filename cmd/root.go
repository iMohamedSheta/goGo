package cmd

import (
	"fmt"
	"imohamedsheta/gocrud/config"
	"imohamedsheta/gocrud/enums"
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
		showHowToUse()
	},
}

func init() {
	// Disable the "completion" command
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}

// registerCommands registers all commands
func registerCommands(root *cobra.Command) {
	for _, cmd := range config.RegisteredCommands {
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

// showHowToUse shows how to use the application
func showHowToUse() {
	fmt.Println("\n" + enums.Blue.Value() + "📌 Usage:" + enums.Reset.Value())
	fmt.Println("  " + enums.Green.Value() + "▶ To start the server:    go run . serve" + enums.Reset.Value())
	fmt.Println("  " + enums.Green.Value() + "▶ To run CLI commands:    go run . <command>" + enums.Reset.Value() + "\n")
}
