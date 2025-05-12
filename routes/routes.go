package routes

import (
	"github.com/Chamanthra/TaskManagmentSystem/controllers"
	"github.com/Chamanthra/TaskManagmentSystem/middlewares"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)

	protected := r.Group("/tasks")
	protected.Use(middlewares.AuthMiddleware())
	{
		protected.POST("/", controllers.CreateTask)
		protected.GET("/", controllers.GetTasks)
		protected.PUT("/:id", controllers.UpdateTask)
		protected.DELETE("/:id", controllers.DeleteTask)
	}

	return r
}
