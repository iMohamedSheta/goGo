package main

import (
	"github.com/iMohamedSheta/xapp/bootstrap"
	_ "github.com/iMohamedSheta/xapp/docs"
)

func main() {
	// load the application
	bootstrap.Load()

	// run the application
	bootstrap.Run()
}
