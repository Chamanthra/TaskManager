package main

import (
	"github.com/Chamanthra/TaskManager/config"
	"github.com/Chamanthra/TaskManager/migrations"
	"github.com/Chamanthra/TaskManager/models"
	"github.com/Chamanthra/TaskManager/routes"
	"github.com/Chamanthra/TaskManager/workers"
)

func main() {
	config.ConnectDatabase()

	// AutoMigrate all models
	config.DB.AutoMigrate(
		&models.User{},
		&models.Task{},
		&models.Role{},
		&models.Notification{},
		&models.Comment{},
		&models.File{},
	)

	// Seed initial data
	if err := migrations.InitRoles(config.DB); err != nil {
		panic("Failed to seed roles: " + err.Error())
	}

	// Start notification worker
	go workers.StartNotificationWorker()

	r := routes.SetupRouter()
	r.Run(":8080")
}
