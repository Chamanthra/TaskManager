package controllers

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/Chamanthra/TaskManager/config"
	"github.com/Chamanthra/TaskManager/models"
	"github.com/Chamanthra/TaskManager/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Register(c *gin.Context) {
	var input struct {
		UserName string `json:"user_name" binding:"required"`
		Password string `json:"password" binding:"required"`
		RoleID   uint   `json:"role_id"` // Optional
		Role     string `json:"role"`    // Optional
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input format", "details": err.Error()})
		return
	}

	input.UserName = strings.TrimSpace(input.UserName)
	input.Password = strings.TrimSpace(input.Password)

	if input.UserName == "" || input.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username and password are required"})
		return
	}

	// Check if username already exists
	var existing models.User
	if err := config.DB.Where("user_name = ?", input.UserName).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
		return
	}

	// Find role: by ID or by name
	var role models.Role
	if input.RoleID != 0 {
		if err := config.DB.First(&role, input.RoleID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID"})
			return
		}
	} else {
		// Default to "user" if role name not provided
		if input.Role == "" {
			input.Role = "user"
		}
		if err := config.DB.Where("role = ?", input.Role).First(&role).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role name"})
			return
		}
	}

	// Hash password
	hashed, err := utils.HashPassword(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Create user
	user := models.User{
		UserName:  input.UserName,
		Password:  hashed,
		RoleID:    role.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Registration failed", "details": err.Error()})
		return
	}

	// Load role data for response
	if err := config.DB.Preload("Role").First(&user, user.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user details", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Registration successful",
		"user": gin.H{
			"id":       user.ID,
			"username": user.UserName,
			"role":     user.Role.Name,
		},
	})
}

// Login handles user authentication
func Login(c *gin.Context) {
	var input struct {
		UserName string `json:"user_name" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input format", "details": err.Error()})
		return
	}

	// Lookup user with role preloaded
	var user models.User
	if err := config.DB.Preload("Role").Where("user_name = ?", input.UserName).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error", "details": err.Error()})
		}
		return
	}

	// Check password
	if !utils.CheckPassword(input.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// Generate JWT with role name
	token, err := utils.GenerateJWT(user.ID, user.Role.Name) // Use Role.Name instead of Role struct
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.UserName,
			"role":     user.Role.Name, // Return the role name instead of the struct
		},
	})
}
