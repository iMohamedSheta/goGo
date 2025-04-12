package cmd

import "github.com/spf13/cobra"

var SayHelloCommand = &cobra.Command{
	Use:   "say:hello",
	Short: "say hello to the world",
	Run:   handle(),
}

func handle() func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		println("Hello, World")
	}
}
