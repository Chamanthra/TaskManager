package controllers

import (
	"net/http"

	"github.com/Chamanthra/TaskManagmentSystem/config"
	"github.com/Chamanthra/TaskManagmentSystem/models"
	"github.com/gin-gonic/gin"
)

func CreateTask(c *gin.Context) {
	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	config.DB.Create(&task)
	c.JSON(http.StatusOK, task)
}

func GetTasks(c *gin.Context) {
	var tasks []models.Task
	config.DB.Find(&tasks)
	c.JSON(http.StatusOK, tasks)
}

func UpdateTask(c *gin.Context) {
	id := c.Param("id")
	var task models.Task
	config.DB.First(&task, id)

	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	config.DB.Save(&task)
	c.JSON(http.StatusOK, task)
}

func DeleteTask(c *gin.Context) {
	id := c.Param("id")
	var task models.Task
	config.DB.Delete(&task, id)
	c.JSON(http.StatusOK, gin.H{"message": "Task deleted"})
}
