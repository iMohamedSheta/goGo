package bootstrap

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/iMohamedSheta/xapp/database"
	"github.com/iMohamedSheta/xapp/pkg/cmd"
	"github.com/iMohamedSheta/xapp/pkg/config"
	"github.com/iMohamedSheta/xapp/pkg/enums"
	"github.com/iMohamedSheta/xapp/pkg/logger"
	"github.com/iMohamedSheta/xapp/pkg/validate"
	"github.com/iMohamedSheta/xapp/routes"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
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
	config.LoadAll()
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

	shutdown_timeout := config.App.GetDuration("app.shutdown_timeout", 10*time.Second)

	url := config.App.GetString("app.url", "localhost")

	port := config.App.GetString("app.port", "8080")

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
	grpcEnabled := config.App.GetBool("app.grpc.enabled", false)

	var grpcServer *grpc.Server
	var grpcListener net.Listener

	if grpcEnabled {
		grpcUrl := config.App.GetString("app.grpc.url", "localhost")
		grpcPort := config.App.GetString("app.grpc.port", "50051")
		grpcServer, grpcListener = startGrpcServer(grpcUrl, grpcPort)
	}

	<-quit

	log.Println(enums.Yellow.Value() + "Shutting down server..." + enums.Reset.Value())

	// timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), shutdown_timeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal(enums.Red.Value() + "Forced to shutdown: " + err.Error() + enums.Reset.Value())
	}

	if grpcEnabled {
		// Graceful shutdown for gRPC
		go func() {
			grpcServer.GracefulStop()
			grpcListener.Close()
		}()

	}

	log.Println(enums.Green.Value() + "Server exited properly" + enums.Reset.Value())
}

func startGrpcServer(url string, port string) (*grpc.Server, net.Listener) {

	listener, err := net.Listen("tcp", url+":"+port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)

	routes.LoadGrpcRoutes(grpcServer)

	go func() {
		log.Println(enums.Yellow.Value() + "Starting gRPC server on tcp://" + url + ":" + port + enums.Reset.Value())
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatalf("Failed to serve gRPC: %v", err)
		}
	}()

	return grpcServer, listener
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
