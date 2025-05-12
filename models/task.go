package models

import "time"

type Task struct {
	ID          uint   `gorm:"primaryKey"`
	Title       string `gorm:"not null"`
	Description string
	Status      string `gorm:"default:pending"`
	Priority    string
	Category    string
	DueDate     time.Time
	UserID      uint
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
