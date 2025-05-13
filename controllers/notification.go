package controllers

import (
	"net/http"

	"github.com/Chamanthra/TaskManagmentSystem/config"
	"github.com/Chamanthra/TaskManagmentSystem/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// GetUserNotifications gets all notifications for the current user
func GetUserNotifications(c *gin.Context) {
	claims := c.MustGet("claims").(jwt.MapClaims)
	userID := uint(claims["user_id"].(float64))

	var notifications []models.Notification
	if err := config.DB.Where("user_id = ?", userID).Order("created_at desc").Find(&notifications).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get notifications"})
		return
	}

	c.JSON(http.StatusOK, notifications)
}

// MarkNotificationAsRead marks a notification as read
func MarkNotificationAsRead(c *gin.Context) {
	claims := c.MustGet("claims").(jwt.MapClaims)
	userID := uint(claims["user_id"].(float64))

	notificationID := c.Param("id")
	var notification models.Notification

	if err := config.DB.First(&notification, notificationID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Notification not found"})
		return
	}

	// Verify ownership
	if notification.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only mark your own notifications as read"})
		return
	}

	if err := config.DB.Model(&notification).Update("status", "read").Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update notification"})
		return
	}

	c.JSON(http.StatusOK, notification)
}
