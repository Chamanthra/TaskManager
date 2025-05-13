package controllers

import (
	"net/http"

	"github.com/Chamanthra/TaskManagmentSystem/config"
	"github.com/Chamanthra/TaskManagmentSystem/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
)

// In controllers/task.go - CreateTask function
func CreateTask(c *gin.Context) {
	// Get user ID from JWT claims
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User information not found"})
		return
	}

	jwtClaims := claims.(jwt.MapClaims)
	userID := uint(jwtClaims["user_id"].(float64))

	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set the user ID for the task
	task.UserID = userID

	if err := config.DB.Create(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, task)
}

// In controllers/task.go - GetTasks function
func GetTasks(c *gin.Context) {
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User information not found"})
		return
	}

	jwtClaims := claims.(jwt.MapClaims)
	userID := uint(jwtClaims["user_id"].(float64))
	role := jwtClaims["role"].(string)

	var tasks []models.Task
	var query *gorm.DB

	if role == "admin" {
		// Admin can see all tasks
		query = config.DB.Preload("User") // Include user information
	} else {
		// Regular users can only see their own tasks
		query = config.DB.Where("user_id = ?", userID)
	}

	if err := query.Find(&tasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tasks", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

// In controllers/task.go - UpdateTask function
func UpdateTask(c *gin.Context) {
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User information not found"})
		return
	}

	jwtClaims := claims.(jwt.MapClaims)
	userID := uint(jwtClaims["user_id"].(float64))
	role := jwtClaims["role"].(string)

	id := c.Param("id")
	var task models.Task

	if err := config.DB.First(&task, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	// Check ownership (unless admin)
	if role != "admin" && task.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only update your own tasks"})
		return
	}

	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Ensure user can't change ownership
	task.UserID = userID

	if err := config.DB.Save(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, task)
}

// Similar changes for DeleteTask
func DeleteTask(c *gin.Context) {
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User information not found"})
		return
	}

	jwtClaims := claims.(jwt.MapClaims)
	userID := uint(jwtClaims["user_id"].(float64))
	role := jwtClaims["role"].(string)

	id := c.Param("id")
	var task models.Task

	if err := config.DB.First(&task, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	// Check ownership (unless admin)
	if role != "admin" && task.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only delete your own tasks"})
		return
	}

	if err := config.DB.Delete(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task deleted"})
}

// GetAllTasks gets all tasks (admin only)
func GetAllTasks(c *gin.Context) {
	claims := c.MustGet("claims").(jwt.MapClaims)
	role := claims["role"].(string)

	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only admins can view all tasks"})
		return
	}

	var tasks []models.Task
	if err := config.DB.Preload("User").Find(&tasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tasks"})
		return
	}

	c.JSON(http.StatusOK, tasks)
}
