package main

import (
	"imohamedsheta/gocrud/bootstrap"
	_ "imohamedsheta/gocrud/docs"
)

func main() {
	// load the application
	bootstrap.Load()

	// data := map[string]interface{}{
	// 	"name":     "imohamedsheta",
	// 	"email":    "imohamedsheta@gmail.com",
	// 	"password": "123456",
	// }

	// support.DD(query.Table("users").Where("id", "=", 15).Where("username", "=", "imohamedsheta").UpdateSql(data))
	// Run the application
	bootstrap.Run()
}
