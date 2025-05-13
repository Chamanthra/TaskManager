package models

import (
	"time"

	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	Title       string    `json:"title" gorm:"not null;size:255"`
	Description string    `json:"description" gorm:"type:text"`
	Status      string    `json:"status" gorm:"type:enum('todo','in_progress','done','archived');default:'todo'"`
	Priority    string    `json:"priority" gorm:"type:enum('low','medium','high','critical');default:'medium'"`
	DueDate     time.Time `json:"due_date"`
	Notified    bool      `json:"notified" gorm:"default:false"` // Track if notification was sent
	UserID      uint      `json:"user_id" gorm:"index"`

	// Relationships
	User     User      `gorm:"foreignKey:UserID"`
	Comments []Comment `gorm:"foreignKey:TaskID"`
	Files    []File    `gorm:"foreignKey:TaskID"`
}
