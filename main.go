package main

import (
	"imohamedsheta/gocrud/bootstrap"
	_ "imohamedsheta/gocrud/docs"
)

func main() {
	// load the application
	bootstrap.Load()

	// support.DD(query.TodosTable().Get("title", "status"))
	// Run the application
	bootstrap.Run()
}
