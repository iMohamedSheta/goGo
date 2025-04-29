package bootstrap

import (
	"imohamedsheta/gocrud/app/rules"
	cmd "imohamedsheta/gocrud/cmd/example"
	"imohamedsheta/gocrud/config"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/cobra"
)

/*
	This file is used to register all the custom stuff in the application
	like the commands, configuration files, validations rules, ...etc
*/

// Load the config files
func loadConfig() {
	config.LoadAppConfig()
	config.LoadDatabaseConfig()
}

// Register new validations rules
var registeredRules = map[string]validator.Func{
	"unique_db": rules.Unique,
}

// Register new commands
var registeredCommands = []*cobra.Command{
	cmd.SayHelloCommand,
}
