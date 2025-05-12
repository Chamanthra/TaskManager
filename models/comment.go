package models

import "time"

type Comment struct {
	ID        uint   `gorm:"primaryKey"`
	Content   string `gorm:"not null"`
	TaskID    uint
	CreatedAt time.Time
}
