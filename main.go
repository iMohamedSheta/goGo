package main

import (
	"imohamedsheta/gocrud/bootstrap"
	_ "imohamedsheta/gocrud/docs"
)

func main() {

	// load the application
	bootstrap.Load()
	// support.DD(config.AppConfig.Get("app"))
	// run the application
	bootstrap.Run()

	// s.DD(session.NewSession())

}
