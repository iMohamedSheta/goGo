package main

import (
	"imohamedsheta/gocrud/bootstrap"
	"imohamedsheta/gocrud/pkg/session"
	s "imohamedsheta/gocrud/pkg/support"
)

func main() {

	// load the application
	bootstrap.Load()
	// support.DD(config.AppConfig.Get("app"))
	// run the application
	// bootstrap.Run()

	s.DD(session.NewSession())

}
