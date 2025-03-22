package config

import (
	migrationCmd "imohamedsheta/gocrud/cmd/migration"

	"github.com/spf13/cobra"
)

/*
Here we register all the commands that we want to use in our application
*/
var RegisteredCommands = []*cobra.Command{
	migrationCmd.MigrateCommand,
}
