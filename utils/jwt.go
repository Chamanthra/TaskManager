package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// GenerateJWT creates a JWT token for a user
func GenerateJWT(userID uint, role string) (string, error) {
	// Create a new JWT token with user information and expiration time
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(24 * time.Hour).Unix(), // 24 hours expiration
	})

	// Get JWT secret from environment variable
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", errors.New("JWT_SECRET is not set in environment")
	}

	// Sign and return the token
	return token.SignedString([]byte(secret))
}

// HashPassword hashes a password using bcrypt
func HashPassword(password string) (string, error) {
	// Hash the password with bcrypt using cost of 14 (higher value means more security)
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", err // Return the error if bcrypt fails
	}
	return string(bytes), nil
}

// CheckPassword compares a plain password to a bcrypt hash
func CheckPassword(password, hash string) bool {
	// Compare the hashed password with the plain-text password
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil // Return true if passwords match
}
