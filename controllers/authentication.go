package controllers

import (
	"net/http"

	"github.com/Chamanthra/TaskManagmentSystem/config"
	"github.com/Chamanthra/TaskManagmentSystem/models"
	"github.com/Chamanthra/TaskManagmentSystem/utils"
	"github.com/gin-gonic/gin"
)

// Register handles user registration
func Register(c *gin.Context) {
	var input models.User

	// Validate input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	// Hash the password before storing
	input.Password = utils.HashPassword(input.Password)

	// Store user in DB
	if err := config.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Registration successful"})
}

// Login handles user authentication and JWT issuance
func Login(c *gin.Context) {
	var input models.User

	// Validate input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	var user models.User

	// Find user by username
	if err := config.DB.Where("user_name = ?", input.UserName).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Check password
	if !utils.CheckPassword(input.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(user.ID, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
