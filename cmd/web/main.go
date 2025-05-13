package main

import (
	"github.com/Chamanthra/TaskManagmentSystem/config"
	"github.com/Chamanthra/TaskManagmentSystem/migrations"
	"github.com/Chamanthra/TaskManagmentSystem/models"
	"github.com/Chamanthra/TaskManagmentSystem/routes"
	"github.com/Chamanthra/TaskManagmentSystem/workers"
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
