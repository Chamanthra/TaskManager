package models

import "time"

type Notification struct {
	ID        uint   `gorm:"primaryKey"`
	Message   string `gorm:"not null"`
	Status    string `gorm:"default:unread"`
	CreatedAt time.Time
}
