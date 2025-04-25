package main

import (
	"imohamedsheta/gocrud/bootstrap"
	_ "imohamedsheta/gocrud/docs"
)

func main() {
	// load the application
	bootstrap.Load()

	// Run the application
	bootstrap.Run()
}
