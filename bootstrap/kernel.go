package bootstrap

import (
	"imohamedsheta/gocrud/config"
	"imohamedsheta/gocrud/database"
	"imohamedsheta/gocrud/pkg/cmd"
	"imohamedsheta/gocrud/pkg/enums"
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
	// loadDatabaseConnection()
}

func loadDatabaseConnection() {
	// Connect to the database
	database.Connect()
}

func loadEnvConfig() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal(enums.Red.Value() + "Error loading .env file" + enums.Reset.Value())
	}

	log.Println(enums.Green.Value() + "Kernel: Loaded .env file" + enums.Reset.Value())
}

func loadConfig() {
	config.LoadDatabaseConfig()
	config.LoadAppConfig()
}

// Run the Application (CLI or HTTP server)
func Run() {
	// If the command is "serve" then we'll start the HTTP server
	if len(os.Args) > 1 && os.Args[1] == "serve" {
		startHttpServer()
	}

	// Otherwise, we'll run the CLI
	cmd.Execute()
}

// Start the HTTP server
func startHttpServer() {
	log.Println(enums.Green.Value() + "Starting HTTP server on :7777..." + enums.Reset.Value())
	if err := http.ListenAndServe(":7777", routes.RegisterRoutes()); err != nil {
		log.Fatal(enums.Red.Value() + err.Error() + enums.Reset.Value())
	}
}
