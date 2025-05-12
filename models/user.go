package models

import "time"

type User struct {
	ID        uint   `gorm:"primaryKey"`
	UserName  string `gorm:"unique;not null"`
	Password  string `gorm:"not null"`
	Role      string `gorm:"default:user"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
