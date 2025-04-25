package bootstrap

import (
	"imohamedsheta/gocrud/database"
	"imohamedsheta/gocrud/pkg/cmd"
	"imohamedsheta/gocrud/pkg/config"
	"imohamedsheta/gocrud/pkg/enums"
	"imohamedsheta/gocrud/pkg/logger"
	"imohamedsheta/gocrud/pkg/validate"
	"imohamedsheta/gocrud/routes"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

/*
*
	This file is used to bootstrap the application.
*
*/

// Load Application
func Load() {
	loadEnvConfig()
	loadConfig()
	loadDatabaseConnection()
	logger.LoadLogger()
	loadValidation()
}

// Run the Application (CLI or HTTP server)
func Run() {
	// If the command is "serve" then we'll start the HTTP server
	if len(os.Args) > 1 && os.Args[1] == "serve" {
		startHttpServer()
	}

	// Otherwise, we'll run the CLI
	executeCommand()
}

// Start the HTTP server
func startHttpServer() {

	url := config.App.Get("app.url").(string)
	port := config.App.Get("app.port").(string)

	log.Println(enums.Green.Value() + "Starting HTTP server on http://" + url + ":" + port + " ..." + enums.Reset.Value())

	if err := http.ListenAndServe(url+":"+port, routes.RegisterRoutes()); err != nil {
		log.Fatal(enums.Red.Value() + err.Error() + enums.Reset.Value())
	}
}

func executeCommand() {
	// Register all commands
	registerCommands()

	// Execute the command
	cmd.Execute()
}

// connect to the database
func loadDatabaseConnection() {
	// Connect to the database
	database.Connect()
}

// Load the environment variables
func loadEnvConfig() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal(enums.Red.Value() + "Error loading .env file" + enums.Reset.Value())
	}

	log.Println(enums.Green.Value() + "Kernel: Loaded .env file" + enums.Reset.Value())
}

// Load the validation rules
func loadValidation() {
	validator := validate.Validator()

	for tag, rule := range registeredRules {
		err := validator.RegisterValidation(tag, rule)

		if err != nil {
			logger.Log().Error("Error registering validation rule: " + tag + ":  error: " + err.Error())
		}
	}
}

// registerCommands registers all commands
func registerCommands() {
	rootCmd := cmd.Command()
	for _, cmd := range registeredCommands {
		rootCmd.AddCommand(cmd)
	}
}
