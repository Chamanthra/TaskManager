package controllers

import (
	"errors"
	"net/http"

	"github.com/Chamanthra/TaskManagmentSystem/config"
	"github.com/Chamanthra/TaskManagmentSystem/models"
	"github.com/Chamanthra/TaskManagmentSystem/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

// GetUsers retrieves all users (admin only)
func GetUsers(c *gin.Context) {
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User information not found"})
		return
	}

	jwtClaims := claims.(jwt.MapClaims)
	role := jwtClaims["role"].(string)

	// Only allow admins to list users
	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only admins can list users"})
		return
	}

	var users []models.User
	if err := config.DB.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users", "details": err.Error()})
		return
	}

	// Don't return password hashes
	var sanitizedUsers []gin.H
	for _, user := range users {
		sanitizedUsers = append(sanitizedUsers, gin.H{
			"id":         user.ID,
			"username":   user.UserName,
			"role":       user.Role,
			"created_at": user.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, sanitizedUsers)
}

// DeleteUser handles user deletion (admin only)
func DeleteUser(c *gin.Context) {
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User information not found"})
		return
	}

	jwtClaims := claims.(jwt.MapClaims)
	role := jwtClaims["role"].(string)

	// Only allow admins to delete users
	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only admins can delete users"})
		return
	}

	userID := c.Param("id")
	var user models.User

	// First check if user exists
	if err := config.DB.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error", "details": err.Error()})
		}
		return
	}

	// Delete the user
	if err := config.DB.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// GetUserProfile returns the current user's profile
func GetUserProfile(c *gin.Context) {
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User information not found"})
		return
	}

	jwtClaims := claims.(jwt.MapClaims)
	userID := uint(jwtClaims["user_id"].(float64))

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error", "details": err.Error()})
		}
		return
	}

	// Return user profile without sensitive information
	c.JSON(http.StatusOK, gin.H{
		"id":         user.ID,
		"username":   user.UserName,
		"role":       user.Role,
		"created_at": user.CreatedAt,
	})
}

// UpdateUserProfile allows users to update their own profile
func UpdateUserProfile(c *gin.Context) {
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User information not found"})
		return
	}

	jwtClaims := claims.(jwt.MapClaims)
	userID := uint(jwtClaims["user_id"].(float64))

	var input struct {
		UserName string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input format", "details": err.Error()})
		return
	}

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Update username if provided
	if input.UserName != "" {
		// Check if username is already taken
		var existingUser models.User
		if err := config.DB.Where("user_name = ? AND id != ?", input.UserName, userID).First(&existingUser).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Username already taken"})
			return
		}
		user.UserName = input.UserName
	}

	// Update password if provided
	if input.Password != "" {
		hashed, err := utils.HashPassword(input.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}
		user.Password = hashed
	}

	if err := config.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Profile updated successfully",
		"user": gin.H{
			"id":       user.ID,
			"username": user.UserName,
			"role":     user.Role,
		},
	})
}
