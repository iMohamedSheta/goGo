package main

import (
	"imohamedsheta/gocrud/bootstrap"
)

func main() {

	// load the application
	bootstrap.Load()
	// support.DD(config.AppConfig.Get("app"))
	// run the application
	bootstrap.Run()

}
