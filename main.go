package main

import (
	"imohamedsheta/gocrud/bootstrap"
	_ "imohamedsheta/gocrud/docs"
	"imohamedsheta/gocrud/pkg/logger"
	"imohamedsheta/gocrud/pkg/support"
	"imohamedsheta/gocrud/query"
	"time"
)

func main() {
	// load the application
	bootstrap.Load()
	users, err := query.UsersTable().Get("email", "name")
	if err != nil {
		logger.Log().Error("Error fetching users: " + err.Error())
		return
	}

	user := query.User{
		Email:     "imohamedsheta@gmail.com",
		Name:      "Imohamed Sheta",
		CreatedAt: time.Now(),
	}
	err = query.UsersTable().Insert(user)

	if err != nil {
		support.DD(err)
	}
	support.DD(users)
}
