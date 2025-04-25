package cmd

import (
	"fmt"
	"imohamedsheta/gocrud/pkg/logger"
	s "imohamedsheta/gocrud/pkg/support"
	"os"

	"github.com/spf13/cobra"
)

// Root command
var rootCmd = &cobra.Command{
	Use:   "",
	Short: "Command runner",
	Long:  "A simple CLI tool for running application commands",
	Run: func(cmd *cobra.Command, args []string) {
		err := cmd.Help()
		if err != nil {
			logger.Log().Error(err.Error())
		}
		s.PrintHowToUseApp()
	},
}

func Command() *cobra.Command {
	return rootCmd
}

func init() {
	// Disable the "completion" command
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}

// // registerCommands registers all commands
// func registerCommands(root *cobra.Command) {
// 	for _, cmd := range bootstrap.RegisteredCommands() {
// 		root.AddCommand(cmd)
// 	}
// }

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
