package config

import (
	migrationCmd "imohamedsheta/gocrud/cmd/migration"

	"github.com/spf13/cobra"
)

/*
Here we register all the commands that we want to use in our application
*/
var registeredCommands = []*cobra.Command{
	migrationCmd.MigrateCommand,
}

/*
GetRegisteredCommands returns all the registered commands
*/
func GetRegisteredCommands() []*cobra.Command {
	return registeredCommands
}
