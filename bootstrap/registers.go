package bootstrap

import (
	"github.com/iMohamedSheta/xapp/app/rules"
	cmd "github.com/iMohamedSheta/xapp/cmd/example"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/cobra"
)

/*
	This file is used to register all the custom stuff in the application
	like the commands, validations rules, ...etc
*/

// Register new validations rules
var registeredRules = map[string]validator.Func{
	"unique_db": rules.Unique,
}

// Register new commands
var registeredCommands = []*cobra.Command{
	cmd.SayHelloCommand,
}
