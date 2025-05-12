package main

import (
	"github.com/Chamanthra/TaskManagmentSystem/config"
	"github.com/Chamanthra/TaskManagmentSystem/routes"
)

func main() {
	config.ConnectDatabase()
	r := routes.SetupRouter()
	r.Run(":8080")
}
