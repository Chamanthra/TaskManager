package controllers

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/Chamanthra/TaskManagmentSystem/config"
	"github.com/Chamanthra/TaskManagmentSystem/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

const uploadDir = "./uploads"

func init() {
	os.MkdirAll(uploadDir, os.ModePerm)
}

// UploadFile handles file uploads for tasks
func UploadFile(c *gin.Context) {
	claims := c.MustGet("claims").(jwt.MapClaims)
	userID := uint(claims["user_id"].(float64))
	role := claims["role"].(string)

	taskID, err := strconv.Atoi(c.Param("taskId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	// Check if task exists
	var task models.Task
	if err := config.DB.First(&task, taskID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	// Check ownership (unless admin)
	if role != "admin" && task.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only upload files to your own tasks"})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File upload error"})
		return
	}

	// Save file
	filename := filepath.Join(uploadDir, strconv.Itoa(int(taskID))+"_"+file.Filename)
	if err := c.SaveUploadedFile(file, filename); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	// Create file record
	fileRecord := models.File{
		FilePath: filename,
		FileName: file.Filename,
		FileType: file.Header.Get("Content-Type"),
		FileSize: file.Size,
		TaskID:   uint(taskID),
		UserID:   userID,
	}

	if err := config.DB.Create(&fileRecord).Error; err != nil {
		os.Remove(filename) // Clean up if DB fails
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create file record"})
		return
	}

	// Create notification for task owner
	if task.UserID != userID { // Don't notify yourself
		notification := models.Notification{
			Message: "New file uploaded to your task",
			UserID:  task.UserID,
			TaskID:  task.ID,
			Type:    "file_upload",
		}
		config.DB.Create(&notification)
	}

	c.JSON(http.StatusCreated, fileRecord)
}

// DownloadFile handles file downloads
func DownloadFile(c *gin.Context) {
	claims := c.MustGet("claims").(jwt.MapClaims)
	userID := uint(claims["user_id"].(float64))
	role := claims["role"].(string)

	fileID, err := strconv.Atoi(c.Param("fileId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file ID"})
		return
	}

	var file models.File
	if err := config.DB.Preload("Task").First(&file, fileID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	// Check ownership (unless admin)
	if role != "admin" && file.Task.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only download files from your own tasks"})
		return
	}

	if _, err := os.Stat(file.FilePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found on server"})
		return
	}

	c.FileAttachment(file.FilePath, file.FileName)
}

// DeleteFile handles file deletion
func DeleteFile(c *gin.Context) {
	claims := c.MustGet("claims").(jwt.MapClaims)
	userID := uint(claims["user_id"].(float64))
	role := claims["role"].(string)

	fileID, err := strconv.Atoi(c.Param("fileId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file ID"})
		return
	}

	var file models.File
	if err := config.DB.Preload("Task").First(&file, fileID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	// Check ownership (admin can delete any, users can only delete their own)
	if role != "admin" && file.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only delete your own files"})
		return
	}

	// Delete file from filesystem
	if err := os.Remove(file.FilePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete file from storage"})
		return
	}

	// Delete record from database
	if err := config.DB.Delete(&file).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete file record"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File deleted"})
}
