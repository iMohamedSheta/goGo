package config

import (
	cmd "imohamedsheta/gocrud/cmd/example"

	"github.com/spf13/cobra"
)

/*
Here we register all the commands that we want to use in our application
*/
var registeredCommands = []*cobra.Command{
	cmd.SayHelloCommand,
}

/*
GetRegisteredCommands returns all the registered commands
*/
func GetRegisteredCommands() []*cobra.Command {
	return registeredCommands
}
