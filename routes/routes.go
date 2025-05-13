package routes

import (
	"github.com/Chamanthra/TaskManager/controllers"
	"github.com/Chamanthra/TaskManager/middlewares"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Health check route
	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "OK",
		})
	})

	// User authentication routes
	r.POST("/api/register", controllers.Register)
	r.POST("/api/login", controllers.Login)

	protected := r.Group("/api")
	protected.Use(middlewares.AuthMiddleware())
	{
		// Task routes
		taskRoutes := protected.Group("/tasks")
		{
			taskRoutes.POST("/", controllers.CreateTask)
			taskRoutes.GET("/", controllers.GetTasks)
			taskRoutes.PUT("/:id", controllers.UpdateTask)
			taskRoutes.DELETE("/:id", controllers.DeleteTask)

			// Task comments
			taskRoutes.POST("/:taskId/comments", controllers.AddComment)
			taskRoutes.GET("/:taskId/comments", controllers.GetTaskComments)
			taskRoutes.DELETE("/comments/:commentId", controllers.DeleteComment)

			// Task files
			taskRoutes.POST("/:taskId/files", controllers.UploadFile)
			taskRoutes.GET("/files/:fileId", controllers.DownloadFile)
			taskRoutes.DELETE("/files/:fileId", controllers.DeleteFile)
		}

		// Notification routes
		protected.GET("/notifications", controllers.GetUserNotifications)
		protected.PUT("/notifications/:id/read", controllers.MarkNotificationAsRead)

		// User profile routes
		protected.GET("/profile", controllers.GetUserProfile)
		protected.PUT("/profile", controllers.UpdateUserProfile)

		// Admin-only routes
		adminRoutes := protected.Group("/admin")
		{
			adminRoutes.GET("/users", controllers.GetUsers)
			adminRoutes.DELETE("/users/:id", controllers.DeleteUser)
			adminRoutes.GET("/tasks/all", controllers.GetAllTasks)
		}
	}

	return r
}
