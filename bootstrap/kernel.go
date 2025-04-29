package bootstrap

import (
	"context"
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
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

/*
	Package bootstrap kernel initializes and runs the core application lifecycle.

	This file is the main entry point for bootstrapping the application. It:
	- Loads essential configurations (environment variables, app settings, DB, logger, validation rules).
	- Decides whether to start the HTTP server or run a CLI command based on input arguments.
	- Starts the HTTP server with graceful shutdown capabilities on system signals (e.g., SIGTERM).
	- Registers and executes CLI commands when applicable.

	Modules involved:
	- config: loads and provides access to application configuration.
	- logger: sets up the global logging system.
	- database: initializes the database connection.
	- validate: registers custom validation rules.
	- routes: provides the HTTP router.
	- cmd: manages CLI commands.

	This file should be invoked from `main.go` via `bootstrap.Load()` and `bootstrap.Run()`.
*/

// Load application (env, config, DB, logger, validation)
func Load() {
	loadEnvConfig()
	loadConfig()
	loadDatabaseConnection()
	logger.LoadLogger()
	loadValidation()
}

// Run the application (serve HTTP or execute CLI command)
func Run() {
	if len(os.Args) > 1 && os.Args[1] == "serve" {
		startHttpServer()
	} else {
		executeCommand()
	}
}

// Start the HTTP server with graceful shutdown
func startHttpServer() {
	shutdown_timeout := config.App.Get("app.shutdown_timeout").(time.Duration)
	url := config.App.Get("app.url").(string)
	port := config.App.Get("app.port").(string)

	srv := &http.Server{
		Addr:    url + ":" + port,
		Handler: routes.RegisterRoutes(),
	}

	// Start the server in a goroutine
	go func() {
		log.Println(enums.Green.Value() + "Starting HTTP server on http://" + url + ":" + port + " ..." + enums.Reset.Value())
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(enums.Red.Value() + "Server error: " + err.Error() + enums.Reset.Value())
		}
	}()

	// Listen for OS signals for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println(enums.Yellow.Value() + "Shutting down server..." + enums.Reset.Value())

	// timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), shutdown_timeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal(enums.Red.Value() + "Forced to shutdown: " + err.Error() + enums.Reset.Value())
	}

	log.Println(enums.Green.Value() + "Server exited properly" + enums.Reset.Value())
}

// Execute CLI command
func executeCommand() {
	registerCommands()
	cmd.Execute()
}

// Load environment variables
func loadEnvConfig() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(enums.Red.Value() + "Error loading .env file" + enums.Reset.Value())
	}
	log.Println(enums.Green.Value() + "Kernel: Loaded .env file" + enums.Reset.Value())
}

// Connect to the database
func loadDatabaseConnection() {
	database.Connect()
}

// Load validation rules
func loadValidation() {
	validator := validate.Validator()
	for tag, rule := range registeredRules {
		if err := validator.RegisterValidation(tag, rule); err != nil {
			logger.Log().Error("Error registering validation rule: " + tag + ": " + err.Error())
		}
	}
}

// Register CLI commands
func registerCommands() {
	rootCmd := cmd.Command()
	for _, c := range registeredCommands {
		rootCmd.AddCommand(c)
	}
}
